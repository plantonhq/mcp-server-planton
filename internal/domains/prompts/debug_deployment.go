package prompts

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// DebugDeploymentPrompt returns the prompt definition for diagnosing a failed
// infrastructure deployment. It guides the LLM through a multi-step sequence
// spanning the stackjob, cloudresource, and graph domains.
func DebugDeploymentPrompt() *mcp.Prompt {
	return &mcp.Prompt{
		Name: "debug_failed_deployment",
		Description: "Investigate and diagnose a failed infrastructure deployment (stack job). " +
			"Walks through error retrieval, AI-generated fix recommendations, IaC resource state " +
			"inspection, and upstream dependency analysis.",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "resource_id",
				Description: "Cloud resource or service ID whose deployment failed. If omitted, the conversation should start by identifying the resource.",
			},
			{
				Name:        "stack_job_id",
				Description: "Specific stack job ID to investigate. If omitted, the latest job for the resource is used.",
			},
		},
	}
}

// DebugDeploymentHandler returns the prompt handler that builds the diagnostic
// guidance message.
func DebugDeploymentHandler() mcp.PromptHandler {
	return func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		text := buildDebugDeploymentText(
			req.Params.Arguments["resource_id"],
			req.Params.Arguments["stack_job_id"],
		)
		return domains.PromptResult(
			"Diagnose a failed infrastructure deployment",
			domains.UserMessage(text),
		), nil
	}
}

func buildDebugDeploymentText(resourceID, stackJobID string) string {
	var b strings.Builder

	b.WriteString("I need help diagnosing a failed infrastructure deployment on Planton Cloud.")

	if resourceID != "" {
		b.WriteString("\n\nThe affected resource is: ")
		b.WriteString(resourceID)
	}
	if stackJobID != "" {
		b.WriteString("\n\nThe stack job to investigate is: ")
		b.WriteString(stackJobID)
	}

	b.WriteString(`

Recommended diagnostic approach:

1. Locate the failed stack job. `)

	if stackJobID != "" {
		b.WriteString("Use get_stack_job to retrieve job ")
		b.WriteString(stackJobID)
		b.WriteString(" and examine its status and error output.")
	} else if resourceID != "" {
		b.WriteString("Use get_latest_stack_job with resource ID ")
		b.WriteString(resourceID)
		b.WriteString(" to find the most recent job, then use get_stack_job to get its full details.")
	} else {
		b.WriteString("Ask me which resource or stack job to investigate, or use list_stack_jobs to find recent failures.")
	}

	b.WriteString(`

2. Once you have the failed stack job, call get_error_resolution_recommendation with the stack job ID. This returns an AI-generated analysis of the error with a suggested fix — it often identifies the root cause faster than manual inspection.

3. Review the IaC resource state by calling find_iac_resources_by_stack_job. This shows every Pulumi or Terraform resource and its status (created, failed, pending, deleted). Look for resources stuck in a failed or pending state.

4. If the error points to a configuration issue, use get_stack_job_input to retrieve the exact input that was fed to the IaC engine. This is credential-free and safe to inspect. Compare it against expected values.

5. Check whether upstream dependencies are healthy by calling get_dependencies on the resource. A failed dependency (e.g. a VPC or credential) can cause cascading deployment failures.

6. Summarize your findings: what failed, the likely root cause, and concrete next steps to resolve it. If a rerun might help, mention rerun_stack_job as an option.`)

	return b.String()
}
