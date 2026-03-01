package flowcontrolpolicy

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	flowcontrolpolicyv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/flowcontrolpolicy/v1"
	"google.golang.org/grpc"
)

// Get retrieves a flow control policy by ID via
// FlowControlPolicyQueryController.Get.
func Get(ctx context.Context, serverAddress, policyID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := flowcontrolpolicyv1.NewFlowControlPolicyQueryControllerClient(conn)
			resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: policyID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("flow control policy %q", policyID))
			}
			return domains.MarshalJSON(resp)
		})
}

// GetBySelector retrieves a flow control policy by its scope selector via
// FlowControlPolicyQueryController.GetBySelector.
func GetBySelector(ctx context.Context, serverAddress, selectorKind, selectorID string) (string, error) {
	kind, err := selectorKindResolver.Resolve(selectorKind)
	if err != nil {
		return "", err
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := flowcontrolpolicyv1.NewFlowControlPolicyQueryControllerClient(conn)
			resp, err := client.GetBySelector(ctx, &apiresource.ApiResourceSelector{
				Kind: kind,
				Id:   selectorID,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf(
					"flow control policy for %s %q", selectorKind, selectorID))
			}
			return domains.MarshalJSON(resp)
		})
}
