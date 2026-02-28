// Package tektonpipeline provides the MCP tools for the ServiceHub
// TektonPipeline domain, backed by the TektonPipelineQueryController and
// TektonPipelineCommandController (ai.planton.servicehub.tektonpipeline.v1)
// RPCs on the Planton backend.
//
// Three tools are exposed:
//   - get_tekton_pipeline:    retrieve a Tekton pipeline template by ID or org+slug
//   - apply_tekton_pipeline:  create or update a Tekton pipeline template (idempotent)
//   - delete_tekton_pipeline: remove a Tekton pipeline template
package tektonpipeline

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// get_tekton_pipeline
// ---------------------------------------------------------------------------

// GetTektonPipelineInput defines the parameters for the get_tekton_pipeline
// tool. Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type GetTektonPipelineInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The Tekton pipeline ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Tekton pipeline slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// GetTool returns the MCP tool definition for get_tekton_pipeline.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_tekton_pipeline",
		Description: "Retrieve the full details of a Tekton pipeline template by its ID or by org+slug. " +
			"A Tekton pipeline is a reusable CI/CD pipeline definition that services reference " +
			"to orchestrate their build and deployment stages. " +
			"Returns the complete pipeline including metadata, spec (description, git_repo, yaml_content, overview_markdown), " +
			"and audit status. " +
			"The output JSON can be modified and passed to apply_tekton_pipeline for updates.",
	}
}

// GetHandler returns the typed tool handler for get_tekton_pipeline.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetTektonPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetTektonPipelineInput) (*mcp.CallToolResult, any, error) {
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
// apply_tekton_pipeline
// ---------------------------------------------------------------------------

// ApplyTektonPipelineInput defines the parameters for the
// apply_tekton_pipeline tool.
type ApplyTektonPipelineInput struct {
	TektonPipeline map[string]any `json:"tekton_pipeline" jsonschema:"required,The full TektonPipeline resource as a JSON object. Must include 'api_version' ('service-hub.planton.ai/v1'), 'kind' ('TektonPipeline'), 'metadata' (with 'name' and 'org'), and 'spec' (with 'selector' and optionally 'description', 'git_repo', 'yaml_content', 'overview_markdown'). The output of get_tekton_pipeline can be modified and passed directly here."`
}

// ApplyTool returns the MCP tool definition for apply_tekton_pipeline.
func ApplyTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "apply_tekton_pipeline",
		Description: "Create or update a Tekton pipeline template (idempotent). " +
			"Accepts the full TektonPipeline resource as a JSON object. " +
			"A Tekton pipeline defines a reusable CI/CD pipeline that services reference " +
			"for orchestrating build and deployment stages. " +
			"For new pipelines, provide api_version, kind, metadata (name, org), " +
			"and spec (selector, and optionally description, git_repo, yaml_content). " +
			"For updates, retrieve the pipeline with get_tekton_pipeline, modify desired fields, and pass here. " +
			"Returns the applied pipeline with server-assigned ID and audit information.",
	}
}

// ApplyHandler returns the typed tool handler for apply_tekton_pipeline.
func ApplyHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ApplyTektonPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ApplyTektonPipelineInput) (*mcp.CallToolResult, any, error) {
		if len(input.TektonPipeline) == 0 {
			return nil, nil, fmt.Errorf("'tekton_pipeline' is required and must be a non-empty JSON object")
		}
		text, err := Apply(ctx, serverAddress, input.TektonPipeline)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// delete_tekton_pipeline
// ---------------------------------------------------------------------------

// DeleteTektonPipelineInput defines the parameters for the
// delete_tekton_pipeline tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Slug path: set both 'org' and 'slug'.
type DeleteTektonPipelineInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The Tekton pipeline ID. Mutually exclusive with org+slug."`
	Org  string `json:"org,omitempty"  jsonschema:"Organization identifier for slug-based lookup. Must be paired with 'slug'."`
	Slug string `json:"slug,omitempty" jsonschema:"Tekton pipeline slug (name) for lookup within an organization. Must be paired with 'org'."`
}

// DeleteTool returns the MCP tool definition for delete_tekton_pipeline.
func DeleteTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_tekton_pipeline",
		Description: "Delete a Tekton pipeline template from the platform. " +
			"WARNING: Services referencing this pipeline will lose their CI/CD pipeline definition. " +
			"Ensure no services reference this pipeline before deleting. " +
			"Identify the pipeline by ID or by org+slug.",
	}
}

// DeleteHandler returns the typed tool handler for delete_tekton_pipeline.
func DeleteHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteTektonPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteTektonPipelineInput) (*mcp.CallToolResult, any, error) {
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
		return fmt.Errorf("provide either 'id' or both 'org' and 'slug' to identify the Tekton pipeline")
	}
}
