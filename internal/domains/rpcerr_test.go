package domains

import (
	"strings"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestClassifyCode_NotFound(t *testing.T) {
	msg := classifyCode(codes.NotFound, "cloud resource \"my-vpc\"", "not found")
	if !strings.Contains(msg, "not found") {
		t.Fatalf("expected 'not found', got: %s", msg)
	}
	if !strings.Contains(msg, "my-vpc") {
		t.Fatalf("expected resource desc in message, got: %s", msg)
	}
}

func TestClassifyCode_PermissionDenied(t *testing.T) {
	msg := classifyCode(codes.PermissionDenied, "cloud resource \"my-vpc\"", "denied")
	if !strings.Contains(msg, "Permission denied") {
		t.Fatalf("expected 'Permission denied', got: %s", msg)
	}
}

func TestClassifyCode_Unauthenticated(t *testing.T) {
	msg := classifyCode(codes.Unauthenticated, "resource", "bad token")
	if !strings.Contains(msg, "Authentication failed") {
		t.Fatalf("expected 'Authentication failed', got: %s", msg)
	}
}

func TestClassifyCode_Unavailable(t *testing.T) {
	msg := classifyCode(codes.Unavailable, "resource", "connection refused")
	if !strings.Contains(msg, "unavailable") {
		t.Fatalf("expected 'unavailable', got: %s", msg)
	}
}

func TestClassifyCode_DeadlineExceeded(t *testing.T) {
	msg := classifyCode(codes.DeadlineExceeded, "resource", "timeout")
	if !strings.Contains(msg, "timed out") {
		t.Fatalf("expected 'timed out', got: %s", msg)
	}
}

func TestClassifyCode_InvalidArgument(t *testing.T) {
	grpcMsg := "field 'name' is required"
	msg := classifyCode(codes.InvalidArgument, "resource", grpcMsg)
	if msg != grpcMsg {
		t.Fatalf("expected grpc message passthrough, got: %s", msg)
	}
}

func TestClassifyCode_Unknown(t *testing.T) {
	msg := classifyCode(codes.Internal, "resource", "something broke")
	if !strings.Contains(msg, "unexpected error") {
		t.Fatalf("expected 'unexpected error', got: %s", msg)
	}
	if !strings.Contains(msg, "something broke") {
		t.Fatalf("expected grpc message in output, got: %s", msg)
	}
}
