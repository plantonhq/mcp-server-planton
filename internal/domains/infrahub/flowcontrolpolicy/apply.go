package flowcontrolpolicy

import (
	"context"
	"fmt"

	apiresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	flowcontrolpolicyv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/infrahub/flowcontrolpolicy/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ApplyInput holds the validated parameters for an apply operation.
type ApplyInput struct {
	PolicyID                              string
	Name                                  string
	SelectorKind                          string
	SelectorID                            string
	IsManual                              bool
	DisableOnLifecycleEvents              bool
	SkipRefresh                           bool
	PreviewBeforeUpdateOrDestroy          bool
	PauseBetweenPreviewAndUpdateOrDestroy bool
}

// Apply creates or updates a flow control policy via
// FlowControlPolicyCommandController.Apply.
//
// The handler constructs the full FlowControlPolicy protobuf message from the
// typed input. If policy_id is set, the backend treats it as an update;
// otherwise it creates a new policy.
func Apply(ctx context.Context, serverAddress string, input ApplyInput) (string, error) {
	kind, err := selectorKindResolver.Resolve(input.SelectorKind)
	if err != nil {
		return "", err
	}

	policy := &flowcontrolpolicyv1.FlowControlPolicy{
		ApiVersion: "infra-hub.planton.ai/v1",
		Kind:       "FlowControlPolicy",
		Metadata: &apiresource.ApiResourceMetadata{
			Id:   input.PolicyID,
			Name: input.Name,
		},
		Spec: &flowcontrolpolicyv1.FlowControlPolicySpec{
			Selector: &apiresource.ApiResourceSelector{
				Kind: kind,
				Id:   input.SelectorID,
			},
			FlowControl: &flowcontrolpolicyv1.StackJobFlowControl{
				IsManual:                              input.IsManual,
				DisableOnLifecycleEvents:              input.DisableOnLifecycleEvents,
				SkipRefresh:                           input.SkipRefresh,
				PreviewBeforeUpdateOrDestroy:          input.PreviewBeforeUpdateOrDestroy,
				PauseBetweenPreviewAndUpdateOrDestroy: input.PauseBetweenPreviewAndUpdateOrDestroy,
			},
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := flowcontrolpolicyv1.NewFlowControlPolicyCommandControllerClient(conn)
			resp, err := client.Apply(ctx, policy)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf(
					"flow control policy for %s %q", input.SelectorKind, input.SelectorID))
			}
			return domains.MarshalJSON(resp)
		})
}
