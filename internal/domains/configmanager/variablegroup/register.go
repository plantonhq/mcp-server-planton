package variablegroup

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all variable group tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, UpsertEntryTool(), UpsertEntryHandler(serverAddress))
	mcp.AddTool(srv, DeleteEntryTool(), DeleteEntryHandler(serverAddress))
	mcp.AddTool(srv, RefreshEntryTool(), RefreshEntryHandler(serverAddress))
	mcp.AddTool(srv, RefreshAllTool(), RefreshAllHandler(serverAddress))
	mcp.AddTool(srv, ResolveEntryTool(), ResolveEntryHandler(serverAddress))
}
