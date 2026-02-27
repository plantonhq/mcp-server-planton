package cloudresource

import (
	"strings"
	"testing"
)

func TestExtractKindFromCloudObject_Valid(t *testing.T) {
	co := map[string]any{"kind": "AwsVpc"}
	kind, err := extractKindFromCloudObject(co)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kind != "AwsVpc" {
		t.Fatalf("expected AwsVpc, got %q", kind)
	}
}

func TestExtractKindFromCloudObject_MissingKey(t *testing.T) {
	co := map[string]any{"api_version": "v1"}
	_, err := extractKindFromCloudObject(co)
	if err == nil {
		t.Fatal("expected error for missing kind key")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Fatalf("expected 'missing' in error, got: %v", err)
	}
}

func TestExtractKindFromCloudObject_NonString(t *testing.T) {
	co := map[string]any{"kind": 42}
	_, err := extractKindFromCloudObject(co)
	if err == nil {
		t.Fatal("expected error for non-string kind")
	}
	if !strings.Contains(err.Error(), "non-empty string") {
		t.Fatalf("expected type error, got: %v", err)
	}
}

func TestExtractKindFromCloudObject_EmptyString(t *testing.T) {
	co := map[string]any{"kind": ""}
	_, err := extractKindFromCloudObject(co)
	if err == nil {
		t.Fatal("expected error for empty kind string")
	}
}

func TestResolveKind_Known(t *testing.T) {
	kind, err := resolveKind("AwsVpc")
	if err != nil {
		t.Fatalf("unexpected error for known kind AwsVpc: %v", err)
	}
	if kind == 0 {
		t.Fatal("expected non-zero enum value for AwsVpc")
	}
}

func TestResolveKind_Unknown(t *testing.T) {
	_, err := resolveKind("CompletelyFakeKind")
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
	if !strings.Contains(err.Error(), "unknown cloud resource kind") {
		t.Fatalf("expected 'unknown' in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "catalog") {
		t.Fatalf("expected catalog hint in error, got: %v", err)
	}
}
