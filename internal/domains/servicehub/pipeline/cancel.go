package pipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	pipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/pipeline/v1"
	"google.golang.org/grpc"
)

// Cancel stops a running pipeline via the
// PipelineCommandController.Cancel RPC.
//
// Cancellation is signal-based and graceful: during the build stage, Tekton
// PipelineRun resources are deleted and running build pods are terminated;
// during the deploy stage, the current deployment task receives a cancellation
// signal and remaining tasks are skipped.
//
// Returns the full updated Pipeline after cancellation.
func Cancel(ctx context.Context, serverAddress, pipelineID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := pipelinev1.NewPipelineCommandControllerClient(conn)
			resp, err := client.Cancel(ctx, &pipelinev1.PipelineId{Value: pipelineID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("cancel pipeline %q", pipelineID))
			}
			return domains.MarshalJSON(resp)
		})
}
