package secretbackend

import (
	"context"

	"google.golang.org/grpc"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	secretbackendv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/secretbackend/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// Delete removes a secret backend by ID or by org+slug.
// When using org+slug, the backend is first resolved to get its ID.
func Delete(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveSecretBackendID(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}

			desc := describeSecretBackend(id, org, slug)
			client := secretbackendv1.NewSecretBackendCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &apiresource.ApiResourceDeleteInput{
				ResourceId: resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			RedactSecretBackend(resp)
			return domains.MarshalJSON(resp)
		})
}
