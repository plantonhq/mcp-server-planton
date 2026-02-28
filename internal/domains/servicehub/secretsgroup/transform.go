package secretsgroup

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	secretsgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/secretsgroup/v1"
	"google.golang.org/grpc"
)

// Transform batch-resolves $secrets-group/ references in environment
// variable maps via the SecretsGroupQueryController.Transform RPC.
//
// Values starting with $secrets-group/ are resolved to their actual
// values. Literal values pass through unchanged. The response includes
// both successfully transformed entries and any entries that failed
// resolution with error details.
//
// WARNING: Resolved values are returned in PLAINTEXT.
func Transform(ctx context.Context, serverAddress, org string, entries map[string]string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := secretsgroupv1.NewSecretsGroupQueryControllerClient(conn)
			resp, err := client.Transform(ctx, &secretsgroupv1.TransformSecretKeysRequest{
				Org:     org,
				Entries: entries,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("transform secrets in org %q", org))
			}
			return domains.MarshalJSON(resp)
		})
}
