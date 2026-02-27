// Package variable provides the MCP tools for the Variable domain, backed by
// the VariableQueryController and VariableCommandController RPCs
// (ai.planton.configmanager.variable.v1) on the Planton backend.
//
// Five tools are exposed:
//   - list_variables:    paginated listing with org/env filters
//   - get_variable:      retrieve by ID or by org+scope+slug
//   - apply_variable:    create or update with explicit parameters
//   - delete_variable:   remove by ID or by org+scope+slug
//   - resolve_variable:  quick value lookup returning plain string
package variable

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	variablev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/variable/v1"
)

// ---------------------------------------------------------------------------
// list_variables
// ---------------------------------------------------------------------------

// ListVariablesInput defines the parameters for the list_variables tool.
type ListVariablesInput struct {
	Org      string `json:"org"                jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
	Env      string `json:"env,omitempty"       jsonschema:"Environment slug to filter by. When omitted, variables across all environments and organization-scoped variables are returned."`
	PageNum  int32  `json:"page_num,omitempty"  jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize int32  `json:"page_size,omitempty" jsonschema:"Number of results per page. Defaults to 20."`
}

// ListTool returns the MCP tool definition for list_variables.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_variables",
		Description: "List configuration variables within an organization. " +
			"Variables hold plaintext configuration values scoped to either an organization or a specific environment. " +
			"Optionally filter by environment slug. " +
			"Use get_variable with an ID from the results to retrieve full details, " +
			"or resolve_variable for a quick value lookup by slug.",
	}
}

// ListHandler returns the typed tool handler for list_variables.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListVariablesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListVariablesInput) (*mcp.CallToolResult, any, error) {
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
// get_variable
// ---------------------------------------------------------------------------

// GetVariableInput defines the parameters for the get_variable tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set 'org', 'scope', and 'slug'.
type GetVariableInput struct {
	ID    string `json:"id,omitempty"    jsonschema:"The variable ID. Mutually exclusive with org+scope+slug."`
	Org   string `json:"org,omitempty"   jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'scope' and 'slug'."`
	Scope string `json:"scope,omitempty" jsonschema:"Variable scope for slug-based lookup. Must be 'organization' or 'environment'. Must be paired with 'org' and 'slug'."`
	Slug  string `json:"slug,omitempty"  jsonschema:"Variable slug for lookup within an organization and scope. Must be paired with 'org' and 'scope'."`
}

// GetTool returns the MCP tool definition for get_variable.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_variable",
		Description: "Retrieve the full details of a configuration variable by its ID or by org+scope+slug. " +
			"Returns the complete variable including metadata, spec (scope, value, description), and audit status. " +
			"Variables are uniquely identified within (org, scope, slug). " +
			"Use resolve_variable instead if you only need the plain string value.",
	}
}

// GetHandler returns the typed tool handler for get_variable.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetVariableInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetVariableInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Scope, input.Slug); err != nil {
			return nil, nil, err
		}
		var scope variablev1.VariableSpec_Scope
		if input.Scope != "" {
			var err error
			scope, err = scopeResolver.Resolve(input.Scope)
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
// apply_variable
// ---------------------------------------------------------------------------

// ApplyVariableInput defines the parameters for the apply_variable tool.
type ApplyVariableInput struct {
	Name        string `json:"name"               jsonschema:"required,Display name for the variable. Also used to derive the slug if not already set."`
	Org         string `json:"org"                jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
	Scope       string `json:"scope"              jsonschema:"required,Variable scope. Must be 'organization' (shared across all environments) or 'environment' (scoped to a specific environment)."`
	Env         string `json:"env,omitempty"       jsonschema:"Environment slug. Required when scope is 'environment'. Ignored when scope is 'organization'."`
	Description string `json:"description,omitempty" jsonschema:"Human-readable description of what this variable is used for."`
	Value       string `json:"value"              jsonschema:"required,The variable value (plaintext string)."`
}

// ApplyTool returns the MCP tool definition for apply_variable.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_variable",
		Description: "Create or update a configuration variable (idempotent). " +
			"Variables hold plaintext configuration values scoped to an organization or environment. " +
			"The scope determines the uniqueness key: variables are unique within (org, scope, slug). " +
			"When scope is 'environment', the env parameter is required. " +
			"Returns the applied variable with server-assigned ID and audit information.",
	}
}

// ApplyHandler returns the typed tool handler for apply_variable.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyVariableInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyVariableInput) (*mcp.CallToolResult, any, error) {
		if input.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Scope == "" {
			return nil, nil, fmt.Errorf("'scope' is required")
		}
		scope, err := scopeResolver.Resolve(input.Scope)
		if err != nil {
			return nil, nil, err
		}
		if err := validateScopeEnv(input.Scope, input.Env); err != nil {
			return nil, nil, err
		}
		if input.Value == "" {
			return nil, nil, fmt.Errorf("'value' is required")
		}
		text, err := Apply(ctx, serverAddress, ApplyInput{
			Name:        input.Name,
			Org:         input.Org,
			Scope:       scope,
			Env:         input.Env,
			Description: input.Description,
			Value:       input.Value,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_variable
// ---------------------------------------------------------------------------

// DeleteVariableInput defines the parameters for the delete_variable tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set 'org', 'scope', and 'slug'.
type DeleteVariableInput struct {
	ID    string `json:"id,omitempty"    jsonschema:"The variable ID. Mutually exclusive with org+scope+slug."`
	Org   string `json:"org,omitempty"   jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'scope' and 'slug'."`
	Scope string `json:"scope,omitempty" jsonschema:"Variable scope for slug-based lookup. Must be 'organization' or 'environment'. Must be paired with 'org' and 'slug'."`
	Slug  string `json:"slug,omitempty"  jsonschema:"Variable slug for lookup within an organization and scope. Must be paired with 'org' and 'scope'."`
}

// DeleteTool returns the MCP tool definition for delete_variable.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_variable",
		Description: "Delete a configuration variable from the platform. " +
			"This permanently removes the variable record. " +
			"Identify the variable by ID or by org+scope+slug.",
	}
}

// DeleteHandler returns the typed tool handler for delete_variable.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteVariableInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteVariableInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Scope, input.Slug); err != nil {
			return nil, nil, err
		}
		var scope variablev1.VariableSpec_Scope
		if input.Scope != "" {
			var err error
			scope, err = scopeResolver.Resolve(input.Scope)
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
// resolve_variable
// ---------------------------------------------------------------------------

// ResolveVariableInput defines the parameters for the resolve_variable tool.
type ResolveVariableInput struct {
	Org   string `json:"org"   jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
	Scope string `json:"scope" jsonschema:"required,Variable scope. Must be 'organization' or 'environment'."`
	Slug  string `json:"slug"  jsonschema:"required,Variable slug within the organization and scope."`
}

// ResolveTool returns the MCP tool definition for resolve_variable.
func ResolveTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "resolve_variable",
		Description: "Resolve a configuration variable's value by org+scope+slug. " +
			"Returns only the plain string value — no metadata, no spec wrapper. " +
			"This is the fastest way to look up a variable's current value. " +
			"Use get_variable instead if you need the full resource with metadata and audit information.",
	}
}

// ResolveHandler returns the typed tool handler for resolve_variable.
func ResolveHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ResolveVariableInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ResolveVariableInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.Scope == "" {
			return nil, nil, fmt.Errorf("'scope' is required")
		}
		scope, err := scopeResolver.Resolve(input.Scope)
		if err != nil {
			return nil, nil, err
		}
		if input.Slug == "" {
			return nil, nil, fmt.Errorf("'slug' is required")
		}
		text, err := Resolve(ctx, serverAddress, input.Org, scope, input.Slug)
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
		return fmt.Errorf("provide either 'id' or all of 'org', 'scope', and 'slug' to identify the variable")
	}
}

// validateScopeEnv checks that 'env' is provided when scope is "environment"
// and warns if provided when scope is "organization".
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
