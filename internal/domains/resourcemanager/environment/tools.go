// Package environment provides the MCP tools for the Environment domain,
// backed by the EnvironmentCommandController and EnvironmentQueryController
// RPCs (ai.planton.resourcemanager.environment.v1) on the Planton backend.
//
// Five tools are exposed:
//   - list_environments:   retrieve environments the caller can access within an organization
//   - get_environment:     retrieve a single environment by ID or by org+slug
//   - create_environment:  provision a new environment within an organization
//   - update_environment:  modify an existing environment (read-modify-write)
//   - delete_environment:  remove an environment by ID
package environment

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// list_environments
// ---------------------------------------------------------------------------

// ListEnvironmentsInput defines the parameters for the list_environments tool.
type ListEnvironmentsInput struct {
	Org string `json:"org" jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
}

// ListTool returns the MCP tool definition for list_environments.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_environments",
		Description: "List environments the caller can access within an organization. " +
			"Returns only environments where the caller has at least view permission. " +
			"Use list_organizations first to discover available organization identifiers.",
	}
}

// ListHandler returns the typed tool handler for list_environments.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListEnvironmentsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListEnvironmentsInput) (*mcp.CallToolResult, any, error) {
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

// ---------------------------------------------------------------------------
// get_environment
// ---------------------------------------------------------------------------

// GetEnvironmentInput defines the parameters for the get_environment tool.
// Supports dual-resolution: provide either env_id alone, or both org and slug.
type GetEnvironmentInput struct {
	EnvID string `json:"env_id,omitempty" jsonschema:"Environment ID. Provide this OR both org and slug."`
	Org   string `json:"org,omitempty"    jsonschema:"Organization identifier. Required when using slug-based lookup."`
	Slug  string `json:"slug,omitempty"   jsonschema:"Environment slug. Required when using slug-based lookup."`
}

// GetTool returns the MCP tool definition for get_environment.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_environment",
		Description: "Get an environment from Planton Cloud. " +
			"Supports two lookup modes: by env_id, or by org + slug. " +
			"Returns the full environment including metadata (id, slug, name, org), spec (description), and status. " +
			"Use list_environments to discover environment IDs and slugs within an organization.",
	}
}

// GetHandler returns the typed tool handler for get_environment.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetEnvironmentInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetEnvironmentInput) (*mcp.CallToolResult, any, error) {
		var text string
		var err error

		switch {
		case input.EnvID != "":
			text, err = Get(ctx, serverAddress, input.EnvID)
		case input.Org != "" && input.Slug != "":
			text, err = GetByOrgBySlug(ctx, serverAddress, input.Org, input.Slug)
		default:
			return nil, nil, fmt.Errorf("provide either 'env_id' or both 'org' and 'slug'")
		}

		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// create_environment
// ---------------------------------------------------------------------------

// CreateEnvironmentInput defines the parameters for the create_environment
// tool. Both org and slug are required; the server assigns the environment ID.
type CreateEnvironmentInput struct {
	Org         string `json:"org"                   jsonschema:"required,Organization identifier that the environment belongs to. Use list_organizations to discover available organizations."`
	Slug        string `json:"slug"                  jsonschema:"required,URL-friendly identifier (2-15 lowercase chars; letters, digits, hyphens; must start with a letter). Immutable after creation."`
	Name        string `json:"name,omitempty"         jsonschema:"Human-readable display name for the environment."`
	Description string `json:"description,omitempty"  jsonschema:"Short description of the environment."`
}

// CreateTool returns the MCP tool definition for create_environment.
func CreateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_environment",
		Description: "Create a new environment within an organization on Planton Cloud. " +
			"Requires a unique slug (2-15 lowercase chars, letters/digits/hyphens, starts with a letter) " +
			"and the organization identifier. " +
			"Returns the created environment with its server-assigned ID.",
	}
}

// CreateHandler returns the typed tool handler for create_environment.
func CreateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CreateEnvironmentInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CreateEnvironmentInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Slug == "" {
			return nil, nil, fmt.Errorf("'slug' is required")
		}
		text, err := Create(ctx, serverAddress, input.Org, input.Slug, input.Name, input.Description)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// update_environment
// ---------------------------------------------------------------------------

// UpdateEnvironmentInput defines the parameters for the update_environment
// tool. Only env_id is required; all other fields are optional updates.
// Empty strings are treated as "no change".
type UpdateEnvironmentInput struct {
	EnvID       string `json:"env_id"                 jsonschema:"required,Environment ID of the environment to update."`
	Name        string `json:"name,omitempty"          jsonschema:"New display name. Leave empty to keep the current value."`
	Description string `json:"description,omitempty"   jsonschema:"New description. Leave empty to keep the current value."`
}

// UpdateTool returns the MCP tool definition for update_environment.
func UpdateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "update_environment",
		Description: "Update an existing environment on Planton Cloud. " +
			"Provide the env_id and any fields to change — omitted fields are left unchanged. " +
			"Requires environment update permission. " +
			"Returns the updated environment.",
	}
}

// UpdateHandler returns the typed tool handler for update_environment.
func UpdateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpdateEnvironmentInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UpdateEnvironmentInput) (*mcp.CallToolResult, any, error) {
		if input.EnvID == "" {
			return nil, nil, fmt.Errorf("'env_id' is required")
		}
		fields := UpdateFields{
			Name:        input.Name,
			Description: input.Description,
		}
		text, err := Update(ctx, serverAddress, input.EnvID, fields)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_environment
// ---------------------------------------------------------------------------

// DeleteEnvironmentInput defines the parameters for the delete_environment tool.
type DeleteEnvironmentInput struct {
	EnvID string `json:"env_id" jsonschema:"required,Environment ID of the environment to delete."`
}

// DeleteTool returns the MCP tool definition for delete_environment.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_environment",
		Description: "Delete an environment from Planton Cloud. " +
			"Requires environment delete permission. " +
			"WARNING: This permanently removes the environment and triggers cascading cleanup " +
			"of all resources deployed to it, including stack-modules, microservices, secrets, and clusters. " +
			"This operation is not reversible.",
	}
}

// DeleteHandler returns the typed tool handler for delete_environment.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteEnvironmentInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteEnvironmentInput) (*mcp.CallToolResult, any, error) {
		if input.EnvID == "" {
			return nil, nil, fmt.Errorf("'env_id' is required")
		}
		text, err := Delete(ctx, serverAddress, input.EnvID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
