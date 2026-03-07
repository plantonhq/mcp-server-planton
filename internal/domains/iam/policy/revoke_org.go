package policy

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	iampolicyv2 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/iampolicy/v2"
	"google.golang.org/grpc"
)

// RevokeOrgAccess removes ALL access a user has to an organization and its
// child resources, via IamPolicyV2CommandController.RevokeOrgAccess.
func RevokeOrgAccess(ctx context.Context, serverAddress, identityAccountID, organizationID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := iampolicyv2.NewIamPolicyV2CommandControllerClient(conn)
			resp, err := client.RevokeOrgAccess(ctx, &iampolicyv2.RevokeOrgAccessInput{
				IdentityAccountId: identityAccountID,
				OrganizationId:    organizationID,
			})
			if err != nil {
				return "", domains.RPCError(err,
					fmt.Sprintf("revoking org access for user %q from org %q", identityAccountID, organizationID))
			}
			return domains.MarshalJSON(resp)
		})
}
