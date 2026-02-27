// Package audit — see doc.go for overview.
//
// This file defines the MCP tool surface: input structs, tool definitions,
// and typed handlers for the three audit tools.
package audit

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// list_resource_versions
// ---------------------------------------------------------------------------

// ListResourceVersionsInput defines the parameters for the
// list_resource_versions tool.
type ListResourceVersionsInput struct {
	ResourceID string `json:"resource_id"        jsonschema:"required,The ID of the resource to retrieve version history for."`
	Kind       string `json:"kind"               jsonschema:"required,Platform resource kind. Common values: cloud_resource, infra_project, infra_chart, infra_pipeline, variable, secret, environment, organization, service, stack_job."`
	PageNum    int32  `json:"page_num,omitempty"  jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty" jsonschema:"Number of results per page. Defaults to 20."`
}

// ListTool returns the MCP tool definition for list_resource_versions.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_resource_versions",
		Description: "List the version history for a specific platform resource. " +
			"Requires the resource ID and its kind (e.g. cloud_resource, infra_project, variable). " +
			"Returns a paginated list of versions with metadata, event type, and timestamps. " +
			"Each version entry includes an ID that can be passed to get_resource_version for full details and diffs. " +
			"Use get_resource_version_count first if you only need to know whether changes exist.",
	}
}

// ListHandler returns the typed tool handler for list_resource_versions.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListResourceVersionsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListResourceVersionsInput) (*mcp.CallToolResult, any, error) {
		if input.ResourceID == "" {
			return nil, nil, fmt.Errorf("'resource_id' is required")
		}
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}

		kind, err := apiResourceKindResolver.Resolve(input.Kind)
		if err != nil {
			return nil, nil, err
		}

		text, err := List(ctx, serverAddress, ListInput{
			Kind:       kind,
			ResourceID: input.ResourceID,
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
// get_resource_version
// ---------------------------------------------------------------------------

// GetResourceVersionInput defines the parameters for the
// get_resource_version tool.
type GetResourceVersionInput struct {
	VersionID   string `json:"version_id"            jsonschema:"required,The resource version ID. Obtain from list_resource_versions results."`
	ContextSize int32  `json:"context_size,omitempty" jsonschema:"Number of surrounding lines to include in the unified diff. Analogous to git diff -U<n>. Defaults to 3."`
}

// GetTool returns the MCP tool definition for get_resource_version.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_resource_version",
		Description: "Retrieve a specific resource version with full change details. " +
			"Returns the original and new state as YAML, a unified diff, the event type (create, update, delete), " +
			"linked stack job ID, and cloud object version details when applicable. " +
			"The context_size parameter controls diff context lines (default 3, like git diff -U3). " +
			"Use list_resource_versions to discover version IDs.",
	}
}

// GetHandler returns the typed tool handler for get_resource_version.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetResourceVersionInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetResourceVersionInput) (*mcp.CallToolResult, any, error) {
		if input.VersionID == "" {
			return nil, nil, fmt.Errorf("'version_id' is required")
		}

		text, err := Get(ctx, serverAddress, input.VersionID, input.ContextSize)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_resource_version_count
// ---------------------------------------------------------------------------

// GetResourceVersionCountInput defines the parameters for the
// get_resource_version_count tool.
type GetResourceVersionCountInput struct {
	ResourceID string `json:"resource_id" jsonschema:"required,The ID of the resource to count versions for."`
	Kind       string `json:"kind"        jsonschema:"required,Platform resource kind. Common values: cloud_resource, infra_project, infra_chart, infra_pipeline, variable, secret, environment, organization, service, stack_job."`
}

// CountTool returns the MCP tool definition for get_resource_version_count.
func CountTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_resource_version_count",
		Description: "Get the number of versions that exist for a specific platform resource. " +
			"This is a lightweight query — no version data is transferred. " +
			"Use to quickly check whether a resource has any change history, " +
			"or to estimate pagination before calling list_resource_versions.",
	}
}

// CountHandler returns the typed tool handler for get_resource_version_count.
func CountHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetResourceVersionCountInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetResourceVersionCountInput) (*mcp.CallToolResult, any, error) {
		if input.ResourceID == "" {
			return nil, nil, fmt.Errorf("'resource_id' is required")
		}
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}

		kind, err := apiResourceKindResolver.Resolve(input.Kind)
		if err != nil {
			return nil, nil, err
		}

		text, err := Count(ctx, serverAddress, kind, input.ResourceID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
