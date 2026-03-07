// Package connection provides MCP tools and resources for managing Planton
// provider connections across cloud providers and third-party services.
//
// A connection object follows the OpenMCF envelope format:
//
//	{
//	  "api_version": "connect.planton.ai/v1",
//	  "kind": "AwsProviderConnection",      // PascalCase type
//	  "metadata": { "name": "…", "org": "…" },
//	  "spec": { … }                         // provider-specific fields
//	}
//
// Five tools are exposed:
//   - apply_connection:      create or update a connection (kind-dispatched)
//   - get_connection:        retrieve a connection by ID or org+slug (kind-dispatched)
//   - delete_connection:     delete a connection by ID (kind-dispatched)
//   - search_connections:    search connections in an organization via ConnectSearchQueryController
//   - check_connection_slug: check slug availability for a connection kind
//
// Two MCP resources are registered:
//   - connection-types://catalog         — JSON catalog of all supported connection types
//   - connection-schema://{kind}         — per-type JSON schema for spec fields
package connection
