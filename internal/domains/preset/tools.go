// Package preset provides the MCP tools for the CloudObjectPreset domain,
// backed by the InfraHubSearchQueryController RPCs
// (ai.planton.search.v1.infrahub) and the CloudObjectPresetQueryController
// RPCs (ai.planton.infrahub.cloudobjectpreset.v1) on the Planton backend.
//
// Two tools are exposed:
//   - search_cloud_object_presets: search for preset templates (official + org-scoped)
//   - get_cloud_object_preset: retrieve full preset content by ID
package preset

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// search_cloud_object_presets
// ---------------------------------------------------------------------------

// SearchCloudObjectPresetsInput defines the parameters for the
// search_cloud_object_presets tool. All fields are optional — an empty input
// returns all official presets.
type SearchCloudObjectPresetsInput struct {
	Org        string `json:"org,omitempty"         jsonschema:"Organization identifier. When provided, results include both official and organization-specific presets. When omitted, only official presets are returned."`
	Kind       string `json:"kind,omitempty"        jsonschema:"PascalCase cloud resource kind to filter by (e.g. AwsEksCluster). Read cloud-resource-kinds://catalog for valid kinds."`
	SearchText string `json:"search_text,omitempty" jsonschema:"Free-text search query to filter presets by name or description."`
}

// SearchTool returns the MCP tool definition for search_cloud_object_presets.
func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_cloud_object_presets",
		Description: "Search for cloud object preset templates. " +
			"Presets are pre-configured cloud resource manifests that can be used as starting points for apply_cloud_resource. " +
			"When 'org' is provided, results include both official platform presets and organization-specific presets. " +
			"When 'org' is omitted, only official presets are returned. " +
			"Use get_cloud_object_preset with the preset ID from the results to retrieve the full YAML content.",
	}
}

// SearchHandler returns the typed tool handler for search_cloud_object_presets.
func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchCloudObjectPresetsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchCloudObjectPresetsInput) (*mcp.CallToolResult, any, error) {
		text, err := Search(ctx, serverAddress, SearchInput{
			Org:        input.Org,
			Kind:       input.Kind,
			SearchText: input.SearchText,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_cloud_object_preset
// ---------------------------------------------------------------------------

// GetCloudObjectPresetInput defines the parameters for the
// get_cloud_object_preset tool.
type GetCloudObjectPresetInput struct {
	ID string `json:"id" jsonschema:"required,The preset ID obtained from search_cloud_object_presets results."`
}

// GetTool returns the MCP tool definition for get_cloud_object_preset.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_cloud_object_preset",
		Description: "Get the full content of a cloud object preset by ID. " +
			"Returns the complete preset including YAML manifest content, markdown documentation, " +
			"cloud resource kind, rank, and provider metadata. " +
			"Use the YAML content as a template for apply_cloud_resource, replacing placeholder values as needed.",
	}
}

// GetHandler returns the typed tool handler for get_cloud_object_preset.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetCloudObjectPresetInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetCloudObjectPresetInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}

		text, err := Get(ctx, serverAddress, input.ID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
