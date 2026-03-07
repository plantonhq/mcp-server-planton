package search

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/rpc"
)

// Register adds all cross-domain search tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	// api-resource search (3 tools)
	mcp.AddTool(srv, SearchByTextTool(), SearchByTextHandler(serverAddress))
	mcp.AddTool(srv, SearchByKindTool(), SearchByKindHandler(serverAddress))
	mcp.AddTool(srv, GetByOrgKindNameTool(), GetByOrgKindNameHandler(serverAddress))

	// connect search (3 tools)
	mcp.AddTool(srv, SearchConnectionsTool(), SearchConnectionsHandler(serverAddress))
	mcp.AddTool(srv, GetConnectionsByEnvTool(), GetConnectionsByEnvHandler(serverAddress))
	mcp.AddTool(srv, SearchRunnerRegistrationsTool(), SearchRunnerRegistrationsHandler(serverAddress))

	// infrahub search (2 tools — search_infra_projects lives in infraproject package)
	mcp.AddTool(srv, SearchIacModulesTool(), SearchIacModulesHandler(serverAddress))
	mcp.AddTool(srv, LookupCloudResourceTool(), LookupCloudResourceHandler(serverAddress))

	// resource-manager search (2 tools)
	mcp.AddTool(srv, GetContextHierarchyTool(), GetContextHierarchyHandler(serverAddress))
	mcp.AddTool(srv, SearchQuickActionsTool(), SearchQuickActionsHandler(serverAddress))

	// service-hub search (1 tool)
	mcp.AddTool(srv, SearchInfraChartsTool(), SearchInfraChartsHandler(serverAddress))
}

const (
	defaultPageNum  int32 = 1
	defaultPageSize int32 = 20
)

// buildPageInfo converts 1-based page_num and page_size into a proto PageInfo
// with 0-based Num, applying sensible defaults.
func buildPageInfo(pageNum, pageSize int32) *rpc.PageInfo {
	if pageNum <= 0 {
		pageNum = defaultPageNum
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	return &rpc.PageInfo{
		Num:  pageNum - 1,
		Size: pageSize,
	}
}
