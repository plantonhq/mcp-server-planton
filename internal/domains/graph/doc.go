// Package graph provides MCP tools for dependency intelligence and impact
// analysis, backed by the GraphQueryController RPCs
// (ai.planton.graph.v1) on the Planton backend.
//
// Seven tools are exposed:
//   - get_organization_graph:    full resource topology for an organization
//   - get_environment_graph:     everything deployed in a specific environment
//   - get_service_graph:         service-centric subgraph with deployments
//   - get_cloud_resource_graph:  resource-centric dependency view
//   - get_dependencies:          upstream — what does a resource depend on?
//   - get_dependents:            downstream — what depends on a resource?
//   - get_impact_analysis:       if a resource changes or is deleted, what breaks?
package graph
