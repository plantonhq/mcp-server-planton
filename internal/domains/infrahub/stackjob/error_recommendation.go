package stackjob

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
	"google.golang.org/grpc"
)

// GetErrorRecommendation retrieves an AI-generated recommendation for resolving
// a specific error from a failed stack job via the
// StackJobQueryController.GetErrorResolutionRecommendation RPC.
//
// The backend forwards the error message to an LLM and returns a plain-text
// recommendation. Authorization is skipped for this RPC — any authenticated
// user may call it.
func GetErrorRecommendation(ctx context.Context, serverAddress, stackJobID, errorMessage string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobQueryControllerClient(conn)
			resp, err := client.GetErrorResolutionRecommendation(ctx, &stackjobv1.GetErrorResolutionRecommendationInput{
				StackJobId:   stackJobID,
				ErrorMessage: errorMessage,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("error recommendation for stack job %q", stackJobID))
			}
			return resp.GetValue(), nil
		})
}
