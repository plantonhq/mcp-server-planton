package cloudresource

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Get retrieves a cloud resource via the CloudResourceQueryController.
//
// Two identification paths are supported (delegated to [resolveResource]):
//   - ID path: fetches by CloudResourceId directly.
//   - Slug path: resolves the PascalCase kind to the proto enum, then fetches
//     by (org, env, kind, slug).
//
// The caller must validate the ResourceIdentifier before calling this function.
func Get(ctx context.Context, serverAddress string, id ResourceIdentifier) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			cr, err := resolveResource(ctx, conn, id)
			if err != nil {
				return "", err
			}
			return domains.MarshalJSON(cr)
		})
}
