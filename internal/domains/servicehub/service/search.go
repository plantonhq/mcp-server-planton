package service

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource/apiresourcekind"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	apiresourcesearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/apiresource"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// SearchInput holds the validated parameters for searching services.
type SearchInput struct {
	Org        string
	SearchText string
	PageNum    int32
	PageSize   int32
}

// Search queries services via the
// ApiResourceSearchQueryController.SearchByKind RPC, filtered to the
// "service" API resource kind.
//
// The search supports free-text filtering, org scoping (required), and
// pagination. Results are lightweight search records — use get_service with
// an ID from the results to retrieve the full service.
func Search(ctx context.Context, serverAddress string, input SearchInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	req := &apiresourcesearch.SearchApiResourcesByKindInput{
		Org:             input.Org,
		ApiResourceKind: apiresourcekind.ApiResourceKind_service,
		SearchText:      input.SearchText,
		PageInfo: &rpc.PageInfo{
			Num:  pageNum - 1,
			Size: pageSize,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := apiresourcesearch.NewApiResourceSearchQueryControllerClient(conn)
			resp, err := client.SearchByKind(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("services in org %q", input.Org))
			}
			return domains.MarshalJSON(resp)
		})
}
