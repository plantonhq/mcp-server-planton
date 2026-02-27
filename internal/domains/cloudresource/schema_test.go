package cloudresource

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseSchemaURI_Valid(t *testing.T) {
	kind, err := parseSchemaURI("cloud-resource-schema://AwsEksCluster")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kind != "AwsEksCluster" {
		t.Fatalf("expected AwsEksCluster, got %q", kind)
	}
}

func TestParseSchemaURI_WrongScheme(t *testing.T) {
	_, err := parseSchemaURI("https://AwsEksCluster")
	if err == nil {
		t.Fatal("expected error for wrong scheme")
	}
	if !strings.Contains(err.Error(), "unexpected URI scheme") {
		t.Fatalf("expected scheme error, got: %v", err)
	}
}

func TestParseSchemaURI_MissingKind(t *testing.T) {
	_, err := parseSchemaURI("cloud-resource-schema://")
	if err == nil {
		t.Fatal("expected error for missing kind")
	}
	if !strings.Contains(err.Error(), "missing kind") {
		t.Fatalf("expected missing kind error, got: %v", err)
	}
}

func TestParseSchemaURI_MalformedURI(t *testing.T) {
	_, err := parseSchemaURI("://")
	if err == nil {
		t.Fatal("expected error for malformed URI")
	}
}

func TestLoadRegistry_ReturnsNonEmpty(t *testing.T) {
	reg, err := loadRegistry()
	if err != nil {
		t.Fatalf("loadRegistry failed: %v", err)
	}
	if len(reg) == 0 {
		t.Fatal("expected non-empty registry")
	}
}

func TestLoadRegistry_ContainsKnownKind(t *testing.T) {
	reg, err := loadRegistry()
	if err != nil {
		t.Fatalf("loadRegistry failed: %v", err)
	}
	entry, ok := reg["AwsVpc"]
	if !ok {
		t.Fatal("expected AwsVpc in registry")
	}
	if entry.CloudProvider == "" {
		t.Fatal("expected non-empty cloud provider for AwsVpc")
	}
	if entry.SchemaFile == "" {
		t.Fatal("expected non-empty schema file for AwsVpc")
	}
}

func TestLoadProviderSchema_KnownKind(t *testing.T) {
	data, err := loadProviderSchema("AwsVpc")
	if err != nil {
		t.Fatalf("loadProviderSchema(AwsVpc) failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty schema data")
	}
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("schema is not valid JSON: %v", err)
	}
}

func TestLoadProviderSchema_UnknownKind(t *testing.T) {
	_, err := loadProviderSchema("CompletelyFakeKind")
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
	if !strings.Contains(err.Error(), "no schema found") {
		t.Fatalf("expected 'no schema found' in error, got: %v", err)
	}
}

func TestBuildKindCatalog_ValidJSON(t *testing.T) {
	data, err := buildKindCatalog()
	if err != nil {
		t.Fatalf("buildKindCatalog failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty catalog")
	}

	var cat kindCatalog
	if err := json.Unmarshal(data, &cat); err != nil {
		t.Fatalf("catalog is not valid JSON: %v", err)
	}

	if cat.TotalKinds == 0 {
		t.Fatal("expected non-zero total_kinds")
	}
	if len(cat.Providers) == 0 {
		t.Fatal("expected non-empty providers map")
	}
	if cat.SchemaURITemplate != "cloud-resource-schema://{kind}" {
		t.Fatalf("unexpected schema_uri_template: %q", cat.SchemaURITemplate)
	}
}

func TestBuildKindCatalog_KindsSorted(t *testing.T) {
	data, err := buildKindCatalog()
	if err != nil {
		t.Fatalf("buildKindCatalog failed: %v", err)
	}

	var cat kindCatalog
	if err := json.Unmarshal(data, &cat); err != nil {
		t.Fatalf("catalog is not valid JSON: %v", err)
	}

	for provider, entry := range cat.Providers {
		for i := 1; i < len(entry.Kinds); i++ {
			if entry.Kinds[i] < entry.Kinds[i-1] {
				t.Errorf("provider %q kinds not sorted: %q comes after %q",
					provider, entry.Kinds[i], entry.Kinds[i-1])
			}
		}
	}
}
