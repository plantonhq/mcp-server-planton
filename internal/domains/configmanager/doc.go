// Package configmanager is the domain root for configuration lifecycle tools,
// backed by the ConfigManager RPCs on the Planton backend.
//
// Three sub-packages expose 11 MCP tools:
//
//   - variable/       (5 tools) — plaintext configuration values scoped to org or env
//   - secret/         (4 tools) — encrypted secret metadata scoped to org or env
//   - secretversion/  (2 tools) — immutable encrypted key-value payloads
//
// Security boundary: agents can write secret values (create_secret_version)
// but cannot read decrypted data back. This follows the same principle as
// AD-01 (credential exclusion) — reading secrets is a human trust boundary.
package configmanager
