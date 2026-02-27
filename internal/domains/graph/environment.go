package graph

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	graphv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/graph/v1"
	"google.golang.org/grpc"
)

// EnvironmentGraphInput holds the validated parameters for querying an
// environment-scoped resource graph.
type EnvironmentGraphInput struct {
	EnvID                   string
	NodeTypes               []graphv1.GraphNode_Type
	IncludeTopologicalOrder bool
}

// GetEnvironmentGraph retrieves all resources deployed in a specific
// environment via the GraphQueryController.GetEnvironmentGraph RPC.
func GetEnvironmentGraph(ctx context.Context, serverAddress string, input EnvironmentGraphInput) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := graphv1.NewGraphQueryControllerClient(conn)
			resp, err := client.GetEnvironmentGraph(ctx, &graphv1.GetEnvironmentGraphInput{
				EnvId:                   input.EnvID,
				NodeTypes:               input.NodeTypes,
				IncludeTopologicalOrder: input.IncludeTopologicalOrder,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("environment graph for %q", input.EnvID))
			}
			return domains.MarshalJSON(resp)
		})
}
