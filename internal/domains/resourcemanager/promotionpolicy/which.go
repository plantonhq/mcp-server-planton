package promotionpolicy

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	promotionpolicyv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/resourcemanager/promotionpolicy/v1"
	"google.golang.org/grpc"
)

// WhichPolicy resolves the effective promotion policy for a given scope via
// PromotionPolicyQueryController.WhichPolicy.
//
// The backend applies inheritance: if an org-specific policy exists, it is
// returned; otherwise the platform default is returned. This makes it the
// right tool for answering "what promotion rules actually apply here?"
func WhichPolicy(ctx context.Context, serverAddress, selectorKind, selectorID string) (string, error) {
	kind, err := selectorKindResolver.Resolve(selectorKind)
	if err != nil {
		return "", err
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := promotionpolicyv1.NewPromotionPolicyQueryControllerClient(conn)
			resp, err := client.WhichPolicy(ctx, &apiresource.ApiResourceSelector{
				Kind: kind,
				Id:   selectorID,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf(
					"effective promotion policy for %s %q", selectorKind, selectorID))
			}
			return domains.MarshalJSON(resp)
		})
}
