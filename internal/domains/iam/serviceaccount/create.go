package serviceaccount

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	serviceaccountv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/serviceaccount/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// Create provisions a new service account via ServiceAccountCommandController.Create.
func Create(ctx context.Context, serverAddress, org, displayName, description string) (string, error) {
	sa := &serviceaccountv1.ServiceAccount{
		ApiVersion: "iam.planton.ai/v1",
		Kind:       "ServiceAccount",
		Metadata:   &apiresource.ApiResourceMetadata{Org: org, Name: displayName},
		Spec: &serviceaccountv1.ServiceAccountSpec{
			DisplayName: displayName,
			Description: description,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := serviceaccountv1.NewServiceAccountCommandControllerClient(conn)
			resp, err := client.Create(ctx, sa)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("service account %q in org %q", displayName, org))
			}
			return domains.MarshalJSON(resp)
		})
}
