package cloudresource

import (
	"context"

	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"github.com/plantoncloud/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Delete removes a cloud resource via the CloudResourceCommandController.
//
// Two identification paths are supported:
//   - ID path: calls Delete(ApiResourceDeleteInput) directly with the given ID.
//   - Slug path: first resolves the composite key (kind, org, env, slug) to a
//     resource ID via the query controller, then calls Delete. Both calls share
//     a single gRPC connection.
//
// The caller must validate the ResourceIdentifier before calling this function.
func Delete(ctx context.Context, serverAddress string, id ResourceIdentifier) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveResourceID(ctx, conn, id)
			if err != nil {
				return "", err
			}

			desc := describeIdentifier(id)
			cmdClient := cloudresourcev1.NewCloudResourceCommandControllerClient(conn)
			deleted, err := cmdClient.Delete(ctx, &apiresource.ApiResourceDeleteInput{
				ResourceId: resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(deleted)
		})
}
