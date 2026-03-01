package organization

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	organizationv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1"
	"google.golang.org/grpc"
)

// UpdateFields holds the optional fields that can be modified on an existing
// organization. Zero-value (empty string) means "leave unchanged".
type UpdateFields struct {
	Name         string
	Description  string
	ContactEmail string
	LogoURL      string
}

// Update performs a read-modify-write on an existing organization.
//
// The function first fetches the current organization by ID via the query
// controller, applies any non-empty fields from UpdateFields, then writes the
// result back via the command controller's Update RPC. Both calls share a
// single gRPC connection within one WithConnection callback.
func Update(ctx context.Context, serverAddress, orgID string, fields UpdateFields) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			queryClient := organizationv1.NewOrganizationQueryControllerClient(conn)
			org, err := queryClient.Get(ctx, &organizationv1.OrganizationId{Value: orgID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("organization %q", orgID))
			}

			if fields.Name != "" {
				org.Metadata.Name = fields.Name
			}
			if org.Spec == nil {
				org.Spec = &organizationv1.OrganizationSpec{}
			}
			if fields.Description != "" {
				org.Spec.Description = fields.Description
			}
			if fields.ContactEmail != "" {
				org.Spec.ContactEmail = fields.ContactEmail
			}
			if fields.LogoURL != "" {
				org.Spec.LogoUrl = fields.LogoURL
			}

			cmdClient := organizationv1.NewOrganizationCommandControllerClient(conn)
			resp, err := cmdClient.Update(ctx, org)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("organization %q", orgID))
			}
			return domains.MarshalJSON(resp)
		})
}
