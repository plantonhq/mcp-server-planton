package cloudresource

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all cloudresource tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, ListTool(), ListHandler(serverAddress))
	mcp.AddTool(srv, DestroyTool(), DestroyHandler(serverAddress))
	mcp.AddTool(srv, PurgeTool(), PurgeHandler(serverAddress))
	mcp.AddTool(srv, CheckSlugAvailabilityTool(), CheckSlugAvailabilityHandler(serverAddress))
	mcp.AddTool(srv, ListLocksTool(), ListLocksHandler(serverAddress))
	mcp.AddTool(srv, RemoveLocksTool(), RemoveLocksHandler(serverAddress))
	mcp.AddTool(srv, RenameTool(), RenameHandler(serverAddress))
	mcp.AddTool(srv, GetEnvVarMapTool(), GetEnvVarMapHandler(serverAddress))
	mcp.AddTool(srv, ResolveValueReferencesTool(), ResolveValueReferencesHandler(serverAddress))
}

// RegisterResources adds the cloudresource MCP resources and resource templates.
func RegisterResources(srv *mcp.Server) {
	srv.AddResource(KindCatalogResource(), KindCatalogHandler())
	srv.AddResourceTemplate(SchemaTemplate(), SchemaHandler())
}
