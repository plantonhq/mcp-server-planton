package pipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	pipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/pipeline/v1"
	"google.golang.org/grpc"
)

// Rerun re-executes a previously executed pipeline via the
// PipelineCommandController.Rerun RPC.
//
// The pipeline is re-run using the same configuration (service, branch,
// commit) as the original execution. Returns the newly created pipeline.
func Rerun(ctx context.Context, serverAddress, pipelineID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := pipelinev1.NewPipelineCommandControllerClient(conn)
			resp, err := client.Rerun(ctx, &pipelinev1.PipelineId{Value: pipelineID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("rerun pipeline %q", pipelineID))
			}
			return domains.MarshalJSON(resp)
		})
}
