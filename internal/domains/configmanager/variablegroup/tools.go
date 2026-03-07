package variablegroup

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	variablegroupv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/variablegroup/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// apply_variable_group
// ---------------------------------------------------------------------------

// ApplyVariableGroupInput defines the parameters for the apply_variable_group tool.
type ApplyVariableGroupInput struct {
	GroupObject map[string]any `json:"group_object" jsonschema:"required,Full VariableGroup object in OpenMCF envelope format."`
}

// ApplyTool returns the MCP tool definition for apply_variable_group.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_variable_group",
		Description: "Create or update a variable group (idempotent). " +
			"A variable group bundles related configuration variables into a single named, " +
			"scoped resource. Each entry has a name, value, optional description, and optional " +
			"external source reference. " +
			"Pass the full VariableGroup object as an OpenMCF envelope.",
	}
}

// ApplyHandler returns the typed tool handler for apply_variable_group.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyVariableGroupInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyVariableGroupInput) (*mcp.CallToolResult, any, error) {
		if input.GroupObject == nil {
			return nil, nil, fmt.Errorf("'group_object' is required")
		}
		text, err := Apply(ctx, serverAddress, input.GroupObject)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_variable_group
// ---------------------------------------------------------------------------

// GetVariableGroupInput defines the parameters for the get_variable_group tool.
type GetVariableGroupInput struct {
	ID    string `json:"id,omitempty"    jsonschema:"The variable group ID. Mutually exclusive with org+scope+slug."`
	Org   string `json:"org,omitempty"   jsonschema:"Organization identifier for slug-based lookup."`
	Scope string `json:"scope,omitempty" jsonschema:"Group scope: 'organization' or 'environment'."`
	Slug  string `json:"slug,omitempty"  jsonschema:"Variable group slug within the organization and scope."`
}

// GetTool returns the MCP tool definition for get_variable_group.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_variable_group",
		Description: "Retrieve a variable group by ID or by org+scope+slug. " +
			"Returns the full group including all entries with their current values. " +
			"Variable groups are uniquely identified within (org, scope, slug).",
	}
}

// GetHandler returns the typed tool handler for get_variable_group.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetVariableGroupInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetVariableGroupInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Scope, input.Slug); err != nil {
			return nil, nil, err
		}
		var scope variablegroupv1.VariableGroupSpec_Scope
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
// delete_variable_group
// ---------------------------------------------------------------------------

// DeleteVariableGroupInput defines the parameters for the delete_variable_group tool.
type DeleteVariableGroupInput struct {
	ID    string `json:"id,omitempty"    jsonschema:"The variable group ID. Mutually exclusive with org+scope+slug."`
	Org   string `json:"org,omitempty"   jsonschema:"Organization identifier for slug-based lookup."`
	Scope string `json:"scope,omitempty" jsonschema:"Group scope: 'organization' or 'environment'."`
	Slug  string `json:"slug,omitempty"  jsonschema:"Variable group slug within the organization and scope."`
}

// DeleteTool returns the MCP tool definition for delete_variable_group.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_variable_group",
		Description: "Delete a variable group and all its entries by ID or by org+scope+slug. " +
			"This is destructive and cannot be undone.",
	}
}

// DeleteHandler returns the typed tool handler for delete_variable_group.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteVariableGroupInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteVariableGroupInput) (*mcp.CallToolResult, any, error) {
		if err := validateIdentification(input.ID, input.Org, input.Scope, input.Slug); err != nil {
			return nil, nil, err
		}
		var scope variablegroupv1.VariableGroupSpec_Scope
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
// shared validation
// ---------------------------------------------------------------------------

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
		return fmt.Errorf("provide either 'id' or all of 'org', 'scope', and 'slug' to identify the variable group")
	}
}

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
