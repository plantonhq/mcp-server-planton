package environment

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	environmentv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/environment/v1"
	"google.golang.org/grpc"
)

// Delete removes an environment by ID via the
// EnvironmentCommandController.Delete RPC.
//
// Deleting an environment triggers cascading cleanup of all resources deployed
// to it, including stack-modules, microservices, secrets, and clusters.
func Delete(ctx context.Context, serverAddress, envID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := environmentv1.NewEnvironmentCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &environmentv1.EnvironmentId{Value: envID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("environment %q", envID))
			}
			return domains.MarshalJSON(resp)
		})
}
