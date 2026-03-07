package runner

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all runner registration tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, SearchTool(), SearchHandler(serverAddress))
	mcp.AddTool(srv, GenerateCredentialsTool(), GenerateCredentialsHandler(serverAddress))
	mcp.AddTool(srv, RegenerateCredentialsTool(), RegenerateCredentialsHandler(serverAddress))
}
