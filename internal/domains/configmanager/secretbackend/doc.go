// Package secretbackend provides MCP tools for managing SecretBackend
// resources, which define where and how encrypted secret data is stored.
//
// Four tools are exposed:
//   - apply_secret_backend:  create or update a backend configuration
//   - get_secret_backend:    retrieve by ID or by org+slug
//   - list_secret_backends:  list all backends in an organization
//   - delete_secret_backend: delete by ID or by org+slug
//
// Sensitive credential fields (tokens, access keys, client secrets) are
// redacted before returning responses to the agent. The backend also
// applies mask-on-write, so this is defense-in-depth.
package secretbackend
