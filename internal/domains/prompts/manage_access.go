package prompts

import (
	"context"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ManageAccessPrompt returns the prompt definition for IAM discovery, audit,
// and policy management. It guides the LLM through the Planton IAM model —
// principals, policies, roles, teams, and API keys — providing a structured
// approach to understanding and modifying access controls.
func ManageAccessPrompt() *mcp.Prompt {
	return &mcp.Prompt{
		Name: "manage_access",
		Description: "Discover, audit, and manage access controls for platform resources. " +
			"Guides through principal discovery, permission checks, role lookup, " +
			"and IAM policy creation or revocation.",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "org_id",
				Description: "Organization ID to scope the access review. If omitted, available organizations are listed.",
			},
			{
				Name:        "resource_id",
				Description: "Specific resource ID to review access for. If omitted, the review covers the organization level.",
			},
		},
	}
}

// ManageAccessHandler returns the prompt handler that builds the access
// management guidance message.
func ManageAccessHandler() mcp.PromptHandler {
	return func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		text := buildManageAccessText(
			req.Params.Arguments["org_id"],
			req.Params.Arguments["resource_id"],
		)
		return domains.PromptResult(
			"Discover, audit, and manage access controls",
			domains.UserMessage(text),
		), nil
	}
}

func buildManageAccessText(orgID, resourceID string) string {
	var b strings.Builder

	b.WriteString("I need help managing access controls")
	if resourceID != "" {
		b.WriteString(" for resource ")
		b.WriteString(resourceID)
	}
	if orgID != "" {
		b.WriteString(" in organization ")
		b.WriteString(orgID)
	}
	b.WriteString(".")

	if orgID == "" {
		b.WriteString("\n\nIf no organization is specified, start by using list_organizations to identify the target organization.")
	}

	b.WriteString(`

The Planton IAM model has these key concepts:
- Principals: identity accounts (users) and teams that can be granted access
- IAM Policies: bindings that grant a principal a specific role on a specific resource
- IAM Roles: permission sets scoped to a resource kind (e.g. 'environment.admin', 'service.viewer')
- Teams: groups of identity accounts that can be granted access as a unit

Recommended approach:

1. Discover who has access. Use list_principals to see all principals (users and teams) with access. Filter by principal_kind ('identity_account' or 'team') and optionally by iam_role to narrow results.

2. For a specific resource, use list_resource_access to see all IAM policies attached to it — this shows which principals have which roles on that resource.

3. To verify a specific permission, use check_authorization with the principal ID, resource ID, and the permission to check. This returns a yes/no answer with the effective policy chain.

4. To understand what roles are available, use list_iam_roles_for_resource_kind with the target resource kind. This returns all assignable roles with their permission descriptions.

5. To grant access, use create_iam_policy with the principal ID, role ID, and resource ID. To revoke a specific policy, use delete_iam_policy with the policy ID.

6. For bulk changes, upsert_iam_policies performs a declarative sync — it takes a list of desired policies and reconciles them with the current state.

7. Use revoke_org_access only as a last resort — it removes ALL access for a principal across the entire organization. This is a nuclear option and should be used with extreme caution.`)

	return b.String()
}
