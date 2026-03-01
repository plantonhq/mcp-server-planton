package identity

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	identityaccountv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/iam/identityaccount/v1"
	"google.golang.org/grpc"
)

// invitationStatusResolver maps user-friendly status strings to the proto enum.
var invitationStatusResolver = domains.NewEnumResolver[identityaccountv1.UserInvitationStatusType](
	identityaccountv1.UserInvitationStatusType_value,
	"invitation status",
	"user_invitation_status_type_unspecified",
)

// ListInvitations retrieves user invitations for an organization filtered
// by status, via UserInvitationQueryController.FindByOrgByStatus.
func ListInvitations(ctx context.Context, serverAddress, org, status string) (string, error) {
	statusEnum, err := invitationStatusResolver.Resolve(status)
	if err != nil {
		return "", err
	}
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := identityaccountv1.NewUserInvitationQueryControllerClient(conn)
			resp, err := client.FindByOrgByStatus(ctx, &identityaccountv1.FindUserInvitationsByOrgByStatusInput{
				Org:    org,
				Status: statusEnum,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("invitations for org %q with status %q", org, status))
			}
			return domains.MarshalJSON(resp)
		})
}
