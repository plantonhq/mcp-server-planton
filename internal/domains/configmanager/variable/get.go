package variable

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	variablev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/variable/v1"
	"google.golang.org/grpc"
)

// Get retrieves a variable by ID or by org+scope+slug via the
// VariableQueryController RPCs.
//
// Two identification paths are supported:
//   - ID path: calls Get(VariableId) directly.
//   - Slug path: calls GetByOrgByScopeBySlug(VariableScopeSlugRequest).
func Get(ctx context.Context, serverAddress, id, org string, scope variablev1.VariableSpec_Scope, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			v, err := resolveVariable(ctx, conn, id, org, scope, slug)
			if err != nil {
				return "", err
			}
			return domains.MarshalJSON(v)
		})
}

// resolveVariable fetches the full Variable proto by ID or by org+scope+slug.
// Used by operations that need the full resource before acting on it.
func resolveVariable(ctx context.Context, conn *grpc.ClientConn, id, org string, scope variablev1.VariableSpec_Scope, slug string) (*variablev1.Variable, error) {
	client := variablev1.NewVariableQueryControllerClient(conn)

	if id != "" {
		resp, err := client.Get(ctx, &variablev1.VariableId{Value: id})
		if err != nil {
			return nil, domains.RPCError(err, fmt.Sprintf("variable %q", id))
		}
		return resp, nil
	}

	resp, err := client.GetByOrgByScopeBySlug(ctx, &variablev1.VariableScopeSlugRequest{
		Org:   org,
		Scope: scope,
		Slug:  slug,
	})
	if err != nil {
		return nil, domains.RPCError(err, fmt.Sprintf("variable %q (scope=%s) in org %q", slug, scope, org))
	}
	return resp, nil
}

// resolveVariableID resolves identification inputs to a system-assigned
// variable ID string. When an ID is already provided it is returned directly.
// Otherwise the variable is fetched by org+scope+slug and its metadata ID is
// extracted.
func resolveVariableID(ctx context.Context, conn *grpc.ClientConn, id, org string, scope variablev1.VariableSpec_Scope, slug string) (string, error) {
	if id != "" {
		return id, nil
	}

	v, err := resolveVariable(ctx, conn, id, org, scope, slug)
	if err != nil {
		return "", err
	}

	resourceID := v.GetMetadata().GetId()
	if resourceID == "" {
		return "", fmt.Errorf("resolved variable %q (scope=%s) in org %q but it has no ID — this indicates a backend issue", slug, scope, org)
	}
	return resourceID, nil
}

// describeVariable returns a human-readable description for error messages.
func describeVariable(id, org string, scope variablev1.VariableSpec_Scope, slug string) string {
	if id != "" {
		return fmt.Sprintf("variable %q", id)
	}
	return fmt.Sprintf("variable %q (scope=%s) in org %q", slug, scope, org)
}
