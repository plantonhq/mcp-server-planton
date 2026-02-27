package infraproject

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	infraprojectv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infraproject/v1"
	"google.golang.org/grpc"
)

// CheckSlugAvailability checks whether a slug is available for use within
// an organization via the InfraProjectQueryController.CheckSlugAvailability RPC.
//
// InfraProject slugs are scoped to the organization only (unlike cloud
// resources which are scoped to org+env+kind).
func CheckSlugAvailability(ctx context.Context, serverAddress, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infraprojectv1.NewInfraProjectQueryControllerClient(conn)
			resp, err := client.CheckSlugAvailability(ctx, &infraprojectv1.InfraProjectSlugAvailabilityCheckRequest{
				Org:  org,
				Slug: slug,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("slug availability for %q in org %q", slug, org))
			}
			return domains.MarshalJSON(resp)
		})
}
