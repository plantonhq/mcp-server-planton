package serviceaccount

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	serviceaccountv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/serviceaccount/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// Delete removes a service account by ID. This cascades: all API keys
// are revoked, authorization tuples are removed, and the backing
// identity account is deleted.
func Delete(ctx context.Context, serverAddress, id string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := serviceaccountv1.NewServiceAccountCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &serviceaccountv1.ServiceAccountId{Value: id})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("service account %q", id))
			}
			return domains.MarshalJSON(resp)
		})
}
