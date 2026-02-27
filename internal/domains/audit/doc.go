// Package audit provides MCP tools for resource version history and change
// tracking, backed by the ApiResourceVersionQueryController RPCs
// (ai.planton.audit.apiresourceversion.v1) on the Planton backend.
//
// Three tools are exposed:
//   - list_resource_versions:     paginated change history for a specific resource
//   - get_resource_version:       full version with YAML states and unified diff
//   - get_resource_version_count: lightweight count of how many versions exist
//
// Unlike domain-specific packages (infrachart, variable, etc.) that hard-code
// their ApiResourceKind, audit accepts the kind from the caller because it is
// a cross-cutting concern that operates on any platform resource type.
package audit
