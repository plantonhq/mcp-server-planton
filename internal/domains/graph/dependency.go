package graph

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	graphv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/graph/v1"
	"google.golang.org/grpc"
)

// DependencyInput holds the validated parameters shared by GetDependencies
// and GetDependents queries.
type DependencyInput struct {
	ResourceID        string
	MaxDepth          int32
	RelationshipTypes []graphv1.GraphRelationship_Type
}

// GetDependencies retrieves all resources that the given resource depends on
// (upstream traversal) via the GraphQueryController.GetDependencies RPC.
func GetDependencies(ctx context.Context, serverAddress string, input DependencyInput) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := graphv1.NewGraphQueryControllerClient(conn)
			resp, err := client.GetDependencies(ctx, &graphv1.GetDependenciesInput{
				ResourceId:        input.ResourceID,
				MaxDepth:          input.MaxDepth,
				RelationshipTypes: input.RelationshipTypes,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("dependencies for resource %q", input.ResourceID))
			}
			return domains.MarshalJSON(resp)
		})
}

// GetDependents retrieves all resources that depend on the given resource
// (downstream traversal) via the GraphQueryController.GetDependents RPC.
func GetDependents(ctx context.Context, serverAddress string, input DependencyInput) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := graphv1.NewGraphQueryControllerClient(conn)
			resp, err := client.GetDependents(ctx, &graphv1.GetDependentsInput{
				ResourceId:        input.ResourceID,
				MaxDepth:          input.MaxDepth,
				RelationshipTypes: input.RelationshipTypes,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("dependents of resource %q", input.ResourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
