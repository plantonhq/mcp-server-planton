package infraproject

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	infrahubsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/infrahub"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// SearchInput holds the validated parameters for searching infra projects.
type SearchInput struct {
	Org        string
	Env        string
	SearchText string
	PageNum    int32
	PageSize   int32
}

// Search queries infra projects via the
// InfraHubSearchQueryController.SearchInfraProjects RPC.
//
// The search supports free-text filtering, org scoping (required), optional
// environment filtering, and pagination. Results are lightweight search
// records — use get_infra_project with an ID from the results to retrieve
// the full project.
func Search(ctx context.Context, serverAddress string, input SearchInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	req := &infrahubsearch.SearchInfraProjectsRequest{
		Org:        input.Org,
		Env:        input.Env,
		SearchText: input.SearchText,
		PageInfo: &rpc.PageInfo{
			Num:  pageNum - 1,
			Size: pageSize,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrahubsearch.NewInfraHubSearchQueryControllerClient(conn)
			resp, err := client.SearchInfraProjects(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("infra projects in org %q", input.Org))
			}
			return domains.MarshalJSON(resp)
		})
}
