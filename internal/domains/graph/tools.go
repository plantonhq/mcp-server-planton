package graph

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	graphv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/graph/v1"
)

// ---------------------------------------------------------------------------
// get_organization_graph
// ---------------------------------------------------------------------------

// GetOrganizationGraphToolInput defines the parameters for the
// get_organization_graph tool.
type GetOrganizationGraphToolInput struct {
	Org                     string   `json:"org"                                  jsonschema:"required,Organization identifier. Use list_organizations to discover available organizations."`
	Envs                    []string `json:"envs,omitempty"                       jsonschema:"Environment slugs to restrict the graph to. When omitted all environments are included."`
	NodeTypes               []string `json:"node_types,omitempty"                 jsonschema:"Node types to include. Valid values: organization, environment, service, cloud_resource, credential, infra_project. When omitted all types are included."`
	IncludeTopologicalOrder bool     `json:"include_topological_order,omitempty"  jsonschema:"When true the response includes a topological ordering of node IDs (roots first). Useful for determining deployment order."`
	MaxDepth                int32    `json:"max_depth,omitempty"                  jsonschema:"Maximum depth for relationship traversal. 0 or omitted means unlimited."`
}

// GetOrganizationGraphTool returns the MCP tool definition for get_organization_graph.
func GetOrganizationGraphTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_organization_graph",
		Description: "Retrieve the complete resource topology graph for an organization. " +
			"Returns all nodes (organizations, environments, services, cloud resources, credentials, infra projects) " +
			"and their relationships (depends_on, uses_credential, deployed_as, etc.). " +
			"Use this as the starting point for understanding an organization's infrastructure landscape. " +
			"Optionally filter by environments or node types to focus the graph. " +
			"Use get_cloud_resource_graph or get_service_graph to drill into specific resources.",
	}
}

// GetOrganizationGraphHandler returns the typed tool handler for get_organization_graph.
func GetOrganizationGraphHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetOrganizationGraphToolInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetOrganizationGraphToolInput) (*mcp.CallToolResult, any, error) {
		if input.Org == "" {
			return nil, nil, fmt.Errorf("'org' is required")
		}
		nodeTypes, err := resolveNodeTypes(input.NodeTypes)
		if err != nil {
			return nil, nil, err
		}
		text, err := GetOrganizationGraph(ctx, serverAddress, OrganizationGraphInput{
			Org:                     input.Org,
			Envs:                    input.Envs,
			NodeTypes:               nodeTypes,
			IncludeTopologicalOrder: input.IncludeTopologicalOrder,
			MaxDepth:                input.MaxDepth,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_environment_graph
// ---------------------------------------------------------------------------

// GetEnvironmentGraphToolInput defines the parameters for the
// get_environment_graph tool.
type GetEnvironmentGraphToolInput struct {
	EnvID                   string   `json:"env_id"                               jsonschema:"required,Environment identifier to retrieve the graph for."`
	NodeTypes               []string `json:"node_types,omitempty"                 jsonschema:"Node types to include. Valid values: organization, environment, service, cloud_resource, credential, infra_project. When omitted all types are included."`
	IncludeTopologicalOrder bool     `json:"include_topological_order,omitempty"  jsonschema:"When true the response includes a topological ordering of node IDs (roots first)."`
}

// GetEnvironmentGraphTool returns the MCP tool definition for get_environment_graph.
func GetEnvironmentGraphTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_environment_graph",
		Description: "Retrieve the resource graph scoped to a specific environment. " +
			"Returns the environment node, its parent organization, and all resources " +
			"deployed in or belonging to the environment with their relationships. " +
			"Use this to understand what is deployed in an environment (e.g. staging, production). " +
			"Use list_environments to discover available environment identifiers.",
	}
}

// GetEnvironmentGraphHandler returns the typed tool handler for get_environment_graph.
func GetEnvironmentGraphHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetEnvironmentGraphToolInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetEnvironmentGraphToolInput) (*mcp.CallToolResult, any, error) {
		if input.EnvID == "" {
			return nil, nil, fmt.Errorf("'env_id' is required")
		}
		nodeTypes, err := resolveNodeTypes(input.NodeTypes)
		if err != nil {
			return nil, nil, err
		}
		text, err := GetEnvironmentGraph(ctx, serverAddress, EnvironmentGraphInput{
			EnvID:                   input.EnvID,
			NodeTypes:               nodeTypes,
			IncludeTopologicalOrder: input.IncludeTopologicalOrder,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_service_graph
// ---------------------------------------------------------------------------

// GetServiceGraphToolInput defines the parameters for the
// get_service_graph tool.
type GetServiceGraphToolInput struct {
	ServiceID         string   `json:"service_id"                  jsonschema:"required,Service identifier. Service IDs appear as node IDs in organization graph results."`
	Envs              []string `json:"envs,omitempty"              jsonschema:"Environment slugs to restrict results to. When omitted all environments are included."`
	IncludeUpstream   bool     `json:"include_upstream,omitempty"  jsonschema:"When true the graph includes upstream dependencies — resources the service depends on."`
	IncludeDownstream bool     `json:"include_downstream,omitempty" jsonschema:"When true the graph includes downstream dependents — resources that depend on the service."`
	MaxDepth          int32    `json:"max_depth,omitempty"         jsonschema:"Maximum depth for dependency traversal. 0 or omitted means unlimited."`
}

// GetServiceGraphTool returns the MCP tool definition for get_service_graph.
func GetServiceGraphTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_service_graph",
		Description: "Retrieve a service-centric subgraph showing the service and all related resources. " +
			"Returns the service node, its cloud resource deployments per environment, " +
			"and optionally upstream dependencies and downstream dependents. " +
			"Use this to understand where a service is deployed and what infrastructure it uses. " +
			"Service IDs can be discovered from get_organization_graph results.",
	}
}

// GetServiceGraphHandler returns the typed tool handler for get_service_graph.
func GetServiceGraphHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetServiceGraphToolInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetServiceGraphToolInput) (*mcp.CallToolResult, any, error) {
		if input.ServiceID == "" {
			return nil, nil, fmt.Errorf("'service_id' is required")
		}
		text, err := GetServiceGraph(ctx, serverAddress, ServiceGraphInput{
			ServiceID:         input.ServiceID,
			Envs:              input.Envs,
			IncludeUpstream:   input.IncludeUpstream,
			IncludeDownstream: input.IncludeDownstream,
			MaxDepth:          input.MaxDepth,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_cloud_resource_graph
// ---------------------------------------------------------------------------

// GetCloudResourceGraphToolInput defines the parameters for the
// get_cloud_resource_graph tool.
type GetCloudResourceGraphToolInput struct {
	CloudResourceID   string `json:"cloud_resource_id"            jsonschema:"required,Cloud resource ID to retrieve the graph for. Use get_cloud_resource to look up the ID if needed."`
	IncludeUpstream   bool   `json:"include_upstream,omitempty"   jsonschema:"When true the graph includes upstream dependencies — resources this cloud resource depends on."`
	IncludeDownstream bool   `json:"include_downstream,omitempty" jsonschema:"When true the graph includes downstream dependents — resources that depend on this cloud resource."`
	MaxDepth          int32  `json:"max_depth,omitempty"          jsonschema:"Maximum depth for dependency traversal. 0 or omitted means unlimited."`
}

// GetCloudResourceGraphTool returns the MCP tool definition for get_cloud_resource_graph.
func GetCloudResourceGraphTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_cloud_resource_graph",
		Description: "Retrieve a cloud-resource-centric subgraph. " +
			"Returns the cloud resource node at the center, services deployed as it, " +
			"credentials it uses, and all connected nodes and relationships. " +
			"Enable include_upstream and include_downstream to traverse dependencies " +
			"beyond the immediate neighbors. " +
			"Use this to understand the full context of a specific cloud resource.",
	}
}

// GetCloudResourceGraphHandler returns the typed tool handler for get_cloud_resource_graph.
func GetCloudResourceGraphHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetCloudResourceGraphToolInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetCloudResourceGraphToolInput) (*mcp.CallToolResult, any, error) {
		if input.CloudResourceID == "" {
			return nil, nil, fmt.Errorf("'cloud_resource_id' is required")
		}
		text, err := GetCloudResourceGraph(ctx, serverAddress, CloudResourceGraphInput{
			CloudResourceID:   input.CloudResourceID,
			IncludeUpstream:   input.IncludeUpstream,
			IncludeDownstream: input.IncludeDownstream,
			MaxDepth:          input.MaxDepth,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_dependencies
// ---------------------------------------------------------------------------

// GetDependenciesToolInput defines the parameters for the get_dependencies tool.
type GetDependenciesToolInput struct {
	ResourceID        string   `json:"resource_id"                   jsonschema:"required,Resource ID to find dependencies for. Can be any resource type — cloud resource, service, credential, etc."`
	MaxDepth          int32    `json:"max_depth,omitempty"           jsonschema:"Maximum depth for dependency traversal. 0 or omitted means unlimited."`
	RelationshipTypes []string `json:"relationship_types,omitempty"  jsonschema:"Relationship types to include. Valid values: belongs_to_org, belongs_to_env, deployed_as, uses_credential, depends_on, runs_on, managed_by, uses, service_depends_on, owned_by. When omitted all types are included."`
}

// GetDependenciesTool returns the MCP tool definition for get_dependencies.
func GetDependenciesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_dependencies",
		Description: "Find all resources that a given resource depends on (upstream traversal). " +
			"For example, an EKS cluster might depend on a VPC and an IAM credential. " +
			"Useful for understanding deployment prerequisites and resource ordering. " +
			"Optionally filter by relationship type (e.g. only depends_on or uses_credential). " +
			"Use get_dependents for the reverse direction (what depends on this resource).",
	}
}

// GetDependenciesHandler returns the typed tool handler for get_dependencies.
func GetDependenciesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetDependenciesToolInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetDependenciesToolInput) (*mcp.CallToolResult, any, error) {
		if input.ResourceID == "" {
			return nil, nil, fmt.Errorf("'resource_id' is required")
		}
		relTypes, err := resolveRelationshipTypes(input.RelationshipTypes)
		if err != nil {
			return nil, nil, err
		}
		text, err := GetDependencies(ctx, serverAddress, DependencyInput{
			ResourceID:        input.ResourceID,
			MaxDepth:          input.MaxDepth,
			RelationshipTypes: relTypes,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_dependents
// ---------------------------------------------------------------------------

// GetDependentsToolInput defines the parameters for the get_dependents tool.
type GetDependentsToolInput struct {
	ResourceID        string   `json:"resource_id"                   jsonschema:"required,Resource ID to find dependents for. Can be any resource type — cloud resource, service, credential, etc."`
	MaxDepth          int32    `json:"max_depth,omitempty"           jsonschema:"Maximum depth for dependency traversal. 0 or omitted means unlimited."`
	RelationshipTypes []string `json:"relationship_types,omitempty"  jsonschema:"Relationship types to include. Valid values: belongs_to_org, belongs_to_env, deployed_as, uses_credential, depends_on, runs_on, managed_by, uses, service_depends_on, owned_by. When omitted all types are included."`
}

// GetDependentsTool returns the MCP tool definition for get_dependents.
func GetDependentsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_dependents",
		Description: "Find all resources that depend on a given resource (downstream traversal). " +
			"For example, a VPC might have EKS clusters, RDS instances, and other resources depending on it. " +
			"Use this before deleting or modifying a resource to understand what might be affected. " +
			"For a comprehensive impact report with counts and breakdown, use get_impact_analysis instead. " +
			"Use get_dependencies for the reverse direction (what does this resource depend on).",
	}
}

// GetDependentsHandler returns the typed tool handler for get_dependents.
func GetDependentsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetDependentsToolInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetDependentsToolInput) (*mcp.CallToolResult, any, error) {
		if input.ResourceID == "" {
			return nil, nil, fmt.Errorf("'resource_id' is required")
		}
		relTypes, err := resolveRelationshipTypes(input.RelationshipTypes)
		if err != nil {
			return nil, nil, err
		}
		text, err := GetDependents(ctx, serverAddress, DependencyInput{
			ResourceID:        input.ResourceID,
			MaxDepth:          input.MaxDepth,
			RelationshipTypes: relTypes,
		})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_impact_analysis
// ---------------------------------------------------------------------------

// GetImpactAnalysisToolInput defines the parameters for the
// get_impact_analysis tool.
type GetImpactAnalysisToolInput struct {
	ResourceID string `json:"resource_id"            jsonschema:"required,Resource ID to analyze impact for. Can be any resource type — cloud resource, service, credential, etc."`
	ChangeType string `json:"change_type,omitempty"  jsonschema:"Type of change being analyzed. Valid values: delete, update. When omitted the server analyzes the general impact without change-type-specific logic."`
}

// GetImpactAnalysisTool returns the MCP tool definition for get_impact_analysis.
func GetImpactAnalysisTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_impact_analysis",
		Description: "Analyze the impact of modifying or deleting a resource. " +
			"Returns directly affected resources, transitively affected resources, " +
			"total affected count, and a breakdown of affected resources by type. " +
			"Use this before destructive operations (delete, destroy) to understand the blast radius. " +
			"Specify change_type as 'delete' or 'update' for change-specific analysis.",
	}
}

// GetImpactAnalysisHandler returns the typed tool handler for get_impact_analysis.
func GetImpactAnalysisHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetImpactAnalysisToolInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetImpactAnalysisToolInput) (*mcp.CallToolResult, any, error) {
		if input.ResourceID == "" {
			return nil, nil, fmt.Errorf("'resource_id' is required")
		}
		var changeType graphv1.GetImpactAnalysisInput_ChangeType
		if input.ChangeType != "" {
			var err error
			changeType, err = resolveChangeType(input.ChangeType)
			if err != nil {
				return nil, nil, err
			}
		}
		text, err := GetImpactAnalysis(ctx, serverAddress, input.ResourceID, changeType)
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
