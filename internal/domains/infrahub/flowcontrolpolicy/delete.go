package flowcontrolpolicy

import (
	"context"
	"fmt"

	apiresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	flowcontrolpolicyv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/infrahub/flowcontrolpolicy/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// Delete removes a flow control policy by ID via
// FlowControlPolicyCommandController.Delete.
func Delete(ctx context.Context, serverAddress, policyID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := flowcontrolpolicyv1.NewFlowControlPolicyCommandControllerClient(conn)
			resp, err := client.Delete(ctx, &apiresource.ApiResourceId{Value: policyID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("flow control policy %q", policyID))
			}
			return domains.MarshalJSON(resp)
		})
}
