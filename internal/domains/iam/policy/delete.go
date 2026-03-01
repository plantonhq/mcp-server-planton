package policy

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	iampolicyv2 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/iam/iampolicy/v2"
	"google.golang.org/grpc"
)

// Delete revokes a specific access grant via
// IamPolicyV2CommandController.Delete. This is idempotent.
func Delete(ctx context.Context, serverAddress, principalKind, principalID, resourceKind, resourceID, relation string) (string, error) {
	spec := &iampolicyv2.IamPolicySpec{
		Principal: &iampolicyv2.ApiResourceRef{Kind: principalKind, Id: principalID},
		Resource:  &iampolicyv2.ApiResourceRef{Kind: resourceKind, Id: resourceID},
		Relation:  relation,
	}
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := iampolicyv2.NewIamPolicyV2CommandControllerClient(conn)
			resp, err := client.Delete(ctx, spec)
			if err != nil {
				return "", domains.RPCError(err,
					fmt.Sprintf("IAM policy: %s:%s -> %s -> %s:%s", principalKind, principalID, relation, resourceKind, resourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
