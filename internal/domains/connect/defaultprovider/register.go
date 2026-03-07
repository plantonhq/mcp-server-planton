package defaultprovider

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all default provider connection tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, GetOrgDefaultTool(), GetOrgDefaultHandler(serverAddress))
	mcp.AddTool(srv, GetEnvDefaultTool(), GetEnvDefaultHandler(serverAddress))
	mcp.AddTool(srv, ResolveTool(), ResolveHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, DeleteOrgDefaultTool(), DeleteOrgDefaultHandler(serverAddress))
	mcp.AddTool(srv, DeleteEnvDefaultTool(), DeleteEnvDefaultHandler(serverAddress))
}
