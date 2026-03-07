package serviceaccount

import (
	"context"

	"google.golang.org/grpc"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	serviceaccountv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/serviceaccount/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// List returns all service accounts in an organization.
func List(ctx context.Context, serverAddress, org string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := serviceaccountv1.NewServiceAccountQueryControllerClient(conn)
			resp, err := client.FindByOrg(ctx, &apiresource.ApiResourceId{Value: org})
			if err != nil {
				return "", domains.RPCError(err, "service accounts")
			}
			return domains.MarshalJSON(resp)
		})
}
