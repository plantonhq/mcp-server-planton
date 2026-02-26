// Package server initializes the MCP server, registers all tools, and exposes
// the transport entry points (STDIO and HTTP).
//
// The server is stateless — all per-request state (API key, gRPC connection)
// is derived from the context that the transport injects. This means the same
// mcp.Server instance can safely serve both STDIO and HTTP concurrently.
package server

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/plantoncloud/mcp-server-planton/internal/auth"
	"github.com/plantoncloud/mcp-server-planton/internal/config"
	"github.com/plantoncloud/mcp-server-planton/internal/domains/cloudresource"
)

// Server wraps an mcp.Server with Planton-specific configuration.
type Server struct {
	mcp    *mcp.Server
	config *config.Config
}

// New creates a configured MCP server with all Planton tools registered.
func New(cfg *config.Config) *Server {
	srv := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-server-planton",
			Version: version(),
		},
		nil,
	)

	registerTools(srv, cfg.ServerAddress)
	registerResources(srv)

	return &Server{
		mcp:    srv,
		config: cfg,
	}
}

// registerTools wires up every domain tool. The serverAddress is captured in
// each handler's closure so that tool handlers can create gRPC connections
// without reaching back into the config layer.
func registerTools(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, cloudresource.ApplyTool(), cloudresource.ApplyHandler(serverAddress))

	// Phase 4 will add:
	//   mcp.AddTool(srv, cloudresource.DeleteTool(), cloudresource.DeleteHandler(serverAddress))
	//   mcp.AddTool(srv, cloudresource.GetTool(), cloudresource.GetHandler(serverAddress))

	slog.Info("tools registered", "count", 1)
}

// registerResources wires up MCP resources and resource templates. The kind
// catalog lets agents discover available kinds; the schema template lets them
// fetch the full spec definition for a specific kind.
func registerResources(srv *mcp.Server) {
	srv.AddResource(cloudresource.KindCatalogResource(), cloudresource.KindCatalogHandler())
	srv.AddResourceTemplate(cloudresource.SchemaTemplate(), cloudresource.SchemaHandler())

	slog.Info("resources registered",
		"static", []string{"cloud-resource-kinds://catalog"},
		"templates", []string{"cloud-resource-schema://{kind}"},
	)
}

// ServeStdio runs the MCP server over stdin/stdout until the client
// disconnects or the context is cancelled.
//
// In STDIO mode the API key is loaded once from the environment at startup
// (validated during config loading) and injected into the base context.
// Every tool handler can then retrieve it via auth.APIKey(ctx).
func (s *Server) ServeStdio(ctx context.Context) error {
	ctx = auth.WithAPIKey(ctx, s.config.APIKey)
	return s.mcp.Run(ctx, &mcp.StdioTransport{})
}

// version returns the server version. This is set at build time via ldflags
// or falls back to "dev".
func version() string {
	if buildVersion != "" {
		return buildVersion
	}
	return "dev"
}

// buildVersion is populated at link time:
//
//	go build -ldflags "-X github.com/plantoncloud/mcp-server-planton/internal/server.buildVersion=v1.0.0"
var buildVersion string
