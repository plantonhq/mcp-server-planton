package connection

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Apply creates or updates a connection by dispatching to the per-type gRPC
// command controller identified by the "kind" field in the connection object.
func Apply(ctx context.Context, serverAddress string, connectionObject map[string]any) (string, error) {
	kind, err := extractKind(connectionObject)
	if err != nil {
		return "", err
	}

	d, ok := dispatchers[kind]
	if !ok {
		return "", fmt.Errorf("unsupported connection kind %q — supported: %s", kind, supportedKinds())
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resp, err := d.apply(ctx, conn, connectionObject)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("connection %q", kind))
			}
			return domains.MarshalJSON(resp)
		})
}
