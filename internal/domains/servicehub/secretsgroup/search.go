package secretsgroup

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	servicehubsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/servicehub"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// SearchInput holds the validated parameters for searching secret entries.
type SearchInput struct {
	Org        string
	SearchText string
	PageNum    int32
	PageSize   int32
}

// Search queries secret entries across all secrets groups in an org via the
// ServiceHubSearchQueryController.SearchSecrets RPC.
//
// Results are individual secret entries (not whole groups), each annotated
// with its parent group name and ID. Supports free-text filtering and
// pagination.
func Search(ctx context.Context, serverAddress string, input SearchInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	req := &servicehubsearch.SearchConfigEntriesRequest{
		Org:        input.Org,
		SearchText: input.SearchText,
		PageInfo: &rpc.PageInfo{
			Num:  pageNum - 1,
			Size: pageSize,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := servicehubsearch.NewServiceHubSearchQueryControllerClient(conn)
			resp, err := client.SearchSecrets(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("secret entries in org %q", input.Org))
			}
			return domains.MarshalJSON(resp)
		})
}
