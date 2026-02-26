// proto2schema converts OpenMCF provider Protocol Buffer definitions to
// JSON schemas for code generation and MCP resource template discovery.
//
// This is Stage 1 of the two-stage codegen pipeline adapted from
// Stigmer's architecture. Stage 2 (generator) consumes these schemas
// to produce Go input types with ToProto() methods.
//
// Usage:
//
//	go run tools/codegen/proto2schema/main.go --all
//	go run tools/codegen/proto2schema/main.go --provider aws/awsalb
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	all := flag.Bool("all", false, "Generate schemas for all providers")
	provider := flag.String("provider", "", "Generate schema for a single provider (e.g., aws/awsalb)")
	openmcfAPIsDir := flag.String("openmcf-apis-dir", "", "Path to openmcf/apis directory (default: $SCM_ROOT/github.com/plantonhq/openmcf/apis)")
	outputDir := flag.String("output-dir", "tools/codegen/schemas", "Output directory for generated schemas")
	flag.Parse()

	if !*all && *provider == "" {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, "  proto2schema --all [--openmcf-apis-dir <dir>] [--output-dir <dir>]")
		fmt.Fprintln(os.Stderr, "  proto2schema --provider aws/awsalb [--openmcf-apis-dir <dir>] [--output-dir <dir>]")
		os.Exit(1)
	}

	if *openmcfAPIsDir == "" {
		*openmcfAPIsDir = defaultOpenMCFAPIsDir()
	}

	if _, err := os.Stat(*openmcfAPIsDir); os.IsNotExist(err) {
		log.Fatalf("OpenMCF APIs directory not found: %s\n"+
			"Ensure the openmcf repo is cloned, or set SCM_ROOT / --openmcf-apis-dir.", *openmcfAPIsDir)
	}

	importPaths := buildImportPaths(*openmcfAPIsDir)

	p := NewParser(importPaths)
	if err := p.Init(); err != nil {
		log.Fatalf("Parser init failed: %v", err)
	}

	if *all {
		runAll(p, *openmcfAPIsDir, *outputDir)
	} else {
		runSingle(p, *provider, *outputDir)
	}
}

// runAll discovers all provider directories and generates schemas for each.
func runAll(p *Parser, openmcfAPIsDir, outputDir string) {
	providerBaseDir := filepath.Join(openmcfAPIsDir, "org", "openmcf", "provider")

	cloudDirs, err := os.ReadDir(providerBaseDir)
	if err != nil {
		log.Fatalf("Cannot read provider directory %s: %v", providerBaseDir, err)
	}

	var schemas []*ProviderSchema
	var parseErrors []string

	for _, cloudDir := range cloudDirs {
		if !cloudDir.IsDir() || strings.HasPrefix(cloudDir.Name(), "_") {
			continue
		}
		cloudName := cloudDir.Name()

		resourceDirs, err := os.ReadDir(filepath.Join(providerBaseDir, cloudName))
		if err != nil {
			log.Printf("Warning: cannot read %s/%s: %v", providerBaseDir, cloudName, err)
			continue
		}

		for _, resDir := range resourceDirs {
			if !resDir.IsDir() {
				continue
			}
			resourceName := resDir.Name()

			apiPath := filepath.Join(providerBaseDir, cloudName, resourceName, "v1", "api.proto")
			if _, err := os.Stat(apiPath); os.IsNotExist(err) {
				continue
			}

			schema, err := p.ParseProvider(cloudName, resourceName)
			if err != nil {
				parseErrors = append(parseErrors, fmt.Sprintf("%s/%s: %v", cloudName, resourceName, err))
				continue
			}

			schemas = append(schemas, schema)
			fmt.Printf("  %s/%s -> %s\n", cloudName, resourceName, schema.Kind)
		}
	}

	sortSchemasByKind(schemas)

	for _, schema := range schemas {
		if err := writeProviderSchema(schema, outputDir); err != nil {
			log.Printf("Error writing schema for %s: %v", schema.Kind, err)
		}
	}

	if err := writeRegistry(schemas, outputDir); err != nil {
		log.Printf("Error writing registry: %v", err)
	}

	metadata, err := p.ParseMetadata()
	if err != nil {
		log.Printf("Error parsing metadata: %v", err)
	} else {
		if err := writeMetadataSchema(metadata, outputDir); err != nil {
			log.Printf("Error writing metadata schema: %v", err)
		}
	}

	fmt.Printf("\nGenerated %d provider schemas in %s\n", len(schemas), filepath.Join(outputDir, "providers"))
	if len(parseErrors) > 0 {
		fmt.Printf("Errors (%d):\n", len(parseErrors))
		for _, e := range parseErrors {
			fmt.Printf("  - %s\n", e)
		}
	}
}

// runSingle generates a schema for a single provider specified as "cloud/resource".
func runSingle(p *Parser, provider, outputDir string) {
	parts := strings.SplitN(provider, "/", 2)
	if len(parts) != 2 {
		log.Fatalf("Invalid provider format: %q (expected cloud/resource, e.g., aws/awsalb)", provider)
	}

	schema, err := p.ParseProvider(parts[0], parts[1])
	if err != nil {
		log.Fatalf("Error parsing provider %s: %v", provider, err)
	}

	if err := writeProviderSchema(schema, outputDir); err != nil {
		log.Fatalf("Error writing schema: %v", err)
	}

	fmt.Printf("Generated schema: providers/%s/%s.json\n", schema.CloudProvider, strings.ToLower(schema.Kind))
}

// defaultOpenMCFAPIsDir returns the default path to the openmcf/apis directory
// using the SCM_ROOT convention ($HOME/scm/github.com/{org}/{repo}/apis).
func defaultOpenMCFAPIsDir() string {
	scmRoot := os.Getenv("SCM_ROOT")
	if scmRoot == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Cannot determine home directory: %v", err)
		}
		scmRoot = filepath.Join(home, "scm")
	}
	return filepath.Join(scmRoot, "github.com", "plantonhq", "openmcf", "apis")
}

// buildImportPaths constructs the proto import paths needed for parsing.
// Includes the openmcf/apis directory and the buf module cache (for buf.validate).
func buildImportPaths(openmcfAPIsDir string) []string {
	paths := []string{openmcfAPIsDir}

	home, err := os.UserHomeDir()
	if err != nil {
		return paths
	}

	bufCachePath := filepath.Join(home, ".cache", "buf", "v3", "modules", "b5",
		"buf.build", "bufbuild", "protovalidate")

	entries, err := os.ReadDir(bufCachePath)
	if err != nil {
		log.Printf("Warning: buf cache not found at %s", bufCachePath)
		log.Printf("Run 'buf dep update' in the openmcf repo to populate the cache.")
		return paths
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		filesPath := filepath.Join(bufCachePath, entry.Name(), "files")
		if _, err := os.Stat(filesPath); err == nil {
			paths = append([]string{filesPath}, paths...)
			break
		}
	}

	return paths
}

// writeJSON marshals v as indented JSON and writes it to path.
func writeJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}
