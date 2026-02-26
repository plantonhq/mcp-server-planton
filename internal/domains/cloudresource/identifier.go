package cloudresource

import (
	"context"
	"fmt"
	"strings"

	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"github.com/plantoncloud/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ResourceIdentifier provides two mutually exclusive ways to identify a cloud
// resource: by system-assigned ID, or by the composite key (kind, org, env, slug).
//
// Exactly one path must be fully specified:
//   - ID path: only ID is set
//   - Slug path: all of Kind, Org, Env, Slug are set
type ResourceIdentifier struct {
	ID   string
	Kind string
	Org  string
	Env  string
	Slug string
}

// validateIdentifier checks that exactly one identification path is fully
// specified. Returns nil when valid.
func validateIdentifier(id ResourceIdentifier) error {
	hasID := id.ID != ""
	slugFields := [4]string{id.Kind, id.Org, id.Env, id.Slug}
	populated := 0
	for _, f := range slugFields {
		if f != "" {
			populated++
		}
	}

	switch {
	case hasID && populated > 0:
		return fmt.Errorf("provide either 'id' alone or all of 'kind', 'org', 'env', and 'slug' — not both")
	case hasID:
		return nil
	case populated == 4:
		return nil
	case populated > 0:
		var missing []string
		for _, pair := range []struct {
			val, name string
		}{
			{id.Kind, "kind"},
			{id.Org, "org"},
			{id.Env, "env"},
			{id.Slug, "slug"},
		} {
			if pair.val == "" {
				missing = append(missing, pair.name)
			}
		}
		return fmt.Errorf(
			"when using the slug path, all of 'kind', 'org', 'env', and 'slug' are required — missing: %s",
			strings.Join(missing, ", "),
		)
	default:
		return fmt.Errorf("provide either 'id' or all of 'kind', 'org', 'env', and 'slug' to identify the cloud resource")
	}
}

// describeIdentifier returns a human-readable description of the resource for
// use in error messages and log entries.
func describeIdentifier(id ResourceIdentifier) string {
	if id.ID != "" {
		return fmt.Sprintf("cloud resource %q", id.ID)
	}
	return fmt.Sprintf("cloud resource %q (kind=%s) in org %q env %q", id.Slug, id.Kind, id.Org, id.Env)
}

// resolveResourceID resolves a ResourceIdentifier to a system-assigned resource
// ID string. If the identifier already carries an ID, it is returned directly.
// Otherwise the slug-path fields are used to look up the resource via the query
// controller and the ID is extracted from the response metadata.
//
// Errors returned are already user-friendly (kind validation errors are passed
// through; gRPC lookup errors are classified via [domains.RPCError]).
func resolveResourceID(ctx context.Context, conn *grpc.ClientConn, id ResourceIdentifier) (string, error) {
	if id.ID != "" {
		return id.ID, nil
	}

	kind, err := resolveKind(id.Kind)
	if err != nil {
		return "", err
	}

	desc := describeIdentifier(id)
	client := cloudresourcev1.NewCloudResourceQueryControllerClient(conn)
	cr, err := client.GetByOrgByEnvByKindBySlug(ctx, &cloudresourcev1.CloudResourceByOrgByEnvByKindBySlugRequest{
		Org:               id.Org,
		Env:               id.Env,
		CloudResourceKind: kind,
		Slug:              id.Slug,
	})
	if err != nil {
		return "", domains.RPCError(err, desc)
	}

	resourceID := cr.GetMetadata().GetId()
	if resourceID == "" {
		return "", fmt.Errorf("resolved %s but it has no ID — this indicates a backend issue", desc)
	}
	return resourceID, nil
}
