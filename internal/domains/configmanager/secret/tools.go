// Package secret provides the MCP tools for the Secret domain, backed by
// the SecretQueryController and SecretCommandController RPCs
// (ai.planton.configmanager.secret.v1) on the Planton backend.
//
// Four tools are exposed:
//   - list_secrets:   paginated listing with org/env filters
//   - get_secret:     retrieve metadata by ID or by org+scope+slug
//   - apply_secret:   create or update metadata with explicit parameters
//   - delete_secret:  remove secret and all its versions
//
// Secret values are managed through the secretversion package. This package
// only manages secret metadata (name, scope, description, backend).
package secret

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	secretv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secret/v1"
)

// ---------------------------------------------------------------------------
// list_secrets
// ---------------------------------------------------------------------------

// ListSecretsInput defines the parameters for the list_secrets tool.
type ListSecretsInput struct {
	Org      string `json:"org"                jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
	Env      string `json:"env,omitempty"       jsonschema:"Environment slug to filter by. When omitted, secrets across all environments and organization-scoped secrets are returned."`
	PageNum  int32  `json:"page_num,omitempty"  jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize int32  `json:"page_size,omitempty" jsonschema:"Number of results per page. Defaults to 20."`
}

// ListTool returns the MCP tool definition for list_secrets.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_secrets",
		Description: "List secrets within an organization. " +
			"Secrets are metadata containers for encrypted key-value pairs scoped to an organization or environment. " +
			"Only metadata is returned — no secret values are exposed. " +
			"Optionally filter by environment slug. " +
			"Use get_secret with an ID from the results to retrieve full metadata, " +
			"or list_secret_versions to see version history for a specific secret.",
	}
}

// ListHandler returns the typed tool handler for list_secrets.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListSecretsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListSecretsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := List(ctx, serverAddress, ListInput{
			Org:      input.Org,
			Env:      input.Env,
			PageNum:  input.PageNum,
			PageSize: input.PageSize,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_secret
// ---------------------------------------------------------------------------

// GetSecretInput defines the parameters for the get_secret tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set 'org', 'scope', and 'slug'.
type GetSecretInput struct {
	ID    string `json:"id,omitempty"    jsonschema:"The secret ID. Mutually exclusive with org+scope+slug."`
	Org   string `json:"org,omitempty"   jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'scope' and 'slug'."`
	Scope string `json:"scope,omitempty" jsonschema:"Secret scope for slug-based lookup. Must be 'organization' or 'environment'. Must be paired with 'org' and 'slug'."`
	Slug  string `json:"slug,omitempty"  jsonschema:"Secret slug for lookup within an organization and scope. Must be paired with 'org' and 'scope'."`
}

// GetTool returns the MCP tool definition for get_secret.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_secret",
		Description: "Retrieve the full metadata of a secret by its ID or by org+scope+slug. " +
			"Returns the secret including metadata, spec (scope, description, backend), and audit status. " +
			"No secret values are exposed — use list_secret_versions to see version history " +
			"and create_secret_version to add new values. " +
			"Secrets are uniquely identified within (org, scope, slug).",
	}
}

// GetHandler returns the typed tool handler for get_secret.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetSecretInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetSecretInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Scope, input.Slug); err != nil {
			return nil, nil, err
		}
		var scope secretv1.SecretSpec_Scope
		if input.Scope != "" {
			var err error
			scope, err = resolveScope(input.Scope)
			if err != nil {
				return nil, nil, err
			}
		}
		text, err := Get(ctx, serverAddress, input.ID, input.Org, scope, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// apply_secret
// ---------------------------------------------------------------------------

// ApplySecretInput defines the parameters for the apply_secret tool.
type ApplySecretInput struct {
	Name        string `json:"name"               jsonschema:"required,Display name for the secret. Also used to derive the slug if not already set."`
	Org         string `json:"org"                jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
	Scope       string `json:"scope"              jsonschema:"required,Secret scope. Must be 'organization' (shared across all environments) or 'environment' (scoped to a specific environment)."`
	Env         string `json:"env,omitempty"       jsonschema:"Environment slug. Required when scope is 'environment'. Ignored when scope is 'organization'."`
	Description string `json:"description,omitempty" jsonschema:"Human-readable description of what this secret stores."`
	Backend     string `json:"backend,omitempty"   jsonschema:"Slug of the SecretBackend resource to use for encryption. When omitted, the organization's default backend is used."`
}

// ApplyTool returns the MCP tool definition for apply_secret.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_secret",
		Description: "Create or update a secret's metadata (idempotent). " +
			"Secrets are metadata containers for encrypted key-value pairs. " +
			"This tool manages the secret record only — use create_secret_version to store actual secret values. " +
			"The scope determines the uniqueness key: secrets are unique within (org, scope, slug). " +
			"When scope is 'environment', the env parameter is required. " +
			"The backend cannot be changed after creation.",
	}
}

// ApplyHandler returns the typed tool handler for apply_secret.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplySecretInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplySecretInput) (*mcp.CallToolResult, any, error) {
		if input.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Scope == "" {
			return nil, nil, fmt.Errorf("'scope' is required")
		}
		scope, err := resolveScope(input.Scope)
		if err != nil {
			return nil, nil, err
		}
		if err := validateScopeEnv(input.Scope, input.Env); err != nil {
			return nil, nil, err
		}
		text, err := Apply(ctx, serverAddress, ApplyInput{
			Name:        input.Name,
			Org:         input.Org,
			Scope:       scope,
			Env:         input.Env,
			Description: input.Description,
			Backend:     input.Backend,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_secret
// ---------------------------------------------------------------------------

// DeleteSecretInput defines the parameters for the delete_secret tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set 'org', 'scope', and 'slug'.
type DeleteSecretInput struct {
	ID    string `json:"id,omitempty"    jsonschema:"The secret ID. Mutually exclusive with org+scope+slug."`
	Org   string `json:"org,omitempty"   jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'scope' and 'slug'."`
	Scope string `json:"scope,omitempty" jsonschema:"Secret scope for slug-based lookup. Must be 'organization' or 'environment'. Must be paired with 'org' and 'slug'."`
	Slug  string `json:"slug,omitempty"  jsonschema:"Secret slug for lookup within an organization and scope. Must be paired with 'org' and 'scope'."`
}

// DeleteTool returns the MCP tool definition for delete_secret.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_secret",
		Description: "Delete a secret and ALL its versions from the platform. " +
			"WARNING: This is a destructive operation that permanently removes the secret record " +
			"AND destroys all encrypted version data stored in the backend. This cannot be undone. " +
			"Use list_secret_versions first to understand what will be destroyed. " +
			"Identify the secret by ID or by org+scope+slug.",
	}
}

// DeleteHandler returns the typed tool handler for delete_secret.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteSecretInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteSecretInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Scope, input.Slug); err != nil {
			return nil, nil, err
		}
		var scope secretv1.SecretSpec_Scope
		if input.Scope != "" {
			var err error
			scope, err = resolveScope(input.Scope)
			if err != nil {
				return nil, nil, err
			}
		}
		text, err := Delete(ctx, serverAddress, input.ID, input.Org, scope, input.Slug)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// shared validation
// ---------------------------------------------------------------------------

// validateIdentification checks that exactly one identification path is
// provided: either 'id' alone, or all of 'org', 'scope', and 'slug'.
func validateIdentification(id, org, scope, slug string) error {
	hasID := id != ""
	slugFields := [3]string{org, scope, slug}
	slugCount := 0
	for _, f := range slugFields {
		if f != "" {
			slugCount++
		}
	}

	switch {
	case hasID && slugCount > 0:
		return fmt.Errorf("provide either 'id' alone or all of 'org', 'scope', and 'slug' — not both paths")
	case hasID:
		return nil
	case slugCount == 3:
		return nil
	case slugCount > 0:
		missing := make([]string, 0, 3)
		if org == "" {
			missing = append(missing, "'org'")
		}
		if scope == "" {
			missing = append(missing, "'scope'")
		}
		if slug == "" {
			missing = append(missing, "'slug'")
		}
		return fmt.Errorf("when not using 'id', all of 'org', 'scope', and 'slug' are required — missing: %s",
			joinMissing(missing))
	default:
		return fmt.Errorf("provide either 'id' or all of 'org', 'scope', and 'slug' to identify the secret")
	}
}

// validateScopeEnv checks that 'env' is provided when scope is "environment".
func validateScopeEnv(scope, env string) error {
	if scope == "environment" && env == "" {
		return fmt.Errorf("'env' is required when scope is 'environment'")
	}
	return nil
}

// joinMissing joins a slice of missing field names for error messages.
func joinMissing(fields []string) string {
	switch len(fields) {
	case 0:
		return ""
	case 1:
		return fields[0]
	default:
		return fields[0] + " and " + joinMissing(fields[1:])
	}
}
