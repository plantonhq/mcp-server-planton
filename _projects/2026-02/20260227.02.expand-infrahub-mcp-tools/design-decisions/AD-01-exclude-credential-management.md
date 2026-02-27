# AD-01: Exclude Credential Management from MCP Server

**Date**: 2026-02-27
**Status**: Accepted
**Decision**: Do not expose credential CRUD (Connect domain) through the MCP server

## Context

The initial architectural analysis identified the Connect domain (AWS, GCP, Azure credentials, runner registrations) as a potential Tier 2 addition to the MCP server. This was reconsidered based on security and design concerns.

## Decision

Credential lifecycle management (create, update, delete, read-with-secrets) is excluded from the MCP server tool surface.

## Rationale

1. **Security boundary**: Exposing credential CRUD through MCP means an AI agent could read access keys, secret keys, and session tokens. Even with API-side redaction, the `apply_credential` path would require the agent to *write* secrets into the system. Credentials are a human-trust-boundary operation.

2. **Not architecturally needed**: When a cloud resource is applied, it already references its credential connection by slug in the resource spec. The platform resolves the correct credential at deployment time. The AI agent doesn't need to manage credentials — it just references them.

3. **Failure mode is handled elsewhere**: If a credential is missing, the `check_stack_job_essentials` preflight (Phase 3B) will surface the issue. The user can then create the credential through the web UI or CLI — both of which have proper authentication and audit controls.

## Future Consideration

A read-only `list_credentials` tool (returning only slugs/names, never secret values) may be added later if agents frequently need to help users pick the right credential slug when composing resource specs. This would be a minimal, safe addition that doesn't cross the security boundary.
