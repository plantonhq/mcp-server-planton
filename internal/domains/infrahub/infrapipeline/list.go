package infrapipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	infrapipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrapipeline/v1"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// ListInput holds the validated filter values for listing infra pipelines.
type ListInput struct {
	Org            string
	InfraProjectID string
	PageNum        int32
	PageSize       int32
}

// List queries infra pipelines via the
// InfraPipelineQueryController.ListByFilters RPC.
//
// Pipelines can be scoped to an organization (required) and optionally
// narrowed to a specific infra project. Results are paginated.
func List(ctx context.Context, serverAddress string, input ListInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	req := &infrapipelinev1.ListInfraPipelinesByFiltersInput{
		Org:            input.Org,
		InfraProjectId: input.InfraProjectID,
		PageInfo: &rpc.PageInfo{
			Num:  pageNum - 1,
			Size: pageSize,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrapipelinev1.NewInfraPipelineQueryControllerClient(conn)
			resp, err := client.ListByFilters(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("infra pipelines in org %q", input.Org))
			}
			return domains.MarshalJSON(resp)
		})
}
