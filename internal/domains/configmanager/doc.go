// Package configmanager is the domain root for configuration lifecycle tools,
// backed by the ConfigManager RPCs on the Planton backend.
//
// Six sub-packages expose 23 MCP tools:
//
//   - variable/       (5 tools) — plaintext configuration values scoped to org or env
//   - variablegroup/  (8 tools) — grouped configuration variables with entry-level operations
//   - secret/         (4 tools) — encrypted secret metadata scoped to org or env
//   - secretbackend/  (4 tools) — storage backend configuration for secrets (apply, get, list, delete)
//   - secretversion/  (2 tools) — immutable encrypted key-value payloads
//
// Security boundary: agents can write secret values (create_secret_version)
// but cannot read decrypted data back. Secret backend credentials (tokens,
// access keys) are redacted before returning to the agent.
package configmanager
