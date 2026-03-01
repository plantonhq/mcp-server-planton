package environment

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	environmentv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/environment/v1"
	organizationv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1"
	"google.golang.org/grpc"
)

// Get retrieves a single environment by ID via the
// EnvironmentQueryController.Get RPC.
func Get(ctx context.Context, serverAddress, envID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := environmentv1.NewEnvironmentQueryControllerClient(conn)
			resp, err := client.Get(ctx, &environmentv1.EnvironmentId{Value: envID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("environment %q", envID))
			}
			return domains.MarshalJSON(resp)
		})
}

// GetByOrgBySlug retrieves a single environment by organization and slug via
// the EnvironmentQueryController.GetByOrgBySlug RPC.
func GetByOrgBySlug(ctx context.Context, serverAddress, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := environmentv1.NewEnvironmentQueryControllerClient(conn)
			resp, err := client.GetByOrgBySlug(ctx, &organizationv1.ByOrgBySlugRequest{
				Org:  org,
				Slug: slug,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("environment %q in org %q", slug, org))
			}
			return domains.MarshalJSON(resp)
		})
}
