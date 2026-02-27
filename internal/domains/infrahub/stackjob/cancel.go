package stackjob

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
	"google.golang.org/grpc"
)

// Cancel gracefully cancels a running stack job via the
// StackJobCommandController.Cancel RPC.
//
// Cancellation is signal-based and graceful: the currently executing IaC
// operation completes fully, remaining operations are skipped and marked
// as cancelled. Infrastructure created by completed operations remains
// (no automatic rollback). The resource lock is released, allowing queued
// stack jobs to proceed.
//
// Returns the StackJob resource; status updates asynchronously as
// cancellation completes.
func Cancel(ctx context.Context, serverAddress, stackJobID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobCommandControllerClient(conn)
			resp, err := client.Cancel(ctx, &stackjobv1.StackJobId{Value: stackJobID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("cancel stack job %q", stackJobID))
			}
			return domains.MarshalJSON(resp)
		})
}
