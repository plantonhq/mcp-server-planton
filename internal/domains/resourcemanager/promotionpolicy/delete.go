package promotionpolicy

import (
	"context"
	"fmt"

	apiresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	promotionpolicyv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/resourcemanager/promotionpolicy/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Delete removes a promotion policy by ID via
// PromotionPolicyCommandController.Delete.
func Delete(ctx context.Context, serverAddress, policyID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := promotionpolicyv1.NewPromotionPolicyCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &apiresource.ApiResourceId{Value: policyID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("promotion policy %q", policyID))
			}
			return domains.MarshalJSON(resp)
		})
}
