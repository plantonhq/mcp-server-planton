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
	"github.com/plantonhq/mcp-server-planton/internal/auth"
	"github.com/plantonhq/mcp-server-planton/internal/config"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/cloudresource"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/infrachart"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/infraproject"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/preset"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/stackjob"
	"github.com/plantonhq/mcp-server-planton/internal/domains/resourcemanager/environment"
	"github.com/plantonhq/mcp-server-planton/internal/domains/resourcemanager/organization"
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
	mcp.AddTool(srv, cloudresource.GetTool(), cloudresource.GetHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.DeleteTool(), cloudresource.DeleteHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.ListTool(), cloudresource.ListHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.DestroyTool(), cloudresource.DestroyHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.CheckSlugAvailabilityTool(), cloudresource.CheckSlugAvailabilityHandler(serverAddress))

	mcp.AddTool(srv, stackjob.GetTool(), stackjob.GetHandler(serverAddress))
	mcp.AddTool(srv, stackjob.GetLatestTool(), stackjob.GetLatestHandler(serverAddress))
	mcp.AddTool(srv, stackjob.ListTool(), stackjob.ListHandler(serverAddress))

	mcp.AddTool(srv, organization.ListTool(), organization.ListHandler(serverAddress))

	mcp.AddTool(srv, environment.ListTool(), environment.ListHandler(serverAddress))

	mcp.AddTool(srv, preset.SearchTool(), preset.SearchHandler(serverAddress))
	mcp.AddTool(srv, preset.GetTool(), preset.GetHandler(serverAddress))

	mcp.AddTool(srv, infrachart.ListTool(), infrachart.ListHandler(serverAddress))
	mcp.AddTool(srv, infrachart.GetTool(), infrachart.GetHandler(serverAddress))
	mcp.AddTool(srv, infrachart.BuildTool(), infrachart.BuildHandler(serverAddress))

	mcp.AddTool(srv, infraproject.SearchTool(), infraproject.SearchHandler(serverAddress))
	mcp.AddTool(srv, infraproject.GetTool(), infraproject.GetHandler(serverAddress))
	mcp.AddTool(srv, infraproject.ApplyTool(), infraproject.ApplyHandler(serverAddress))
	mcp.AddTool(srv, infraproject.DeleteTool(), infraproject.DeleteHandler(serverAddress))
	mcp.AddTool(srv, infraproject.CheckSlugTool(), infraproject.CheckSlugHandler(serverAddress))
	mcp.AddTool(srv, infraproject.UndeployTool(), infraproject.UndeployHandler(serverAddress))

	mcp.AddTool(srv, cloudresource.ListLocksTool(), cloudresource.ListLocksHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.RemoveLocksTool(), cloudresource.RemoveLocksHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.RenameTool(), cloudresource.RenameHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.GetEnvVarMapTool(), cloudresource.GetEnvVarMapHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.ResolveValueReferencesTool(), cloudresource.ResolveValueReferencesHandler(serverAddress))

	slog.Info("tools registered", "count", 27, "tools", []string{
		"apply_cloud_resource",
		"get_cloud_resource",
		"delete_cloud_resource",
		"list_cloud_resources",
		"destroy_cloud_resource",
		"check_slug_availability",
		"get_stack_job",
		"get_latest_stack_job",
		"list_stack_jobs",
		"list_organizations",
		"list_environments",
		"search_cloud_object_presets",
		"get_cloud_object_preset",
		"list_infra_charts",
		"get_infra_chart",
		"build_infra_chart",
		"search_infra_projects",
		"get_infra_project",
		"apply_infra_project",
		"delete_infra_project",
		"check_infra_project_slug",
		"undeploy_infra_project",
		"list_cloud_resource_locks",
		"remove_cloud_resource_locks",
		"rename_cloud_resource",
		"get_env_var_map",
		"resolve_value_references",
	})
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
//	go build -ldflags "-X github.com/plantonhq/mcp-server-planton/internal/server.buildVersion=v1.0.0"
var buildVersion string
