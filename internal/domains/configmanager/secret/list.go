package secret

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource/apiresourcekind"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
	secretv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secret/v1"
	"google.golang.org/grpc"
)

const (
	defaultPageNum  = 1
	defaultPageSize = 20
)

// ListInput holds the validated filter values for listing secrets.
type ListInput struct {
	Org      string
	Env      string
	PageNum  int32
	PageSize int32
}

// List queries secrets via the SecretQueryController.Find RPC.
//
// The Find RPC supports pagination and optional org/env scoping. The
// ApiResourceKind is hard-coded to secret since this tool is exclusively
// for secret discovery. Only metadata is returned — no secret values.
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
		Kind: apiresourcekind.ApiResourceKind_secret,
		Org:  input.Org,
		Env:  input.Env,
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := secretv1.NewSecretQueryControllerClient(conn)
			resp, err := client.Find(ctx, req)
			if err != nil {
				return "", domains.RPCError(err, "secrets")
			}
			return domains.MarshalJSON(resp)
		})
}
