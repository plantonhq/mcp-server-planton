package domains

import (
	"strings"
	"testing"
)

func TestResolveKind_Known(t *testing.T) {
	kind, err := ResolveKind("AwsVpc")
	if err != nil {
		t.Fatalf("unexpected error for known kind AwsVpc: %v", err)
	}
	if kind == 0 {
		t.Fatal("expected non-zero enum value for AwsVpc")
	}
}

func TestResolveKind_Unknown(t *testing.T) {
	_, err := ResolveKind("CompletelyFakeKind")
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
	if !strings.Contains(err.Error(), "unknown cloud resource kind") {
		t.Fatalf("expected 'unknown cloud resource kind' in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "catalog") {
		t.Fatalf("expected catalog hint in error, got: %v", err)
	}
}

func TestResolveKind_Empty(t *testing.T) {
	_, err := ResolveKind("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}
