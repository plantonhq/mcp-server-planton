package stackjob

import (
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// resolveOperationType
// ---------------------------------------------------------------------------

func TestResolveOperationType_ValidValues(t *testing.T) {
	valid := []string{"init", "refresh", "update_preview", "update", "destroy_preview", "destroy"}
	for _, v := range valid {
		got, err := resolveOperationType(v)
		if err != nil {
			t.Errorf("resolveOperationType(%q) unexpected error: %v", v, err)
		}
		if got == 0 {
			t.Errorf("resolveOperationType(%q) returned unspecified (0)", v)
		}
	}
}

func TestResolveOperationType_Unknown(t *testing.T) {
	_, err := resolveOperationType("bogus")
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

func TestResolveOperationType_Empty(t *testing.T) {
	_, err := resolveOperationType("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}

// ---------------------------------------------------------------------------
// resolveExecutionStatus
// ---------------------------------------------------------------------------

func TestResolveExecutionStatus_ValidValues(t *testing.T) {
	valid := []string{"queued", "running", "completed", "awaiting_approval"}
	for _, v := range valid {
		got, err := resolveExecutionStatus(v)
		if err != nil {
			t.Errorf("resolveExecutionStatus(%q) unexpected error: %v", v, err)
		}
		if got == 0 {
			t.Errorf("resolveExecutionStatus(%q) returned unspecified (0)", v)
		}
	}
}

func TestResolveExecutionStatus_Unknown(t *testing.T) {
	_, err := resolveExecutionStatus("bogus")
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

func TestResolveExecutionStatus_Empty(t *testing.T) {
	_, err := resolveExecutionStatus("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}

// ---------------------------------------------------------------------------
// resolveExecutionResult
// ---------------------------------------------------------------------------

func TestResolveExecutionResult_ValidValues(t *testing.T) {
	valid := []string{"tbd", "succeeded", "failed", "cancelled", "skipped", "discovered"}
	for _, v := range valid {
		got, err := resolveExecutionResult(v)
		if err != nil {
			t.Errorf("resolveExecutionResult(%q) unexpected error: %v", v, err)
		}
		if got == 0 {
			t.Errorf("resolveExecutionResult(%q) returned unspecified (0)", v)
		}
	}
}

func TestResolveExecutionResult_Unknown(t *testing.T) {
	_, err := resolveExecutionResult("bogus")
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

func TestResolveExecutionResult_Empty(t *testing.T) {
	_, err := resolveExecutionResult("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}
