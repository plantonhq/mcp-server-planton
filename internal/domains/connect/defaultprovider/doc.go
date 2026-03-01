// Package defaultprovider provides MCP tools for managing DefaultProviderConnection
// resources, which bind a cloud provider credential as the default for an
// organization or environment.
//
// Four tools are exposed:
//   - apply_default_provider_connection:    create or update
//   - get_default_provider_connection:      retrieve by ID
//   - resolve_default_provider_connection:  resolve the effective default for an org/provider/env
//   - delete_default_provider_connection:   delete by ID
package defaultprovider
