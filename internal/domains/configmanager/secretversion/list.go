package secretversion

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	secretversionv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secretversion/v1"
	"google.golang.org/grpc"
)

// List retrieves all versions of a secret via the
// SecretVersionQueryController.ListBySecret RPC.
//
// Returns metadata only (secret ID, timestamps, backend version ID) —
// the data field is intentionally omitted for performance and security.
// Use this to inspect version history before creating a new version.
func List(ctx context.Context, serverAddress, secretID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := secretversionv1.NewSecretVersionQueryControllerClient(conn)
			resp, err := client.ListBySecret(ctx, &secretversionv1.SecretVersionsBySecretInput{
				SecretId: secretID,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("secret versions for secret %q", secretID))
			}
			return domains.MarshalJSON(resp)
		})
}
