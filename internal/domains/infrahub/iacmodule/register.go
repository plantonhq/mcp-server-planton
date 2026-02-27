package iacmodule

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all iacmodule tools to the MCP server.
func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, SearchTool(), SearchHandler(serverAddress))
	mcp.AddTool(srv, GetTool(), GetHandler(serverAddress))
}
