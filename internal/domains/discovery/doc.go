// Package discovery provides platform-wide MCP resources that help agents
// navigate the Planton API surface. Unlike domain-specific resources
// (e.g. cloud-resource-kinds, credential-types), these resources span all
// bounded contexts and serve as a top-level index of the platform.
//
// One MCP resource is exposed:
//   - api-resource-kinds://catalog: all platform API resource types grouped by domain
package discovery
