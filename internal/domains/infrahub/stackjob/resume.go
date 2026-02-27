package stackjob

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
	"google.golang.org/grpc"
)

// Resume approves and resumes a stack job that is in the awaiting_approval
// state via the StackJobCommandController.Resume RPC.
//
// Stack jobs enter awaiting_approval when a flow control policy requires
// manual approval before IaC execution proceeds. This function unblocks
// such jobs, allowing them to continue with their remaining operations.
//
// Returns the full updated StackJob after the resume is initiated.
func Resume(ctx context.Context, serverAddress, stackJobID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobCommandControllerClient(conn)
			resp, err := client.Resume(ctx, &stackjobv1.StackJobId{Value: stackJobID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("resume stack job %q", stackJobID))
			}
			return domains.MarshalJSON(resp)
		})
}
