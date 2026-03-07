package policy

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource/apiresourcekind"
	rpc "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/rpc"
	iampolicyv2 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/iampolicy/v2"
	"google.golang.org/grpc"
)

// principalKindResolver maps user-supplied principal kind strings to the enum.
var principalKindResolver = domains.NewEnumResolver[apiresourcekind.ApiResourceKind](
	apiresourcekind.ApiResourceKind_value,
	"principal kind",
	"api_resource_kind_unspecified",
)

// ListPrincipals retrieves principals (identity accounts or teams) with
// access to an organization or environment, via
// IamPolicyV2QueryController.ListPrincipals.
func ListPrincipals(ctx context.Context, serverAddress, orgID, env, principalKind string, pageNumber, pageSize int32) (string, error) {
	pk, err := principalKindResolver.Resolve(principalKind)
	if err != nil {
		return "", err
	}

	input := &iampolicyv2.ListPrincipalsInput{
		OrgId:         orgID,
		Env:           env,
		PrincipalKind: pk,
	}
	if pageNumber > 0 || pageSize > 0 {
		input.PageInfo = &rpc.PageInfo{
			Num:  pageNumber,
			Size: pageSize,
		}
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := iampolicyv2.NewIamPolicyV2QueryControllerClient(conn)
			resp, err := client.ListPrincipals(ctx, input)
			if err != nil {
				return "", domains.RPCError(err,
					fmt.Sprintf("principals of kind %q in org %q", principalKind, orgID))
			}
			return domains.MarshalJSON(resp)
		})
}
