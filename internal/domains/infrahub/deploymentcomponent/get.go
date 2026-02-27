package deploymentcomponent

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	deploymentcomponentv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/deploymentcomponent/v1"
	"google.golang.org/grpc"
)

// Get retrieves a deployment component by ID or by cloud resource kind via the
// DeploymentComponentQueryController RPCs.
//
// Exactly one identification path is used:
//   - When id is non-empty, calls Get with ApiResourceId.
//   - When kind is non-empty, resolves the PascalCase kind string to a
//     CloudResourceKind enum value and calls GetByCloudResourceKind.
func Get(ctx context.Context, serverAddress, id, kind string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := deploymentcomponentv1.NewDeploymentComponentQueryControllerClient(conn)

			if id != "" {
				resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: id})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("deployment component %q", id))
				}
				return domains.MarshalJSON(resp)
			}

			kindEnum, err := domains.ResolveKind(kind)
			if err != nil {
				return "", err
			}
			resp, err := client.GetByCloudResourceKind(ctx, &cloudresourcev1.CloudResourceKindRequest{Value: kindEnum})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("deployment component for kind %q", kind))
			}
			return domains.MarshalJSON(resp)
		})
}
