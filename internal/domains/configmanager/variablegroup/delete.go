package variablegroup

import (
	"context"

	"google.golang.org/grpc"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	variablegroupv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/variablegroup/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// Delete removes a variable group by ID or by org+scope+slug.
// When using org+scope+slug, the group is first resolved to get its ID.
func Delete(ctx context.Context, serverAddress, id, org string, scope variablegroupv1.VariableGroupSpec_Scope, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveVariableGroupID(ctx, conn, id, org, scope, slug)
			if err != nil {
				return "", err
			}

			desc := describeVariableGroup(id, org, scope, slug)
			client := variablegroupv1.NewVariableGroupCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &apiresource.ApiResourceDeleteInput{
				ResourceId: resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(resp)
		})
}
