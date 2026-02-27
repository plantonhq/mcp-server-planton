package cloudresource

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"google.golang.org/grpc"
)

// Destroy tears down the cloud infrastructure for a resource via the
// CloudResourceCommandController.Destroy RPC while keeping the resource
// record on the platform.
//
// Two identification paths are supported:
//   - ID path: fetches the full resource by ID, then calls Destroy.
//   - Slug path: fetches the full resource by (kind, org, env, slug), then
//     calls Destroy. Both calls share a single gRPC connection.
//
// The caller must validate the ResourceIdentifier before calling this function.
func Destroy(ctx context.Context, serverAddress string, id ResourceIdentifier) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			cr, err := resolveResource(ctx, conn, id)
			if err != nil {
				return "", err
			}

			desc := describeIdentifier(id)
			cmdClient := cloudresourcev1.NewCloudResourceCommandControllerClient(conn)
			result, err := cmdClient.Destroy(ctx, cr)
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
