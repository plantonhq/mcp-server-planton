// Package environment provides the MCP tools for the Environment domain,
// backed by the EnvironmentQueryController RPCs
// (ai.planton.resourcemanager.environment.v1) on the Planton backend.
//
// One tool is exposed:
//   - list_environments: retrieve environments the caller can access within an organization
package environment

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ListEnvironmentsInput defines the parameters for the list_environments tool.
type ListEnvironmentsInput struct {
	Org string `json:"org" jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
}

// ListTool returns the MCP tool definition for list_environments.
func ListTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_environments",
		Description: "List environments the caller can access within an organization. " +
			"Returns only environments where the caller has at least view permission. " +
			"Use list_organizations first to discover available organization identifiers.",
	}
}

// ListHandler returns the typed tool handler for list_environments.
func ListHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListEnvironmentsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListEnvironmentsInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		text, err := List(ctx, serverAddress, input.Org)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
