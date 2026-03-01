package apikey

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all API key tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, CreateTool(), CreateHandler(serverAddress))
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
}
