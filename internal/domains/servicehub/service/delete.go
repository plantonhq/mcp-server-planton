package service

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	servicev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/service/v1"
	"google.golang.org/grpc"
)

// Delete removes a service record via the
// ServiceCommandController.Delete RPC.
//
// Two identification paths are supported:
//   - ID path: calls Delete directly with the given ID.
//   - Slug path: first resolves org+slug to a service ID via the query
//     controller, then calls Delete. Both calls share a single gRPC connection.
//
// This removes the database record and disconnects webhooks. It does NOT
// delete deployed cloud resources (microservices, databases, etc.).
func Delete(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveServiceID(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}

			desc := describeService(id, org, slug)
			client := servicev1.NewServiceCommandControllerClient(conn)
			deleted, err := client.Delete(ctx, &apiresource.ApiResourceDeleteInput{
				ResourceId: resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(deleted)
		})
}
