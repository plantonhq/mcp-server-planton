package policy

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all IAM policy v2 tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, CreateTool(), CreateHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, UpsertTool(), UpsertHandler(serverAddress))
	mcp.AddTool(srv, RevokeOrgAccessTool(), RevokeOrgAccessHandler(serverAddress))
	mcp.AddTool(srv, ListResourceAccessTool(), ListResourceAccessHandler(serverAddress))
	mcp.AddTool(srv, CheckAuthorizationTool(), CheckAuthorizationHandler(serverAddress))
	mcp.AddTool(srv, ListPrincipalsTool(), ListPrincipalsHandler(serverAddress))
}
