// Package variablesgroup provides the MCP tools for the ServiceHub
// VariablesGroup domain, backed by the VariablesGroupQueryController,
// VariablesGroupCommandController (ai.planton.servicehub.variablesgroup.v1),
// and ServiceHubSearchQueryController (ai.planton.search.v1.servicehub) RPCs.
//
// Eight tools are exposed:
//   - search_variables:        entry-level search across all groups in an org
//   - get_variables_group:     retrieve a group by ID or org+slug
//   - apply_variables_group:   create or update a group (idempotent)
//   - delete_variables_group:  remove a group
//   - upsert_variable:         add or update a single variable entry
//   - delete_variable:         remove a single variable entry
//   - get_variable_value:      resolve a specific variable's value
//   - transform_variables:     batch-resolve $variables-group/ references
package variablesgroup

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// search_variables
// ---------------------------------------------------------------------------

// SearchVariablesInput defines the parameters for the search_variables tool.
type SearchVariablesInput struct {
	Org        string `json:"org"                   jsonschema:"required,Organization identifier to search within. Use list_organizations to discover available organizations."`
	SearchText string `json:"search_text,omitempty"  jsonschema:"Free-text search query to filter variable entries by name, value, or description."`
	PageNum    int32  `json:"page_num,omitempty"     jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty"    jsonschema:"Number of results per page. Defaults to 20."`
}

// SearchTool returns the MCP tool definition for search_variables.
func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_variables",
		Description: "Search variable entries across all variables groups within an organization. " +
			"Returns individual variable entries (not whole groups) with their group context — " +
			"useful for finding where a specific variable like DATABASE_HOST is defined. " +
			"Each result includes the group name, group ID, variable name, value, and description. " +
			"Use get_variables_group with the group ID from results to see the full group.",
	}
}

// SearchHandler returns the typed tool handler for search_variables.
func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchVariablesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchVariablesInput) (*mcp.CallToolResult, any, error) {
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
// get_variables_group
// ---------------------------------------------------------------------------

// GetVariablesGroupInput defines the parameters for the get_variables_group
// tool. Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type GetVariablesGroupInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The variables group ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Variables group slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// GetTool returns the MCP tool definition for get_variables_group.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_variables_group",
		Description: "Retrieve the full details of a variables group by its ID or by org+slug. " +
			"A variables group is a named collection of environment variables that services reference " +
			"for shared configuration (e.g., database hosts, API URLs). " +
			"Returns the complete group including metadata, all variable entries (names, values, descriptions), " +
			"and audit status. " +
			"The output JSON can be modified and passed to apply_variables_group for updates.",
	}
}

// GetHandler returns the typed tool handler for get_variables_group.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetVariablesGroupInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetVariablesGroupInput) (*mcp.CallToolResult, any, error) {
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
// apply_variables_group
// ---------------------------------------------------------------------------

// ApplyVariablesGroupInput defines the parameters for the
// apply_variables_group tool.
type ApplyVariablesGroupInput struct {
	VariablesGroup map[string]any `json:"variables_group" jsonschema:"required,The full VariablesGroup resource as a JSON object. Must include 'api_version' ('service-hub.planton.ai/v1'), 'kind' ('VariablesGroup'), 'metadata' (with 'name' and 'org'), and 'spec' (with optional 'description' and 'entries'). The output of get_variables_group can be modified and passed directly here."`
}

// ApplyTool returns the MCP tool definition for apply_variables_group.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_variables_group",
		Description: "Create or update a variables group (idempotent). " +
			"Accepts the full VariablesGroup resource as a JSON object. " +
			"A variables group holds a collection of named environment variables that services " +
			"can reference for shared configuration. " +
			"For new groups, provide api_version, kind, metadata (name, org), and spec (description, entries). " +
			"For updates, retrieve the group with get_variables_group, modify desired fields, and pass here. " +
			"WARNING: This replaces all entries — to modify a single variable without affecting others, " +
			"use upsert_variable instead. " +
			"Returns the applied group with server-assigned ID and audit information.",
	}
}

// ApplyHandler returns the typed tool handler for apply_variables_group.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyVariablesGroupInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyVariablesGroupInput) (*mcp.CallToolResult, any, error) {
		if len(input.VariablesGroup) == 0 {
			return nil, nil, fmt.Errorf("'variables_group' is required and must be a non-empty JSON object")
		}
		text, err := Apply(ctx, serverAddress, input.VariablesGroup)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_variables_group
// ---------------------------------------------------------------------------

// DeleteVariablesGroupInput defines the parameters for the
// delete_variables_group tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type DeleteVariablesGroupInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The variables group ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Variables group slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// DeleteTool returns the MCP tool definition for delete_variables_group.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_variables_group",
		Description: "Delete a variables group from the platform. " +
			"WARNING: Ensure no services reference this group before deleting — " +
			"services using $variables-group/ references to this group will fail during deployment. " +
			"Identify the group by ID or by org+slug.",
	}
}

// DeleteHandler returns the typed tool handler for delete_variables_group.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteVariablesGroupInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteVariablesGroupInput) (*mcp.CallToolResult, any, error) {
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
// upsert_variable
// ---------------------------------------------------------------------------

// UpsertVariableInput defines the parameters for the upsert_variable tool.
// The target group can be identified by group_id alone, or by org+group_slug.
type UpsertVariableInput struct {
	GroupID   string         `json:"group_id,omitempty"    jsonschema:"The variables group ID. Mutually exclusive with org+group_slug."`
	Org       string         `json:"org,omitempty"         jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'group_slug'."`
	GroupSlug string         `json:"group_slug,omitempty"  jsonschema:"Variables group slug (name) for lookup within an organization. Must be paired with 'org'."`
	Entry     map[string]any `json:"entry"                 jsonschema:"required,The variable entry as a JSON object with 'name' (required), optional 'description', and either 'value' (literal string) or 'value_from' (reference object). If a variable with the same name exists it will be updated; otherwise it will be added."`
}

// UpsertTool returns the MCP tool definition for upsert_variable.
func UpsertTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "upsert_variable",
		Description: "Add or update a single variable in a variables group. " +
			"If a variable with the same name already exists, it is updated; otherwise it is added. " +
			"This is safer than apply_variables_group when modifying a single variable, " +
			"because it does not affect other entries in the group. " +
			"Identify the target group by group_id or by org+group_slug. " +
			"The entry must include 'name' and either 'value' (a literal string) or " +
			"'value_from' (a reference to another resource's field). " +
			"Returns the updated variables group.",
	}
}

// UpsertHandler returns the typed tool handler for upsert_variable.
func UpsertHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpsertVariableInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UpsertVariableInput) (*mcp.CallToolResult, any, error) {
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
// delete_variable
// ---------------------------------------------------------------------------

// DeleteVariableInput defines the parameters for the delete_variable tool.
// The target group can be identified by group_id alone, or by org+group_slug.
type DeleteVariableInput struct {
	GroupID   string `json:"group_id,omitempty"    jsonschema:"The variables group ID. Mutually exclusive with org+group_slug."`
	Org       string `json:"org,omitempty"         jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'group_slug'."`
	GroupSlug string `json:"group_slug,omitempty"  jsonschema:"Variables group slug (name) for lookup within an organization. Must be paired with 'org'."`
	EntryName string `json:"entry_name"            jsonschema:"required,Name of the variable to remove from the group."`
}

// DeleteEntryTool returns the MCP tool definition for delete_variable.
func DeleteEntryTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_variable",
		Description: "Remove a single variable from a variables group. " +
			"Other variables in the group remain unchanged. " +
			"Identify the target group by group_id or by org+group_slug. " +
			"WARNING: Services referencing this variable via $variables-group/ will fail during deployment.",
	}
}

// DeleteEntryHandler returns the typed tool handler for delete_variable.
func DeleteEntryHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteVariableInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteVariableInput) (*mcp.CallToolResult, any, error) {
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
// get_variable_value
// ---------------------------------------------------------------------------

// GetVariableValueInput defines the parameters for the get_variable_value
// tool.
type GetVariableValueInput struct {
	Org       string `json:"org"        jsonschema:"required,Organization identifier."`
	GroupName string `json:"group_name" jsonschema:"required,Name (slug) of the variables group."`
	EntryName string `json:"entry_name" jsonschema:"required,Name of the variable whose value to retrieve."`
}

// GetValueTool returns the MCP tool definition for get_variable_value.
func GetValueTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_variable_value",
		Description: "Retrieve the resolved value of a specific variable from a variables group. " +
			"If the variable uses a value_from reference, the reference is resolved to its current value. " +
			"This is a convenient shortcut — instead of fetching the entire group and finding the entry, " +
			"you can look up a single value by org, group name, and variable name.",
	}
}

// GetValueHandler returns the typed tool handler for get_variable_value.
func GetValueHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetVariableValueInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetVariableValueInput) (*mcp.CallToolResult, any, error) {
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
// transform_variables
// ---------------------------------------------------------------------------

// TransformVariablesInput defines the parameters for the
// transform_variables tool.
type TransformVariablesInput struct {
	Org     string            `json:"org"     jsonschema:"required,Organization identifier. Determines which variables groups are available for reference resolution."`
	Entries map[string]string `json:"entries" jsonschema:"required,Map of environment variable names to values. Values starting with $variables-group/ will be resolved to their actual values. Other values pass through unchanged."`
}

// TransformTool returns the MCP tool definition for transform_variables.
func TransformTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "transform_variables",
		Description: "Batch-resolve $variables-group/ references in a set of environment variables. " +
			"Accepts a map of key-value pairs where values may contain $variables-group/<group-name>/<entry-name> references. " +
			"References are resolved to their actual values; literal values pass through unchanged. " +
			"Returns two maps: successfully transformed entries and any entries that failed resolution " +
			"(with error messages explaining why). " +
			"Useful for debugging configuration issues or previewing resolved values before deployment.",
	}
}

// TransformHandler returns the typed tool handler for transform_variables.
func TransformHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *TransformVariablesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *TransformVariablesInput) (*mcp.CallToolResult, any, error) {
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
		return fmt.Errorf("provide either 'id' or both 'org' and 'slug' to identify the variables group")
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
		return fmt.Errorf("provide either 'group_id' or both 'org' and 'group_slug' to identify the variables group")
	}
}
