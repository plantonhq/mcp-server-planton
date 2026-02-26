package cloudresource

import (
	"strings"
	"testing"
)

func TestValidateIdentifier_IDOnly(t *testing.T) {
	err := validateIdentifier(ResourceIdentifier{ID: "res-abc123"})
	if err != nil {
		t.Fatalf("expected nil error for ID-only path, got: %v", err)
	}
}

func TestValidateIdentifier_SlugPathComplete(t *testing.T) {
	err := validateIdentifier(ResourceIdentifier{
		Kind: "AwsVpc",
		Org:  "acme",
		Env:  "prod",
		Slug: "my-vpc",
	})
	if err != nil {
		t.Fatalf("expected nil error for complete slug path, got: %v", err)
	}
}

func TestValidateIdentifier_BothPathsError(t *testing.T) {
	err := validateIdentifier(ResourceIdentifier{
		ID:   "res-abc123",
		Kind: "AwsVpc",
		Org:  "acme",
		Env:  "prod",
		Slug: "my-vpc",
	})
	if err == nil {
		t.Fatal("expected error when both ID and slug fields are set")
	}
	if !strings.Contains(err.Error(), "not both") {
		t.Fatalf("expected 'not both' in error, got: %v", err)
	}
}

func TestValidateIdentifier_IDWithPartialSlug(t *testing.T) {
	err := validateIdentifier(ResourceIdentifier{
		ID:  "res-abc123",
		Org: "acme",
	})
	if err == nil {
		t.Fatal("expected error when ID is set alongside a slug field")
	}
}

func TestValidateIdentifier_PartialSlugMissing(t *testing.T) {
	tests := []struct {
		name    string
		id      ResourceIdentifier
		missing []string
	}{
		{
			name:    "missing kind",
			id:      ResourceIdentifier{Org: "acme", Env: "prod", Slug: "my-vpc"},
			missing: []string{"kind"},
		},
		{
			name:    "missing org",
			id:      ResourceIdentifier{Kind: "AwsVpc", Env: "prod", Slug: "my-vpc"},
			missing: []string{"org"},
		},
		{
			name:    "missing env and slug",
			id:      ResourceIdentifier{Kind: "AwsVpc", Org: "acme"},
			missing: []string{"env", "slug"},
		},
		{
			name:    "only kind set",
			id:      ResourceIdentifier{Kind: "AwsVpc"},
			missing: []string{"org", "env", "slug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIdentifier(tt.id)
			if err == nil {
				t.Fatal("expected error for partial slug path")
			}
			for _, field := range tt.missing {
				if !strings.Contains(err.Error(), field) {
					t.Errorf("expected missing field %q in error, got: %v", field, err)
				}
			}
		})
	}
}

func TestValidateIdentifier_AllEmpty(t *testing.T) {
	err := validateIdentifier(ResourceIdentifier{})
	if err == nil {
		t.Fatal("expected error when all fields are empty")
	}
	if !strings.Contains(err.Error(), "provide either") {
		t.Fatalf("expected guidance message, got: %v", err)
	}
}

func TestDescribeIdentifier_IDPath(t *testing.T) {
	desc := describeIdentifier(ResourceIdentifier{ID: "res-abc123"})
	if !strings.Contains(desc, "res-abc123") {
		t.Fatalf("expected ID in description, got: %s", desc)
	}
}

func TestDescribeIdentifier_SlugPath(t *testing.T) {
	desc := describeIdentifier(ResourceIdentifier{
		Kind: "AwsVpc",
		Org:  "acme",
		Env:  "prod",
		Slug: "my-vpc",
	})
	if !strings.Contains(desc, "my-vpc") {
		t.Fatalf("expected slug in description, got: %s", desc)
	}
	if !strings.Contains(desc, "AwsVpc") {
		t.Fatalf("expected kind in description, got: %s", desc)
	}
	if !strings.Contains(desc, "acme") {
		t.Fatalf("expected org in description, got: %s", desc)
	}
	if !strings.Contains(desc, "prod") {
		t.Fatalf("expected env in description, got: %s", desc)
	}
}
