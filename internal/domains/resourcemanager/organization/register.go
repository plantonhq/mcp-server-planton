package organization

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all organization tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, CreateTool(), CreateHandler(serverAddress))
	mcp.AddTool(srv, UpdateTool(), UpdateHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
}
