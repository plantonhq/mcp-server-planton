package role

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all IAM role tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, ListForResourceKindTool(), ListForResourceKindHandler(serverAddress))
}
