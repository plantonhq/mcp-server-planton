# Design Decision 001: Remove Credential Field Redaction

**Date**: 2026-03-07
**Status**: Accepted
**Decision Maker**: User (explicit choice)

## Context

The old `connect/credential` package included `redact.go` — a module that masked sensitive fields (API keys, secrets, tokens) in the JSON response returned by `get_credential`. Each credential type in `registry.go` declared a `sensitiveFields []string` list, and the `redactFields` function replaced matching values with `"*****"`.

## The Change

The new proto contracts fundamentally changed how secrets are handled:

- **Old model**: Credential specs contained plaintext secret values (e.g., `access_key_secret: "ABCD1234"`).
- **New model**: Credential specs use `ConnectionFieldSecretRef` wrappers that store a **secret slug** (e.g., `secret_slug: "my-aws-access-key"`) instead of the actual secret value. The real secret is resolved at runtime by the platform.

## Decision

**Remove redaction entirely.** Secret slugs are identifiers, not sensitive data. Redacting them would hide useful information from MCP tool consumers without any security benefit.

## Alternatives Considered

1. **Keep redaction with updated field paths** — Would require mapping new nested `ConnectionFieldSecretRef` paths. Rejected because the slugs themselves are not sensitive.
2. **Redact only `ConnectionFieldSecretRef` fields** — Same issue: slugs are references, not secrets.
3. **Remove redaction entirely** — Selected. Clean, simple, and correct for the new data model.

## Impact

- `redact.go` deleted entirely
- `registry.go` no longer declares `sensitiveFields` per connection type
- `get.go` no longer calls `redactFields` after fetching a connection
- MCP tool consumers now see the full connection object including secret slug references
