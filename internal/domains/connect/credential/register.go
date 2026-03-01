package credential

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all credential tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ApplyTool(), ApplyHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
	mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))
	mcp.AddTool(srv, SearchTool(), SearchHandler(serverAddress))
	mcp.AddTool(srv, CheckSlugTool(), CheckSlugHandler(serverAddress))
}

// RegisterResources adds the credential MCP resources and resource templates.
func RegisterResources(srv *mcp.Server) {
	srv.AddResource(CatalogResource(), CatalogHandler())
	srv.AddResourceTemplate(SchemaTemplate(), SchemaHandler())
}
