package prompts

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ExploreInfrastructurePrompt returns the prompt definition for a top-down
// discovery of an organization's infrastructure. It guides the LLM from
// organization selection through graph visualization to per-resource drill-down.
func ExploreInfrastructurePrompt() *mcp.Prompt {
	return &mcp.Prompt{
		Name: "explore_infrastructure",
		Description: "Get a comprehensive overview of an organization's infrastructure topology. " +
			"Walks through organization discovery, graph visualization, environment summary, " +
			"and health checks for key resources.",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "org_id",
				Description: "Organization ID to explore. If omitted, available organizations are listed first.",
			},
		},
	}
}

// ExploreInfrastructureHandler returns the prompt handler that builds the
// exploration guidance message.
func ExploreInfrastructureHandler() mcp.PromptHandler {
	return func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		text := buildExploreInfrastructureText(req.Params.Arguments["org_id"])
		return domains.PromptResult(
			"Explore an organization's infrastructure topology",
			domains.UserMessage(text),
		), nil
	}
}

func buildExploreInfrastructureText(orgID string) string {
	var b strings.Builder

	b.WriteString("I want to understand the infrastructure landscape")
	if orgID != "" {
		b.WriteString(" for organization ")
		b.WriteString(orgID)
	}
	b.WriteString(".")

	b.WriteString(`

For context, the platform's full catalog of resource types is available at the api-resource-kinds://catalog MCP resource.

Recommended exploration approach:

1. `)

	if orgID != "" {
		b.WriteString("Use get_organization to confirm the organization exists and get its details.")
	} else {
		b.WriteString("Use list_organizations to show available organizations and help me choose one.")
	}

	b.WriteString(`

2. Use get_organization_graph to retrieve the full resource topology. This returns all nodes (environments, services, cloud resources, credentials, infra projects) and their relationships. It is the single best tool for getting a big-picture view.

3. Summarize the topology:
   - How many environments exist and their names
   - How many cloud resources, services, and credentials are deployed
   - How resources are distributed across environments
   - Notable dependency patterns or architectural observations

4. For each environment, briefly describe what is deployed there. Use get_environment_graph if a deeper per-environment view is needed.

5. Spot-check health: for a few key cloud resources, call get_latest_stack_job to check whether the most recent deployment succeeded or failed. Flag any resources in a failed state.

6. If I want to drill deeper into a specific resource, use get_cloud_resource_graph or get_service_graph to show its dependencies and dependents.`)

	return b.String()
}
