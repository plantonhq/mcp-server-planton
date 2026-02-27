package cloudresource

import (
	"strings"
	"testing"

	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestBuildCloudResource_Valid(t *testing.T) {
	co := map[string]any{
		"api_version": "org.openmcf/v1",
		"kind":        "AwsVpc",
		"metadata": map[string]any{
			"name": "my-vpc",
			"org":  "acme",
			"env":  "prod",
		},
		"spec": map[string]any{"cidr": "10.0.0.0/16"},
	}
	normalized, err := structpb.NewStruct(map[string]any{"cidr": "10.0.0.0/16"})
	if err != nil {
		t.Fatalf("structpb.NewStruct: %v", err)
	}

	cr, err := buildCloudResource(co, "AwsVpc", normalized)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cr.ApiVersion != cloudResourceAPIVersion {
		t.Errorf("api_version = %q, want %q", cr.ApiVersion, cloudResourceAPIVersion)
	}
	if cr.Kind != cloudResourceKindConst {
		t.Errorf("kind = %q, want %q", cr.Kind, cloudResourceKindConst)
	}
	if cr.Metadata == nil {
		t.Fatal("expected non-nil metadata")
	}
	if cr.Metadata.Name != "my-vpc" {
		t.Errorf("metadata.name = %q, want %q", cr.Metadata.Name, "my-vpc")
	}
	if cr.Spec == nil {
		t.Fatal("expected non-nil spec")
	}
	if cr.Spec.CloudObject == nil {
		t.Fatal("expected non-nil spec.cloud_object")
	}
}

func TestBuildCloudResource_InvalidKind(t *testing.T) {
	co := map[string]any{
		"metadata": map[string]any{
			"name": "my-vpc",
			"org":  "acme",
			"env":  "prod",
		},
	}
	normalized, _ := structpb.NewStruct(map[string]any{})

	_, err := buildCloudResource(co, "CompletelyFakeKind", normalized)
	if err == nil {
		t.Fatal("expected error for invalid kind")
	}
	if !strings.Contains(err.Error(), "unknown cloud resource kind") {
		t.Fatalf("expected kind error, got: %v", err)
	}
}

func TestBuildCloudResource_MissingMetadata(t *testing.T) {
	co := map[string]any{
		"spec": map[string]any{"cidr": "10.0.0.0/16"},
	}
	normalized, _ := structpb.NewStruct(map[string]any{})

	_, err := buildCloudResource(co, "AwsVpc", normalized)
	if err == nil {
		t.Fatal("expected error for missing metadata")
	}
	if !strings.Contains(err.Error(), "metadata") {
		t.Fatalf("expected metadata error, got: %v", err)
	}
}

func TestDescribeResource_WithMetadata(t *testing.T) {
	cr := &cloudresourcev1.CloudResource{
		Metadata: &apiresource.ApiResourceMetadata{
			Name: "my-vpc",
			Org:  "acme",
			Env:  "prod",
		},
	}
	desc := describeResource(cr)
	if !strings.Contains(desc, "my-vpc") {
		t.Errorf("expected name in description, got: %s", desc)
	}
	if !strings.Contains(desc, "acme") {
		t.Errorf("expected org in description, got: %s", desc)
	}
}

func TestDescribeResource_NilMetadata(t *testing.T) {
	cr := &cloudresourcev1.CloudResource{}
	desc := describeResource(cr)
	if desc != "cloud resource" {
		t.Errorf("expected 'cloud resource', got: %s", desc)
	}
}

func TestDescribeResource_FallsBackToSlug(t *testing.T) {
	cr := &cloudresourcev1.CloudResource{
		Metadata: &apiresource.ApiResourceMetadata{
			Slug: "my-slug",
			Org:  "acme",
			Env:  "prod",
		},
	}
	desc := describeResource(cr)
	if !strings.Contains(desc, "my-slug") {
		t.Errorf("expected slug in description, got: %s", desc)
	}
}
