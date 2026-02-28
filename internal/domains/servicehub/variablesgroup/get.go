package variablesgroup

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	variablesgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/variablesgroup/v1"
	"google.golang.org/grpc"
)

// Get retrieves a variables group by ID or by org+slug via the
// VariablesGroupQueryController RPCs.
//
// Two identification paths are supported:
//   - ID path: calls Get(VariablesGroupId) directly.
//   - Slug path: calls GetByOrgBySlug(ApiResourceByOrgBySlugRequest).
func Get(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			group, err := resolveGroup(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}
			return domains.MarshalJSON(group)
		})
}

// resolveGroup fetches the full VariablesGroup proto by ID or by org+slug.
func resolveGroup(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (*variablesgroupv1.VariablesGroup, error) {
	client := variablesgroupv1.NewVariablesGroupQueryControllerClient(conn)

	if id != "" {
		resp, err := client.Get(ctx, &variablesgroupv1.VariablesGroupId{Value: id})
		if err != nil {
			return nil, domains.RPCError(err, fmt.Sprintf("variables group %q", id))
		}
		return resp, nil
	}

	resp, err := client.GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{
		Org:  org,
		Slug: slug,
	})
	if err != nil {
		return nil, domains.RPCError(err, fmt.Sprintf("variables group %q in org %q", slug, org))
	}
	return resp, nil
}

// resolveGroupID resolves identification inputs to a system-assigned
// variables group ID string. When an ID is already provided it is returned
// directly. Otherwise the group is fetched by org+slug and its metadata ID
// is extracted.
func resolveGroupID(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (string, error) {
	if id != "" {
		return id, nil
	}

	group, err := resolveGroup(ctx, conn, id, org, slug)
	if err != nil {
		return "", err
	}

	resourceID := group.GetMetadata().GetId()
	if resourceID == "" {
		return "", fmt.Errorf("resolved variables group %q in org %q but it has no ID — this indicates a backend issue", slug, org)
	}
	return resourceID, nil
}

// describeGroup returns a human-readable description of the variables group
// for use in error messages.
func describeGroup(id, org, slug string) string {
	if id != "" {
		return fmt.Sprintf("variables group %q", id)
	}
	return fmt.Sprintf("variables group %q in org %q", slug, org)
}
