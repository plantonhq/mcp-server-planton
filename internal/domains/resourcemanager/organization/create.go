package organization

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	organizationv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1"
	"google.golang.org/grpc"
)

// Create provisions a new organization via the
// OrganizationCommandController.Create RPC.
//
// The caller supplies user-facing fields; the function assembles the full
// Organization proto with the required api_version and kind constants.
func Create(ctx context.Context, serverAddress, slug, name, description, contactEmail string) (string, error) {
	org := &organizationv1.Organization{
		ApiVersion: "resource-manager.planton.ai/v1",
		Kind:       "Organization",
		Metadata:   &apiresource.ApiResourceMetadata{Slug: slug, Name: name},
		Spec: &organizationv1.OrganizationSpec{
			Description:  description,
			ContactEmail: contactEmail,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := organizationv1.NewOrganizationCommandControllerClient(conn)
			resp, err := client.Create(ctx, org)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("organization %q", slug))
			}
			return domains.MarshalJSON(resp)
		})
}
