package serviceaccount

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	serviceaccountv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/serviceaccount/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// Get retrieves a service account by ID.
func Get(ctx context.Context, serverAddress, id string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := serviceaccountv1.NewServiceAccountQueryControllerClient(conn)
			resp, err := client.Get(ctx, &serviceaccountv1.ServiceAccountId{Value: id})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("service account %q", id))
			}
			return domains.MarshalJSON(resp)
		})
}
