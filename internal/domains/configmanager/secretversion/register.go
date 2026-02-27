package secretversion

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all secretversion tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, CreateTool(), CreateHandler(serverAddress))
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
}
