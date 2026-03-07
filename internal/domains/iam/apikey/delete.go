package apikey

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apikeyv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/apikey/v1"
	"google.golang.org/grpc"
)

// Delete permanently revokes and removes an API key via
// ApiKeyCommandController.Delete.
func Delete(ctx context.Context, serverAddress, apiKeyID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := apikeyv1.NewApiKeyCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &apikeyv1.ApiKeyId{Value: apiKeyID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("API key %q", apiKeyID))
			}
			return domains.MarshalJSON(resp)
		})
}
