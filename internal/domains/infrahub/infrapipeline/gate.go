package infrapipeline

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/workflow"
	infrapipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrapipeline/v1"
	"google.golang.org/grpc"
)

// ResolveEnvGate approves or rejects a manual gate for an entire deployment
// environment within an infra pipeline via the
// InfraPipelineCommandController.ResolveEnvironmentManualGate RPC.
//
// The decision string must be "approve" or "reject".
func ResolveEnvGate(ctx context.Context, serverAddress, pipelineID, env, decision string) (string, error) {
	d, err := resolveDecision(decision)
	if err != nil {
		return "", err
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrapipelinev1.NewInfraPipelineCommandControllerClient(conn)
			_, err := client.ResolveEnvironmentManualGate(ctx, &infrapipelinev1.ResolveInfraPipelineEnvironmentManualGateRequest{
				InfraPipelineId:    pipelineID,
				Env:                env,
				ManualGateDecision: d,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("resolve env gate for pipeline %q env %q", pipelineID, env))
			}
			return fmt.Sprintf("Manual gate for environment %q in pipeline %q resolved: %s.", env, pipelineID, decision), nil
		})
}

// ResolveNodeGate approves or rejects a manual gate for a specific DAG node
// within an infra pipeline via the
// InfraPipelineCommandController.ResolveNodeManualGate RPC.
//
// The nodeID format is "{CloudResourceKind}/{slug}"
// (e.g. "KubernetesOpenFga/fga-gcp-dev").
// The decision string must be "approve" or "reject".
func ResolveNodeGate(ctx context.Context, serverAddress, pipelineID, env, nodeID, decision string) (string, error) {
	d, err := resolveDecision(decision)
	if err != nil {
		return "", err
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrapipelinev1.NewInfraPipelineCommandControllerClient(conn)
			_, err := client.ResolveNodeManualGate(ctx, &infrapipelinev1.ResolveInfraPipelineNodeManualGateRequest{
				InfraPipelineId:    pipelineID,
				Env:                env,
				NodeId:             nodeID,
				ManualGateDecision: d,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("resolve node gate for pipeline %q node %q in env %q", pipelineID, nodeID, env))
			}
			return fmt.Sprintf("Manual gate for node %q in environment %q of pipeline %q resolved: %s.", nodeID, env, pipelineID, decision), nil
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
