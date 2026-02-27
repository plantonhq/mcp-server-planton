package cloudresource

import (
	"strings"
	"testing"
)

func TestResolveKinds_Nil(t *testing.T) {
	kinds, err := resolveKinds(nil)
	if err != nil {
		t.Fatalf("unexpected error for nil input: %v", err)
	}
	if kinds != nil {
		t.Fatalf("expected nil slice for nil input, got %v", kinds)
	}
}

func TestResolveKinds_Empty(t *testing.T) {
	kinds, err := resolveKinds([]string{})
	if err != nil {
		t.Fatalf("unexpected error for empty input: %v", err)
	}
	if kinds != nil {
		t.Fatalf("expected nil slice for empty input, got %v", kinds)
	}
}

func TestResolveKinds_ValidSingle(t *testing.T) {
	kinds, err := resolveKinds([]string{"AwsVpc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kinds) != 1 {
		t.Fatalf("expected 1 kind, got %d", len(kinds))
	}
	if kinds[0] == 0 {
		t.Fatal("expected non-zero enum value for AwsVpc")
	}
}

func TestResolveKinds_ValidMultiple(t *testing.T) {
	kinds, err := resolveKinds([]string{"AwsVpc", "GcpDnsZone"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kinds) != 2 {
		t.Fatalf("expected 2 kinds, got %d", len(kinds))
	}
	for i, k := range kinds {
		if k == 0 {
			t.Fatalf("expected non-zero enum value at index %d", i)
		}
	}
}

func TestResolveKinds_UnknownKind(t *testing.T) {
	_, err := resolveKinds([]string{"CompletelyFakeKind"})
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
	if !strings.Contains(err.Error(), "unknown cloud resource kind") {
		t.Fatalf("expected 'unknown cloud resource kind' in error, got: %v", err)
	}
}

func TestResolveKinds_MixedValidAndInvalid(t *testing.T) {
	_, err := resolveKinds([]string{"AwsVpc", "CompletelyFakeKind"})
	if err == nil {
		t.Fatal("expected error when one kind is invalid")
	}
	if !strings.Contains(err.Error(), "CompletelyFakeKind") {
		t.Fatalf("expected invalid kind name in error, got: %v", err)
	}
}
