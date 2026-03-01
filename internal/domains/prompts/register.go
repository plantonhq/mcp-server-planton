package prompts

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Register adds all cross-domain MCP prompts to the server.
func Register(srv *mcp.Server) {
	srv.AddPrompt(DebugDeploymentPrompt(), DebugDeploymentHandler())
	srv.AddPrompt(AssessImpactPrompt(), AssessImpactHandler())
	srv.AddPrompt(ExploreInfrastructurePrompt(), ExploreInfrastructureHandler())
	srv.AddPrompt(ProvisionResourcePrompt(), ProvisionResourceHandler())
	srv.AddPrompt(ManageAccessPrompt(), ManageAccessHandler())
}
