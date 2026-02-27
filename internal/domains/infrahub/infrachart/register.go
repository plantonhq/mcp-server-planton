package infrachart

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all infrachart tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, BuildTool(), BuildHandler(serverAddress))
}
