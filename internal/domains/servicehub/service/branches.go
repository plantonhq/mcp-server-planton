package service

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	servicev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/service/v1"
	"google.golang.org/grpc"
)

// ListBranches lists all Git branches in the repository connected to the
// service via the ServiceQueryController.ListBranches RPC.
//
// Two identification paths are supported:
//   - ID path: uses the given ID directly.
//   - Slug path: first resolves org+slug to a service ID via the query
//     controller, then calls ListBranches.
func ListBranches(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveServiceID(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}

			desc := describeService(id, org, slug)
			client := servicev1.NewServiceQueryControllerClient(conn)
			result, err := client.ListBranches(ctx, &servicev1.ServiceId{Value: resourceID})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
