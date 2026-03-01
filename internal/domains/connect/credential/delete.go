package credential

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Delete removes a credential by ID, dispatching to the per-type gRPC
// service identified by kind.
func Delete(ctx context.Context, serverAddress, kind, id string) (string, error) {
	d, ok := dispatchers[kind]
	if !ok {
		return "", fmt.Errorf("unsupported credential kind %q — supported: %s", kind, supportedKinds())
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resp, err := d.del(ctx, conn, id)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("credential %q (kind=%s)", id, kind))
			}
			return domains.MarshalJSON(resp)
		})
}
