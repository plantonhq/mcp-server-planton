package serviceaccount

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// create_service_account
// ---------------------------------------------------------------------------

// CreateServiceAccountInput defines the parameters for create_service_account.
type CreateServiceAccountInput struct {
	Org         string `json:"org"                    jsonschema:"required,Organization identifier."`
	DisplayName string `json:"display_name"           jsonschema:"required,Human-readable display name."`
	Description string `json:"description,omitempty"   jsonschema:"Optional description of the service account's purpose."`
}

// CreateTool returns the MCP tool definition for create_service_account.
func CreateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_service_account",
		Description: "Create a new service account for programmatic API access. " +
			"A backing machine identity is auto-created. " +
			"After creation, use create_service_account_key to generate an API key.",
	}
}

// CreateHandler returns the typed tool handler for create_service_account.
func CreateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CreateServiceAccountInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CreateServiceAccountInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.DisplayName == "" {
			return nil, nil, fmt.Errorf("'display_name' is required")
		}
		text, err := Create(ctx, serverAddress, input.Org, input.DisplayName, input.Description)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_service_account
// ---------------------------------------------------------------------------

// GetServiceAccountInput defines the parameters for get_service_account.
type GetServiceAccountInput struct {
	ID string `json:"id" jsonschema:"required,Service account ID."`
}

// GetTool returns the MCP tool definition for get_service_account.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_service_account",
		Description: "Retrieve a service account by ID. Returns metadata, display name, and description.",
	}
}

// GetHandler returns the typed tool handler for get_service_account.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetServiceAccountInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetServiceAccountInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := Get(ctx, serverAddress, input.ID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// update_service_account
// ---------------------------------------------------------------------------

// UpdateServiceAccountInput defines the parameters for update_service_account.
type UpdateServiceAccountInput struct {
	ID          string `json:"id"                     jsonschema:"required,Service account ID."`
	DisplayName string `json:"display_name,omitempty"  jsonschema:"New display name. Leave empty to keep current."`
	Description string `json:"description,omitempty"   jsonschema:"New description. Leave empty to keep current."`
}

// UpdateTool returns the MCP tool definition for update_service_account.
func UpdateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "update_service_account",
		Description: "Update a service account's display name and/or description. " +
			"Only provided fields are changed — omitted fields are left as-is.",
	}
}

// UpdateHandler returns the typed tool handler for update_service_account.
func UpdateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpdateServiceAccountInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UpdateServiceAccountInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := Update(ctx, serverAddress, input.ID, UpdateFields{
			DisplayName: input.DisplayName,
			Description: input.Description,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_service_account
// ---------------------------------------------------------------------------

// DeleteServiceAccountInput defines the parameters for delete_service_account.
type DeleteServiceAccountInput struct {
	ID string `json:"id" jsonschema:"required,Service account ID to delete."`
}

// DeleteTool returns the MCP tool definition for delete_service_account.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_service_account",
		Description: "Delete a service account and cascade: all API keys are revoked, " +
			"authorization tuples are removed, and the backing identity is deleted. " +
			"This cannot be undone.",
	}
}

// DeleteHandler returns the typed tool handler for delete_service_account.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteServiceAccountInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteServiceAccountInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := Delete(ctx, serverAddress, input.ID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_service_accounts
// ---------------------------------------------------------------------------

// ListServiceAccountsInput defines the parameters for list_service_accounts.
type ListServiceAccountsInput struct {
	Org string `json:"org" jsonschema:"required,Organization identifier."`
}

// ListTool returns the MCP tool definition for list_service_accounts.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_service_accounts",
		Description: "List all service accounts in an organization. " +
			"Returns each account's ID, display name, and description.",
	}
}

// ListHandler returns the typed tool handler for list_service_accounts.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListServiceAccountsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListServiceAccountsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := List(ctx, serverAddress, input.Org)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
