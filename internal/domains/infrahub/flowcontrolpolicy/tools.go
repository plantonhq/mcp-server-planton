package flowcontrolpolicy

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource/apiresourcekind"
)

// selectorKindResolver maps user-supplied selector kind strings
// (e.g. "organization", "environment", "platform", or a cloud resource kind)
// to the ApiResourceKind enum.
var selectorKindResolver = domains.NewEnumResolver[apiresourcekind.ApiResourceKind](
	apiresourcekind.ApiResourceKind_value,
	"selector kind",
	"api_resource_kind_unspecified",
)

// ---------------------------------------------------------------------------
// apply_flow_control_policy
// ---------------------------------------------------------------------------

// ApplyFlowControlPolicyInput defines the parameters for the
// apply_flow_control_policy tool.
type ApplyFlowControlPolicyInput struct {
	PolicyID                              string `json:"policy_id,omitempty"                                  jsonschema:"Existing policy ID for update. Omit to create a new policy."`
	Name                                  string `json:"name,omitempty"                                       jsonschema:"Human-readable name for the policy. Required when creating a new policy."`
	SelectorKind                          string `json:"selector_kind"                                        jsonschema:"required,Resource kind the policy targets (organization, environment, platform, or a cloud resource kind like GcpCloudSqlDatabase)."`
	SelectorID                            string `json:"selector_id"                                          jsonschema:"required,ID of the target resource."`
	IsManual                              bool   `json:"is_manual,omitempty"                                   jsonschema:"Require manual approval before every stack job execution. Stack jobs are still created on lifecycle events but wait for approval."`
	DisableOnLifecycleEvents              bool   `json:"disable_on_lifecycle_events,omitempty"                  jsonschema:"Disable automatic stack jobs on create/update/delete/restore lifecycle events. Stack jobs can still be created manually."`
	SkipRefresh                           bool   `json:"skip_refresh,omitempty"                                jsonschema:"Skip the refresh step during preview, update, and destroy operations."`
	PreviewBeforeUpdateOrDestroy          bool   `json:"preview_before_update_or_destroy,omitempty"             jsonschema:"Run a preview step before update or destroy operations. Defaults to true in the backend."`
	PauseBetweenPreviewAndUpdateOrDestroy bool   `json:"pause_between_preview_and_update_or_destroy,omitempty"  jsonschema:"Pause for approval after preview completes but before update/destroy executes. Only effective when preview_before_update_or_destroy is true."`
}

// ApplyTool returns the MCP tool definition for apply_flow_control_policy.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_flow_control_policy",
		Description: "Create or update a flow control policy that governs how stack jobs " +
			"(Pulumi/Terraform operations) execute for a given scope. " +
			"Flow control policies can require manual approval, disable automatic " +
			"lifecycle-triggered jobs, skip refresh, or add preview-before-update pauses. " +
			"Policies can target an organization, environment, platform, or a specific " +
			"cloud resource. Provide policy_id to update an existing policy; omit it to " +
			"create a new one.",
	}
}

// ApplyHandler returns the typed tool handler for apply_flow_control_policy.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyFlowControlPolicyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyFlowControlPolicyInput) (*mcp.CallToolResult, any, error) {
		if input.SelectorKind == "" {
			return nil, nil, fmt.Errorf("'selector_kind' is required")
		}
		if input.SelectorID == "" {
			return nil, nil, fmt.Errorf("'selector_id' is required")
		}

		text, err := Apply(ctx, serverAddress, ApplyInput{
			PolicyID:                              input.PolicyID,
			Name:                                  input.Name,
			SelectorKind:                          input.SelectorKind,
			SelectorID:                            input.SelectorID,
			IsManual:                              input.IsManual,
			DisableOnLifecycleEvents:              input.DisableOnLifecycleEvents,
			SkipRefresh:                           input.SkipRefresh,
			PreviewBeforeUpdateOrDestroy:          input.PreviewBeforeUpdateOrDestroy,
			PauseBetweenPreviewAndUpdateOrDestroy: input.PauseBetweenPreviewAndUpdateOrDestroy,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_flow_control_policy
// ---------------------------------------------------------------------------

// GetFlowControlPolicyInput defines the parameters for the
// get_flow_control_policy tool.
type GetFlowControlPolicyInput struct {
	PolicyID     string `json:"policy_id,omitempty"     jsonschema:"Policy ID for direct lookup. If provided, selector_kind and selector_id are ignored."`
	SelectorKind string `json:"selector_kind,omitempty" jsonschema:"Resource kind (organization, environment, platform, or a cloud resource kind). Use with selector_id to look up the policy assigned to that scope."`
	SelectorID   string `json:"selector_id,omitempty"   jsonschema:"ID of the target resource. Use with selector_kind."`
}

// GetTool returns the MCP tool definition for get_flow_control_policy.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_flow_control_policy",
		Description: "Retrieve a flow control policy by its ID or by its scope selector. " +
			"Provide policy_id for a direct lookup, or provide selector_kind and selector_id " +
			"to find the policy assigned to a specific organization, environment, platform, " +
			"or cloud resource.",
	}
}

// GetHandler returns the typed tool handler for get_flow_control_policy.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetFlowControlPolicyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetFlowControlPolicyInput) (*mcp.CallToolResult, any, error) {
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
// delete_flow_control_policy
// ---------------------------------------------------------------------------

// DeleteFlowControlPolicyInput defines the parameters for the
// delete_flow_control_policy tool.
type DeleteFlowControlPolicyInput struct {
	PolicyID string `json:"policy_id" jsonschema:"required,The flow control policy ID to delete."`
}

// DeleteTool returns the MCP tool definition for delete_flow_control_policy.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_flow_control_policy",
		Description: "Delete a flow control policy by ID. " +
			"Once deleted, stack jobs for the targeted scope revert to default behavior " +
			"(or inherit from a higher-level policy if one exists). " +
			"Use get_flow_control_policy to find the policy ID first.",
	}
}

// DeleteHandler returns the typed tool handler for delete_flow_control_policy.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteFlowControlPolicyInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteFlowControlPolicyInput) (*mcp.CallToolResult, any, error) {
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
