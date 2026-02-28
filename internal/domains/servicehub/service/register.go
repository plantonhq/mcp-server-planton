package service

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all service tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, SearchTool(), SearchHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, DisconnectGitRepoTool(), DisconnectGitRepoHandler(serverAddress))
	mcp.AddTool(srv, ConfigureWebhookTool(), ConfigureWebhookHandler(serverAddress))
	mcp.AddTool(srv, ListBranchesTool(), ListBranchesHandler(serverAddress))
}
