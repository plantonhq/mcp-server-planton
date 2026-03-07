package stackjob

import (
	"context"
	"fmt"

	stackjobv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/infrahub/stackjob/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Rerun re-runs a previously executed stack job via the
// StackJobCommandController.Rerun RPC.
//
// Returns the full updated StackJob after the rerun is initiated.
func Rerun(ctx context.Context, serverAddress, stackJobID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobCommandControllerClient(conn)
			resp, err := client.Rerun(ctx, &stackjobv1.StackJobId{Value: stackJobID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("rerun stack job %q", stackJobID))
			}
			return domains.MarshalJSON(resp)
		})
}
