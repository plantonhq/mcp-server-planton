package infrapipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	infrapipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrapipeline/v1"
	"google.golang.org/grpc"
)

// Cancel stops a running infra pipeline via the
// InfraPipelineCommandController.Cancel RPC.
//
// Returns the full updated InfraPipeline after cancellation.
func Cancel(ctx context.Context, serverAddress, pipelineID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrapipelinev1.NewInfraPipelineCommandControllerClient(conn)
			resp, err := client.Cancel(ctx, &infrapipelinev1.InfraPipelineId{Value: pipelineID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("cancel infra pipeline %q", pipelineID))
			}
			return domains.MarshalJSON(resp)
		})
}
