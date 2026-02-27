package variable

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all variable tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, ResolveTool(), ResolveHandler(serverAddress))
}
