package credential

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Apply creates or updates a credential of any supported type.
//
// The credentialObject is the full OpenMCF envelope provided by the agent:
//
//	{ api_version, kind, metadata: { name, org }, spec: { ... } }
//
// The kind field determines which per-type gRPC service is called.
func Apply(ctx context.Context, serverAddress string, credentialObject map[string]any) (string, error) {
	kind, err := extractKind(credentialObject)
	if err != nil {
		return "", err
	}

	d, ok := dispatchers[kind]
	if !ok {
		return "", fmt.Errorf("unsupported credential kind %q — supported: %s", kind, supportedKinds())
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resp, err := d.apply(ctx, conn, credentialObject)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("credential %q", kind))
			}
			return domains.MarshalJSON(resp)
		})
}
