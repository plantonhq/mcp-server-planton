package team

import (
	"context"
	"fmt"

	teamv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/team/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Get retrieves a team by ID via TeamQueryController.Get.
func Get(ctx context.Context, serverAddress, teamID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := teamv1.NewTeamQueryControllerClient(conn)
			resp, err := client.Get(ctx, &teamv1.TeamId{Value: teamID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("team %q", teamID))
			}
			return domains.MarshalJSON(resp)
		})
}
