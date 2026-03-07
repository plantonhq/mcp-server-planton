package policy

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	iampolicyv2 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/iampolicy/v2"
	"google.golang.org/grpc"
)

// CheckAuthorization checks whether a principal is authorized for a
// specific relation on a resource, via
// IamPolicyV2QueryController.CheckAuthorization.
func CheckAuthorization(ctx context.Context, serverAddress, principalKind, principalID, resourceKind, resourceID, relation string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := iampolicyv2.NewIamPolicyV2QueryControllerClient(conn)
			resp, err := client.CheckAuthorization(ctx, &iampolicyv2.CheckAuthorizationInput{
				Policy: &iampolicyv2.IamPolicySpec{
					Principal: &iampolicyv2.ApiResourceRef{Kind: principalKind, Id: principalID},
					Resource:  &iampolicyv2.ApiResourceRef{Kind: resourceKind, Id: resourceID},
					Relation:  relation,
				},
			})
			if err != nil {
				return "", domains.RPCError(err,
					fmt.Sprintf("authorization check: %s:%s -> %s -> %s:%s",
						principalKind, principalID, relation, resourceKind, resourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
