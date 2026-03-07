package identity

import (
	"context"
	"fmt"

	identityaccountv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/identityaccount/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Invite creates a user invitation for the given email to join an
// organization with the specified IAM roles, via
// UserInvitationCommandController.Create.
func Invite(ctx context.Context, serverAddress, org, email string, iamRoleIDs []string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := identityaccountv1.NewUserInvitationCommandControllerClient(conn)
			resp, err := client.Create(ctx, &identityaccountv1.CreateUserInvitationInput{
				Org:        org,
				Email:      email,
				IamRoleIds: iamRoleIDs,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("user invitation for %q in org %q", email, org))
			}
			return domains.MarshalJSON(resp)
		})
}
