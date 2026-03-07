package identity

import (
	"context"
	"fmt"

	identityaccountv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/identityaccount/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Get retrieves an identity account by ID via
// IdentityAccountQueryController.Get.
func Get(ctx context.Context, serverAddress, id string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := identityaccountv1.NewIdentityAccountQueryControllerClient(conn)
			resp, err := client.Get(ctx, &identityaccountv1.IdentityAccountId{Value: id})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("identity account %q", id))
			}
			return domains.MarshalJSON(resp)
		})
}

// GetByEmail retrieves an identity account by email via
// IdentityAccountQueryController.GetByEmail.
func GetByEmail(ctx context.Context, serverAddress, email string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := identityaccountv1.NewIdentityAccountQueryControllerClient(conn)
			resp, err := client.GetByEmail(ctx, &identityaccountv1.IdentityAccountEmail{Value: email})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("identity account with email %q", email))
			}
			return domains.MarshalJSON(resp)
		})
}
