package parse

import (
	"strings"
	"testing"
)

func TestValidateHeader_Valid(t *testing.T) {
	co := map[string]any{
		"api_version": "org.openmcf/v1",
		"kind":        "AwsVpc",
	}
	if err := ValidateHeader(co, "org.openmcf/v1", "AwsVpc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateHeader_MissingAPIVersion(t *testing.T) {
	co := map[string]any{"kind": "AwsVpc"}
	err := ValidateHeader(co, "org.openmcf/v1", "AwsVpc")
	if err == nil {
		t.Fatal("expected error for missing api_version")
	}
	if !strings.Contains(err.Error(), "api_version") {
		t.Fatalf("expected 'api_version' in error, got: %v", err)
	}
}

func TestValidateHeader_WrongAPIVersion(t *testing.T) {
	co := map[string]any{
		"api_version": "wrong/v2",
		"kind":        "AwsVpc",
	}
	err := ValidateHeader(co, "org.openmcf/v1", "AwsVpc")
	if err == nil {
		t.Fatal("expected error for wrong api_version")
	}
	if !strings.Contains(err.Error(), "expected api_version") {
		t.Fatalf("expected mismatch error, got: %v", err)
	}
}

func TestValidateHeader_MissingKind(t *testing.T) {
	co := map[string]any{"api_version": "org.openmcf/v1"}
	err := ValidateHeader(co, "org.openmcf/v1", "AwsVpc")
	if err == nil {
		t.Fatal("expected error for missing kind")
	}
	if !strings.Contains(err.Error(), "kind") {
		t.Fatalf("expected 'kind' in error, got: %v", err)
	}
}

func TestValidateHeader_WrongKind(t *testing.T) {
	co := map[string]any{
		"api_version": "org.openmcf/v1",
		"kind":        "GcpVpc",
	}
	err := ValidateHeader(co, "org.openmcf/v1", "AwsVpc")
	if err == nil {
		t.Fatal("expected error for wrong kind")
	}
	if !strings.Contains(err.Error(), "expected kind") {
		t.Fatalf("expected mismatch error, got: %v", err)
	}
}

func TestExtractSpecMap_Valid(t *testing.T) {
	co := map[string]any{
		"spec": map[string]any{"cidr": "10.0.0.0/16"},
	}
	spec, err := ExtractSpecMap(co)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec["cidr"] != "10.0.0.0/16" {
		t.Fatalf("expected cidr field, got: %v", spec)
	}
}

func TestExtractSpecMap_Missing(t *testing.T) {
	co := map[string]any{"kind": "AwsVpc"}
	_, err := ExtractSpecMap(co)
	if err == nil {
		t.Fatal("expected error for missing spec")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Fatalf("expected 'missing' in error, got: %v", err)
	}
}

func TestExtractSpecMap_NonObject(t *testing.T) {
	co := map[string]any{"spec": "not-a-map"}
	_, err := ExtractSpecMap(co)
	if err == nil {
		t.Fatal("expected error for non-object spec")
	}
	if !strings.Contains(err.Error(), "must be an object") {
		t.Fatalf("expected type error, got: %v", err)
	}
}

func TestRebuildCloudObject_PreservesFields(t *testing.T) {
	original := map[string]any{
		"api_version": "org.openmcf/v1",
		"kind":        "AwsVpc",
		"metadata":    map[string]any{"name": "test"},
		"spec":        map[string]any{"old": "value"},
	}
	normalized := map[string]any{"new": "value"}

	result, err := RebuildCloudObject(original, normalized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fields := result.GetFields()
	if fields["api_version"].GetStringValue() != "org.openmcf/v1" {
		t.Error("api_version not preserved")
	}
	if fields["kind"].GetStringValue() != "AwsVpc" {
		t.Error("kind not preserved")
	}

	specFields := fields["spec"].GetStructValue().GetFields()
	if specFields["new"].GetStringValue() != "value" {
		t.Error("normalized spec not applied")
	}
	if _, ok := specFields["old"]; ok {
		t.Error("old spec field should not be present")
	}
}

func TestRebuildCloudObject_DoesNotMutateOriginal(t *testing.T) {
	original := map[string]any{
		"api_version": "org.openmcf/v1",
		"spec":        map[string]any{"old": "value"},
	}
	normalized := map[string]any{"new": "value"}

	_, err := RebuildCloudObject(original, normalized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	origSpec := original["spec"].(map[string]any)
	if _, ok := origSpec["new"]; ok {
		t.Error("original map was mutated")
	}
}
