// Package deploymentcomponent provides the MCP tools for the DeploymentComponent
// domain, backed by the InfraHubSearchQueryController
// (ai.planton.search.v1.infrahub) and DeploymentComponentQueryController
// (ai.planton.infrahub.deploymentcomponent.v1) RPCs on the Planton backend.
//
// Two tools are exposed:
//   - search_deployment_components: browse the cloud resource type catalog
//   - get_deployment_component:     retrieve full component details by ID or by kind
package deploymentcomponent

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// ---------------------------------------------------------------------------
// search_deployment_components
// ---------------------------------------------------------------------------

// SearchDeploymentComponentsInput defines the parameters for the
// search_deployment_components tool.
type SearchDeploymentComponentsInput struct {
	SearchText string `json:"search_text,omitempty" jsonschema:"Free-text search query to filter deployment components by name or description."`
	Provider   string `json:"provider,omitempty"    jsonschema:"Cloud provider to filter by (e.g. aws, gcp, azure, confluent, snowflake). When omitted, components from all providers are returned."`
	PageNum    int32  `json:"page_num,omitempty"    jsonschema:"Page number (1-based). Defaults to 1."`
	PageSize   int32  `json:"page_size,omitempty"   jsonschema:"Number of results per page. Defaults to 20."`
}

// SearchTool returns the MCP tool definition for search_deployment_components.
func SearchTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "search_deployment_components",
		Description: "Browse the deployment component catalog. " +
			"Deployment components represent the types of cloud resources that can be provisioned on the platform " +
			"(e.g. AwsEksCluster, GcpCloudRunService, ConfluentKafkaCluster). " +
			"Use the optional 'provider' filter to narrow results to a specific cloud provider. " +
			"Returns lightweight search records — use get_deployment_component with an ID or kind " +
			"from the results to retrieve the full component definition.",
	}
}

// SearchHandler returns the typed tool handler for search_deployment_components.
func SearchHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *SearchDeploymentComponentsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *SearchDeploymentComponentsInput) (*mcp.CallToolResult, any, error) {
		text, err := Search(ctx, serverAddress, SearchInput{
			SearchText: input.SearchText,
			Provider:   input.Provider,
			PageNum:    input.PageNum,
			PageSize:   input.PageSize,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_deployment_component
// ---------------------------------------------------------------------------

// GetDeploymentComponentInput defines the parameters for the
// get_deployment_component tool.
// Exactly one identification path must be provided:
//   - ID path: set 'id' alone.
//   - Kind path: set 'kind' alone (PascalCase, e.g. AwsEksCluster).
type GetDeploymentComponentInput struct {
	ID   string `json:"id,omitempty"   jsonschema:"The deployment component ID. Mutually exclusive with 'kind'."`
	Kind string `json:"kind,omitempty" jsonschema:"PascalCase cloud resource kind (e.g. AwsEksCluster). Read cloud-resource-kinds://catalog for valid kinds. Mutually exclusive with 'id'."`
}

// GetTool returns the MCP tool definition for get_deployment_component.
func GetTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_deployment_component",
		Description: "Retrieve the full details of a deployment component by its ID or by cloud resource kind. " +
			"A deployment component defines a type of cloud resource that can be provisioned, including its " +
			"supported IaC modules, provider, and configuration schema. " +
			"Use search_deployment_components to discover component IDs and kinds, or pass a known kind " +
			"string directly (e.g. AwsEksCluster). " +
			"To find which IaC modules can provision a component, use search_iac_modules with the 'kind' filter.",
	}
}

// GetHandler returns the typed tool handler for get_deployment_component.
func GetHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetDeploymentComponentInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetDeploymentComponentInput) (*mcp.CallToolResult, any, error) {
		if err := validateGetInput(input.ID, input.Kind); err != nil {
			return nil, nil, err
		}
		text, err := Get(ctx, serverAddress, input.ID, input.Kind)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// validateGetInput checks that exactly one identification path is provided:
// either 'id' alone or 'kind' alone.
func validateGetInput(id, kind string) error {
	hasID := id != ""
	hasKind := kind != ""

	switch {
	case hasID && hasKind:
		return fmt.Errorf("provide either 'id' or 'kind' — not both")
	case hasID || hasKind:
		return nil
	default:
		return fmt.Errorf("provide either 'id' or 'kind' to identify the deployment component")
	}
}
