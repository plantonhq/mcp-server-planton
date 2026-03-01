package infrapipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	infrapipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrapipeline/v1"
	"google.golang.org/grpc"
)

// Delete removes an infra pipeline record by ID via the
// InfraPipelineCommandController.Delete RPC.
//
// Returns the deleted InfraPipeline including its final status.
func Delete(ctx context.Context, serverAddress, pipelineID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrapipelinev1.NewInfraPipelineCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &apiresource.ApiResourceDeleteInput{
				ResourceId: pipelineID,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("infra pipeline %q", pipelineID))
			}
			return domains.MarshalJSON(resp)
		})
}
