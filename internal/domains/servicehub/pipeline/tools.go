// Package pipeline provides the MCP tools for the ServiceHub Pipeline domain,
// backed by the PipelineQueryController and PipelineCommandController
// RPCs (ai.planton.servicehub.pipeline.v1) on the Planton backend.
//
// Nine tools are exposed:
//   - list_pipelines:        list pipelines by org with optional service/env filters
//   - get_pipeline:          retrieve a pipeline by ID
//   - get_last_pipeline:     most recent pipeline for a service
//   - run_pipeline:          trigger a pipeline run for a branch/commit
//   - rerun_pipeline:        re-run a previously executed pipeline
//   - cancel_pipeline:       cancel a running pipeline
//   - resolve_pipeline_gate: approve/reject a deployment manual gate
//   - list_pipeline_files:   discover Tekton pipeline YAMLs in the service repo
//   - update_pipeline_file:  modify a pipeline file in the service repo
package pipeline

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// list_pipelines
// ---------------------------------------------------------------------------

// ListPipelinesInput defines the parameters for the list_pipelines tool.
type ListPipelinesInput struct {
	Org       string   `json:"org"                  jsonschema:"required,Organization identifier to scope results. Use list_organizations to discover available organizations."`
	ServiceID string   `json:"service_id,omitempty" jsonschema:"Service ID to filter pipelines by. When omitted, all pipelines in the organization are returned."`
	Envs      []string `json:"envs,omitempty"       jsonschema:"Environment names to filter by. Returns pipelines where any deployment task targets one of the specified environments."`
	PageNum   int32    `json:"page_num,omitempty"   jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize  int32    `json:"page_size,omitempty"  jsonschema:"Number of results per page. Defaults to 20."`
}

// ListTool returns the MCP tool definition for list_pipelines.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_pipelines",
		Description: "List CI/CD pipelines within an organization. " +
			"Pipelines represent build-and-deploy runs triggered by Git pushes or manual actions on services. " +
			"Optionally filter by service ID and/or environment names. " +
			"Returns a paginated list. Use get_pipeline with a pipeline ID from the results to retrieve full details.",
	}
}

// ListHandler returns the typed tool handler for list_pipelines.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListPipelinesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListPipelinesInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := List(ctx, serverAddress, ListInput{
			Org:       input.Org,
			ServiceID: input.ServiceID,
			Envs:      input.Envs,
			PageNum:   input.PageNum,
			PageSize:  input.PageSize,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_pipeline
// ---------------------------------------------------------------------------

// GetPipelineInput defines the parameters for the get_pipeline tool.
type GetPipelineInput struct {
	ID string `json:"id" jsonschema:"required,The pipeline ID."`
}

// GetTool returns the MCP tool definition for get_pipeline.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_pipeline",
		Description: "Retrieve the full details of a CI/CD pipeline by its ID. " +
			"Returns the complete pipeline including status, build stage, deployment tasks, " +
			"timestamps, and any errors. " +
			"Use this to check pipeline progress after triggering a run or to inspect a specific pipeline.",
	}
}

// GetHandler returns the typed tool handler for get_pipeline.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetPipelineInput) (*mcp.CallToolResult, any, error) {
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
// get_last_pipeline
// ---------------------------------------------------------------------------

// GetLastPipelineInput defines the parameters for the get_last_pipeline tool.
type GetLastPipelineInput struct {
	ServiceID string `json:"service_id" jsonschema:"required,The service ID to look up the most recent pipeline for."`
}

// GetLastTool returns the MCP tool definition for get_last_pipeline.
func GetLastTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_last_pipeline",
		Description: "Retrieve the most recent CI/CD pipeline for a service. " +
			"This is the primary tool to check whether a run_pipeline operation completed successfully. " +
			"Returns the full pipeline including status, build stage, deployment tasks, " +
			"timestamps, and any errors.",
	}
}

// GetLastHandler returns the typed tool handler for get_last_pipeline.
func GetLastHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetLastPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetLastPipelineInput) (*mcp.CallToolResult, any, error) {
		if input.ServiceID == "" {
			return nil, nil, fmt.Errorf("'service_id' is required")
		}
		text, err := GetLatest(ctx, serverAddress, input.ServiceID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// run_pipeline
// ---------------------------------------------------------------------------

// RunPipelineInput defines the parameters for the run_pipeline tool.
type RunPipelineInput struct {
	ServiceID string `json:"service_id"           jsonschema:"required,The service ID to run the pipeline for."`
	Branch    string `json:"branch"               jsonschema:"required,Git branch name to run the pipeline on."`
	CommitSHA string `json:"commit_sha,omitempty" jsonschema:"Git commit SHA to deploy. When omitted, the branch HEAD is used."`
}

// RunTool returns the MCP tool definition for run_pipeline.
func RunTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "run_pipeline",
		Description: "Trigger a new CI/CD pipeline run for a service on a specific Git branch. " +
			"Optionally provide a commit SHA to deploy that exact commit; " +
			"when omitted, the branch HEAD is used. " +
			"Use get_last_pipeline with the service ID to monitor the triggered pipeline.",
	}
}

// RunHandler returns the typed tool handler for run_pipeline.
func RunHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RunPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RunPipelineInput) (*mcp.CallToolResult, any, error) {
		if input.ServiceID == "" {
			return nil, nil, fmt.Errorf("'service_id' is required")
		}
		if input.Branch == "" {
			return nil, nil, fmt.Errorf("'branch' is required")
		}
		text, err := Run(ctx, serverAddress, input.ServiceID, input.Branch, input.CommitSHA)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// rerun_pipeline
// ---------------------------------------------------------------------------

// RerunPipelineInput defines the parameters for the rerun_pipeline tool.
type RerunPipelineInput struct {
	ID string `json:"id" jsonschema:"required,The pipeline ID to re-run."`
}

// RerunTool returns the MCP tool definition for rerun_pipeline.
func RerunTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "rerun_pipeline",
		Description: "Re-run a previously executed CI/CD pipeline using the same service, branch, and commit " +
			"configuration as the original execution. " +
			"Useful for retrying pipelines that failed due to transient issues. " +
			"Returns the newly created pipeline.",
	}
}

// RerunHandler returns the typed tool handler for rerun_pipeline.
func RerunHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RerunPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RerunPipelineInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := Rerun(ctx, serverAddress, input.ID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// cancel_pipeline
// ---------------------------------------------------------------------------

// CancelPipelineInput defines the parameters for the cancel_pipeline tool.
type CancelPipelineInput struct {
	ID string `json:"id" jsonschema:"required,The pipeline ID to cancel."`
}

// CancelTool returns the MCP tool definition for cancel_pipeline.
func CancelTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cancel_pipeline",
		Description: "Cancel a running CI/CD pipeline. " +
			"During the build stage, Tekton PipelineRun resources are deleted and running build pods are terminated. " +
			"During the deploy stage, the current deployment task receives a cancellation signal " +
			"and remaining tasks are skipped. " +
			"Returns the updated pipeline with its final status after cancellation.",
	}
}

// CancelHandler returns the typed tool handler for cancel_pipeline.
func CancelHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CancelPipelineInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CancelPipelineInput) (*mcp.CallToolResult, any, error) {
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
// resolve_pipeline_gate
// ---------------------------------------------------------------------------

// ResolvePipelineGateInput defines the parameters for the
// resolve_pipeline_gate tool.
type ResolvePipelineGateInput struct {
	PipelineID         string `json:"pipeline_id"          jsonschema:"required,The pipeline ID containing the manual gate."`
	DeploymentTaskName string `json:"deployment_task_name" jsonschema:"required,Name of the deployment task with the pending manual gate. Visible in the get_pipeline output."`
	Decision           string `json:"decision"             jsonschema:"required,The gate decision. Must be approve or reject."`
}

// ResolveGateTool returns the MCP tool definition for
// resolve_pipeline_gate.
func ResolveGateTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "resolve_pipeline_gate",
		Description: "Approve or reject a manual gate for a deployment task within a CI/CD pipeline. " +
			"Manual gates pause pipeline execution at a specific deployment task until a human " +
			"or agent explicitly approves or rejects. " +
			"WARNING: Approving a gate for a production deployment task will deploy to production. " +
			"Use get_pipeline to inspect which deployment tasks have pending gates.",
	}
}

// ResolveGateHandler returns the typed tool handler for
// resolve_pipeline_gate.
func ResolveGateHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ResolvePipelineGateInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ResolvePipelineGateInput) (*mcp.CallToolResult, any, error) {
		if input.PipelineID == "" {
			return nil, nil, fmt.Errorf("'pipeline_id' is required")
		}
		if input.DeploymentTaskName == "" {
			return nil, nil, fmt.Errorf("'deployment_task_name' is required")
		}
		if input.Decision == "" {
			return nil, nil, fmt.Errorf("'decision' is required")
		}
		text, err := ResolveGate(ctx, serverAddress, input.PipelineID, input.DeploymentTaskName, input.Decision)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_pipeline_files
// ---------------------------------------------------------------------------

// ListPipelineFilesInput defines the parameters for the
// list_pipeline_files tool.
type ListPipelineFilesInput struct {
	ServiceID string `json:"service_id"       jsonschema:"required,The service ID whose repository will be scanned for pipeline files."`
	Branch    string `json:"branch,omitempty" jsonschema:"Git branch to scan. When omitted, the service's default branch is used."`
}

// ListFilesTool returns the MCP tool definition for list_pipeline_files.
func ListFilesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_pipeline_files",
		Description: "List Tekton pipeline files in the Git repository connected to a service. " +
			"Discovers pipeline YAMLs under Planton conventions (.planton/, .tekton/, tekton/) " +
			"and returns their paths, content, and blob SHAs. " +
			"Use this to inspect current pipeline configuration before making changes with update_pipeline_file.",
	}
}

// ListFilesHandler returns the typed tool handler for list_pipeline_files.
func ListFilesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListPipelineFilesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListPipelineFilesInput) (*mcp.CallToolResult, any, error) {
		if input.ServiceID == "" {
			return nil, nil, fmt.Errorf("'service_id' is required")
		}
		text, err := ListFiles(ctx, serverAddress, input.ServiceID, input.Branch)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// update_pipeline_file
// ---------------------------------------------------------------------------

// UpdatePipelineFileInput defines the parameters for the
// update_pipeline_file tool.
type UpdatePipelineFileInput struct {
	ServiceID       string `json:"service_id"                  jsonschema:"required,The service ID whose repository will be updated."`
	Path            string `json:"path"                        jsonschema:"required,Path to write relative to repository root (e.g. .planton/pipeline.yaml)."`
	Content         string `json:"content"                     jsonschema:"required,New file content (plain text, typically YAML)."`
	ExpectedBaseSHA string `json:"expected_base_sha,omitempty" jsonschema:"When set, the write is rejected if the current blob SHA differs. Use the sha from list_pipeline_files for optimistic locking."`
	CommitMessage   string `json:"commit_message,omitempty"    jsonschema:"Custom commit message. When omitted, a default message is generated."`
	Branch          string `json:"branch,omitempty"            jsonschema:"Target branch. When omitted, the service's default branch is used."`
}

// UpdateFileTool returns the MCP tool definition for update_pipeline_file.
func UpdateFileTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "update_pipeline_file",
		Description: "Create or update a pipeline file in the Git repository connected to a service. " +
			"Commits the change directly to the specified branch (or the service's default branch). " +
			"For safe concurrent editing, provide expected_base_sha from list_pipeline_files — " +
			"the write is rejected if the file has been modified since. " +
			"Returns the new blob SHA, commit SHA, and branch name.",
	}
}

// UpdateFileHandler returns the typed tool handler for update_pipeline_file.
func UpdateFileHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpdatePipelineFileInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UpdatePipelineFileInput) (*mcp.CallToolResult, any, error) {
		if input.ServiceID == "" {
			return nil, nil, fmt.Errorf("'service_id' is required")
		}
		if input.Path == "" {
			return nil, nil, fmt.Errorf("'path' is required")
		}
		if input.Content == "" {
			return nil, nil, fmt.Errorf("'content' is required")
		}
		text, err := UpdateFile(ctx, serverAddress, input.ServiceID, input.Path, input.Content, input.ExpectedBaseSHA, input.CommitMessage, input.Branch)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
