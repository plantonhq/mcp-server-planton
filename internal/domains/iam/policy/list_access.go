package policy

import (
	"context"
	"fmt"

	iampolicyv2 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/iampolicy/v2"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ListResourceAccess returns all principals and their access grants for a
// specific resource, via IamPolicyV2QueryController.ListResourceAccessByPrincipal.
func ListResourceAccess(ctx context.Context, serverAddress, resourceKind, resourceID string, includeInherited bool) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := iampolicyv2.NewIamPolicyV2QueryControllerClient(conn)
			resp, err := client.ListResourceAccessByPrincipal(ctx, &iampolicyv2.ListResourceAccessInput{
				Resource:         &iampolicyv2.ApiResourceRef{Kind: resourceKind, Id: resourceID},
				IncludeInherited: includeInherited,
			})
			if err != nil {
				return "", domains.RPCError(err,
					fmt.Sprintf("access list for %s:%s", resourceKind, resourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
