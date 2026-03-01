// Package organization provides the MCP tools for the Organization domain,
// backed by the OrganizationCommandController and OrganizationQueryController
// RPCs (ai.planton.resourcemanager.organization.v1) on the Planton backend.
//
// Five tools are exposed:
//   - list_organizations:  retrieve all organizations the caller is a member of
//   - get_organization:    retrieve a single organization by ID
//   - create_organization: provision a new organization
//   - update_organization: modify an existing organization (read-modify-write)
//   - delete_organization: remove an organization by ID
package organization

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// list_organizations
// ---------------------------------------------------------------------------

// ListOrganizationsInput defines the parameters for the list_organizations
// tool. This tool takes no input — the server returns organizations scoped to
// the authenticated caller's membership.
type ListOrganizationsInput struct{}

// ListTool returns the MCP tool definition for list_organizations.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_organizations",
		Description: "List all organizations the caller is a member of. " +
			"Returns the full organization objects including id, name, and slug. " +
			"Use this as the first step to discover the operating context before working with cloud resources or environments.",
	}
}

// ListHandler returns the typed tool handler for list_organizations.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListOrganizationsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ *ListOrganizationsInput) (*mcp.CallToolResult, any, error) {
		text, err := List(ctx, serverAddress)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_organization
// ---------------------------------------------------------------------------

// GetOrganizationInput defines the parameters for the get_organization tool.
type GetOrganizationInput struct {
	OrgID string `json:"org_id" jsonschema:"required,Organization ID. Use list_organizations to discover available IDs."`
}

// GetTool returns the MCP tool definition for get_organization.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_organization",
		Description: "Get an organization by ID from Planton Cloud. " +
			"Returns the full organization including metadata (id, slug, name), spec (description, contact email), and status. " +
			"Use list_organizations first to discover organization IDs.",
	}
}

// GetHandler returns the typed tool handler for get_organization.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetOrganizationInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetOrganizationInput) (*mcp.CallToolResult, any, error) {
		if input.OrgID == "" {
			return nil, nil, fmt.Errorf("'org_id' is required")
		}
		text, err := Get(ctx, serverAddress, input.OrgID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// create_organization
// ---------------------------------------------------------------------------

// CreateOrganizationInput defines the parameters for the create_organization
// tool. Only slug is required; the server assigns the organization ID.
type CreateOrganizationInput struct {
	Slug         string `json:"slug"                    jsonschema:"required,URL-friendly identifier (2-15 lowercase chars; letters, digits, hyphens; must start with a letter). Immutable after creation."`
	Name         string `json:"name,omitempty"           jsonschema:"Human-readable display name for the organization."`
	Description  string `json:"description,omitempty"    jsonschema:"Short description of the organization."`
	ContactEmail string `json:"contact_email,omitempty"  jsonschema:"Primary contact email for billing and administrative communications."`
}

// CreateTool returns the MCP tool definition for create_organization.
func CreateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "create_organization",
		Description: "Create a new organization on Planton Cloud. " +
			"Requires a unique slug (2-15 lowercase chars, letters/digits/hyphens, starts with a letter). " +
			"Any authenticated user can create an organization. " +
			"Returns the created organization with its server-assigned ID.",
	}
}

// CreateHandler returns the typed tool handler for create_organization.
func CreateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CreateOrganizationInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CreateOrganizationInput) (*mcp.CallToolResult, any, error) {
		if input.Slug == "" {
			return nil, nil, fmt.Errorf("'slug' is required")
		}
		text, err := Create(ctx, serverAddress, input.Slug, input.Name, input.Description, input.ContactEmail)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// update_organization
// ---------------------------------------------------------------------------

// UpdateOrganizationInput defines the parameters for the update_organization
// tool. Only org_id is required; all other fields are optional updates.
// Empty strings are treated as "no change".
type UpdateOrganizationInput struct {
	OrgID        string `json:"org_id"                   jsonschema:"required,Organization ID of the organization to update."`
	Name         string `json:"name,omitempty"            jsonschema:"New display name. Leave empty to keep the current value."`
	Description  string `json:"description,omitempty"     jsonschema:"New description. Leave empty to keep the current value."`
	ContactEmail string `json:"contact_email,omitempty"   jsonschema:"New contact email. Leave empty to keep the current value."`
	LogoURL      string `json:"logo_url,omitempty"        jsonschema:"New logo URL. Leave empty to keep the current value."`
}

// UpdateTool returns the MCP tool definition for update_organization.
func UpdateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "update_organization",
		Description: "Update an existing organization on Planton Cloud. " +
			"Provide the org_id and any fields to change — omitted fields are left unchanged. " +
			"Requires organization update permission. " +
			"Returns the updated organization.",
	}
}

// UpdateHandler returns the typed tool handler for update_organization.
func UpdateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpdateOrganizationInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UpdateOrganizationInput) (*mcp.CallToolResult, any, error) {
		if input.OrgID == "" {
			return nil, nil, fmt.Errorf("'org_id' is required")
		}
		fields := UpdateFields{
			Name:         input.Name,
			Description:  input.Description,
			ContactEmail: input.ContactEmail,
			LogoURL:      input.LogoURL,
		}
		text, err := Update(ctx, serverAddress, input.OrgID, fields)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_organization
// ---------------------------------------------------------------------------

// DeleteOrganizationInput defines the parameters for the delete_organization tool.
type DeleteOrganizationInput struct {
	OrgID string `json:"org_id" jsonschema:"required,Organization ID of the organization to delete."`
}

// DeleteTool returns the MCP tool definition for delete_organization.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_organization",
		Description: "Delete an organization from Planton Cloud. " +
			"Requires organization delete permission. " +
			"WARNING: This permanently removes the organization and is not reversible.",
	}
}

// DeleteHandler returns the typed tool handler for delete_organization.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteOrganizationInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteOrganizationInput) (*mcp.CallToolResult, any, error) {
		if input.OrgID == "" {
			return nil, nil, fmt.Errorf("'org_id' is required")
		}
		text, err := Delete(ctx, serverAddress, input.OrgID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
