package organization

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	organizationv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1"
	"google.golang.org/grpc"
)

// Delete removes an organization by ID via the
// OrganizationCommandController.Delete RPC.
func Delete(ctx context.Context, serverAddress, orgID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := organizationv1.NewOrganizationCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &organizationv1.OrganizationId{Value: orgID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("organization %q", orgID))
			}
			return domains.MarshalJSON(resp)
		})
}
