package organization

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all organization tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
}
