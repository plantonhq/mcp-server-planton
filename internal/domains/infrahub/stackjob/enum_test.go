package stackjob

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// operationTypeResolver
// ---------------------------------------------------------------------------

func TestOperationTypeResolver_ValidValues(t *testing.T) {
	valid := []string{"init", "refresh", "update_preview", "update", "destroy_preview", "destroy"}
	for _, v := range valid {
		got, err := operationTypeResolver.Resolve(v)
		if err != nil {
			t.Errorf("operationTypeResolver.Resolve(%q) unexpected error: %v", v, err)
		}
		if got == 0 {
			t.Errorf("operationTypeResolver.Resolve(%q) returned unspecified (0)", v)
		}
	}
}

func TestOperationTypeResolver_Unknown(t *testing.T) {
	_, err := operationTypeResolver.Resolve("bogus")
	if err == nil {
		t.Fatal("expected error for unknown operation type")
	}
	if !strings.Contains(err.Error(), "unknown stack job operation type") {
		t.Fatalf("expected 'unknown stack job operation type' in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "update") {
		t.Fatal("error should list valid values")
	}
}

func TestOperationTypeResolver_Empty(t *testing.T) {
	_, err := operationTypeResolver.Resolve("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}

// ---------------------------------------------------------------------------
// executionStatusResolver
// ---------------------------------------------------------------------------

func TestExecutionStatusResolver_ValidValues(t *testing.T) {
	valid := []string{"queued", "running", "completed", "awaiting_approval"}
	for _, v := range valid {
		got, err := executionStatusResolver.Resolve(v)
		if err != nil {
			t.Errorf("executionStatusResolver.Resolve(%q) unexpected error: %v", v, err)
		}
		if got == 0 {
			t.Errorf("executionStatusResolver.Resolve(%q) returned unspecified (0)", v)
		}
	}
}

func TestExecutionStatusResolver_Unknown(t *testing.T) {
	_, err := executionStatusResolver.Resolve("bogus")
	if err == nil {
		t.Fatal("expected error for unknown execution status")
	}
	if !strings.Contains(err.Error(), "unknown execution status") {
		t.Fatalf("expected 'unknown execution status' in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "running") {
		t.Fatal("error should list valid values")
	}
}

func TestExecutionStatusResolver_Empty(t *testing.T) {
	_, err := executionStatusResolver.Resolve("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}

// ---------------------------------------------------------------------------
// executionResultResolver
// ---------------------------------------------------------------------------

func TestExecutionResultResolver_ValidValues(t *testing.T) {
	valid := []string{"tbd", "succeeded", "failed", "cancelled", "skipped", "discovered"}
	for _, v := range valid {
		got, err := executionResultResolver.Resolve(v)
		if err != nil {
			t.Errorf("executionResultResolver.Resolve(%q) unexpected error: %v", v, err)
		}
		if got == 0 {
			t.Errorf("executionResultResolver.Resolve(%q) returned unspecified (0)", v)
		}
	}
}

func TestExecutionResultResolver_Unknown(t *testing.T) {
	_, err := executionResultResolver.Resolve("bogus")
	if err == nil {
		t.Fatal("expected error for unknown execution result")
	}
	if !strings.Contains(err.Error(), "unknown execution result") {
		t.Fatalf("expected 'unknown execution result' in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "failed") {
		t.Fatal("error should list valid values")
	}
}

func TestExecutionResultResolver_Empty(t *testing.T) {
	_, err := executionResultResolver.Resolve("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}
