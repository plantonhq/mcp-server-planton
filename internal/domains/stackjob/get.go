package stackjob

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
	"google.golang.org/grpc"
)

// Get retrieves a single stack job by its ID via the
// StackJobQueryController.Get RPC.
func Get(ctx context.Context, serverAddress, stackJobID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobQueryControllerClient(conn)
			resp, err := client.Get(ctx, &stackjobv1.StackJobId{Value: stackJobID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("stack job %q", stackJobID))
			}
			return domains.MarshalJSON(resp)
		})
}
