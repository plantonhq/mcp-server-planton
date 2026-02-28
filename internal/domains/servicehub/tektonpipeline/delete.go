package tektonpipeline

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"

	tektonpipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/tektonpipeline/v1"
)

// Delete removes a Tekton pipeline template via the
// TektonPipelineCommandController.Delete RPC.
//
// The delete RPC requires the full TektonPipeline entity as input, so this
// function first resolves the pipeline by ID or org+slug via the query
// controller, then passes the resolved entity to Delete. Both calls share a
// single gRPC connection.
func Delete(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			pipeline, err := resolvePipeline(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}

			desc := describePipeline(id, org, slug)
			client := tektonpipelinev1.NewTektonPipelineCommandControllerClient(conn)
			deleted, err := client.Delete(ctx, pipeline)
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(deleted)
		})
}
