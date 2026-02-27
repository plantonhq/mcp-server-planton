// Package organization provides the MCP tools for the Organization domain,
// backed by the OrganizationQueryController RPCs
// (ai.planton.resourcemanager.organization.v1) on the Planton backend.
//
// One tool is exposed:
//   - list_organizations: retrieve all organizations the caller is a member of
package organization

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ListOrganizationsInput defines the parameters for the list_organizations
// tool. This tool takes no input — the server returns organizations scoped to
// the authenticated caller's membership.
type ListOrganizationsInput struct{}

// ListTool returns the MCP tool definition for list_organizations.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_organizations",
		Description: "List all organizations the caller is a member of. " +
			"Returns the full organization objects including id, name, and slug. " +
			"Use this as the first step to discover the operating context before working with cloud resources or environments.",
	}
}

// ListHandler returns the typed tool handler for list_organizations.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListOrganizationsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ *ListOrganizationsInput) (*mcp.CallToolResult, any, error) {
		text, err := List(ctx, serverAddress)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
