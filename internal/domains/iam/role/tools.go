package role

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// get_iam_role
// ---------------------------------------------------------------------------

// GetIamRoleInput defines the parameters for the get_iam_role tool.
type GetIamRoleInput struct {
	RoleID string `json:"role_id" jsonschema:"required,IAM role ID (e.g. iamr-xxx)."`
}

// GetTool returns the MCP tool definition for get_iam_role.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_iam_role",
		Description: "Get an IAM role definition by ID. " +
			"Returns the role including its code, name, description, resource kind, and principal type. " +
			"Useful for understanding what a role grants before assigning it.",
	}
}

// GetHandler returns the typed tool handler for get_iam_role.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetIamRoleInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetIamRoleInput) (*mcp.CallToolResult, any, error) {
		if input.RoleID == "" {
			return nil, nil, fmt.Errorf("'role_id' is required")
		}
		text, err := Get(ctx, serverAddress, input.RoleID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_iam_roles_for_resource_kind
// ---------------------------------------------------------------------------

// ListIamRolesForResourceKindInput defines the parameters for the
// list_iam_roles_for_resource_kind tool.
type ListIamRolesForResourceKindInput struct {
	ResourceKind  string `json:"resource_kind"  jsonschema:"required,API resource kind (e.g. 'organization', 'environment', 'aws_credential')."`
	PrincipalType string `json:"principal_type" jsonschema:"required,Principal type (e.g. 'identity_account', 'team')."`
}

// ListForResourceKindTool returns the MCP tool definition for
// list_iam_roles_for_resource_kind.
func ListForResourceKindTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_iam_roles_for_resource_kind",
		Description: "List all IAM roles available for a given resource kind and principal type. " +
			"For example, to see what roles can be granted to users on an organization, " +
			"use resource_kind='organization' and principal_type='identity_account'. " +
			"Returns role IDs needed for create_iam_policy and invite_member.",
	}
}

// ListForResourceKindHandler returns the typed tool handler for
// list_iam_roles_for_resource_kind.
func ListForResourceKindHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListIamRolesForResourceKindInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListIamRolesForResourceKindInput) (*mcp.CallToolResult, any, error) {
		if input.ResourceKind == "" {
			return nil, nil, fmt.Errorf("'resource_kind' is required")
		}
		if input.PrincipalType == "" {
			return nil, nil, fmt.Errorf("'principal_type' is required")
		}
		text, err := ListForResourceKind(ctx, serverAddress, input.ResourceKind, input.PrincipalType)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
