package infrapipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	plantongrpc "github.com/plantonhq/mcp-server-planton/internal/grpc"
	infrapipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrapipeline/v1"
	tekton "github.com/plantonhq/planton/apis/stubs/go/ai/planton/integration/tekton"
	"google.golang.org/grpc"
)

const maxLogEntries = 1000

// GetLogs collects Tekton task log entries for an infra pipeline by internally
// calling the streaming GetLogStream RPC and draining results until the
// stream closes (completed job) or the collect timeout expires (running job).
func GetLogs(ctx context.Context, serverAddress, pipelineID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			streamCtx, cancel := context.WithTimeout(ctx, plantongrpc.StreamCollectTimeout)
			defer cancel()

			client := infrapipelinev1.NewInfraPipelineQueryControllerClient(conn)
			stream, err := client.GetLogStream(streamCtx, &infrapipelinev1.InfraPipelineId{Value: pipelineID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("infra pipeline logs %q", pipelineID))
			}

			text, _, err := domains.DrainStream(stream, maxLogEntries, formatLogEntry)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("infra pipeline logs %q", pipelineID))
			}
			return text, nil
		})
}

func formatLogEntry(e *tekton.TektonTaskLogEntry) string {
	return fmt.Sprintf("[%s] %s", e.GetTaskName(), e.GetLogMessage())
}
