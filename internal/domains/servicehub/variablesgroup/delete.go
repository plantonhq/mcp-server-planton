package variablesgroup

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	variablesgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/variablesgroup/v1"
	"google.golang.org/grpc"
)

// Delete removes a variables group via the
// VariablesGroupCommandController.Delete RPC.
//
// Two identification paths are supported:
//   - ID path: calls Delete directly with the given ID.
//   - Slug path: first resolves org+slug to a group ID via the query
//     controller, then calls Delete.
func Delete(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveGroupID(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}

			desc := describeGroup(id, org, slug)
			client := variablesgroupv1.NewVariablesGroupCommandControllerClient(conn)
			deleted, err := client.Delete(ctx, &apiresource.ApiResourceDeleteInput{
				ResourceId: resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(deleted)
		})
}
