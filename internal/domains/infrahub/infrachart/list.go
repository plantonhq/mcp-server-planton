package infrachart

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource/apiresourcekind"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	infrachartv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrachart/v1"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// ListInput holds the validated filter values for listing infra charts.
type ListInput struct {
	Org      string
	Env      string
	PageNum  int32
	PageSize int32
}

// List queries infra charts via the InfraChartQueryController.Find RPC.
//
// The Find RPC supports pagination and optional org/env scoping. The
// ApiResourceKind is hard-coded to infra_chart since this tool is
// exclusively for infra chart discovery.
func List(ctx context.Context, serverAddress string, input ListInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	// The Find RPC uses 0-based page numbers; our tool API is 1-based
	// for consistency with list_stack_jobs and human convention.
	req := &apiresource.FindApiResourcesRequest{
		Page: &rpc.PageInfo{
			Num:  pageNum - 1,
			Size: pageSize,
		},
		Kind: apiresourcekind.ApiResourceKind_infra_chart,
		Org:  input.Org,
		Env:  input.Env,
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrachartv1.NewInfraChartQueryControllerClient(conn)
			resp, err := client.Find(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, "infra charts")
			}
			return domains.MarshalJSON(resp)
		})
}
