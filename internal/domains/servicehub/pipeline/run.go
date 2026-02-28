package pipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	pipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/pipeline/v1"
	"google.golang.org/grpc"
)

// Run triggers a new pipeline for a service via the
// PipelineCommandController.RunGitCommit RPC.
//
// The branch is always required. When commitSHA is provided, the pipeline
// builds that exact commit; when empty, the pipeline uses the branch HEAD.
//
// The RPC returns google.protobuf.Empty — no pipeline ID is available
// directly. Callers should use GetLatest to check the result.
func Run(ctx context.Context, serverAddress, serviceID, branch, commitSHA string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := pipelinev1.NewPipelineCommandControllerClient(conn)
			_, err := client.RunGitCommit(ctx, &pipelinev1.RunGitCommitPipelineRequest{
				ServiceId: serviceID,
				Branch:    branch,
				CommitSha: commitSHA,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("run pipeline for service %q on branch %q", serviceID, branch))
			}

			msg := fmt.Sprintf("Pipeline triggered for service %q on branch %q.", serviceID, branch)
			if commitSHA != "" {
				msg = fmt.Sprintf("Pipeline triggered for service %q on branch %q at commit %q.", serviceID, branch, commitSHA)
			}
			msg += " Use get_last_pipeline with the service ID to check the result."
			return msg, nil
		})
}
