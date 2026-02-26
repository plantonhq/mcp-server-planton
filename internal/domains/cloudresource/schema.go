package cloudresource

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/plantoncloud/mcp-server-planton/schemas"
)

const schemaScheme = "cloud-resource-schema"

// registryEntry mirrors the JSON structure in providers/registry.json.
type registryEntry struct {
	CloudProvider string `json:"cloudProvider"`
	APIVersion    string `json:"apiVersion"`
	SchemaFile    string `json:"schemaFile"`
}

type registryData struct {
	Providers map[string]registryEntry `json:"providers"`
}

var (
	registryOnce sync.Once
	registryMap  map[string]registryEntry
	registryErr  error
)

// loadRegistry parses providers/registry.json from the embedded FS.
// The result is cached after the first successful load.
func loadRegistry() (map[string]registryEntry, error) {
	registryOnce.Do(func() {
		data, err := schemas.FS.ReadFile("providers/registry.json")
		if err != nil {
			registryErr = fmt.Errorf("reading embedded registry: %w", err)
			return
		}
		var rd registryData
		if err := json.Unmarshal(data, &rd); err != nil {
			registryErr = fmt.Errorf("parsing embedded registry: %w", err)
			return
		}
		registryMap = rd.Providers
	})
	return registryMap, registryErr
}

// loadProviderSchema reads the JSON schema for a specific kind from the
// embedded filesystem.
func loadProviderSchema(kind string) ([]byte, error) {
	reg, err := loadRegistry()
	if err != nil {
		return nil, err
	}
	entry, ok := reg[kind]
	if !ok {
		return nil, fmt.Errorf("no schema found for cloud resource kind %q", kind)
	}
	data, err := schemas.FS.ReadFile("providers/" + entry.SchemaFile)
	if err != nil {
		return nil, fmt.Errorf("reading schema for %q: %w", kind, err)
	}
	return data, nil
}

// parseSchemaURI extracts the kind parameter from a cloud-resource-schema://
// URI. For example, "cloud-resource-schema://AwsEksCluster" returns "AwsEksCluster".
func parseSchemaURI(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("malformed schema URI: %w", err)
	}
	if u.Scheme != schemaScheme {
		return "", fmt.Errorf("unexpected URI scheme %q, expected %q", u.Scheme, schemaScheme)
	}
	// In cloud-resource-schema://AwsEksCluster the standard library parses:
	//   Host = "AwsEksCluster"
	kind := u.Host
	if kind == "" {
		return "", fmt.Errorf("schema URI missing kind: %q", uri)
	}
	return kind, nil
}
