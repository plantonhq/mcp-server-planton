package infraproject

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	infraprojectv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infraproject/v1"
	"google.golang.org/grpc"
)

// Get retrieves an infra project by ID or by org+slug via the
// InfraProjectQueryController RPCs.
//
// Two identification paths are supported:
//   - ID path: calls Get(InfraProjectId) directly.
//   - Slug path: calls GetByOrgBySlug(ApiResourceByOrgBySlugRequest).
func Get(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			project, err := resolveProject(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}
			return domains.MarshalJSON(project)
		})
}

// resolveProject fetches the full InfraProject proto by ID or by org+slug.
// Used by operations that need the full resource (e.g. undeploy passes the
// ID derived from the resolved project).
func resolveProject(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (*infraprojectv1.InfraProject, error) {
	client := infraprojectv1.NewInfraProjectQueryControllerClient(conn)

	if id != "" {
		resp, err := client.Get(ctx, &infraprojectv1.InfraProjectId{Value: id})
		if err != nil {
			return nil, domains.RPCError(err, fmt.Sprintf("infra project %q", id))
		}
		return resp, nil
	}

	resp, err := client.GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{
		Org:  org,
		Slug: slug,
	})
	if err != nil {
		return nil, domains.RPCError(err, fmt.Sprintf("infra project %q in org %q", slug, org))
	}
	return resp, nil
}

// resolveProjectID resolves identification inputs to a system-assigned project
// ID string. When an ID is already provided it is returned directly. Otherwise
// the project is fetched by org+slug and its metadata ID is extracted.
func resolveProjectID(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (string, error) {
	if id != "" {
		return id, nil
	}

	project, err := resolveProject(ctx, conn, id, org, slug)
	if err != nil {
		return "", err
	}

	resourceID := project.GetMetadata().GetId()
	if resourceID == "" {
		return "", fmt.Errorf("resolved infra project %q in org %q but it has no ID — this indicates a backend issue", slug, org)
	}
	return resourceID, nil
}

// describeProject returns a human-readable description of the project for
// use in error messages.
func describeProject(id, org, slug string) string {
	if id != "" {
		return fmt.Sprintf("infra project %q", id)
	}
	return fmt.Sprintf("infra project %q in org %q", slug, org)
}
