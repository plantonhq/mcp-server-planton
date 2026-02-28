// Package tektontask provides the MCP tools for the ServiceHub TektonTask
// domain, backed by the TektonTaskQueryController and
// TektonTaskCommandController (ai.planton.servicehub.tektontask.v1) RPCs on
// the Planton backend.
//
// Three tools are exposed:
//   - get_tekton_task:    retrieve a Tekton task template by ID or org+slug
//   - apply_tekton_task:  create or update a Tekton task template (idempotent)
//   - delete_tekton_task: remove a Tekton task template
package tektontask

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// get_tekton_task
// ---------------------------------------------------------------------------

// GetTektonTaskInput defines the parameters for the get_tekton_task tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type GetTektonTaskInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The Tekton task ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Tekton task slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// GetTool returns the MCP tool definition for get_tekton_task.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_tekton_task",
		Description: "Retrieve the full details of a Tekton task template by its ID or by org+slug. " +
			"A Tekton task is a reusable CI/CD step definition (e.g., git-clone, docker-build) " +
			"that Tekton pipelines reference as individual build steps. " +
			"Returns the complete task including metadata, spec (description, git_repo, yaml_content, overview_markdown), " +
			"and audit status. " +
			"The output JSON can be modified and passed to apply_tekton_task for updates.",
	}
}

// GetHandler returns the typed tool handler for get_tekton_task.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetTektonTaskInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetTektonTaskInput) (*mcp.CallToolResult, any, error) {
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
// apply_tekton_task
// ---------------------------------------------------------------------------

// ApplyTektonTaskInput defines the parameters for the apply_tekton_task tool.
type ApplyTektonTaskInput struct {
	TektonTask map[string]any `json:"tekton_task" jsonschema:"required,The full TektonTask resource as a JSON object. Must include 'api_version' ('service-hub.planton.ai/v1'), 'kind' ('TektonTask'), 'metadata' (with 'name' and 'org'), and 'spec' (with 'selector' and optionally 'description', 'git_repo', 'yaml_content', 'overview_markdown'). The output of get_tekton_task can be modified and passed directly here."`
}

// ApplyTool returns the MCP tool definition for apply_tekton_task.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_tekton_task",
		Description: "Create or update a Tekton task template (idempotent). " +
			"Accepts the full TektonTask resource as a JSON object. " +
			"A Tekton task defines a reusable CI/CD step (e.g., git-clone, docker-build, deploy) " +
			"that Tekton pipelines reference as individual build steps. " +
			"For new tasks, provide api_version, kind, metadata (name, org), " +
			"and spec (selector, and optionally description, git_repo, yaml_content). " +
			"For updates, retrieve the task with get_tekton_task, modify desired fields, and pass here. " +
			"Returns the applied task with server-assigned ID and audit information.",
	}
}

// ApplyHandler returns the typed tool handler for apply_tekton_task.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyTektonTaskInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyTektonTaskInput) (*mcp.CallToolResult, any, error) {
		if len(input.TektonTask) == 0 {
			return nil, nil, fmt.Errorf("'tekton_task' is required and must be a non-empty JSON object")
		}
		text, err := Apply(ctx, serverAddress, input.TektonTask)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_tekton_task
// ---------------------------------------------------------------------------

// DeleteTektonTaskInput defines the parameters for the delete_tekton_task
// tool. Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type DeleteTektonTaskInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The Tekton task ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Tekton task slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// DeleteTool returns the MCP tool definition for delete_tekton_task.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_tekton_task",
		Description: "Delete a Tekton task template from the platform. " +
			"WARNING: Tekton pipelines referencing this task will fail during execution. " +
			"Ensure no pipelines reference this task before deleting. " +
			"Identify the task by ID or by org+slug.",
	}
}

// DeleteHandler returns the typed tool handler for delete_tekton_task.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteTektonTaskInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteTektonTaskInput) (*mcp.CallToolResult, any, error) {
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
		return fmt.Errorf("provide either 'id' or both 'org' and 'slug' to identify the Tekton task")
	}
}
