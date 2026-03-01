package stackjob

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all stackjob tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, GetLatestTool(), GetLatestHandler(serverAddress))
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
	mcp.AddTool(srv, RerunTool(), RerunHandler(serverAddress))
	mcp.AddTool(srv, CancelTool(), CancelHandler(serverAddress))
	mcp.AddTool(srv, ResumeTool(), ResumeHandler(serverAddress))
	mcp.AddTool(srv, CheckEssentialsTool(), CheckEssentialsHandler(serverAddress))
	mcp.AddTool(srv, FindIacResourcesByStackJobTool(), FindIacResourcesByStackJobHandler(serverAddress))
	mcp.AddTool(srv, FindIacResourcesByApiResourceTool(), FindIacResourcesByApiResourceHandler(serverAddress))
	mcp.AddTool(srv, GetStackInputTool(), GetStackInputHandler(serverAddress))
	mcp.AddTool(srv, FindServiceStackJobsByEnvTool(), FindServiceStackJobsByEnvHandler(serverAddress))
	mcp.AddTool(srv, GetErrorRecommendationTool(), GetErrorRecommendationHandler(serverAddress))
}
