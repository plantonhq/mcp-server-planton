package variable

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	variablev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/variable/v1"
	"google.golang.org/grpc"
)

// Delete removes a variable record via the
// VariableCommandController.Delete RPC.
//
// Two identification paths are supported:
//   - ID path: calls Delete directly with the given ID.
//   - Slug path: first resolves org+scope+slug to a variable ID via the
//     query controller, then calls Delete.
func Delete(ctx context.Context, serverAddress, id, org string, scope variablev1.VariableSpec_Scope, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveVariableID(ctx, conn, id, org, scope, slug)
			if err != nil {
				return "", err
			}

			desc := describeVariable(id, org, scope, slug)
			client := variablev1.NewVariableCommandControllerClient(conn)
			deleted, err := client.Delete(ctx, &apiresource.ApiResourceDeleteInput{
				ResourceId: resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(deleted)
		})
}
