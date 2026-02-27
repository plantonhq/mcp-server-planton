package deploymentcomponent

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	infrahubsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/infrahub"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// SearchInput holds the validated parameters for searching deployment components.
type SearchInput struct {
	SearchText string
	Provider   string
	PageNum    int32
	PageSize   int32
}

// Search queries deployment components via the
// InfraHubSearchQueryController.SearchDeploymentComponentsByFilter RPC.
//
// This is a public endpoint — no organization context is required.
// Results are lightweight search records containing component IDs,
// names, and metadata. Use get_deployment_component with an ID or kind
// from the results to retrieve full details.
func Search(ctx context.Context, serverAddress string, input SearchInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	var providers []cloudresourcekind.CloudResourceProvider
	if input.Provider != "" {
		p, err := domains.ResolveProvider(input.Provider)
		if err != nil {
			return "", err
		}
		providers = []cloudresourcekind.CloudResourceProvider{p}
	}

	req := &infrahubsearch.SearchDeploymentComponentsByFilterInput{
		SearchText: input.SearchText,
		PageInfo: &rpc.PageInfo{
			Num:  pageNum - 1,
			Size: pageSize,
		},
		Providers: providers,
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrahubsearch.NewInfraHubSearchQueryControllerClient(conn)
			resp, err := client.SearchDeploymentComponentsByFilter(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, "deployment components")
			}
			return domains.MarshalJSON(resp)
		})
}
