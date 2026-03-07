package promotionpolicy

import (
	"context"
	"fmt"

	apiresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	promotionpolicyv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/resourcemanager/promotionpolicy/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ApplyInput holds the validated parameters for an apply operation.
type ApplyInput struct {
	PolicyID     string
	Name         string
	SelectorKind string
	SelectorID   string
	Strict       bool
	Nodes        []NodeInput
	Edges        []EdgeInput
}

// NodeInput represents a single environment node in the promotion graph.
type NodeInput struct {
	Name string
}

// EdgeInput represents a directed edge between two environment nodes.
type EdgeInput struct {
	From           string
	To             string
	ManualApproval bool
}

// Apply creates or updates a promotion policy via
// PromotionPolicyCommandController.Apply.
//
// The handler constructs the full PromotionPolicy protobuf message from the
// typed input. If policy_id is set, the backend treats it as an update;
// otherwise it creates a new policy.
func Apply(ctx context.Context, serverAddress string, input ApplyInput) (string, error) {
	kind, err := selectorKindResolver.Resolve(input.SelectorKind)
	if err != nil {
		return "", err
	}

	nodes := make([]*promotionpolicyv1.EnvironmentNode, len(input.Nodes))
	for i, n := range input.Nodes {
		nodes[i] = &promotionpolicyv1.EnvironmentNode{Name: n.Name}
	}

	edges := make([]*promotionpolicyv1.PromotionEdge, len(input.Edges))
	for i, e := range input.Edges {
		edges[i] = &promotionpolicyv1.PromotionEdge{
			From:           e.From,
			To:             e.To,
			ManualApproval: e.ManualApproval,
		}
	}

	policy := &promotionpolicyv1.PromotionPolicy{
		ApiVersion: "resource-manager.planton.ai/v1",
		Kind:       "PromotionPolicy",
		Metadata: &apiresource.ApiResourceMetadata{
			Id:   input.PolicyID,
			Name: input.Name,
		},
		Spec: &promotionpolicyv1.PromotionPolicySpec{
			Selector: &apiresource.ApiResourceSelector{
				Kind: kind,
				Id:   input.SelectorID,
			},
			Strict: input.Strict,
			Graph: &promotionpolicyv1.PromotionGraph{
				Nodes: nodes,
				Edges: edges,
			},
		},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := promotionpolicyv1.NewPromotionPolicyCommandControllerClient(conn)
			resp, err := client.Apply(ctx, policy)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf(
					"promotion policy for %s %q", input.SelectorKind, input.SelectorID))
			}
			return domains.MarshalJSON(resp)
		})
}
