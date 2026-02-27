package cloudresource

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	cloudresourcesearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/infrahub/cloudresource"
	"google.golang.org/grpc"
)

// List queries the search index for cloud resources visible to the caller
// within an organization, optionally filtered by environments, kinds, and
// free-text search.
//
// It calls CloudResourceSearchQueryController.GetCloudResourcesCanvasView,
// which returns lightweight search records grouped by environment and kind.
func List(ctx context.Context, serverAddress string, org string, envs []string, searchText string, kinds []cloudresourcekind.CloudResourceKind) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := cloudresourcesearch.NewCloudResourceSearchQueryControllerClient(conn)
			resp, err := client.GetCloudResourcesCanvasView(ctx, &cloudresourcesearch.ExploreCloudResourcesRequest{
				Org:        org,
				Envs:       envs,
				SearchText: searchText,
				Kinds:      kinds,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("cloud resources in org %q", org))
			}
			return domains.MarshalJSON(resp)
		})
}
