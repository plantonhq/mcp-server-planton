package identity

import (
	"context"

	identityaccountv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/identityaccount/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// WhoAmI retrieves the identity account associated with the current
// authentication token via IdentityAccountQueryController.WhoAmI.
func WhoAmI(ctx context.Context, serverAddress string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := identityaccountv1.NewIdentityAccountQueryControllerClient(conn)
			resp, err := client.WhoAmI(ctx, &emptypb.Empty{})
			if err != nil {
				return "", domains.RPCError(err, "current identity account")
			}
			return domains.MarshalJSON(resp)
		})
}
