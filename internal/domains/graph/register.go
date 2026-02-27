package graph

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all graph tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, GetOrganizationGraphTool(), GetOrganizationGraphHandler(serverAddress))
	mcp.AddTool(srv, GetEnvironmentGraphTool(), GetEnvironmentGraphHandler(serverAddress))
	mcp.AddTool(srv, GetServiceGraphTool(), GetServiceGraphHandler(serverAddress))
	mcp.AddTool(srv, GetCloudResourceGraphTool(), GetCloudResourceGraphHandler(serverAddress))
	mcp.AddTool(srv, GetDependenciesTool(), GetDependenciesHandler(serverAddress))
	mcp.AddTool(srv, GetDependentsTool(), GetDependentsHandler(serverAddress))
	mcp.AddTool(srv, GetImpactAnalysisTool(), GetImpactAnalysisHandler(serverAddress))
}
