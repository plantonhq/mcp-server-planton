// Package stackjob provides the MCP tools for the StackJob domain, backed by
// the StackJobQueryController, StackJobCommandController, and
// StackJobEssentialsQueryController RPCs on the Planton backend.
//
// Seven tools are exposed:
//   - get_stack_job:              retrieve a specific stack job by its ID
//   - get_latest_stack_job:       retrieve the most recent stack job for a cloud resource
//   - list_stack_jobs:            query stack jobs by organization with optional filters
//   - rerun_stack_job:            re-run a previously executed stack job
//   - cancel_stack_job:           gracefully cancel a running stack job
//   - resume_stack_job:           approve and resume an awaiting-approval stack job
//   - check_stack_job_essentials: pre-validate deployment prerequisites for a cloud resource kind
package stackjob

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// get_stack_job
// ---------------------------------------------------------------------------

// GetStackJobInput defines the parameters for the get_stack_job tool.
type GetStackJobInput struct {
	ID string `json:"id" jsonschema:"required,The stack job ID."`
}

// GetTool returns the MCP tool definition for get_stack_job.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_stack_job",
		Description: "Retrieve a specific stack job by its ID. " +
			"Returns the full job including operation type, progress, result, timestamps, errors, and IaC resource counts. " +
			"Use when you have a stack job ID from a previous response or from the user.",
	}
}

// GetHandler returns the typed tool handler for get_stack_job.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetStackJobInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetStackJobInput) (*mcp.CallToolResult, any, error) {
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
// get_latest_stack_job
// ---------------------------------------------------------------------------

// GetLatestStackJobInput defines the parameters for the get_latest_stack_job tool.
type GetLatestStackJobInput struct {
	CloudResourceID string `json:"cloud_resource_id" jsonschema:"required,The cloud resource ID to look up the most recent stack job for."`
}

// GetLatestTool returns the MCP tool definition for get_latest_stack_job.
func GetLatestTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_latest_stack_job",
		Description: "Retrieve the most recent stack job for a cloud resource. " +
			"This is the primary tool to check whether an apply_cloud_resource or destroy_cloud_resource operation completed successfully. " +
			"Returns the full stack job including operation type, progress, result, timestamps, errors, and IaC resource counts.",
	}
}

// GetLatestHandler returns the typed tool handler for get_latest_stack_job.
func GetLatestHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetLatestStackJobInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetLatestStackJobInput) (*mcp.CallToolResult, any, error) {
		if input.CloudResourceID == "" {
			return nil, nil, fmt.Errorf("'cloud_resource_id' is required")
		}
		text, err := GetLatest(ctx, serverAddress, input.CloudResourceID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_stack_jobs
// ---------------------------------------------------------------------------

// ListStackJobsInput defines the parameters for the list_stack_jobs tool.
type ListStackJobsInput struct {
	Org               string `json:"org"                          jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
	Env               string `json:"env,omitempty"                jsonschema:"Environment name to filter by."`
	CloudResourceKind string `json:"cloud_resource_kind,omitempty" jsonschema:"PascalCase cloud resource kind to filter by (e.g. AwsEksCluster). Read cloud-resource-kinds://catalog for valid kinds."`
	CloudResourceID   string `json:"cloud_resource_id,omitempty"  jsonschema:"Cloud resource ID to filter by."`
	OperationType     string `json:"operation_type,omitempty"     jsonschema:"Stack job operation type filter. One of: init, refresh, update_preview, update, destroy_preview, destroy."`
	Status            string `json:"status,omitempty"             jsonschema:"Execution status filter. One of: queued, running, completed, awaiting_approval."`
	Result            string `json:"result,omitempty"             jsonschema:"Execution result filter. One of: tbd, succeeded, failed, cancelled, skipped."`
	PageNum           int32  `json:"page_num,omitempty"           jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize          int32  `json:"page_size,omitempty"          jsonschema:"Number of results per page. Defaults to 20."`
}

// ListTool returns the MCP tool definition for list_stack_jobs.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_stack_jobs",
		Description: "List stack jobs matching the given filters. Requires an organization ID. " +
			"Supports filtering by environment, cloud resource kind, resource ID, operation type, execution status, and result. " +
			"Returns a paginated list. " +
			"Use to find failed deployments, audit provisioning history, or discover jobs across resources.",
	}
}

// ListHandler returns the typed tool handler for list_stack_jobs.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListStackJobsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListStackJobsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}

		text, err := List(ctx, serverAddress, ListInput{
			Org:               input.Org,
			Env:               input.Env,
			CloudResourceKind: input.CloudResourceKind,
			CloudResourceID:   input.CloudResourceID,
			OperationType:     input.OperationType,
			Status:            input.Status,
			Result:            input.Result,
			PageNum:           input.PageNum,
			PageSize:          input.PageSize,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// rerun_stack_job
// ---------------------------------------------------------------------------

// RerunStackJobInput defines the parameters for the rerun_stack_job tool.
type RerunStackJobInput struct {
	ID string `json:"id" jsonschema:"required,The stack job ID to re-run."`
}

// RerunTool returns the MCP tool definition for rerun_stack_job.
func RerunTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "rerun_stack_job",
		Description: "Re-run a previously executed stack job. " +
			"Use this to retry a failed deployment without recreating the cloud resource apply. " +
			"The new execution uses the same parameters as the original stack job. " +
			"Returns the updated stack job. Use get_stack_job to monitor progress.",
	}
}

// RerunHandler returns the typed tool handler for rerun_stack_job.
func RerunHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RerunStackJobInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RerunStackJobInput) (*mcp.CallToolResult, any, error) {
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
// cancel_stack_job
// ---------------------------------------------------------------------------

// CancelStackJobInput defines the parameters for the cancel_stack_job tool.
type CancelStackJobInput struct {
	ID string `json:"id" jsonschema:"required,The stack job ID to cancel."`
}

// CancelTool returns the MCP tool definition for cancel_stack_job.
func CancelTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "cancel_stack_job",
		Description: "Gracefully cancel a running stack job. " +
			"Cancellation is not immediate: the currently executing IaC operation " +
			"(e.g. pulumi up, terraform apply) completes fully, then remaining " +
			"operations are skipped and marked as cancelled. Infrastructure created " +
			"by completed operations remains — there is no automatic rollback. " +
			"The resource lock is released, allowing queued stack jobs to proceed. " +
			"The stack job must be in running status. " +
			"Returns the stack job; use get_stack_job to monitor cancellation progress.",
	}
}

// CancelHandler returns the typed tool handler for cancel_stack_job.
func CancelHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CancelStackJobInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CancelStackJobInput) (*mcp.CallToolResult, any, error) {
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
// resume_stack_job
// ---------------------------------------------------------------------------

// ResumeStackJobInput defines the parameters for the resume_stack_job tool.
type ResumeStackJobInput struct {
	ID string `json:"id" jsonschema:"required,The stack job ID to resume."`
}

// ResumeTool returns the MCP tool definition for resume_stack_job.
func ResumeTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "resume_stack_job",
		Description: "Approve and resume a stack job that is awaiting approval. " +
			"Stack jobs enter the awaiting_approval state when a flow control policy " +
			"requires manual approval before IaC execution proceeds. " +
			"This tool unblocks the job, allowing it to continue with its remaining operations. " +
			"To reject instead, use cancel_stack_job. " +
			"Returns the updated stack job.",
	}
}

// ResumeHandler returns the typed tool handler for resume_stack_job.
func ResumeHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ResumeStackJobInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ResumeStackJobInput) (*mcp.CallToolResult, any, error) {
		if input.ID == "" {
			return nil, nil, fmt.Errorf("'id' is required")
		}
		text, err := Resume(ctx, serverAddress, input.ID)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// check_stack_job_essentials
// ---------------------------------------------------------------------------

// CheckEssentialsInput defines the parameters for the
// check_stack_job_essentials tool.
type CheckEssentialsInput struct {
	CloudResourceKind string `json:"cloud_resource_kind" jsonschema:"required,PascalCase cloud resource kind (e.g. AwsEksCluster). Read cloud-resource-kinds://catalog for valid kinds."`
	Org               string `json:"org"                 jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
	Env               string `json:"env,omitempty"        jsonschema:"Environment name. Provide when the resource will be deployed to a specific environment."`
}

// CheckEssentialsTool returns the MCP tool definition for
// check_stack_job_essentials.
func CheckEssentialsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "check_stack_job_essentials",
		Description: "Pre-validate that all prerequisites for running a stack job are in place " +
			"for a given cloud resource kind and organization. " +
			"Returns four preflight checks: iac_module (IaC module resolved), " +
			"backend_credential (state backend configured), flow_control (approval policy resolved), " +
			"and provider_credential (cloud provider credentials available). " +
			"Each check includes a passed flag and any errors. " +
			"Use before apply_cloud_resource to catch missing configuration early.",
	}
}

// CheckEssentialsHandler returns the typed tool handler for
// check_stack_job_essentials.
func CheckEssentialsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *CheckEssentialsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *CheckEssentialsInput) (*mcp.CallToolResult, any, error) {
		if input.CloudResourceKind == "" {
			return nil, nil, fmt.Errorf("'cloud_resource_kind' is required")
		}
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := CheckEssentials(ctx, serverAddress, input.CloudResourceKind, input.Org, input.Env)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
