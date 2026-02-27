package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// buildRegistry creates a Registry from a list of parsed provider schemas.
// The registry maps each kind name to its cloud provider, API version,
// and the relative path to its JSON schema file.
func buildRegistry(schemas []*ProviderSchema) *Registry {
	reg := &Registry{
		Providers: make(map[string]RegistryEntry, len(schemas)),
	}
	for _, s := range schemas {
		reg.Providers[s.Kind] = RegistryEntry{
			CloudProvider: s.CloudProvider,
			APIVersion:    s.APIVersion,
			SchemaFile:    filepath.Join(s.CloudProvider, strings.ToLower(s.Kind)+".json"),
		}
	}
	return reg
}

// writeRegistry writes the provider registry as registry.json.
func writeRegistry(schemas []*ProviderSchema, baseDir string) error {
	dir := filepath.Join(baseDir, "providers")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return writeJSON(filepath.Join(dir, "registry.json"), buildRegistry(schemas))
}

// writeProviderSchema writes a single provider schema to disk.
func writeProviderSchema(schema *ProviderSchema, baseDir string) error {
	dir := filepath.Join(baseDir, "providers", schema.CloudProvider)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	filename := strings.ToLower(schema.Kind) + ".json"
	return writeJSON(filepath.Join(dir, filename), schema)
}

// writeMetadataSchema writes the shared metadata schema to disk.
func writeMetadataSchema(schema *MetadataSchema, baseDir string) error {
	dir := filepath.Join(baseDir, "shared")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return writeJSON(filepath.Join(dir, "metadata.json"), schema)
}

// sortSchemasByKind sorts provider schemas alphabetically by kind name
// for deterministic output ordering.
func sortSchemasByKind(schemas []*ProviderSchema) {
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].Kind < schemas[j].Kind
	})
}
