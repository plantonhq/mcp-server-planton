package infraproject

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all infraproject tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, SearchTool(), SearchHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, CheckSlugTool(), CheckSlugHandler(serverAddress))
	mcp.AddTool(srv, UndeployTool(), UndeployHandler(serverAddress))
}
