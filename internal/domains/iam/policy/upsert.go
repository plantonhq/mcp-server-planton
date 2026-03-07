package policy

import (
	"context"
	"fmt"

	iampolicyv2 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/iampolicy/v2"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Upsert declaratively syncs the relations a principal has on a resource.
// After this call the principal will have EXACTLY the specified relations —
// extra ones are removed and missing ones are added.
func Upsert(ctx context.Context, serverAddress, principalKind, principalID, resourceKind, resourceID string, relations []string) (string, error) {
	input := &iampolicyv2.UpsertIamPoliciesInput{
		Principal: &iampolicyv2.ApiResourceRef{Kind: principalKind, Id: principalID},
		Resource:  &iampolicyv2.ApiResourceRef{Kind: resourceKind, Id: resourceID},
		Relations: relations,
	}
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := iampolicyv2.NewIamPolicyV2CommandControllerClient(conn)
			resp, err := client.Upsert(ctx, input)
			if err != nil {
				return "", domains.RPCError(err,
					fmt.Sprintf("upsert IAM policies for %s:%s on %s:%s", principalKind, principalID, resourceKind, resourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
