package secretsgroup

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	secretsgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/secretsgroup/v1"
	"google.golang.org/grpc"
)

// GetValue retrieves the resolved value of a specific secret from a
// secrets group via the SecretsGroupQueryController.GetValue RPC.
//
// If the secret uses a value_from reference, the backend resolves it to
// the current value. The result is returned as a plain text string.
//
// WARNING: The returned value is in PLAINTEXT.
func GetValue(ctx context.Context, serverAddress, org, groupName, entryName string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := secretsgroupv1.NewSecretsGroupQueryControllerClient(conn)
			resp, err := client.GetValue(ctx, &secretsgroupv1.GetSecretValueRequest{
				Org:       org,
				GroupName: groupName,
				EntryName: entryName,
			})
			if err != nil {
				return "", domains.RPCError(err,
					fmt.Sprintf("secret %q in group %q (org %q)", entryName, groupName, org))
			}

			if resp == nil || resp.GetValue() == "" {
				return fmt.Sprintf("No value found for secret %q in group %q (org %q).",
					entryName, groupName, org), nil
			}
			return resp.GetValue(), nil
		})
}
