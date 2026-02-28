package pipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/workflow"
	pipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/pipeline/v1"
	"google.golang.org/grpc"
)

// ResolveGate approves or rejects a manual gate for a deployment task
// within a pipeline via the PipelineCommandController.ResolveManualGate RPC.
//
// The decision string must be "approve" or "reject".
func ResolveGate(ctx context.Context, serverAddress, pipelineID, deploymentTaskName, decision string) (string, error) {
	d, err := resolveDecision(decision)
	if err != nil {
		return "", err
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := pipelinev1.NewPipelineCommandControllerClient(conn)
			_, err := client.ResolveManualGate(ctx, &pipelinev1.ResolvePipelineManualGateRequest{
				PipelineId:         pipelineID,
				DeploymentTaskName: deploymentTaskName,
				ManualGateDecision: d,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("resolve gate for pipeline %q task %q", pipelineID, deploymentTaskName))
			}
			return fmt.Sprintf("Manual gate for deployment task %q in pipeline %q resolved: %s.", deploymentTaskName, pipelineID, decision), nil
		})
}

// resolveDecision maps the user-facing decision strings "approve" and "reject"
// to the corresponding WorkflowStepManualGateDecision proto enum values.
func resolveDecision(s string) (workflow.WorkflowStepManualGateDecision, error) {
	switch s {
	case "approve":
		return workflow.WorkflowStepManualGateDecision_yes, nil
	case "reject":
		return workflow.WorkflowStepManualGateDecision_no, nil
	default:
		return 0, fmt.Errorf("invalid decision %q — must be \"approve\" or \"reject\"", s)
	}
}
