package infrapipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	infrapipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrapipeline/v1"
	infraprojectv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infraproject/v1"
	"google.golang.org/grpc"
)

// Run triggers a new infra pipeline for a project.
//
// Two execution paths are supported based on the project's source type:
//   - Chart-sourced: when commitSHA is empty, calls
//     InfraPipelineCommandController.RunInfraProjectChartSourcePipeline.
//   - Git-sourced: when commitSHA is provided, calls
//     InfraPipelineCommandController.RunGitCommit with the given SHA.
//
// Both RPCs return the newly created pipeline's ID.
func Run(ctx context.Context, serverAddress, infraProjectID, commitSHA string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrapipelinev1.NewInfraPipelineCommandControllerClient(conn)

			if commitSHA != "" {
				resp, err := client.RunGitCommit(ctx, &infrapipelinev1.RunGitCommitInfraPipelineRequest{
					InfraProjectId: infraProjectID,
					CommitSha:      commitSHA,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("run git-commit pipeline for project %q at %q", infraProjectID, commitSHA))
				}
				return domains.MarshalJSON(resp)
			}

			resp, err := client.RunInfraProjectChartSourcePipeline(ctx, &infraprojectv1.InfraProjectId{Value: infraProjectID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("run chart-source pipeline for project %q", infraProjectID))
			}
			return domains.MarshalJSON(resp)
		})
}
