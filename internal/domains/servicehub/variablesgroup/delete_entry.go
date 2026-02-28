package variablesgroup

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	variablesgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/variablesgroup/v1"
	"google.golang.org/grpc"
)

// DeleteEntry removes a single variable from a variables group via the
// VariablesGroupCommandController.DeleteEntry RPC.
//
// The target group can be identified by group_id directly, or by org+slug
// (which triggers an extra lookup to resolve the group ID).
func DeleteEntry(ctx context.Context, serverAddress, groupID, org, groupSlug, entryName string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resolvedID, err := resolveGroupID(ctx, conn, groupID, org, groupSlug)
			if err != nil {
				return "", err
			}

			client := variablesgroupv1.NewVariablesGroupCommandControllerClient(conn)
			result, err := client.DeleteEntry(ctx, &variablesgroupv1.DeleteVariableRequest{
				GroupId:   resolvedID,
				EntryName: entryName,
			})
			if err != nil {
				desc := describeGroup(groupID, org, groupSlug)
				return "", domains.RPCError(err, fmt.Sprintf("delete variable %q from %s", entryName, desc))
			}
			return domains.MarshalJSON(result)
		})
}
