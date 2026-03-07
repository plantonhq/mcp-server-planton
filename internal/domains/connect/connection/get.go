package connection

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Get retrieves a connection by ID, dispatching to the per-type gRPC query
// controller identified by kind.
func Get(ctx context.Context, serverAddress, kind, id string) (string, error) {
	d, ok := dispatchers[kind]
	if !ok {
		return "", fmt.Errorf("unsupported connection kind %q — supported: %s", kind, supportedKinds())
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resp, err := d.get(ctx, conn, id)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("connection %q (kind=%s)", id, kind))
			}
			return domains.MarshalJSON(resp)
		})
}

// GetByOrgBySlug retrieves a connection by organization and slug, dispatching
// to the per-type gRPC query controller identified by kind.
func GetByOrgBySlug(ctx context.Context, serverAddress, kind, org, slug string) (string, error) {
	d, ok := dispatchers[kind]
	if !ok {
		return "", fmt.Errorf("unsupported connection kind %q — supported: %s", kind, supportedKinds())
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resp, err := d.getByOrgBySlug(ctx, conn, org, slug)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("connection %s/%s (kind=%s)", org, slug, kind))
			}
			return domains.MarshalJSON(resp)
		})
}
