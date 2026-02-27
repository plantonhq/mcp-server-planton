package infrapipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	infrapipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrapipeline/v1"
	"google.golang.org/grpc"
)

// Get retrieves a single infra pipeline by its ID via the
// InfraPipelineQueryController.Get RPC.
func Get(ctx context.Context, serverAddress, pipelineID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrapipelinev1.NewInfraPipelineQueryControllerClient(conn)
			resp, err := client.Get(ctx, &infrapipelinev1.InfraPipelineId{Value: pipelineID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("infra pipeline %q", pipelineID))
			}
			return domains.MarshalJSON(resp)
		})
}
