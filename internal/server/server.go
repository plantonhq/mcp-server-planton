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
	connectcredential "github.com/plantonhq/mcp-server-planton/internal/domains/connect/credential"
	"github.com/plantonhq/mcp-server-planton/internal/domains/discovery"
	connectdefaultprovider "github.com/plantonhq/mcp-server-planton/internal/domains/connect/defaultprovider"
	connectdefaultrunner "github.com/plantonhq/mcp-server-planton/internal/domains/connect/defaultrunner"
	connectgithub "github.com/plantonhq/mcp-server-planton/internal/domains/connect/github"
	connectproviderauth "github.com/plantonhq/mcp-server-planton/internal/domains/connect/providerauth"
	connectrunner "github.com/plantonhq/mcp-server-planton/internal/domains/connect/runner"
	"github.com/plantonhq/mcp-server-planton/internal/domains/configmanager/secret"
	"github.com/plantonhq/mcp-server-planton/internal/domains/configmanager/secretversion"
	"github.com/plantonhq/mcp-server-planton/internal/domains/configmanager/variable"
	"github.com/plantonhq/mcp-server-planton/internal/domains/graph"
	iamapikey "github.com/plantonhq/mcp-server-planton/internal/domains/iam/apikey"
	iamidentity "github.com/plantonhq/mcp-server-planton/internal/domains/iam/identity"
	iampolicy "github.com/plantonhq/mcp-server-planton/internal/domains/iam/policy"
	iamrole "github.com/plantonhq/mcp-server-planton/internal/domains/iam/role"
	iamteam "github.com/plantonhq/mcp-server-planton/internal/domains/iam/team"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/cloudresource"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/deploymentcomponent"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/flowcontrolpolicy"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/iacmodule"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/infrachart"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/infrapipeline"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/infraproject"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/preset"
	"github.com/plantonhq/mcp-server-planton/internal/domains/infrahub/stackjob"
	"github.com/plantonhq/mcp-server-planton/internal/domains/resourcemanager/environment"
	"github.com/plantonhq/mcp-server-planton/internal/domains/resourcemanager/organization"
	"github.com/plantonhq/mcp-server-planton/internal/domains/resourcemanager/promotionpolicy"
	servicehubdnsdomain "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/dnsdomain"
	servicehubpipeline "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/pipeline"
	servicehubsecretsgroup "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/secretsgroup"
	servicehubservice "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/service"
	servicehubtektonpipeline "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/tektonpipeline"
	servicehubtektontask "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/tektontask"
	servicehubvariablesgroup "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/variablesgroup"
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

// registerTools delegates tool registration to each domain package.
// Adding a new tool only requires changes in the owning domain package.
func registerTools(srv *mcp.Server, serverAddress string) {
	cloudresource.Register(srv, serverAddress)
	stackjob.Register(srv, serverAddress)
	organization.Register(srv, serverAddress)
	environment.Register(srv, serverAddress)
	preset.Register(srv, serverAddress)
	infrachart.Register(srv, serverAddress)
	infraproject.Register(srv, serverAddress)
	infrapipeline.Register(srv, serverAddress)
	graph.Register(srv, serverAddress)
	variable.Register(srv, serverAddress)
	secret.Register(srv, serverAddress)
	secretversion.Register(srv, serverAddress)
	audit.Register(srv, serverAddress)
	deploymentcomponent.Register(srv, serverAddress)
	iacmodule.Register(srv, serverAddress)
	servicehubservice.Register(srv, serverAddress)
	servicehubpipeline.Register(srv, serverAddress)
	servicehubvariablesgroup.Register(srv, serverAddress)
	servicehubsecretsgroup.Register(srv, serverAddress)
	servicehubdnsdomain.Register(srv, serverAddress)
	servicehubtektonpipeline.Register(srv, serverAddress)
	servicehubtektontask.Register(srv, serverAddress)
	connectcredential.Register(srv, serverAddress)
	connectgithub.Register(srv, serverAddress)
	connectdefaultprovider.Register(srv, serverAddress)
	connectdefaultrunner.Register(srv, serverAddress)
	connectrunner.Register(srv, serverAddress)
	connectproviderauth.Register(srv, serverAddress)
	iamidentity.Register(srv, serverAddress)
	iamteam.Register(srv, serverAddress)
	iampolicy.Register(srv, serverAddress)
	iamrole.Register(srv, serverAddress)
	iamapikey.Register(srv, serverAddress)
	promotionpolicy.Register(srv, serverAddress)
	flowcontrolpolicy.Register(srv, serverAddress)

	slog.Info("tools registered")
}

// registerResources delegates MCP resource registration to domain packages.
func registerResources(srv *mcp.Server) {
	cloudresource.RegisterResources(srv)
	connectcredential.RegisterResources(srv)
	discovery.RegisterResources(srv)

	slog.Info("resources registered")
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
