# Connect Domain Tool Architecture Decision (DD-01)

**Date**: March 1, 2026

## Summary

Produced the foundational architecture decision for the Connect bounded context, determining how 20+ credential types, GitHub-specific operations, and platform connection resources will be exposed as MCP tools. The decision establishes a generic credential tool pattern (mirroring CloudResource), dedicated tools for GitHub extras and platform connections, MCP-side secret redaction for security, and the correct use of search/v1 RPCs for credential listing.

## Problem Statement

The Connect domain in Planton Cloud manages credentials for 20+ third-party providers and platform connection resources, but the MCP server had zero Connect domain tools. Before implementation could begin, a critical architecture question needed resolution: should each credential type get its own set of tools (60-80 tools), or should a single set of generic tools handle all types through a discriminator?

### Pain Points

- 20+ credential types with identical CRUD operations but different spec fields -- choosing the wrong pattern would either bloat the tool count or create an awkward abstraction
- GitHub has non-CRUD operations (webhook management, installation info, repo listing) that don't fit any generic pattern
- Platform connection resources (DefaultProviderConnection, DefaultRunnerBinding, RunnerRegistration) are NOT credentials but govern which credential/runner is used
- Credential specs contain secrets (AWS keys, GCP service account keys) that must not leak into the LLM context window
- Server-side API has per-type gRPC services (unlike CloudResource's single generic service), requiring a dispatch layer

## Solution

A three-category architecture with distinct tool patterns matched to each category's characteristics:

1. **Standard Credentials (20 types)**: Generic tools with type discriminator + MCP resources for discovery (5 tools + 2 MCP resources)
2. **GitHub Extras**: Dedicated purpose-built tools for non-CRUD operations (5 tools)
3. **Platform Connection Resources**: Dedicated tools per resource type with unique semantics (12 tools)

## Implementation Details

### Generic Credential Pattern

Mirrors the CloudResource approach: `credential-types://catalog` for discovery, `credential-schema://{kind}` for per-type spec schemas, and generic CRUD tools (`apply_credential`, `get_credential`, `delete_credential`, `search_credentials`, `check_credential_slug`).

Server-side dispatch is handled via a type registry mapping credential kind to gRPC client constructors, since credentials have per-type gRPC services unlike CloudResource's single generic service.

### Search RPC Discovery

Investigation revealed that the per-type `QueryController.find` RPCs are backend sync operations (populating search indices), not user-facing endpoints. The correct RPCs for listing are on `ConnectSearchQueryController`:

- `searchCredentialApiResourcesByContext` -- search credentials by org with optional kind/env/text filters
- `findByOrgByProvider` -- find credentials by cloud provider
- `checkConnectionSlugAvailability` -- slug pre-validation
- `searchRunnerRegistrationsByOrgContext` -- runner registration search

This is consistent with every other domain in the codebase (service, cloudresource, variablesgroup, etc.).

### Security Architecture

Based on OWASP MCP Top 10 (MCP01:2025 -- Token Mismanagement and Secret Exposure), the decision establishes MCP-side field redaction as defense-in-depth. Search responses return lightweight records without spec fields (inherently safe). `get_credential` responses have known sensitive fields redacted to `[REDACTED]` before entering the LLM context window. Each credential type declares its sensitive field paths in the dispatch registry.

### Schema Sourcing

Embedded static JSON schemas under `schemas/credentials/`, consistent with the CloudResource pattern (`schemas/providers/`). Simple, no runtime dependencies, and credential types change infrequently.

## Benefits

- **Tool count discipline**: 22 new tools instead of 60-80 with per-type approach (105 -> 127 total vs 105 -> 165-185)
- **Consistent agent UX**: Credential workflow mirrors the proven CloudResource catalog -> schema -> apply pattern
- **Zero-touch extensibility**: New credential types require only a registry entry + schema file, no new tool code
- **Secure by design**: Secret redaction built into the architecture from day one, based on OWASP guidance
- **Correct search pattern**: Uses search/v1 RPCs, consistent with every other domain

## Impact

This decision unblocks T05 (Connect Domain implementation, the largest single task in the gap completion project). It also informs the design of future MCP resources (T12) and establishes the secret redaction pattern that will apply to any future domain with sensitive data.

### Package structure established:

```
internal/domains/connect/
    credential/       -- Generic CRUD (5 tools + 2 MCP resources)
    github/           -- GitHub extras (5 tools)
    defaultprovider/  -- DefaultProviderConnection (4 tools)
    defaultrunner/    -- DefaultRunnerBinding (4 tools)
    runner/           -- RunnerRegistration (4 tools)
```

## Related Work

- **CloudResource pattern**: The architectural precedent for generic tools (`internal/domains/infrahub/cloudresource/`)
- **T01 Gap Analysis**: The master plan that identified the Connect domain as a critical gap
- **T05**: The implementation task this decision unblocks
- **OWASP MCP Top 10**: MCP01:2025 informed the security architecture

---

**Status**: Production Ready
**Timeline**: T02 design task (~2 hours)
