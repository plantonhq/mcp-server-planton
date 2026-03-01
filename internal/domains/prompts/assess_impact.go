package prompts

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// AssessImpactPrompt returns the prompt definition for analyzing the blast
// radius of a change before executing it. It combines graph traversal, impact
// analysis, and version history into a safety-oriented workflow.
func AssessImpactPrompt() *mcp.Prompt {
	return &mcp.Prompt{
		Name: "assess_change_impact",
		Description: "Analyze the blast radius before modifying or deleting a platform resource. " +
			"Combines impact analysis, dependency graph traversal, and version history " +
			"into a risk assessment with a go/no-go recommendation.",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "resource_id",
				Description: "The platform resource ID to assess (e.g. a cloud resource, service, or credential ID).",
				Required:    true,
			},
			{
				Name:        "change_type",
				Description: "Type of planned change: 'delete' or 'update'. Defaults to 'delete' if omitted.",
			},
		},
	}
}

// AssessImpactHandler returns the prompt handler that builds the impact
// assessment guidance message.
func AssessImpactHandler() mcp.PromptHandler {
	return func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		text := buildAssessImpactText(
			req.Params.Arguments["resource_id"],
			req.Params.Arguments["change_type"],
		)
		return domains.PromptResult(
			"Assess the impact of a planned change before executing it",
			domains.UserMessage(text),
		), nil
	}
}

func buildAssessImpactText(resourceID, changeType string) string {
	if changeType == "" {
		changeType = "delete"
	}

	var b strings.Builder

	b.WriteString("Before making a destructive change, I need to understand the full impact.\n\n")
	b.WriteString("Resource to assess: ")
	b.WriteString(resourceID)
	b.WriteString("\nPlanned change: ")
	b.WriteString(changeType)

	b.WriteString(`

Recommended assessment approach:

1. Start with get_impact_analysis — pass the resource ID and change type. This returns the total count of affected resources (directly and transitively), broken down by resource type. It gives the clearest picture of blast radius.

2. Use get_dependents to list every downstream resource that depends on this one. These are the resources that would be affected if this resource is removed or changed. Pay attention to production-critical resources.

3. Use get_dependencies to understand the upstream context — what this resource itself depends on. This helps determine whether the resource can be safely recreated if needed.

4. Check the resource's recent change history with list_resource_versions. Understanding what changed recently can reveal whether this resource was recently modified (and the change might be the problem) or has been stable for a long time (and deletion has broader implications).

5. Present a clear risk assessment:
   - Total number of directly and transitively affected resources
   - Names and types of the most critical affected resources
   - Whether any production environments are in the blast radius
   - A go/no-go recommendation with reasoning
   - If proceeding, any precautions to take (e.g. backup, notify stakeholders, schedule a maintenance window)`)

	return b.String()
}
