// Package infrachart provides the MCP tools for the InfraChart domain, backed
// by the InfraChartQueryController RPCs
// (ai.planton.infrahub.infrachart.v1) on the Planton backend.
//
// Three tools are exposed:
//   - list_infra_charts:  paginated listing of infra charts with org/env filters
//   - get_infra_chart:    retrieve a single infra chart by ID
//   - build_infra_chart:  preview rendered output by applying param overrides
package infrachart

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// list_infra_charts
// ---------------------------------------------------------------------------

// ListInfraChartsInput defines the parameters for the list_infra_charts tool.
// All fields are optional — an empty input returns the first page of all
// infra charts visible to the caller.
type ListInfraChartsInput struct {
	Org      string `json:"org,omitempty"       jsonschema:"Organization identifier to scope results. Use list_organizations to discover available organizations."`
	Env      string `json:"env,omitempty"       jsonschema:"Environment identifier to scope results. Use list_environments to discover available environments."`
	PageNum  int32  `json:"page_num,omitempty"  jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize int32  `json:"page_size,omitempty" jsonschema:"Number of results per page. Defaults to 20."`
}

// ListTool returns the MCP tool definition for list_infra_charts.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_infra_charts",
		Description: "List infra chart templates with optional organization and environment filters. " +
			"Infra charts are reusable infrastructure-as-code templates that define cloud resource compositions. " +
			"Returns a paginated list of charts including metadata, description, and parameters. " +
			"Use get_infra_chart with a chart ID from the results to retrieve full details.",
	}
}

// ListHandler returns the typed tool handler for list_infra_charts.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListInfraChartsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListInfraChartsInput) (*mcp.CallToolResult, any, error) {
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
// get_infra_chart
// ---------------------------------------------------------------------------

// GetInfraChartInput defines the parameters for the get_infra_chart tool.
type GetInfraChartInput struct {
	ID string `json:"id" jsonschema:"required,The infra chart ID obtained from list_infra_charts results."`
}

// GetTool returns the MCP tool definition for get_infra_chart.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_infra_chart",
		Description: "Retrieve the full details of an infra chart by its ID. " +
			"Returns the complete chart including template YAML files, values.yaml, parameter definitions, " +
			"description, and web links. " +
			"Use build_infra_chart to preview the rendered output with custom parameter values.",
	}
}

// GetHandler returns the typed tool handler for get_infra_chart.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInfraChartInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInfraChartInput) (*mcp.CallToolResult, any, error) {
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

// ---------------------------------------------------------------------------
// build_infra_chart
// ---------------------------------------------------------------------------

// BuildInfraChartInput defines the parameters for the build_infra_chart tool.
type BuildInfraChartInput struct {
	ChartID string         `json:"chart_id" jsonschema:"required,The infra chart ID to build. The chart is fetched automatically — you only need to supply parameter overrides."`
	Params  map[string]any `json:"params,omitempty"  jsonschema:"Parameter overrides as a name-to-value map. Keys must match parameter names from the chart's param definitions (visible in get_infra_chart output). Unspecified params keep their chart defaults."`
}

// BuildTool returns the MCP tool definition for build_infra_chart.
func BuildTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "build_infra_chart",
		Description: "Preview the rendered output of an infra chart by applying parameter overrides. " +
			"Fetches the chart by ID, merges the supplied params with the chart defaults, " +
			"and returns the rendered YAML and cloud resource DAG without persisting anything. " +
			"Use this to validate parameter choices before creating an infra project from the chart.",
	}
}

// BuildHandler returns the typed tool handler for build_infra_chart.
func BuildHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *BuildInfraChartInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *BuildInfraChartInput) (*mcp.CallToolResult, any, error) {
		if input.ChartID == "" {
			return nil, nil, fmt.Errorf("'chart_id' is required")
		}
		text, err := Build(ctx, serverAddress, BuildInput{
			ChartID: input.ChartID,
			Params:  input.Params,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
