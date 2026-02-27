package secretversion

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	secretversionv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secretversion/v1"
	"google.golang.org/grpc"
)

// Create stores a new secret version via the
// SecretVersionCommandController.Create RPC.
//
// The data map is encrypted via envelope encryption and stored in the parent
// secret's backend. Version metadata (secret ID, timestamps) is stored in
// the platform database. The encrypted data is NOT stored in the database.
func Create(ctx context.Context, serverAddress, secretID string, data map[string]string) (string, error) {
	sv := &secretversionv1.SecretVersion{
		Metadata: &apiresource.ApiResourceMetadata{},
		Spec: &secretversionv1.SecretVersionSpec{
			SecretId: secretID,
			Data:     data,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := secretversionv1.NewSecretVersionCommandControllerClient(conn)
			result, err := client.Create(ctx, sv)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("secret version for secret %q", secretID))
			}
			return domains.MarshalJSON(result)
		})
}
