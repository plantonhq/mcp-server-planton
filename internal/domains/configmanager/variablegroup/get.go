package variablegroup

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	variablegroupv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/variablegroup/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// Get retrieves a variable group by ID or by org+scope+slug.
func Get(ctx context.Context, serverAddress, id, org string, scope variablegroupv1.VariableGroupSpec_Scope, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			vg, err := resolveVariableGroup(ctx, conn, id, org, scope, slug)
			if err != nil {
				return "", err
			}
			return domains.MarshalJSON(vg)
		})
}

// resolveVariableGroup fetches the full VariableGroup proto by ID or by org+scope+slug.
func resolveVariableGroup(ctx context.Context, conn *grpc.ClientConn, id, org string, scope variablegroupv1.VariableGroupSpec_Scope, slug string) (*variablegroupv1.VariableGroup, error) {
	client := variablegroupv1.NewVariableGroupQueryControllerClient(conn)

	if id != "" {
		resp, err := client.Get(ctx, &variablegroupv1.VariableGroupId{Value: id})
		if err != nil {
			return nil, domains.RPCError(err, fmt.Sprintf("variable group %q", id))
		}
		return resp, nil
	}

	resp, err := client.GetByOrgByScopeBySlug(ctx, &variablegroupv1.VariableGroupScopeSlugRequest{
		Org:   org,
		Scope: scope,
		Slug:  slug,
	})
	if err != nil {
		return nil, domains.RPCError(err, fmt.Sprintf("variable group %q (scope=%s) in org %q", slug, scope, org))
	}
	return resp, nil
}

// resolveVariableGroupID resolves identification inputs to a system-assigned ID.
func resolveVariableGroupID(ctx context.Context, conn *grpc.ClientConn, id, org string, scope variablegroupv1.VariableGroupSpec_Scope, slug string) (string, error) {
	if id != "" {
		return id, nil
	}

	vg, err := resolveVariableGroup(ctx, conn, id, org, scope, slug)
	if err != nil {
		return "", err
	}

	resourceID := vg.GetMetadata().GetId()
	if resourceID == "" {
		return "", fmt.Errorf("resolved variable group %q (scope=%s) in org %q but it has no ID — this indicates a backend issue", slug, scope, org)
	}
	return resourceID, nil
}

// describeVariableGroup returns a human-readable description for error messages.
func describeVariableGroup(id, org string, scope variablegroupv1.VariableGroupSpec_Scope, slug string) string {
	if id != "" {
		return fmt.Sprintf("variable group %q", id)
	}
	return fmt.Sprintf("variable group %q (scope=%s) in org %q", slug, scope, org)
}
