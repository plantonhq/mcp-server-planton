package secretbackend

import (
	"context"

	"google.golang.org/grpc"

	secretbackendv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/secretbackend/v1"
	orgv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/resourcemanager/organization/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// List returns all secret backends for an organization.
func List(ctx context.Context, serverAddress, org string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := secretbackendv1.NewSecretBackendQueryControllerClient(conn)
			resp, err := client.ListByOrg(ctx, &orgv1.OrganizationId{Value: org})
			if err != nil {
				return "", domains.RPCError(err, "secret backends")
			}
			for _, sb := range resp.GetEntries() {
				RedactSecretBackend(sb)
			}
			return domains.MarshalJSON(resp)
		})
}
