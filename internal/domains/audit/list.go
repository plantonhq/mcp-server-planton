package audit

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource/apiresourcekind"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	apiresourceversionv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/audit/apiresourceversion/v1"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// ListInput holds the validated filter values for listing resource versions.
type ListInput struct {
	Kind       apiresourcekind.ApiResourceKind
	ResourceID string
	PageNum    int32
	PageSize   int32
}

// List queries resource versions via the
// ApiResourceVersionQueryController.ListByFilters RPC.
//
// Pagination follows the 1-based convention: the caller provides 1-based
// page numbers, and we convert to 0-based for the proto PageInfo.Num field.
func List(ctx context.Context, serverAddress string, input ListInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	req := &apiresourceversionv1.ListApiResourceVersionsInput{
		PageInfo: &rpc.PageInfo{
			Num:  pageNum - 1,
			Size: pageSize,
		},
		Kind:       input.Kind,
		ResourceId: input.ResourceID,
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := apiresourceversionv1.NewApiResourceVersionQueryControllerClient(conn)
			resp, err := client.ListByFilters(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("resource versions for %s %q", input.Kind, input.ResourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
