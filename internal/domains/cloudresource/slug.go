package cloudresource

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"google.golang.org/grpc"
)

// CheckSlugAvailability checks whether a slug is available for use within the
// scoped composite key (org, env, kind) via the
// CloudResourceQueryController.CheckSlugAvailability RPC.
func CheckSlugAvailability(ctx context.Context, serverAddress string, org, env string, kind cloudresourcekind.CloudResourceKind, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := cloudresourcev1.NewCloudResourceQueryControllerClient(conn)
			resp, err := client.CheckSlugAvailability(ctx, &cloudresourcev1.CloudResourceSlugAvailabilityCheckRequest{
				Org:               org,
				Env:               env,
				CloudResourceKind: kind,
				Slug:              slug,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("slug availability for %q in org %q env %q", slug, org, env))
			}
			return domains.MarshalJSON(resp)
		})
}
