package pipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	pipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/pipeline/v1"
	"google.golang.org/grpc"
)

// UpdateFile modifies a pipeline file in a service's Git repository via the
// PipelineCommandController.UpdateServiceRepoPipelineFile RPC.
//
// The content string is encoded to bytes before sending. When
// expectedBaseSHA is provided, the write is rejected if the current blob SHA
// differs (optimistic locking).
func UpdateFile(ctx context.Context, serverAddress, serviceID, path, content, expectedBaseSHA, commitMessage, branch string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := pipelinev1.NewPipelineCommandControllerClient(conn)
			resp, err := client.UpdateServiceRepoPipelineFile(ctx, &pipelinev1.UpdateServiceRepoPipelineFileInput{
				ServiceId:       serviceID,
				Path:            path,
				Content:         []byte(content),
				ExpectedBaseSha: expectedBaseSHA,
				CommitMessage:   commitMessage,
				Branch:          branch,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("update pipeline file %q for service %q", path, serviceID))
			}
			return domains.MarshalJSON(resp)
		})
}
