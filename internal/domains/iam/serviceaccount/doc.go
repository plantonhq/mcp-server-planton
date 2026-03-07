// Package serviceaccount provides MCP tools for managing ServiceAccount
// resources, which represent machine identities for programmatic API access.
//
// Eight tools are exposed:
//   - create_service_account:      create a new service account
//   - get_service_account:         retrieve by ID
//   - update_service_account:      update display name or description
//   - delete_service_account:      delete and cascade (revokes keys, removes tuples)
//   - list_service_accounts:       list all service accounts in an organization
//   - create_service_account_key:  generate a new API key (sensitive — shown once)
//   - revoke_service_account_key:  revoke a specific API key
//   - list_service_account_keys:   list all API keys for a service account
package serviceaccount
