package service

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	servicev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/service/v1"
	"google.golang.org/grpc"
)

// ConfigureWebhook creates or refreshes the webhook on GitHub/GitLab for the
// given service via the ServiceCommandController.ConfigureWebhook RPC.
//
// Two identification paths are supported:
//   - ID path: uses the given ID directly.
//   - Slug path: first resolves org+slug to a service ID via the query
//     controller, then calls ConfigureWebhook.
//
// Useful for recovering from accidentally deleted webhooks or refreshing
// webhook configuration.
func ConfigureWebhook(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveServiceID(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}

			desc := describeService(id, org, slug)
			client := servicev1.NewServiceCommandControllerClient(conn)
			result, err := client.ConfigureWebhook(ctx, &servicev1.ServiceId{Value: resourceID})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
