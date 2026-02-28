package secretsgroup

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all secrets group tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, SearchTool(), SearchHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, UpsertTool(), UpsertHandler(serverAddress))
	mcp.AddTool(srv, DeleteEntryTool(), DeleteEntryHandler(serverAddress))
	mcp.AddTool(srv, GetValueTool(), GetValueHandler(serverAddress))
	mcp.AddTool(srv, TransformTool(), TransformHandler(serverAddress))
}
