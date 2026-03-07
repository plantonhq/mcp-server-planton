package stackjob

import (
	"context"
	"fmt"

	stackjobv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/infrahub/stackjob/v1"
	servicev1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/servicehub/service/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// FindServiceStackJobsByEnv retrieves the most recent stack job for a service
// in each of its deployed environments via the
// StackJobQueryController.FindServiceStackJobsByEnv RPC.
//
// The response is a map keyed by environment name, where each value is the
// latest stack job for that service-environment pair.
func FindServiceStackJobsByEnv(ctx context.Context, serverAddress, serviceID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobQueryControllerClient(conn)
			resp, err := client.FindServiceStackJobsByEnv(ctx, &servicev1.ServiceId{Value: serviceID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("stack jobs by environment for service %q", serviceID))
			}
			return domains.MarshalJSON(resp)
		})
}
