package cloudresource

import (
	"strings"
	"testing"
)

func validCloudObject() map[string]any {
	return map[string]any{
		"metadata": map[string]any{
			"name": "my-resource",
			"org":  "acme",
			"env":  "prod",
		},
	}
}

func TestExtractMetadata_ValidRequired(t *testing.T) {
	md, err := extractMetadata(validCloudObject())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if md.Name != "my-resource" {
		t.Errorf("name = %q, want %q", md.Name, "my-resource")
	}
	if md.Org != "acme" {
		t.Errorf("org = %q, want %q", md.Org, "acme")
	}
	if md.Env != "prod" {
		t.Errorf("env = %q, want %q", md.Env, "prod")
	}
}

func TestExtractMetadata_MissingMetadataKey(t *testing.T) {
	co := map[string]any{"kind": "AwsVpc"}
	_, err := extractMetadata(co)
	if err == nil {
		t.Fatal("expected error for missing metadata key")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Fatalf("expected 'missing' in error, got: %v", err)
	}
}

func TestExtractMetadata_NonObjectMetadata(t *testing.T) {
	co := map[string]any{"metadata": "not-an-object"}
	_, err := extractMetadata(co)
	if err == nil {
		t.Fatal("expected error for non-object metadata")
	}
	if !strings.Contains(err.Error(), "must be an object") {
		t.Fatalf("expected type error, got: %v", err)
	}
}

func TestExtractMetadata_MissingName(t *testing.T) {
	co := map[string]any{
		"metadata": map[string]any{
			"org": "acme",
			"env": "prod",
		},
	}
	_, err := extractMetadata(co)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "metadata.name") {
		t.Fatalf("expected field-specific error, got: %v", err)
	}
}

func TestExtractMetadata_MissingOrg(t *testing.T) {
	co := map[string]any{
		"metadata": map[string]any{
			"name": "my-resource",
			"env":  "prod",
		},
	}
	_, err := extractMetadata(co)
	if err == nil {
		t.Fatal("expected error for missing org")
	}
	if !strings.Contains(err.Error(), "metadata.org") {
		t.Fatalf("expected field-specific error, got: %v", err)
	}
}

func TestExtractMetadata_MissingEnv(t *testing.T) {
	co := map[string]any{
		"metadata": map[string]any{
			"name": "my-resource",
			"org":  "acme",
		},
	}
	_, err := extractMetadata(co)
	if err == nil {
		t.Fatal("expected error for missing env")
	}
	if !strings.Contains(err.Error(), "metadata.env") {
		t.Fatalf("expected field-specific error, got: %v", err)
	}
}

func TestExtractMetadata_OptionalFieldsPresent(t *testing.T) {
	co := map[string]any{
		"metadata": map[string]any{
			"name":        "my-resource",
			"org":         "acme",
			"env":         "prod",
			"slug":        "my-slug",
			"id":          "res-123",
			"labels":      map[string]any{"team": "platform"},
			"annotations": map[string]any{"note": "important"},
			"tags":        []any{"infra", "production"},
			"version":     map[string]any{"message": "initial"},
		},
	}

	md, err := extractMetadata(co)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if md.Slug != "my-slug" {
		t.Errorf("slug = %q, want %q", md.Slug, "my-slug")
	}
	if md.Id != "res-123" {
		t.Errorf("id = %q, want %q", md.Id, "res-123")
	}
	if md.Labels["team"] != "platform" {
		t.Errorf("labels[team] = %q, want %q", md.Labels["team"], "platform")
	}
	if md.Annotations["note"] != "important" {
		t.Errorf("annotations[note] = %q, want %q", md.Annotations["note"], "important")
	}
	if len(md.Tags) != 2 || md.Tags[0] != "infra" || md.Tags[1] != "production" {
		t.Errorf("tags = %v, want [infra production]", md.Tags)
	}
	if md.Version == nil || md.Version.Message != "initial" {
		t.Errorf("version.message = %v, want %q", md.Version, "initial")
	}
}

func TestExtractMetadata_OptionalFieldsAbsent(t *testing.T) {
	md, err := extractMetadata(validCloudObject())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if md.Slug != "" {
		t.Errorf("slug should be empty, got %q", md.Slug)
	}
	if md.Id != "" {
		t.Errorf("id should be empty, got %q", md.Id)
	}
	if len(md.Labels) != 0 {
		t.Errorf("labels should be empty, got %v", md.Labels)
	}
	if len(md.Tags) != 0 {
		t.Errorf("tags should be empty, got %v", md.Tags)
	}
	if md.Version != nil {
		t.Errorf("version should be nil, got %v", md.Version)
	}
}

func TestToStringMap_SkipsNonStrings(t *testing.T) {
	m := map[string]any{
		"a": "hello",
		"b": 42,
		"c": true,
		"d": "world",
	}
	result := toStringMap(m)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(result), result)
	}
	if result["a"] != "hello" || result["d"] != "world" {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestToStringSlice_SkipsNonStrings(t *testing.T) {
	a := []any{"hello", 42, true, "world"}
	result := toStringSlice(a)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(result), result)
	}
	if result[0] != "hello" || result[1] != "world" {
		t.Fatalf("unexpected result: %v", result)
	}
}
