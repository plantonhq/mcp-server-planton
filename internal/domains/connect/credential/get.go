package credential

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Get retrieves a credential by ID, dispatching to the per-type gRPC service
// identified by kind. Sensitive fields are redacted before the response is
// returned to the agent.
func Get(ctx context.Context, serverAddress, kind, id string) (string, error) {
	d, ok := dispatchers[kind]
	if !ok {
		return "", fmt.Errorf("unsupported credential kind %q — supported: %s", kind, supportedKinds())
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resp, err := d.get(ctx, conn, id)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("credential %q (kind=%s)", id, kind))
			}
			jsonStr, err := domains.MarshalJSON(resp)
			if err != nil {
				return "", err
			}
			return redactFields(jsonStr, d.sensitiveFields)
		})
}

// GetByOrgBySlug retrieves a credential by organization and slug, dispatching
// to the per-type gRPC service identified by kind. Sensitive fields are
// redacted before the response is returned to the agent.
func GetByOrgBySlug(ctx context.Context, serverAddress, kind, org, slug string) (string, error) {
	d, ok := dispatchers[kind]
	if !ok {
		return "", fmt.Errorf("unsupported credential kind %q — supported: %s", kind, supportedKinds())
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resp, err := d.getByOrgBySlug(ctx, conn, org, slug)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("credential %q in org %q (kind=%s)", slug, org, kind))
			}
			jsonStr, err := domains.MarshalJSON(resp)
			if err != nil {
				return "", err
			}
			return redactFields(jsonStr, d.sensitiveFields)
		})
}
