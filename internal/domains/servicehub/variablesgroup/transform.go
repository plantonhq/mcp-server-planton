package variablesgroup

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	variablesgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/variablesgroup/v1"
	"google.golang.org/grpc"
)

// Transform batch-resolves $variables-group/ references in environment
// variable maps via the VariablesGroupQueryController.Transform RPC.
//
// Values starting with $variables-group/ are resolved to their actual
// values. Literal values pass through unchanged. The response includes
// both successfully transformed entries and any entries that failed
// resolution with error details.
func Transform(ctx context.Context, serverAddress, org string, entries map[string]string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := variablesgroupv1.NewVariablesGroupQueryControllerClient(conn)
			resp, err := client.Transform(ctx, &variablesgroupv1.TransformConfigKeysRequest{
				Org:     org,
				Entries: entries,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("transform variables in org %q", org))
			}
			return domains.MarshalJSON(resp)
		})
}
