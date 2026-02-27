package graph

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	graphv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/graph/v1"
	"google.golang.org/grpc"
)

// CloudResourceGraphInput holds the validated parameters for querying a
// cloud-resource-centric subgraph.
type CloudResourceGraphInput struct {
	CloudResourceID  string
	IncludeUpstream  bool
	IncludeDownstream bool
	MaxDepth         int32
}

// GetCloudResourceGraph retrieves a cloud-resource-centric subgraph including
// services deployed as it, credentials it uses, and optionally upstream/
// downstream dependencies via the GraphQueryController.GetCloudResourceGraph RPC.
func GetCloudResourceGraph(ctx context.Context, serverAddress string, input CloudResourceGraphInput) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := graphv1.NewGraphQueryControllerClient(conn)
			resp, err := client.GetCloudResourceGraph(ctx, &graphv1.GetCloudResourceGraphInput{
				CloudResourceId:  input.CloudResourceID,
				IncludeUpstream:  input.IncludeUpstream,
				IncludeDownstream: input.IncludeDownstream,
				MaxDepth:         input.MaxDepth,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("cloud resource graph for %q", input.CloudResourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
