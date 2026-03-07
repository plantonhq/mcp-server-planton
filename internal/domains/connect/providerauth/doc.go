// Package providerauth provides MCP tools for managing provider connection
// authorizations, which control which credentials can be used in which
// environments.
//
// Four tools are exposed:
//   - apply_provider_connection_authorization:    create or update via OpenMCF envelope
//   - get_provider_connection_authorization:      dual-resolution by ID or semantic key
//   - sync_provider_connection_authorization:     reconcile authorization state by semantic key
//   - delete_provider_connection_authorization:   dual-resolution delete by ID or semantic key
package providerauth
