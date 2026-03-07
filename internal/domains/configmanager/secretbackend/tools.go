package secretbackend

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// apply_secret_backend
// ---------------------------------------------------------------------------

// ApplySecretBackendInput defines the parameters for the apply_secret_backend tool.
type ApplySecretBackendInput struct {
	BackendObject map[string]any `json:"backend_object" jsonschema:"required,Full SecretBackend object in OpenMCF envelope format."`
}

// ApplyTool returns the MCP tool definition for apply_secret_backend.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_secret_backend",
		Description: "Create or update a secret backend configuration (idempotent). " +
			"A secret backend defines where encrypted secret data is stored. " +
			"Six backend types are supported: platform (managed OpenBAO), openbao, " +
			"aws_secrets_manager, hashicorp_vault, gcp_secret_manager, and azure_key_vault. " +
			"Pass the full SecretBackend object as an OpenMCF envelope.",
	}
}

// ApplyHandler returns the typed tool handler for apply_secret_backend.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplySecretBackendInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplySecretBackendInput) (*mcp.CallToolResult, any, error) {
		if input.BackendObject == nil {
			return nil, nil, fmt.Errorf("'backend_object' is required")
		}
		text, err := Apply(ctx, serverAddress, input.BackendObject)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_secret_backend
// ---------------------------------------------------------------------------

// GetSecretBackendInput defines the parameters for the get_secret_backend tool.
type GetSecretBackendInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The secret backend ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Secret backend slug within the organization. Must be paired with 'org'."`
}

// GetTool returns the MCP tool definition for get_secret_backend.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_secret_backend",
		Description: "Retrieve a secret backend by ID or by org+slug. " +
			"Returns the backend configuration including type, encryption settings, " +
			"and whether it is the organization's default. " +
			"Sensitive credential fields (tokens, keys) are redacted.",
	}
}

// GetHandler returns the typed tool handler for get_secret_backend.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetSecretBackendInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetSecretBackendInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Slug); err != nil {
			return nil, nil, err
		}
		text, err := Get(ctx, serverAddress, input.ID, input.Org, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_secret_backends
// ---------------------------------------------------------------------------

// ListSecretBackendsInput defines the parameters for the list_secret_backends tool.
type ListSecretBackendsInput struct {
	Org string `json:"org" jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
}

// ListTool returns the MCP tool definition for list_secret_backends.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_secret_backends",
		Description: "List all secret backends configured for an organization. " +
			"Returns each backend's type, slug, and whether it is the default. " +
			"Sensitive credential fields are redacted.",
	}
}

// ListHandler returns the typed tool handler for list_secret_backends.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListSecretBackendsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListSecretBackendsInput) (*mcp.CallToolResult, any, error) {
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
// delete_secret_backend
// ---------------------------------------------------------------------------

// DeleteSecretBackendInput defines the parameters for the delete_secret_backend tool.
type DeleteSecretBackendInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The secret backend ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Secret backend slug within the organization. Must be paired with 'org'."`
}

// DeleteTool returns the MCP tool definition for delete_secret_backend.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_secret_backend",
		Description: "Delete a secret backend by ID or by org+slug. " +
			"WARNING: Deletion will fail if any secrets still reference this backend. " +
			"Migrate secrets to another backend first.",
	}
}

// DeleteHandler returns the typed tool handler for delete_secret_backend.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteSecretBackendInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteSecretBackendInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Slug); err != nil {
			return nil, nil, err
		}
		text, err := Delete(ctx, serverAddress, input.ID, input.Org, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// shared validation
// ---------------------------------------------------------------------------

func validateIdentification(id, org, slug string) error {
	hasID := id != ""
	hasOrg := org != ""
	hasSlug := slug != ""

	switch {
	case hasID && (hasOrg || hasSlug):
		return fmt.Errorf("provide either 'id' alone or both 'org' and 'slug' — not both paths")
	case hasID:
		return nil
	case hasOrg && hasSlug:
		return nil
	case hasOrg || hasSlug:
		if !hasOrg {
			return fmt.Errorf("'org' is required when using slug-based lookup")
		}
		return fmt.Errorf("'slug' is required when using org-based lookup")
	default:
		return fmt.Errorf("provide either 'id' or both 'org' and 'slug' to identify the secret backend")
	}
}
