package graph

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	graphv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/graph/v1"
	"google.golang.org/grpc"
)

// GetImpactAnalysis evaluates the impact of modifying or deleting a resource
// via the GraphQueryController.GetImpactAnalysis RPC.
//
// Returns direct and transitive impacts, total affected count, and a
// breakdown by resource type.
func GetImpactAnalysis(ctx context.Context, serverAddress, resourceID string, changeType graphv1.GetImpactAnalysisInput_ChangeType) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := graphv1.NewGraphQueryControllerClient(conn)
			resp, err := client.GetImpactAnalysis(ctx, &graphv1.GetImpactAnalysisInput{
				ResourceId: resourceID,
				ChangeType: changeType,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("impact analysis for resource %q", resourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
