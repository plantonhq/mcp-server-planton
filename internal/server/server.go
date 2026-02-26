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

// registerResources wires up the URI-addressable resource templates. These
// let MCP clients discover per-kind cloud resource schemas before calling
// the apply tool.
func registerResources(srv *mcp.Server) {
	srv.AddResourceTemplate(cloudresource.SchemaTemplate(), cloudresource.SchemaHandler())

	slog.Info("resource templates registered", "count", 1,
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
