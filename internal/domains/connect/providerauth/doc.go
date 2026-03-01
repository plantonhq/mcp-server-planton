// Package providerauth provides MCP tools for managing provider connection
// authorizations, which control which credentials can be used in which
// environments.
//
// Three tools are exposed:
//   - apply_provider_connection_authorization:    create or update via OpenMCF envelope
//   - get_provider_connection_authorization:      dual-resolution by ID or semantic key
//   - delete_provider_connection_authorization:   remove an authorization by ID
package providerauth
