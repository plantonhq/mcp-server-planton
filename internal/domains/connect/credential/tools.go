package credential

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// apply_credential
// ---------------------------------------------------------------------------

// ApplyCredentialInput defines the parameters for the apply_credential tool.
type ApplyCredentialInput struct {
	CredentialObject map[string]any `json:"credential_object" jsonschema:"required,Full credential object in OpenMCF envelope format: { api_version, kind, metadata: { name, org }, spec: { ... } }. Read credential-types://catalog to discover available kinds, then credential-schema://{kind} to learn the spec fields."`
}

// ApplyTool returns the MCP tool definition for apply_credential.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_credential",
		Description: "Create or update a credential of any supported type (AWS, GCP, Azure, GitHub, Kubernetes, etc.). " +
			"Pass the full credential object as an OpenMCF envelope with api_version, kind, metadata, and spec. " +
			"Workflow: read credential-types://catalog to discover types, then credential-schema://{kind} for the spec schema, then call this tool. " +
			"The kind field in the credential_object determines which provider-specific backend is called.",
	}
}

// ApplyHandler returns the typed tool handler for apply_credential.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyCredentialInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyCredentialInput) (*mcp.CallToolResult, any, error) {
		if input.CredentialObject == nil {
			return nil, nil, fmt.Errorf("'credential_object' is required")
		}
		text, err := Apply(ctx, serverAddress, input.CredentialObject)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_credential
// ---------------------------------------------------------------------------

// GetCredentialInput defines the parameters for the get_credential tool.
// Either id or (org + slug) must be provided. kind is always required.
type GetCredentialInput struct {
	Kind string `json:"kind" jsonschema:"required,PascalCase credential type (e.g. AwsCredential, GcpCredential). Read credential-types://catalog to discover available kinds."`
	ID   string `json:"id,omitempty" jsonschema:"Credential ID. Provide either id or (org + slug) to identify the credential."`
	Org  string `json:"org,omitempty" jsonschema:"Organization ID. Used with slug for lookup by org+slug."`
	Slug string `json:"slug,omitempty" jsonschema:"Credential slug within the organization. Used with org for lookup by org+slug."`
}

// GetTool returns the MCP tool definition for get_credential.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_credential",
		Description: "Retrieve a credential by ID or by organization+slug. " +
			"Requires kind to dispatch to the correct provider backend. " +
			"Provide either 'id' alone, or both 'org' and 'slug'. " +
			"Sensitive fields (secret keys, tokens) are automatically redacted in the response.",
	}
}

// GetHandler returns the typed tool handler for get_credential.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetCredentialInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetCredentialInput) (*mcp.CallToolResult, any, error) {
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}

		var text string
		var err error

		switch {
		case input.ID != "":
			text, err = Get(ctx, serverAddress, input.Kind, input.ID)
		case input.Org != "" && input.Slug != "":
			text, err = GetByOrgBySlug(ctx, serverAddress, input.Kind, input.Org, input.Slug)
		default:
			return nil, nil, fmt.Errorf("provide either 'id' or both 'org' and 'slug'")
		}

		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_credential
// ---------------------------------------------------------------------------

// DeleteCredentialInput defines the parameters for the delete_credential tool.
type DeleteCredentialInput struct {
	Kind string `json:"kind" jsonschema:"required,PascalCase credential type (e.g. AwsCredential, GcpCredential)."`
	ID   string `json:"id" jsonschema:"required,Credential ID to delete."`
}

// DeleteTool returns the MCP tool definition for delete_credential.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_credential",
		Description: "Delete a credential by ID. " +
			"Requires kind to dispatch to the correct provider backend. " +
			"WARNING: This permanently removes the credential. Cloud resources using this credential will lose connectivity.",
	}
}

// DeleteHandler returns the typed tool handler for delete_credential.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteCredentialInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteCredentialInput) (*mcp.CallToolResult, any, error) {
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := Delete(ctx, serverAddress, input.Kind, input.ID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// search_credentials
// ---------------------------------------------------------------------------

// SearchCredentialsInput defines the parameters for the search_credentials tool.
type SearchCredentialsInput struct {
	Org        string   `json:"org" jsonschema:"required,Organization ID to search within."`
	Env        string   `json:"env,omitempty" jsonschema:"Optional environment slug to narrow results."`
	Kinds      []string `json:"kinds,omitempty" jsonschema:"Optional list of PascalCase credential kinds to filter (e.g. ['AwsCredential', 'GcpCredential'])."`
	SearchText string   `json:"search_text,omitempty" jsonschema:"Optional free-text filter on credential name."`
}

// SearchTool returns the MCP tool definition for search_credentials.
func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_credentials",
		Description: "Search credentials within an organization. " +
			"Returns lightweight search records (id, name, kind) without sensitive spec fields. " +
			"Optionally filter by environment, credential kinds, or free-text search. " +
			"Use get_credential to retrieve the full (redacted) credential details.",
	}
}

// SearchHandler returns the typed tool handler for search_credentials.
func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchCredentialsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchCredentialsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := Search(ctx, serverAddress, input.Org, input.Env, input.Kinds, input.SearchText)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// check_credential_slug
// ---------------------------------------------------------------------------

// CheckCredentialSlugInput defines the parameters for the check_credential_slug tool.
type CheckCredentialSlugInput struct {
	Org  string `json:"org" jsonschema:"required,Organization ID to check within."`
	Kind string `json:"kind" jsonschema:"required,PascalCase credential type (e.g. AwsCredential)."`
	Slug string `json:"slug" jsonschema:"required,Slug to check for availability."`
}

// CheckSlugTool returns the MCP tool definition for check_credential_slug.
func CheckSlugTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "check_credential_slug",
		Description: "Check if a credential slug is available within an organization for a given credential kind. " +
			"Returns { is_available: true/false }. " +
			"Use this before apply_credential to verify the slug is unique.",
	}
}

// CheckSlugHandler returns the typed tool handler for check_credential_slug.
func CheckSlugHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CheckCredentialSlugInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CheckCredentialSlugInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}
		if input.Slug == "" {
			return nil, nil, fmt.Errorf("'slug' is required")
		}
		text, err := CheckSlugAvailability(ctx, serverAddress, input.Org, input.Kind, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
