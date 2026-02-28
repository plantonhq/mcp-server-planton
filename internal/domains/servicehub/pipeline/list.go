package pipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	pipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/pipeline/v1"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// ListInput holds the validated filter values for listing pipelines.
type ListInput struct {
	Org       string
	ServiceID string
	Envs      []string
	PageNum   int32
	PageSize  int32
}

// List queries pipelines via the
// PipelineQueryController.ListByFilters RPC.
//
// Pipelines can be scoped to an organization (required) and optionally
// narrowed to a specific service and/or environments. Results are paginated.
func List(ctx context.Context, serverAddress string, input ListInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	req := &pipelinev1.ListPipelinesByFiltersInput{
		Org:       input.Org,
		ServiceId: input.ServiceID,
		Envs:      input.Envs,
		PageInfo: &rpc.PageInfo{
			Num:  pageNum - 1,
			Size: pageSize,
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := pipelinev1.NewPipelineQueryControllerClient(conn)
			resp, err := client.ListByFilters(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("pipelines in org %q", input.Org))
			}
			return domains.MarshalJSON(resp)
		})
}
