package variable

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource/apiresourcekind"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	variablev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/variable/v1"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// ListInput holds the validated filter values for listing variables.
type ListInput struct {
	Org      string
	Env      string
	PageNum  int32
	PageSize int32
}

// List queries variables via the VariableQueryController.Find RPC.
//
// The Find RPC supports pagination and optional org/env scoping. The
// ApiResourceKind is hard-coded to variable since this tool is exclusively
// for variable discovery.
func List(ctx context.Context, serverAddress string, input ListInput) (string, error) {
	pageNum := input.PageNum
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	pageSize := input.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	req := &apiresource.FindApiResourcesRequest{
		Page: &rpc.PageInfo{
			Num:  pageNum - 1,
			Size: pageSize,
		},
		Kind: apiresourcekind.ApiResourceKind_variable,
		Org:  input.Org,
		Env:  input.Env,
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := variablev1.NewVariableQueryControllerClient(conn)
			resp, err := client.Find(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, "variables")
			}
			return domains.MarshalJSON(resp)
		})
}
