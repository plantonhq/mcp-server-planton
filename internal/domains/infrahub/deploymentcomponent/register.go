package deploymentcomponent

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all deploymentcomponent tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, SearchTool(), SearchHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
}
