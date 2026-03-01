package github

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all GitHub extras tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ConfigureWebhookTool(), ConfigureWebhookHandler(serverAddress))
	mcp.AddTool(srv, RemoveWebhookTool(), RemoveWebhookHandler(serverAddress))
	mcp.AddTool(srv, GetInstallationInfoTool(), GetInstallationInfoHandler(serverAddress))
	mcp.AddTool(srv, ListRepositoriesTool(), ListRepositoriesHandler(serverAddress))
	mcp.AddTool(srv, GetInstallationTokenTool(), GetInstallationTokenHandler(serverAddress))
}
