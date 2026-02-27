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
	"github.com/plantonhq/mcp-server-planton/internal/domains/audit"
	"github.com/plantonhq/mcp-server-planton/internal/domains/configmanager/secret"
	"github.com/plantonhq/mcp-server-planton/internal/domains/configmanager/secretversion"
	"github.com/plantonhq/mcp-server-planton/internal/domains/configmanager/variable"
	"github.com/plantonhq/mcp-server-planton/internal/domains/graph"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/cloudresource"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/infrachart"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/infrapipeline"
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
	mcp.AddTool(srv, stackjob.RerunTool(), stackjob.RerunHandler(serverAddress))
	mcp.AddTool(srv, stackjob.CancelTool(), stackjob.CancelHandler(serverAddress))
	mcp.AddTool(srv, stackjob.ResumeTool(), stackjob.ResumeHandler(serverAddress))
	mcp.AddTool(srv, stackjob.CheckEssentialsTool(), stackjob.CheckEssentialsHandler(serverAddress))

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

	mcp.AddTool(srv, infrapipeline.ListTool(), infrapipeline.ListHandler(serverAddress))
	mcp.AddTool(srv, infrapipeline.GetTool(), infrapipeline.GetHandler(serverAddress))
	mcp.AddTool(srv, infrapipeline.GetLatestTool(), infrapipeline.GetLatestHandler(serverAddress))
	mcp.AddTool(srv, infrapipeline.RunTool(), infrapipeline.RunHandler(serverAddress))
	mcp.AddTool(srv, infrapipeline.CancelTool(), infrapipeline.CancelHandler(serverAddress))
	mcp.AddTool(srv, infrapipeline.ResolveEnvGateTool(), infrapipeline.ResolveEnvGateHandler(serverAddress))
	mcp.AddTool(srv, infrapipeline.ResolveNodeGateTool(), infrapipeline.ResolveNodeGateHandler(serverAddress))

	mcp.AddTool(srv, cloudresource.ListLocksTool(), cloudresource.ListLocksHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.RemoveLocksTool(), cloudresource.RemoveLocksHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.RenameTool(), cloudresource.RenameHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.GetEnvVarMapTool(), cloudresource.GetEnvVarMapHandler(serverAddress))
	mcp.AddTool(srv, cloudresource.ResolveValueReferencesTool(), cloudresource.ResolveValueReferencesHandler(serverAddress))

	mcp.AddTool(srv, graph.GetOrganizationGraphTool(), graph.GetOrganizationGraphHandler(serverAddress))
	mcp.AddTool(srv, graph.GetEnvironmentGraphTool(), graph.GetEnvironmentGraphHandler(serverAddress))
	mcp.AddTool(srv, graph.GetServiceGraphTool(), graph.GetServiceGraphHandler(serverAddress))
	mcp.AddTool(srv, graph.GetCloudResourceGraphTool(), graph.GetCloudResourceGraphHandler(serverAddress))
	mcp.AddTool(srv, graph.GetDependenciesTool(), graph.GetDependenciesHandler(serverAddress))
	mcp.AddTool(srv, graph.GetDependentsTool(), graph.GetDependentsHandler(serverAddress))
	mcp.AddTool(srv, graph.GetImpactAnalysisTool(), graph.GetImpactAnalysisHandler(serverAddress))

	mcp.AddTool(srv, variable.ListTool(), variable.ListHandler(serverAddress))
	mcp.AddTool(srv, variable.GetTool(), variable.GetHandler(serverAddress))
	mcp.AddTool(srv, variable.ApplyTool(), variable.ApplyHandler(serverAddress))
	mcp.AddTool(srv, variable.DeleteTool(), variable.DeleteHandler(serverAddress))
	mcp.AddTool(srv, variable.ResolveTool(), variable.ResolveHandler(serverAddress))

	mcp.AddTool(srv, secret.ListTool(), secret.ListHandler(serverAddress))
	mcp.AddTool(srv, secret.GetTool(), secret.GetHandler(serverAddress))
	mcp.AddTool(srv, secret.ApplyTool(), secret.ApplyHandler(serverAddress))
	mcp.AddTool(srv, secret.DeleteTool(), secret.DeleteHandler(serverAddress))

	mcp.AddTool(srv, secretversion.CreateTool(), secretversion.CreateHandler(serverAddress))
	mcp.AddTool(srv, secretversion.ListTool(), secretversion.ListHandler(serverAddress))

	mcp.AddTool(srv, audit.ListTool(), audit.ListHandler(serverAddress))
	mcp.AddTool(srv, audit.GetTool(), audit.GetHandler(serverAddress))
	mcp.AddTool(srv, audit.CountTool(), audit.CountHandler(serverAddress))

	slog.Info("tools registered", "count", 59, "tools", []string{
		"apply_cloud_resource",
		"get_cloud_resource",
		"delete_cloud_resource",
		"list_cloud_resources",
		"destroy_cloud_resource",
		"check_slug_availability",
		"get_stack_job",
		"get_latest_stack_job",
		"list_stack_jobs",
		"rerun_stack_job",
		"cancel_stack_job",
		"resume_stack_job",
		"check_stack_job_essentials",
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
		"list_infra_pipelines",
		"get_infra_pipeline",
		"get_latest_infra_pipeline",
		"run_infra_pipeline",
		"cancel_infra_pipeline",
		"resolve_infra_pipeline_env_gate",
		"resolve_infra_pipeline_node_gate",
		"list_cloud_resource_locks",
		"remove_cloud_resource_locks",
		"rename_cloud_resource",
		"get_env_var_map",
		"resolve_value_references",
		"get_organization_graph",
		"get_environment_graph",
		"get_service_graph",
		"get_cloud_resource_graph",
		"get_dependencies",
		"get_dependents",
		"get_impact_analysis",
		"list_variables",
		"get_variable",
		"apply_variable",
		"delete_variable",
		"resolve_variable",
		"list_secrets",
		"get_secret",
		"apply_secret",
		"delete_secret",
		"create_secret_version",
		"list_secret_versions",
		"list_resource_versions",
		"get_resource_version",
		"get_resource_version_count",
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
