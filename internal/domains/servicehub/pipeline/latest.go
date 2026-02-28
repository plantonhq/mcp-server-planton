package pipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	pipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/pipeline/v1"
	servicev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/service/v1"
	"google.golang.org/grpc"
)

// GetLatest retrieves the most recent pipeline for a service via the
// PipelineQueryController.GetLastPipelineByServiceId RPC.
//
// This is the primary function agents call after run_pipeline to check
// whether the triggered pipeline completed successfully.
func GetLatest(ctx context.Context, serverAddress, serviceID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := pipelinev1.NewPipelineQueryControllerClient(conn)
			resp, err := client.GetLastPipelineByServiceId(ctx, &servicev1.ServiceId{Value: serviceID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("latest pipeline for service %q", serviceID))
			}
			return domains.MarshalJSON(resp)
		})
}
