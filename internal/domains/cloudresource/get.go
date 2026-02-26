package cloudresource

import (
	"context"

	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"github.com/plantoncloud/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Get retrieves a cloud resource via the CloudResourceQueryController.
//
// Two identification paths are supported:
//   - ID path: calls Get(CloudResourceId) directly.
//   - Slug path: resolves the PascalCase kind to the proto enum, then calls
//     GetByOrgByEnvByKindBySlug with all four fields.
//
// The caller must validate the ResourceIdentifier before calling this function.
func Get(ctx context.Context, serverAddress string, id ResourceIdentifier) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			desc := describeIdentifier(id)
			client := cloudresourcev1.NewCloudResourceQueryControllerClient(conn)

			if id.ID != "" {
				cr, err := client.Get(ctx, &cloudresourcev1.CloudResourceId{Value: id.ID})
				if err != nil {
					return "", domains.RPCError(err, desc)
				}
				return domains.MarshalJSON(cr)
			}

			kind, err := resolveKind(id.Kind)
			if err != nil {
				return "", err
			}

			cr, err := client.GetByOrgByEnvByKindBySlug(ctx, &cloudresourcev1.CloudResourceByOrgByEnvByKindBySlugRequest{
				Org:               id.Org,
				Env:               id.Env,
				CloudResourceKind: kind,
				Slug:              id.Slug,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(cr)
		})
}
