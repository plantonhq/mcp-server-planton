package secret

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	secretv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secret/v1"
	"google.golang.org/grpc"
)

// ApplyInput holds the explicit parameters for creating or updating a secret.
type ApplyInput struct {
	Name        string
	Org         string
	Scope       secretv1.SecretSpec_Scope
	Env         string
	Description string
	Backend     string
}

// Apply creates or updates a secret via the
// SecretCommandController.Apply RPC.
//
// The Secret proto is constructed from explicit parameters. Only metadata
// fields are managed — secret values are stored via create_secret_version.
func Apply(ctx context.Context, serverAddress string, input ApplyInput) (string, error) {
	sec := &secretv1.Secret{
		Metadata: &apiresource.ApiResourceMetadata{
			Name: input.Name,
			Org:  input.Org,
			Env:  input.Env,
		},
		Spec: &secretv1.SecretSpec{
			Scope:       input.Scope,
			Description: input.Description,
			Backend:     input.Backend,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := secretv1.NewSecretCommandControllerClient(conn)
			result, err := client.Apply(ctx, sec)
			if err != nil {
				desc := fmt.Sprintf("secret %q in org %q", input.Name, input.Org)
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
