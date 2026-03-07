// Package defaultprovider provides MCP tools for managing DefaultProviderConnection
// resources, which bind a cloud provider credential as the default for an
// organization or environment.
//
// Eight tools are exposed:
//   - apply_default_provider_connection:          create or update
//   - get_default_provider_connection:            retrieve by ID
//   - get_org_default_provider_connection:        retrieve the org-level default by org+provider
//   - get_env_default_provider_connection:        retrieve the env-level default by org+provider+env
//   - resolve_default_provider_connection:        resolve effective default with env→org fallback
//   - delete_default_provider_connection:         delete by ID
//   - delete_org_default_provider_connection:     delete org-level default by org+provider
//   - delete_env_default_provider_connection:     delete env-level default by org+provider+env
package defaultprovider
