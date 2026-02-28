package variablesgroup

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	variablesgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/variablesgroup/v1"
	"google.golang.org/grpc"
)

// GetValue retrieves the resolved value of a specific variable from a
// variables group via the VariablesGroupQueryController.GetValue RPC.
//
// If the variable uses a value_from reference, the backend resolves it to
// the current value. The result is returned as a plain text string.
func GetValue(ctx context.Context, serverAddress, org, groupName, entryName string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := variablesgroupv1.NewVariablesGroupQueryControllerClient(conn)
			resp, err := client.GetValue(ctx, &variablesgroupv1.GetVariableValueRequest{
				Org:       org,
				GroupName: groupName,
				EntryName: entryName,
			})
			if err != nil {
				return "", domains.RPCError(err,
					fmt.Sprintf("variable %q in group %q (org %q)", entryName, groupName, org))
			}

			if resp == nil || resp.GetValue() == "" {
				return fmt.Sprintf("No value found for variable %q in group %q (org %q).",
					entryName, groupName, org), nil
			}
			return resp.GetValue(), nil
		})
}
