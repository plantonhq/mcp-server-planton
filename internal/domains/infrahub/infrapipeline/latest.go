package infrapipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	infrapipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrapipeline/v1"
	infraprojectv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infraproject/v1"
	"google.golang.org/grpc"
)

// GetLatest retrieves the most recent infra pipeline for an infra project via
// the InfraPipelineQueryController.GetLastInfraPipelineByInfraProjectId RPC.
//
// This is the primary function agents call after run_infra_pipeline or
// apply_infra_project to check whether the triggered pipeline completed
// successfully.
func GetLatest(ctx context.Context, serverAddress, infraProjectID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrapipelinev1.NewInfraPipelineQueryControllerClient(conn)
			resp, err := client.GetLastInfraPipelineByInfraProjectId(ctx, &infraprojectv1.InfraProjectId{Value: infraProjectID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("latest infra pipeline for project %q", infraProjectID))
			}
			return domains.MarshalJSON(resp)
		})
}
