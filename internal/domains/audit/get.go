package audit

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresourceversionv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/audit/apiresourceversion/v1"
	"google.golang.org/grpc"
)

const defaultContextSize = 3

// Get retrieves a single resource version by its ID via the
// ApiResourceVersionQueryController.GetByIdWithContextSize RPC.
//
// The context_size parameter controls the number of surrounding lines included
// in the unified diff output, analogous to git diff -U<n>. When zero or
// negative, it defaults to 3 (the standard unified diff default).
func Get(ctx context.Context, serverAddress, versionID string, contextSize int32) (string, error) {
	if contextSize <= 0 {
		contextSize = defaultContextSize
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := apiresourceversionv1.NewApiResourceVersionQueryControllerClient(conn)
			resp, err := client.GetByIdWithContextSize(ctx, &apiresourceversionv1.ApiResourceVersionWithContextSizeInput{
				VersionId:   versionID,
				ContextSize: contextSize,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("resource version %q", versionID))
			}
			return domains.MarshalJSON(resp)
		})
}
