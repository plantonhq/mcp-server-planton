// Package infrapipeline provides the MCP tools for the InfraPipeline domain,
// backed by the InfraPipelineQueryController and InfraPipelineCommandController
// RPCs (ai.planton.infrahub.infrapipeline.v1) on the Planton backend.
//
// Seven tools are exposed:
//   - list_infra_pipelines:              list pipelines by org with optional project filter
//   - get_infra_pipeline:                retrieve a pipeline by ID
//   - get_latest_infra_pipeline:         most recent pipeline for a project
//   - run_infra_pipeline:                trigger a pipeline run (chart-source or git-commit)
//   - cancel_infra_pipeline:             cancel a running pipeline
//   - resolve_infra_pipeline_env_gate:   approve/reject a manual gate for an environment
//   - resolve_infra_pipeline_node_gate:  approve/reject a manual gate for a DAG node
package infrapipeline

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// list_infra_pipelines
// ---------------------------------------------------------------------------

// ListInfraPipelinesInput defines the parameters for the
// list_infra_pipelines tool.
type ListInfraPipelinesInput struct {
	Org            string `json:"org"                       jsonschema:"required,Organization identifier to scope results. Use list_organizations to discover available organizations."`
	InfraProjectID string `json:"infra_project_id,omitempty" jsonschema:"Infra project ID to filter pipelines by. When omitted, all pipelines in the organization are returned."`
	PageNum        int32  `json:"page_num,omitempty"         jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize       int32  `json:"page_size,omitempty"        jsonschema:"Number of results per page. Defaults to 20."`
}

// ListTool returns the MCP tool definition for list_infra_pipelines.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_infra_pipelines",
		Description: "List infra pipelines within an organization. " +
			"Infra pipelines represent deployment runs triggered by apply or run operations on infra projects. " +
			"Optionally filter by infra project ID to see pipelines for a specific project. " +
			"Returns a paginated list. Use get_infra_pipeline with a pipeline ID from the results to retrieve full details.",
	}
}

// ListHandler returns the typed tool handler for list_infra_pipelines.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListInfraPipelinesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListInfraPipelinesInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := List(ctx, serverAddress, ListInput{
			Org:            input.Org,
			InfraProjectID: input.InfraProjectID,
			PageNum:        input.PageNum,
			PageSize:       input.PageSize,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_infra_pipeline
// ---------------------------------------------------------------------------

// GetInfraPipelineInput defines the parameters for the get_infra_pipeline tool.
type GetInfraPipelineInput struct {
	ID string `json:"id" jsonschema:"required,The infra pipeline ID."`
}

// GetTool returns the MCP tool definition for get_infra_pipeline.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_infra_pipeline",
		Description: "Retrieve the full details of an infra pipeline by its ID. " +
			"Returns the complete pipeline including status, environment stages, DAG nodes, " +
			"timestamps, and any errors. " +
			"Use this to check pipeline progress after triggering a run or to inspect a specific pipeline.",
	}
}

// GetHandler returns the typed tool handler for get_infra_pipeline.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInfraPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInfraPipelineInput) (*mcp.CallToolResult, any, error) {
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
// get_latest_infra_pipeline
// ---------------------------------------------------------------------------

// GetLatestInfraPipelineInput defines the parameters for the
// get_latest_infra_pipeline tool.
type GetLatestInfraPipelineInput struct {
	InfraProjectID string `json:"infra_project_id" jsonschema:"required,The infra project ID to look up the most recent pipeline for."`
}

// GetLatestTool returns the MCP tool definition for get_latest_infra_pipeline.
func GetLatestTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_latest_infra_pipeline",
		Description: "Retrieve the most recent infra pipeline for an infra project. " +
			"This is the primary tool to check whether an apply_infra_project or run_infra_pipeline " +
			"operation completed successfully. " +
			"Returns the full pipeline including status, environment stages, DAG nodes, " +
			"timestamps, and any errors.",
	}
}

// GetLatestHandler returns the typed tool handler for get_latest_infra_pipeline.
func GetLatestHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetLatestInfraPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetLatestInfraPipelineInput) (*mcp.CallToolResult, any, error) {
		if input.InfraProjectID == "" {
			return nil, nil, fmt.Errorf("'infra_project_id' is required")
		}
		text, err := GetLatest(ctx, serverAddress, input.InfraProjectID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// run_infra_pipeline
// ---------------------------------------------------------------------------

// RunInfraPipelineInput defines the parameters for the run_infra_pipeline tool.
type RunInfraPipelineInput struct {
	InfraProjectID string `json:"infra_project_id" jsonschema:"required,The infra project ID to run the pipeline for."`
	CommitSHA      string `json:"commit_sha,omitempty" jsonschema:"Git commit SHA to deploy. Required for git-repo sourced projects. Omit for infra-chart sourced projects."`
}

// RunTool returns the MCP tool definition for run_infra_pipeline.
func RunTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "run_infra_pipeline",
		Description: "Trigger a new infra pipeline run for a project. " +
			"For chart-sourced projects, omit commit_sha — the latest chart configuration is used. " +
			"For git-repo sourced projects, provide commit_sha to deploy a specific commit. " +
			"Returns the newly created pipeline ID. " +
			"Use get_infra_pipeline or get_latest_infra_pipeline to monitor progress.",
	}
}

// RunHandler returns the typed tool handler for run_infra_pipeline.
func RunHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RunInfraPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RunInfraPipelineInput) (*mcp.CallToolResult, any, error) {
		if input.InfraProjectID == "" {
			return nil, nil, fmt.Errorf("'infra_project_id' is required")
		}
		text, err := Run(ctx, serverAddress, input.InfraProjectID, input.CommitSHA)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// cancel_infra_pipeline
// ---------------------------------------------------------------------------

// CancelInfraPipelineInput defines the parameters for the
// cancel_infra_pipeline tool.
type CancelInfraPipelineInput struct {
	ID string `json:"id" jsonschema:"required,The infra pipeline ID to cancel."`
}

// CancelTool returns the MCP tool definition for cancel_infra_pipeline.
func CancelTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cancel_infra_pipeline",
		Description: "Cancel a running infra pipeline. " +
			"Returns the updated pipeline with its final status after cancellation.",
	}
}

// CancelHandler returns the typed tool handler for cancel_infra_pipeline.
func CancelHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CancelInfraPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CancelInfraPipelineInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := Cancel(ctx, serverAddress, input.ID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// resolve_infra_pipeline_env_gate
// ---------------------------------------------------------------------------

// ResolveEnvGateInput defines the parameters for the
// resolve_infra_pipeline_env_gate tool.
type ResolveEnvGateInput struct {
	InfraPipelineID string `json:"infra_pipeline_id" jsonschema:"required,The infra pipeline ID containing the manual gate."`
	Env             string `json:"env"               jsonschema:"required,Environment name where the manual gate is pending (e.g. staging or production)."`
	Decision        string `json:"decision"          jsonschema:"required,The gate decision. Must be approve or reject."`
}

// ResolveEnvGateTool returns the MCP tool definition for
// resolve_infra_pipeline_env_gate.
func ResolveEnvGateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "resolve_infra_pipeline_env_gate",
		Description: "Approve or reject a manual gate for an entire deployment environment " +
			"within an infra pipeline. Manual gates pause pipeline execution until a human " +
			"or agent explicitly approves or rejects. " +
			"Use get_infra_pipeline to inspect which environments have pending gates.",
	}
}

// ResolveEnvGateHandler returns the typed tool handler for
// resolve_infra_pipeline_env_gate.
func ResolveEnvGateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ResolveEnvGateInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ResolveEnvGateInput) (*mcp.CallToolResult, any, error) {
		if input.InfraPipelineID == "" {
			return nil, nil, fmt.Errorf("'infra_pipeline_id' is required")
		}
		if input.Env == "" {
			return nil, nil, fmt.Errorf("'env' is required")
		}
		if input.Decision == "" {
			return nil, nil, fmt.Errorf("'decision' is required")
		}
		text, err := ResolveEnvGate(ctx, serverAddress, input.InfraPipelineID, input.Env, input.Decision)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// resolve_infra_pipeline_node_gate
// ---------------------------------------------------------------------------

// ResolveNodeGateInput defines the parameters for the
// resolve_infra_pipeline_node_gate tool.
type ResolveNodeGateInput struct {
	InfraPipelineID string `json:"infra_pipeline_id" jsonschema:"required,The infra pipeline ID containing the manual gate."`
	Env             string `json:"env"               jsonschema:"required,Environment name where the node exists."`
	NodeID          string `json:"node_id"           jsonschema:"required,Node identifier in the format CloudResourceKind/slug (e.g. KubernetesOpenFga/fga-gcp-dev). Visible in the get_infra_pipeline output."`
	Decision        string `json:"decision"          jsonschema:"required,The gate decision. Must be approve or reject."`
}

// ResolveNodeGateTool returns the MCP tool definition for
// resolve_infra_pipeline_node_gate.
func ResolveNodeGateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "resolve_infra_pipeline_node_gate",
		Description: "Approve or reject a manual gate for a specific DAG node " +
			"within an infra pipeline. DAG nodes represent individual cloud resources " +
			"in the deployment graph. Manual gates pause that node's deployment until resolved. " +
			"Use get_infra_pipeline to inspect which nodes have pending gates.",
	}
}

// ResolveNodeGateHandler returns the typed tool handler for
// resolve_infra_pipeline_node_gate.
func ResolveNodeGateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ResolveNodeGateInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ResolveNodeGateInput) (*mcp.CallToolResult, any, error) {
		if input.InfraPipelineID == "" {
			return nil, nil, fmt.Errorf("'infra_pipeline_id' is required")
		}
		if input.Env == "" {
			return nil, nil, fmt.Errorf("'env' is required")
		}
		if input.NodeID == "" {
			return nil, nil, fmt.Errorf("'node_id' is required")
		}
		if input.Decision == "" {
			return nil, nil, fmt.Errorf("'decision' is required")
		}
		text, err := ResolveNodeGate(ctx, serverAddress, input.InfraPipelineID, input.Env, input.NodeID, input.Decision)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
