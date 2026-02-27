package graph

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	graphv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/graph/v1"
	"google.golang.org/grpc"
)

// ServiceGraphInput holds the validated parameters for querying a
// service-centric subgraph.
type ServiceGraphInput struct {
	ServiceID         string
	Envs              []string
	IncludeUpstream   bool
	IncludeDownstream bool
	MaxDepth          int32
}

// GetServiceGraph retrieves a service-centric subgraph including its cloud
// resource deployments and optionally upstream/downstream dependencies via
// the GraphQueryController.GetServiceGraph RPC.
func GetServiceGraph(ctx context.Context, serverAddress string, input ServiceGraphInput) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := graphv1.NewGraphQueryControllerClient(conn)
			resp, err := client.GetServiceGraph(ctx, &graphv1.GetServiceGraphInput{
				ServiceId:         input.ServiceID,
				Envs:              input.Envs,
				IncludeUpstream:   input.IncludeUpstream,
				IncludeDownstream: input.IncludeDownstream,
				MaxDepth:          input.MaxDepth,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("service graph for %q", input.ServiceID))
			}
			return domains.MarshalJSON(resp)
		})
}
