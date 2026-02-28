package pipeline

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all pipeline tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, GetLastTool(), GetLastHandler(serverAddress))
	mcp.AddTool(srv, RunTool(), RunHandler(serverAddress))
	mcp.AddTool(srv, RerunTool(), RerunHandler(serverAddress))
	mcp.AddTool(srv, CancelTool(), CancelHandler(serverAddress))
	mcp.AddTool(srv, ResolveGateTool(), ResolveGateHandler(serverAddress))
	mcp.AddTool(srv, ListFilesTool(), ListFilesHandler(serverAddress))
	mcp.AddTool(srv, UpdateFileTool(), UpdateFileHandler(serverAddress))
}
