package environment

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	environmentv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/environment/v1"
	organizationv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1"
	"google.golang.org/grpc"
)

// List retrieves environments within an organization that the authenticated
// caller is authorized to access via the
// EnvironmentQueryController.FindAuthorized RPC.
//
// Unlike the platform's findByOrg (which returns ALL environments in an org),
// FindAuthorized returns only environments where the caller has at least "get"
// permission, filtered through OpenFGA.
func List(ctx context.Context, serverAddress string, org string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := environmentv1.NewEnvironmentQueryControllerClient(conn)
			resp, err := client.FindAuthorized(ctx, &organizationv1.OrganizationId{Value: org})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("environments in org %q", org))
			}
			return domains.MarshalJSON(resp)
		})
}
