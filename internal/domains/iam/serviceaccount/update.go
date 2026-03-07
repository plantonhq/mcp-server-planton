package serviceaccount

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	serviceaccountv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/serviceaccount/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// UpdateFields holds optional fields for a read-modify-write update.
// Empty strings mean "leave unchanged".
type UpdateFields struct {
	DisplayName string
	Description string
}

// Update performs a read-modify-write on an existing service account.
func Update(ctx context.Context, serverAddress, id string, fields UpdateFields) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			queryClient := serviceaccountv1.NewServiceAccountQueryControllerClient(conn)
			sa, err := queryClient.Get(ctx, &serviceaccountv1.ServiceAccountId{Value: id})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("service account %q", id))
			}

			if sa.Spec == nil {
				sa.Spec = &serviceaccountv1.ServiceAccountSpec{}
			}
			if fields.DisplayName != "" {
				sa.Spec.DisplayName = fields.DisplayName
			}
			if fields.Description != "" {
				sa.Spec.Description = fields.Description
			}

			cmdClient := serviceaccountv1.NewServiceAccountCommandControllerClient(conn)
			resp, err := cmdClient.Update(ctx, sa)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("service account %q", id))
			}
			return domains.MarshalJSON(resp)
		})
}
