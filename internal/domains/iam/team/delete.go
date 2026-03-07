package team

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	teamv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/team/v1"
	"google.golang.org/grpc"
)

// Delete permanently removes a team via TeamCommandController.Delete.
func Delete(ctx context.Context, serverAddress, teamID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := teamv1.NewTeamCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &teamv1.TeamId{Value: teamID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("team %q", teamID))
			}
			return domains.MarshalJSON(resp)
		})
}
