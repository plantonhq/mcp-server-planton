package secret

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	secretv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secret/v1"
	"google.golang.org/grpc"
)

// Delete removes a secret and all its versions via the
// SecretCommandController.Delete RPC.
//
// Two identification paths are supported:
//   - ID path: calls Delete directly with the given ID.
//   - Slug path: first resolves org+scope+slug to a secret ID via the
//     query controller, then calls Delete.
//
// WARNING: This permanently destroys the secret record AND all its versions
// (including encrypted data stored in the backend).
func Delete(ctx context.Context, serverAddress, id, org string, scope secretv1.SecretSpec_Scope, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveSecretID(ctx, conn, id, org, scope, slug)
			if err != nil {
				return "", err
			}

			desc := describeSecret(id, org, scope, slug)
			client := secretv1.NewSecretCommandControllerClient(conn)
			deleted, err := client.Delete(ctx, &apiresource.ApiResourceDeleteInput{
				ResourceId: resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(deleted)
		})
}
