package apikey

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apikeyv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/iam/apikey/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// List retrieves all API keys belonging to the authenticated user
// via ApiKeyQueryController.FindAll.
func List(ctx context.Context, serverAddress string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := apikeyv1.NewApiKeyQueryControllerClient(conn)
			resp, err := client.FindAll(ctx, &emptypb.Empty{})
			if err != nil {
				return "", domains.RPCError(err, "API keys")
			}
			return domains.MarshalJSON(resp)
		})
}
