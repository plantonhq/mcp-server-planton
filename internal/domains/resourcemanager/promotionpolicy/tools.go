package promotionpolicy

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource/apiresourcekind"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// selectorKindResolver maps user-supplied selector kind strings
// (e.g. "organization", "platform") to the ApiResourceKind enum.
var selectorKindResolver = domains.NewEnumResolver[apiresourcekind.ApiResourceKind](
	apiresourcekind.ApiResourceKind_value,
	"selector kind",
	"api_resource_kind_unspecified",
)

// ---------------------------------------------------------------------------
// apply_promotion_policy
// ---------------------------------------------------------------------------

// EnvironmentNodeInput represents a single node in the promotion graph.
type EnvironmentNodeInput struct {
	Name string `json:"name" jsonschema:"required,Logical environment identifier (e.g. dev, staging, prod). Must be unique within the graph."`
}

// PromotionEdgeInput represents a directed edge between two environment nodes.
type PromotionEdgeInput struct {
	From           string `json:"from"                       jsonschema:"required,Source environment name. Must match a node name."`
	To             string `json:"to"                         jsonschema:"required,Target environment name. Must match a node name."`
	ManualApproval bool   `json:"manual_approval,omitempty"  jsonschema:"If true, deployment pauses for explicit approval before promoting along this edge."`
}

// ApplyPromotionPolicyInput defines the parameters for the
// apply_promotion_policy tool.
type ApplyPromotionPolicyInput struct {
	PolicyID     string                 `json:"policy_id,omitempty"    jsonschema:"Existing policy ID for update. Omit to create a new policy."`
	Name         string                 `json:"name,omitempty"         jsonschema:"Human-readable name for the policy. Required when creating a new policy."`
	SelectorKind string                 `json:"selector_kind"          jsonschema:"required,Resource kind the policy targets. Must be organization or platform."`
	SelectorID   string                 `json:"selector_id"            jsonschema:"required,ID of the target resource (organization ID or platform ID)."`
	Strict       bool                   `json:"strict,omitempty"       jsonschema:"When true every environment referenced by a pipeline must appear in the graph or the workflow fails. When false (default) unknown environments are appended alphabetically."`
	Nodes        []EnvironmentNodeInput `json:"nodes"                  jsonschema:"required,Environment nodes in the promotion graph. The graph must be a DAG."`
	Edges        []PromotionEdgeInput   `json:"edges"                  jsonschema:"required,Directed edges defining promotion paths between environment nodes."`
}

// ApplyTool returns the MCP tool definition for apply_promotion_policy.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_promotion_policy",
		Description: "Create or update a promotion policy that defines the order in which " +
			"deployments are promoted through environments (e.g. dev -> staging -> prod). " +
			"The policy is scoped to an organization or the platform. " +
			"Provide policy_id to update an existing policy; omit it to create a new one. " +
			"The graph must be a directed acyclic graph (DAG) — cycles are rejected.",
	}
}

// ApplyHandler returns the typed tool handler for apply_promotion_policy.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyPromotionPolicyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyPromotionPolicyInput) (*mcp.CallToolResult, any, error) {
		if input.SelectorKind == "" {
			return nil, nil, fmt.Errorf("'selector_kind' is required")
		}
		if input.SelectorID == "" {
			return nil, nil, fmt.Errorf("'selector_id' is required")
		}
		if len(input.Nodes) == 0 {
			return nil, nil, fmt.Errorf("'nodes' must contain at least one environment")
		}

		nodes := make([]NodeInput, len(input.Nodes))
		for i, n := range input.Nodes {
			if n.Name == "" {
				return nil, nil, fmt.Errorf("node at index %d has an empty name", i)
			}
			nodes[i] = NodeInput{Name: n.Name}
		}

		edges := make([]EdgeInput, len(input.Edges))
		for i, e := range input.Edges {
			if e.From == "" || e.To == "" {
				return nil, nil, fmt.Errorf("edge at index %d must have both 'from' and 'to'", i)
			}
			edges[i] = EdgeInput{From: e.From, To: e.To, ManualApproval: e.ManualApproval}
		}

		text, err := Apply(ctx, serverAddress, ApplyInput{
			PolicyID:     input.PolicyID,
			Name:         input.Name,
			SelectorKind: input.SelectorKind,
			SelectorID:   input.SelectorID,
			Strict:       input.Strict,
			Nodes:        nodes,
			Edges:        edges,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_promotion_policy
// ---------------------------------------------------------------------------

// GetPromotionPolicyInput defines the parameters for the
// get_promotion_policy tool.
type GetPromotionPolicyInput struct {
	PolicyID     string `json:"policy_id,omitempty"     jsonschema:"Policy ID for direct lookup. If provided, selector_kind and selector_id are ignored."`
	SelectorKind string `json:"selector_kind,omitempty" jsonschema:"Resource kind (organization or platform). Use with selector_id to look up the policy assigned to that scope."`
	SelectorID   string `json:"selector_id,omitempty"   jsonschema:"ID of the target resource (org ID or platform ID). Use with selector_kind."`
}

// GetTool returns the MCP tool definition for get_promotion_policy.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_promotion_policy",
		Description: "Retrieve a promotion policy by its ID or by its scope selector. " +
			"Provide policy_id for a direct lookup, or provide selector_kind and selector_id " +
			"to find the policy assigned to a specific organization or platform. " +
			"To resolve the effective policy with inheritance, use which_promotion_policy instead.",
	}
}

// GetHandler returns the typed tool handler for get_promotion_policy.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetPromotionPolicyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetPromotionPolicyInput) (*mcp.CallToolResult, any, error) {
		if input.PolicyID != "" {
			text, err := Get(ctx, serverAddress, input.PolicyID)
			if err != nil {
				return nil, nil, err
			}
			return domains.TextResult(text)
		}

		if input.SelectorKind == "" || input.SelectorID == "" {
			return nil, nil, fmt.Errorf("provide either 'policy_id' or both 'selector_kind' and 'selector_id'")
		}

		text, err := GetBySelector(ctx, serverAddress, input.SelectorKind, input.SelectorID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// which_promotion_policy
// ---------------------------------------------------------------------------

// WhichPromotionPolicyInput defines the parameters for the
// which_promotion_policy tool.
type WhichPromotionPolicyInput struct {
	SelectorKind string `json:"selector_kind" jsonschema:"required,Resource kind to resolve the effective policy for (e.g. organization or platform)."`
	SelectorID   string `json:"selector_id"   jsonschema:"required,ID of the target resource."`
}

// WhichTool returns the MCP tool definition for which_promotion_policy.
func WhichTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "which_promotion_policy",
		Description: "Resolve the effective promotion policy for a given scope. " +
			"Unlike get_promotion_policy, this applies inheritance: if an org-specific " +
			"policy exists it is returned; otherwise the platform default is returned. " +
			"Use this to answer 'what promotion rules actually apply to this organization?'",
	}
}

// WhichHandler returns the typed tool handler for which_promotion_policy.
func WhichHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *WhichPromotionPolicyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *WhichPromotionPolicyInput) (*mcp.CallToolResult, any, error) {
		if input.SelectorKind == "" {
			return nil, nil, fmt.Errorf("'selector_kind' is required")
		}
		if input.SelectorID == "" {
			return nil, nil, fmt.Errorf("'selector_id' is required")
		}
		text, err := WhichPolicy(ctx, serverAddress, input.SelectorKind, input.SelectorID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_promotion_policy
// ---------------------------------------------------------------------------

// DeletePromotionPolicyInput defines the parameters for the
// delete_promotion_policy tool.
type DeletePromotionPolicyInput struct {
	PolicyID string `json:"policy_id" jsonschema:"required,The promotion policy ID to delete."`
}

// DeleteTool returns the MCP tool definition for delete_promotion_policy.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_promotion_policy",
		Description: "Delete a promotion policy by ID. " +
			"Once deleted, the scope reverts to the platform default policy (if one exists). " +
			"Use get_promotion_policy or which_promotion_policy to find the policy ID first.",
	}
}

// DeleteHandler returns the typed tool handler for delete_promotion_policy.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeletePromotionPolicyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeletePromotionPolicyInput) (*mcp.CallToolResult, any, error) {
		if input.PolicyID == "" {
			return nil, nil, fmt.Errorf("'policy_id' is required")
		}
		text, err := Delete(ctx, serverAddress, input.PolicyID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
