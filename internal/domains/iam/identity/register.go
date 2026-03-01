package identity

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all identity account tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, WhoAmITool(), WhoAmIHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, InviteTool(), InviteHandler(serverAddress))
	mcp.AddTool(srv, ListInvitationsTool(), ListInvitationsHandler(serverAddress))
}
