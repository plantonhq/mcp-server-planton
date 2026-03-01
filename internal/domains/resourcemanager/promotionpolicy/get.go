package promotionpolicy

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	promotionpolicyv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/promotionpolicy/v1"
	"google.golang.org/grpc"
)

// Get retrieves a promotion policy by ID via
// PromotionPolicyQueryController.Get.
func Get(ctx context.Context, serverAddress, policyID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := promotionpolicyv1.NewPromotionPolicyQueryControllerClient(conn)
			resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: policyID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("promotion policy %q", policyID))
			}
			return domains.MarshalJSON(resp)
		})
}

// GetBySelector retrieves a promotion policy by its scope selector via
// PromotionPolicyQueryController.GetBySelector.
func GetBySelector(ctx context.Context, serverAddress, selectorKind, selectorID string) (string, error) {
	kind, err := selectorKindResolver.Resolve(selectorKind)
	if err != nil {
		return "", err
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := promotionpolicyv1.NewPromotionPolicyQueryControllerClient(conn)
			resp, err := client.GetBySelector(ctx, &apiresource.ApiResourceSelector{
				Kind: kind,
				Id:   selectorID,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf(
					"promotion policy for %s %q", selectorKind, selectorID))
			}
			return domains.MarshalJSON(resp)
		})
}
