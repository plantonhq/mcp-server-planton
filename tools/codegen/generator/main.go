// Command generator reads JSON schemas produced by proto2schema (Stage 1) and
// generates Go input types with validation and structpb.Struct conversion for
// each OpenMCF cloud resource provider.
//
// Usage:
//
//	go run ./tools/codegen/generator/ --schemas-dir=schemas --output-dir=gen/infrahub/cloudresource
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ---------------------------------------------------------------------------
// Schema types — match the JSON structure from proto2schema/schema.go
// ---------------------------------------------------------------------------

type ProviderSchema struct {
	Name          string       `json:"name"`
	Kind          string       `json:"kind"`
	CloudProvider string       `json:"cloudProvider"`
	APIVersion    string       `json:"apiVersion"`
	Description   string       `json:"description,omitempty"`
	ProtoPackage  string       `json:"protoPackage"`
	ProtoFiles    ProtoFiles   `json:"protoFiles"`
	Spec          SpecSchema   `json:"spec"`
	NestedTypes   []TypeSchema `json:"nestedTypes,omitempty"`
}

type ProtoFiles struct {
	API  string `json:"api"`
	Spec string `json:"spec"`
}

type SpecSchema struct {
	Name   string         `json:"name"`
	Fields []*FieldSchema `json:"fields"`
}

type TypeSchema struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	ProtoType   string         `json:"protoType,omitempty"`
	Fields      []*FieldSchema `json:"fields"`
}

type FieldSchema struct {
	Name               string      `json:"name"`
	JSONName           string      `json:"jsonName"`
	ProtoField         string      `json:"protoField"`
	Type               TypeSpec    `json:"type"`
	Description        string      `json:"description,omitempty"`
	Required           bool        `json:"required"`
	Validation         *Validation `json:"validation,omitempty"`
	ReferenceKind      string      `json:"referenceKind,omitempty"`
	ReferenceFieldPath string      `json:"referenceFieldPath,omitempty"`
	Default            string      `json:"default,omitempty"`
	RecommendedDefault string      `json:"recommendedDefault,omitempty"`
	OneofGroup         string      `json:"oneofGroup,omitempty"`
}

type TypeSpec struct {
	Kind        string    `json:"kind"`
	KeyType     *TypeSpec `json:"keyType,omitempty"`
	ValueType   *TypeSpec `json:"valueType,omitempty"`
	ElementType *TypeSpec `json:"elementType,omitempty"`
	MessageType string    `json:"messageType,omitempty"`
	EnumType    string    `json:"enumType,omitempty"`
	EnumValues  []string  `json:"enumValues,omitempty"`
}

type Validation struct {
	Required  bool     `json:"required,omitempty"`
	MinLength int      `json:"minLength,omitempty"`
	MaxLength int      `json:"maxLength,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
	Min       int      `json:"min,omitempty"`
	Max       int      `json:"max,omitempty"`
	MinItems  int      `json:"minItems,omitempty"`
	MaxItems  int      `json:"maxItems,omitempty"`
	Enum      []string `json:"enum,omitempty"`
	Unique    bool     `json:"unique,omitempty"`
	Const     string   `json:"const,omitempty"`
}

type Registry struct {
	Providers map[string]RegistryEntry `json:"providers"`
}

type RegistryEntry struct {
	CloudProvider string `json:"cloudProvider"`
	APIVersion    string `json:"apiVersion"`
	SchemaFile    string `json:"schemaFile"`
}

// ---------------------------------------------------------------------------
// Main
// ---------------------------------------------------------------------------

func main() {
	schemasDir := flag.String("schemas-dir", "schemas", "path to schemas directory")
	outputDir := flag.String("output-dir", "gen/infrahub/cloudresource", "path to output directory for generated code")
	goModule := flag.String("module", "github.com/plantonhq/mcp-server-planton", "Go module path")
	flag.Parse()

	reg, err := loadRegistry(filepath.Join(*schemasDir, "providers", "registry.json"))
	if err != nil {
		log.Fatalf("load registry: %v", err)
	}

	// Group providers by cloud
	cloudProviders := groupByCloud(reg)
	clouds := sortedKeys(cloudProviders)

	// Load all provider schemas grouped by cloud
	type cloudGroup struct {
		cloud   string
		schemas []*ProviderSchema
	}
	var groups []cloudGroup
	for _, cloud := range clouds {
		entries := cloudProviders[cloud]
		var schemas []*ProviderSchema
		for _, entry := range entries {
			schemaPath := filepath.Join(*schemasDir, "providers", entry.SchemaFile)
			schema, err := loadProviderSchema(schemaPath)
			if err != nil {
				log.Fatalf("load schema %s: %v", schemaPath, err)
			}
			schemas = append(schemas, schema)
		}
		sort.Slice(schemas, func(i, j int) bool {
			return schemas[i].Kind < schemas[j].Kind
		})
		groups = append(groups, cloudGroup{cloud: cloud, schemas: schemas})
	}

	log.Printf("loaded %d providers across %d cloud platforms", len(reg.Providers), len(clouds))

	// Generate code for each cloud package
	var allProviders []registryItem
	for _, g := range groups {
		items, err := generateCloudPackage(*outputDir, *goModule, g.cloud, g.schemas)
		if err != nil {
			log.Fatalf("generate %s: %v", g.cloud, err)
		}
		allProviders = append(allProviders, items...)
	}

	// Generate the registry
	if err := generateRegistry(*outputDir, *goModule, clouds, allProviders); err != nil {
		log.Fatalf("generate registry: %v", err)
	}

	log.Printf("generated %d provider input types + registry in %s", len(allProviders), *outputDir)
}

// ---------------------------------------------------------------------------
// Schema loading
// ---------------------------------------------------------------------------

func loadRegistry(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &reg, nil
}

func loadProviderSchema(path string) (*ProviderSchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var schema ProviderSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &schema, nil
}

// ---------------------------------------------------------------------------
// Grouping helpers
// ---------------------------------------------------------------------------

type registryItem struct {
	kind     string
	cloud    string
	funcName string // e.g., "ParseAwsEksCluster"
}

func groupByCloud(reg *Registry) map[string][]RegistryEntry {
	m := make(map[string][]RegistryEntry)
	for _, entry := range reg.Providers {
		m[entry.CloudProvider] = append(m[entry.CloudProvider], entry)
	}
	return m
}

func sortedKeys(m map[string][]RegistryEntry) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// toSnakeFileName converts a PascalCase kind name to a snake_case file name.
// e.g., "AwsEksCluster" -> "awsekscluster"
func toSnakeFileName(kind string) string {
	return strings.ToLower(kind)
}
