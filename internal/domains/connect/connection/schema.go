package connection

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"sync"

	"github.com/plantonhq/mcp-server-planton/schemas"
)

const schemaScheme = "connection-schema"

// registryEntry mirrors the JSON structure in connections/registry.json.
type registryEntry struct {
	Provider   string `json:"provider"`
	APIVersion string `json:"apiVersion"`
	SchemaFile string `json:"schemaFile"`
}

type registryData struct {
	Connections map[string]registryEntry `json:"connections"`
}

var (
	registryOnce sync.Once
	registryMap  map[string]registryEntry
	registryErr  error
)

func loadRegistry() (map[string]registryEntry, error) {
	registryOnce.Do(func() {
		data, err := schemas.ConnectionFS.ReadFile("connections/registry.json")
		if err != nil {
			registryErr = fmt.Errorf("reading embedded connection registry: %w", err)
			return
		}
		var rd registryData
		if err := json.Unmarshal(data, &rd); err != nil {
			registryErr = fmt.Errorf("parsing embedded connection registry: %w", err)
			return
		}
		registryMap = rd.Connections
	})
	return registryMap, registryErr
}

func loadConnectionSchema(kind string) ([]byte, error) {
	reg, err := loadRegistry()
	if err != nil {
		return nil, err
	}
	entry, ok := reg[kind]
	if !ok {
		return nil, fmt.Errorf("no schema found for connection kind %q", kind)
	}
	data, err := schemas.ConnectionFS.ReadFile("connections/" + entry.SchemaFile)
	if err != nil {
		return nil, fmt.Errorf("reading schema for %q: %w", kind, err)
	}
	return data, nil
}

func parseSchemaURI(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("malformed schema URI: %w", err)
	}
	if u.Scheme != schemaScheme {
		return "", fmt.Errorf("unexpected URI scheme %q, expected %q", u.Scheme, schemaScheme)
	}
	kind := u.Host
	if kind == "" {
		return "", fmt.Errorf("schema URI missing kind: %q", uri)
	}
	return kind, nil
}

// catalogProviderEntry groups connection types under a provider.
type catalogProviderEntry struct {
	APIVersion string   `json:"api_version"`
	Kinds      []string `json:"kinds"`
}

type connectionCatalog struct {
	SchemaURITemplate string                          `json:"schema_uri_template"`
	TotalTypes        int                             `json:"total_types"`
	Providers         map[string]catalogProviderEntry `json:"providers"`
}

var (
	catalogOnce sync.Once
	catalogJSON []byte
	catalogErr  error
)

func buildConnectionCatalog() ([]byte, error) {
	catalogOnce.Do(func() {
		reg, err := loadRegistry()
		if err != nil {
			catalogErr = fmt.Errorf("building connection catalog: %w", err)
			return
		}

		grouped := make(map[string]*catalogProviderEntry)
		for kind, entry := range reg {
			pe, ok := grouped[entry.Provider]
			if !ok {
				pe = &catalogProviderEntry{APIVersion: entry.APIVersion}
				grouped[entry.Provider] = pe
			}
			pe.Kinds = append(pe.Kinds, kind)
		}

		totalTypes := 0
		providers := make(map[string]catalogProviderEntry, len(grouped))
		for provider, pe := range grouped {
			sort.Strings(pe.Kinds)
			totalTypes += len(pe.Kinds)
			providers[provider] = *pe
		}

		cat := connectionCatalog{
			SchemaURITemplate: schemaScheme + "://{kind}",
			TotalTypes:        totalTypes,
			Providers:         providers,
		}

		catalogJSON, catalogErr = json.Marshal(cat)
		if catalogErr != nil {
			catalogErr = fmt.Errorf("marshaling connection catalog: %w", catalogErr)
		}
	})
	return catalogJSON, catalogErr
}
