package domains

import (
	"context"

	"github.com/plantoncloud/mcp-server-planton/internal/auth"
	plantongrpc "github.com/plantoncloud/mcp-server-planton/internal/grpc"
	"google.golang.org/grpc"
)

// WithConnection creates an authenticated gRPC connection with timeout,
// passes it to fn, and ensures cleanup. This eliminates the repetitive
// connect/auth/timeout/defer pattern in every domain function.
//
// The API key is read from ctx via [auth.APIKey]. The connection targets
// serverAddress using the transport rules in [plantongrpc.NewConnection].
// The context passed to fn has a deadline of [plantongrpc.DefaultRPCTimeout].
func WithConnection(ctx context.Context, serverAddress string,
	fn func(ctx context.Context, conn *grpc.ClientConn) (string, error),
) (string, error) {
	conn, err := plantongrpc.NewConnection(serverAddress, auth.APIKey(ctx))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	rpcCtx, cancel := context.WithTimeout(ctx, plantongrpc.DefaultRPCTimeout)
	defer cancel()

	return fn(rpcCtx, conn)
}
