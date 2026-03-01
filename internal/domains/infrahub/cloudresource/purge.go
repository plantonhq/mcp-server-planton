package cloudresource

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"google.golang.org/grpc"
)

// Purge fully removes a cloud resource by first destroying its IaC-managed
// infrastructure and then deleting the resource record. The backend
// orchestrates this as a Temporal workflow (destroy → wait → delete).
//
// Two identification paths are supported:
//   - ID path: calls Purge(CloudResourceId) directly with the given ID.
//   - Slug path: first resolves the composite key (kind, org, env, slug) to a
//     resource ID via the query controller, then calls Purge. Both calls share
//     a single gRPC connection.
//
// The caller must validate the ResourceIdentifier before calling this function.
func Purge(ctx context.Context, serverAddress string, id ResourceIdentifier) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveResourceID(ctx, conn, id)
			if err != nil {
				return "", err
			}

			desc := describeIdentifier(id)
			cmdClient := cloudresourcev1.NewCloudResourceCommandControllerClient(conn)
			purged, err := cmdClient.Purge(ctx, &cloudresourcev1.CloudResourceId{
				Value: resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(purged)
		})
}
