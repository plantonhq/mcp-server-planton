package prompts

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ProvisionResourcePrompt returns the prompt definition for guided creation
// and deployment of a new cloud resource. It walks through prerequisite
// verification, configuration, deployment, and post-deployment monitoring
// across the connect, cloudresource, stackjob, and discovery domains.
func ProvisionResourcePrompt() *mcp.Prompt {
	return &mcp.Prompt{
		Name: "provision_cloud_resource",
		Description: "Guide through creating and deploying new cloud infrastructure. " +
			"Covers resource type selection, prerequisite checks (credentials, provider connections), " +
			"configuration, deployment, and post-deployment monitoring.",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "kind",
				Description: "Cloud resource kind to provision (e.g. 'AwsEksCluster', 'GcpGkeCluster'). If omitted, available kinds are shown first.",
			},
			{
				Name:        "org_id",
				Description: "Organization to provision in. If omitted, available organizations are listed.",
			},
			{
				Name:        "env_id",
				Description: "Environment to deploy to. If omitted, available environments are listed.",
			},
		},
	}
}

// ProvisionResourceHandler returns the prompt handler that builds the
// provisioning guidance message.
func ProvisionResourceHandler() mcp.PromptHandler {
	return func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		text := buildProvisionResourceText(
			req.Params.Arguments["kind"],
			req.Params.Arguments["org_id"],
			req.Params.Arguments["env_id"],
		)
		return domains.PromptResult(
			"Guide through provisioning new cloud infrastructure",
			domains.UserMessage(text),
		), nil
	}
}

func buildProvisionResourceText(kind, orgID, envID string) string {
	var b strings.Builder

	b.WriteString("I want to provision new cloud infrastructure")
	if kind != "" {
		b.WriteString(" of kind ")
		b.WriteString(kind)
	}
	b.WriteString(".")

	b.WriteString(`

The platform's cloud resource type catalog is available at the cloud-resource-kinds://catalog MCP resource, and per-type configuration schemas are available at cloud-resource-schema://{kind}.

Recommended provisioning approach:

1. Choose the resource type. `)

	if kind != "" {
		b.WriteString("Read cloud-resource-schema://")
		b.WriteString(kind)
		b.WriteString(" to get the configuration schema for this resource type.")
	} else {
		b.WriteString("Read cloud-resource-kinds://catalog to show available resource types and help me choose. Then read the schema for the chosen type.")
	}

	b.WriteString("\n\n2. Identify the target organization and environment. ")

	if orgID != "" && envID != "" {
		b.WriteString("Organization: ")
		b.WriteString(orgID)
		b.WriteString(", Environment: ")
		b.WriteString(envID)
		b.WriteString(".")
	} else if orgID != "" {
		b.WriteString("Organization: ")
		b.WriteString(orgID)
		b.WriteString(". Use list_environments to identify the target environment.")
	} else {
		b.WriteString("Use list_organizations to identify the organization, then list_environments for the environment.")
	}

	b.WriteString(`

3. Verify prerequisites before creating the resource:
   - Use search_credentials to confirm the required provider credentials exist in the organization (e.g. an AWS or GCP credential for the target cloud).
   - Use resolve_default_provider_connection to confirm a provider connection is configured for the target cloud provider and scope. If not, one must be created first with apply_default_provider_connection.

4. Check for pre-approved configurations by calling search_cloud_object_presets with the resource kind. Presets provide vetted defaults that follow organizational standards — use them when available.

5. Help me fill in the configuration based on the schema, then use apply_cloud_resource to create the resource. The tool takes a structured cloud_resource_object with api_version, kind, metadata (org, name, environment), and spec fields.

6. After creation, use get_latest_stack_job to monitor the deployment. Stack jobs represent the infrastructure-as-code execution (Pulumi or Terraform). Check the job status — it progresses through queued, in_progress, and eventually succeeded or failed.

7. If the deployment fails, use get_error_resolution_recommendation for an AI-generated diagnosis. Use find_iac_resources_by_stack_job to see which infrastructure resources were created, failed, or left pending.`)

	return b.String()
}
