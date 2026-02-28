// Package secretsgroup provides the MCP tools for the ServiceHub
// SecretsGroup domain, backed by the SecretsGroupQueryController,
// SecretsGroupCommandController (ai.planton.servicehub.secretsgroup.v1),
// and ServiceHubSearchQueryController (ai.planton.search.v1.servicehub) RPCs.
//
// Eight tools are exposed:
//   - search_secrets:        entry-level search across all groups in an org
//   - get_secrets_group:     retrieve a group by ID or org+slug
//   - apply_secrets_group:   create or update a group (idempotent)
//   - delete_secrets_group:  remove a group
//   - upsert_secret:         add or update a single secret entry
//   - delete_secret:         remove a single secret entry
//   - get_secret_value:      resolve a specific secret's value (PLAINTEXT)
//   - transform_secrets:     batch-resolve $secrets-group/ references
package secretsgroup

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// search_secrets
// ---------------------------------------------------------------------------

// SearchSecretsInput defines the parameters for the search_secrets tool.
type SearchSecretsInput struct {
	Org        string `json:"org"                   jsonschema:"required,Organization identifier to search within. Use list_organizations to discover available organizations."`
	SearchText string `json:"search_text,omitempty"  jsonschema:"Free-text search query to filter secret entries by name or description."`
	PageNum    int32  `json:"page_num,omitempty"     jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty"    jsonschema:"Number of results per page. Defaults to 20."`
}

// SearchTool returns the MCP tool definition for search_secrets.
func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_secrets",
		Description: "Search secret entries across all secrets groups within an organization. " +
			"Returns individual secret entries (not whole groups) with their group context — " +
			"useful for finding where a specific secret like DB_PASSWORD is defined. " +
			"Each result includes the group name, group ID, secret name, and description. " +
			"Use get_secrets_group with the group ID from results to see the full group.",
	}
}

// SearchHandler returns the typed tool handler for search_secrets.
func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchSecretsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchSecretsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := Search(ctx, serverAddress, SearchInput{
			Org:        input.Org,
			SearchText: input.SearchText,
			PageNum:    input.PageNum,
			PageSize:   input.PageSize,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_secrets_group
// ---------------------------------------------------------------------------

// GetSecretsGroupInput defines the parameters for the get_secrets_group
// tool. Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type GetSecretsGroupInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The secrets group ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Secrets group slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// GetTool returns the MCP tool definition for get_secrets_group.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_secrets_group",
		Description: "Retrieve the full details of a secrets group by its ID or by org+slug. " +
			"A secrets group is a named collection of secrets (e.g., API keys, database passwords) " +
			"that services reference for sensitive configuration. " +
			"Returns the complete group including metadata, all secret entries (names, descriptions), " +
			"and audit status. " +
			"The output JSON can be modified and passed to apply_secrets_group for updates.",
	}
}

// GetHandler returns the typed tool handler for get_secrets_group.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetSecretsGroupInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetSecretsGroupInput) (*mcp.CallToolResult, any, error) {
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
// apply_secrets_group
// ---------------------------------------------------------------------------

// ApplySecretsGroupInput defines the parameters for the
// apply_secrets_group tool.
type ApplySecretsGroupInput struct {
	SecretsGroup map[string]any `json:"secrets_group" jsonschema:"required,The full SecretsGroup resource as a JSON object. Must include 'api_version' ('service-hub.planton.ai/v1'), 'kind' ('SecretsGroup'), 'metadata' (with 'name' and 'org'), and 'spec' (with optional 'description' and 'entries'). The output of get_secrets_group can be modified and passed directly here."`
}

// ApplyTool returns the MCP tool definition for apply_secrets_group.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_secrets_group",
		Description: "Create or update a secrets group (idempotent). " +
			"Accepts the full SecretsGroup resource as a JSON object. " +
			"A secrets group holds a collection of named secrets that services " +
			"can reference for sensitive configuration. " +
			"For new groups, provide api_version, kind, metadata (name, org), and spec (description, entries). " +
			"For updates, retrieve the group with get_secrets_group, modify desired fields, and pass here. " +
			"WARNING: This replaces all entries — to modify a single secret without affecting others, " +
			"use upsert_secret instead. " +
			"Returns the applied group with server-assigned ID and audit information.",
	}
}

// ApplyHandler returns the typed tool handler for apply_secrets_group.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplySecretsGroupInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplySecretsGroupInput) (*mcp.CallToolResult, any, error) {
		if len(input.SecretsGroup) == 0 {
			return nil, nil, fmt.Errorf("'secrets_group' is required and must be a non-empty JSON object")
		}
		text, err := Apply(ctx, serverAddress, input.SecretsGroup)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_secrets_group
// ---------------------------------------------------------------------------

// DeleteSecretsGroupInput defines the parameters for the
// delete_secrets_group tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type DeleteSecretsGroupInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The secrets group ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Secrets group slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// DeleteTool returns the MCP tool definition for delete_secrets_group.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_secrets_group",
		Description: "Delete a secrets group from the platform. " +
			"WARNING: Ensure no services reference this group before deleting — " +
			"services using $secrets-group/ references to this group will fail during deployment. " +
			"Identify the group by ID or by org+slug.",
	}
}

// DeleteHandler returns the typed tool handler for delete_secrets_group.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteSecretsGroupInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteSecretsGroupInput) (*mcp.CallToolResult, any, error) {
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
// upsert_secret
// ---------------------------------------------------------------------------

// UpsertSecretInput defines the parameters for the upsert_secret tool.
// The target group can be identified by group_id alone, or by org+group_slug.
type UpsertSecretInput struct {
	GroupID   string         `json:"group_id,omitempty"    jsonschema:"The secrets group ID. Mutually exclusive with org+group_slug."`
	Org       string         `json:"org,omitempty"         jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'group_slug'."`
	GroupSlug string         `json:"group_slug,omitempty"  jsonschema:"Secrets group slug (name) for lookup within an organization. Must be paired with 'org'."`
	Entry     map[string]any `json:"entry"                 jsonschema:"required,The secret entry as a JSON object with 'name' (required), optional 'description', and either 'value' (literal string) or 'value_from' (reference object). If a secret with the same name exists it will be updated; otherwise it will be added."`
}

// UpsertTool returns the MCP tool definition for upsert_secret.
func UpsertTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "upsert_secret",
		Description: "Add or update a single secret in a secrets group. " +
			"If a secret with the same name already exists, it is updated; otherwise it is added. " +
			"This is safer than apply_secrets_group when modifying a single secret, " +
			"because it does not affect other entries in the group. " +
			"Identify the target group by group_id or by org+group_slug. " +
			"The entry must include 'name' and either 'value' (a literal string) or " +
			"'value_from' (a reference to another resource's field). " +
			"Returns the updated secrets group.",
	}
}

// UpsertHandler returns the typed tool handler for upsert_secret.
func UpsertHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpsertSecretInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UpsertSecretInput) (*mcp.CallToolResult, any, error) {
		if err := validateGroupIdentification(input.GroupID, input.Org, input.GroupSlug); err != nil {
			return nil, nil, err
		}
		if len(input.Entry) == 0 {
			return nil, nil, fmt.Errorf("'entry' is required and must be a non-empty JSON object")
		}
		text, err := UpsertEntry(ctx, serverAddress, input.GroupID, input.Org, input.GroupSlug, input.Entry)
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
// The target group can be identified by group_id alone, or by org+group_slug.
type DeleteSecretInput struct {
	GroupID   string `json:"group_id,omitempty"    jsonschema:"The secrets group ID. Mutually exclusive with org+group_slug."`
	Org       string `json:"org,omitempty"         jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'group_slug'."`
	GroupSlug string `json:"group_slug,omitempty"  jsonschema:"Secrets group slug (name) for lookup within an organization. Must be paired with 'org'."`
	EntryName string `json:"entry_name"            jsonschema:"required,Name of the secret to remove from the group."`
}

// DeleteEntryTool returns the MCP tool definition for delete_secret.
func DeleteEntryTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_secret",
		Description: "Remove a single secret from a secrets group. " +
			"Other secrets in the group remain unchanged. " +
			"Identify the target group by group_id or by org+group_slug. " +
			"WARNING: Services referencing this secret via $secrets-group/ will fail during deployment.",
	}
}

// DeleteEntryHandler returns the typed tool handler for delete_secret.
func DeleteEntryHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteSecretInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteSecretInput) (*mcp.CallToolResult, any, error) {
		if err := validateGroupIdentification(input.GroupID, input.Org, input.GroupSlug); err != nil {
			return nil, nil, err
		}
		if input.EntryName == "" {
			return nil, nil, fmt.Errorf("'entry_name' is required")
		}
		text, err := DeleteEntry(ctx, serverAddress, input.GroupID, input.Org, input.GroupSlug, input.EntryName)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_secret_value
// ---------------------------------------------------------------------------

// GetSecretValueInput defines the parameters for the get_secret_value tool.
type GetSecretValueInput struct {
	Org       string `json:"org"        jsonschema:"required,Organization identifier."`
	GroupName string `json:"group_name" jsonschema:"required,Name (slug) of the secrets group."`
	EntryName string `json:"entry_name" jsonschema:"required,Name of the secret whose value to retrieve."`
}

// GetValueTool returns the MCP tool definition for get_secret_value.
func GetValueTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_secret_value",
		Description: "Retrieve the resolved value of a specific secret from a secrets group. " +
			"If the secret uses a value_from reference, the reference is resolved to its current value. " +
			"WARNING: This returns the secret value in PLAINTEXT. " +
			"Only use when the user explicitly requests to see a secret value. " +
			"Never log or display secret values unless specifically asked.",
	}
}

// GetValueHandler returns the typed tool handler for get_secret_value.
func GetValueHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetSecretValueInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetSecretValueInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if input.GroupName == "" {
			return nil, nil, fmt.Errorf("'group_name' is required")
		}
		if input.EntryName == "" {
			return nil, nil, fmt.Errorf("'entry_name' is required")
		}
		text, err := GetValue(ctx, serverAddress, input.Org, input.GroupName, input.EntryName)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// transform_secrets
// ---------------------------------------------------------------------------

// TransformSecretsInput defines the parameters for the
// transform_secrets tool.
type TransformSecretsInput struct {
	Org     string            `json:"org"     jsonschema:"required,Organization identifier. Determines which secrets groups are available for reference resolution."`
	Entries map[string]string `json:"entries" jsonschema:"required,Map of environment variable names to values. Values starting with $secrets-group/ will be resolved to their actual values. Other values pass through unchanged."`
}

// TransformTool returns the MCP tool definition for transform_secrets.
func TransformTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "transform_secrets",
		Description: "Batch-resolve $secrets-group/ references in a set of environment variables. " +
			"Accepts a map of key-value pairs where values may contain $secrets-group/<group-name>/<entry-name> references. " +
			"References are resolved to their actual values; literal values pass through unchanged. " +
			"Returns two maps: successfully transformed entries and any entries that failed resolution " +
			"(with error messages explaining why). " +
			"WARNING: Resolved values are returned in PLAINTEXT. Use only when debugging or when the user " +
			"explicitly requests resolved secret values.",
	}
}

// TransformHandler returns the typed tool handler for transform_secrets.
func TransformHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *TransformSecretsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *TransformSecretsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		if len(input.Entries) == 0 {
			return nil, nil, fmt.Errorf("'entries' is required and must be a non-empty map")
		}
		text, err := Transform(ctx, serverAddress, input.Org, input.Entries)
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
// provided: either 'id' alone, or both 'org' and 'slug'.
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
	case hasOrg:
		return fmt.Errorf("'slug' is required when using 'org' for identification")
	case hasSlug:
		return fmt.Errorf("'org' is required when using 'slug' for identification")
	default:
		return fmt.Errorf("provide either 'id' or both 'org' and 'slug' to identify the secrets group")
	}
}

// validateGroupIdentification checks that exactly one group identification
// path is provided: either 'group_id' alone, or both 'org' and 'group_slug'.
func validateGroupIdentification(groupID, org, groupSlug string) error {
	hasID := groupID != ""
	hasOrg := org != ""
	hasSlug := groupSlug != ""

	switch {
	case hasID && (hasOrg || hasSlug):
		return fmt.Errorf("provide either 'group_id' alone or both 'org' and 'group_slug' — not both paths")
	case hasID:
		return nil
	case hasOrg && hasSlug:
		return nil
	case hasOrg:
		return fmt.Errorf("'group_slug' is required when using 'org' for group identification")
	case hasSlug:
		return fmt.Errorf("'org' is required when using 'group_slug' for group identification")
	default:
		return fmt.Errorf("provide either 'group_id' or both 'org' and 'group_slug' to identify the secrets group")
	}
}
