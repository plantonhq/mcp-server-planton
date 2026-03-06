// Package credential provides MCP tools for managing Planton credentials
// across all supported provider types (AWS, GCP, Azure, GitHub, Kubernetes, etc.).
//
// All credential types share the same OpenMCF envelope structure:
//
//	{ api_version, kind, metadata: { name, org }, spec: { ... } }
//
// and differ only in their spec fields. This package exposes a single set of
// generic tools that dispatch to the correct per-type gRPC service based on
// the credential kind.
//
// Five tools are exposed:
//   - apply_credential:       create or update any credential type
//   - get_credential:         retrieve by ID or by org+slug (with secret redaction)
//   - delete_credential:      delete by ID
//   - search_credentials:     search credentials within an organization
//   - check_credential_slug:  validate slug uniqueness before creation
//
// Two MCP resources are exposed:
//   - credential-types://catalog:       all supported credential types
//   - credential-schema://{kind}:       per-type spec schema
//
// Security: get_credential redacts known sensitive fields (secret keys, tokens,
// certificates) from responses before they enter the LLM context window. See
// redact.go and the sensitive field declarations in registry.go.
package credential
