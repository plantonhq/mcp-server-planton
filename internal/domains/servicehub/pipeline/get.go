package pipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	pipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/pipeline/v1"
	"google.golang.org/grpc"
)

// Get retrieves a single pipeline by its ID via the
// PipelineQueryController.Get RPC.
func Get(ctx context.Context, serverAddress, pipelineID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := pipelinev1.NewPipelineQueryControllerClient(conn)
			resp, err := client.Get(ctx, &pipelinev1.PipelineId{Value: pipelineID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("pipeline %q", pipelineID))
			}
			return domains.MarshalJSON(resp)
		})
}
