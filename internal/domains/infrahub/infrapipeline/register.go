package infrapipeline

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all infrapipeline tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, GetLatestTool(), GetLatestHandler(serverAddress))
	mcp.AddTool(srv, RunTool(), RunHandler(serverAddress))
	mcp.AddTool(srv, CancelTool(), CancelHandler(serverAddress))
	mcp.AddTool(srv, ResolveEnvGateTool(), ResolveEnvGateHandler(serverAddress))
	mcp.AddTool(srv, ResolveNodeGateTool(), ResolveNodeGateHandler(serverAddress))
}
