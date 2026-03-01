package discovery

import "github.com/modelcontextprotocol/go-sdk/mcp"

// RegisterResources adds the platform-wide discovery MCP resources.
func RegisterResources(srv *mcp.Server) {
	srv.AddResource(CatalogResource(), CatalogHandler())
}
