package graph

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	graphv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/graph/v1"
	"google.golang.org/grpc"
)

// OrganizationGraphInput holds the validated parameters for querying an
// organization's resource topology.
type OrganizationGraphInput struct {
	Org                     string
	Envs                    []string
	NodeTypes               []graphv1.GraphNode_Type
	IncludeTopologicalOrder bool
	MaxDepth                int32
}

// GetOrganizationGraph retrieves the complete resource topology for an
// organization via the GraphQueryController.GetOrganizationGraph RPC.
func GetOrganizationGraph(ctx context.Context, serverAddress string, input OrganizationGraphInput) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := graphv1.NewGraphQueryControllerClient(conn)
			resp, err := client.GetOrganizationGraph(ctx, &graphv1.GetOrganizationGraphInput{
				Org:                     input.Org,
				Envs:                    input.Envs,
				NodeTypes:               input.NodeTypes,
				IncludeTopologicalOrder: input.IncludeTopologicalOrder,
				MaxDepth:                input.MaxDepth,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("organization graph for %q", input.Org))
			}
			return domains.MarshalJSON(resp)
		})
}
