package tektontask

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all Tekton task tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
}
