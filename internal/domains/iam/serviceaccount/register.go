package serviceaccount

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all service account tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, CreateTool(), CreateHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, UpdateTool(), UpdateHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
	mcp.AddTool(srv, CreateKeyTool(), CreateKeyHandler(serverAddress))
	mcp.AddTool(srv, RevokeKeyTool(), RevokeKeyHandler(serverAddress))
	mcp.AddTool(srv, ListKeysTool(), ListKeysHandler(serverAddress))
}
