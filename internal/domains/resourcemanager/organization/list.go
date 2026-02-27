package organization

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/protobuf"
	organizationv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1"
	"google.golang.org/grpc"
)

// List retrieves all organizations the authenticated caller is a member of
// via the OrganizationQueryController.FindOrganizations RPC.
//
// No input parameters are required — the server determines membership from
// the caller's identity.
func List(ctx context.Context, serverAddress string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := organizationv1.NewOrganizationQueryControllerClient(conn)
			resp, err := client.FindOrganizations(ctx, &protobuf.CustomEmpty{})
			if err != nil {
				return "", domains.RPCError(err, "organizations")
			}
			return domains.MarshalJSON(resp)
		})
}
