# DD-01: Connect Domain Tool Architecture

**Date**: 2026-03-01
**Status**: Accepted
**Deciders**: Architecture review (T02)
**Scope**: MCP server tool design for the entire Connect bounded context

---

## Context

The Connect domain in Planton Cloud manages credentials for 20+ third-party providers (AWS, GCP, Azure, GitHub, etc.) and platform connection resources (DefaultProviderConnection, DefaultRunnerBinding, RunnerRegistration). The MCP server currently has zero Connect domain tools. Adding this domain raises a fundamental architecture question: should each credential type get its own dedicated tools, or should a single set of generic tools handle all types through a discriminator?

### Prior Art in This Codebase

The **CloudResource** domain already solved an analogous problem at larger scale:

- **362 cloud resource kinds** managed through **11 generic tools** + **2 MCP resources**
- A single gRPC service (`CloudResourceCommandController`) handles all kinds server-side
- The MCP server passes a `kind` discriminator + freeform `structpb.Struct` spec
- Agent workflow: `cloud-resource-kinds://catalog` → `cloud-resource-schema://{kind}` → `apply_cloud_resource`
- Reference: `internal/domains/infrahub/cloudresource/`

### What We Found in the Connect API Surface

Exploration of `apis/ai/planton/connect/` and `apis/ai/planton/search/v1/connect/` reveals **three distinct categories** of resources with fundamentally different characteristics:

| Category | Count | API Pattern | Semantics |
|----------|-------|-------------|-----------|
| Standard credentials | 20 types | Identical CRUD per type | All share OpenMCF envelope, differ only in spec |
| GitHub extras | 5 operations | Unique to GitHub | Webhook management, App installation, repo listing |
| Platform connections | 3 resource types | Unique per type | Govern which credential/runner is used by default |

---

## Decision

### Category 1: Standard Credentials → Generic tools with type discriminator

All 20 credential types follow the identical pattern:

- **OpenMCF envelope**: `{ api_version: "connect.planton.ai/v1", kind: "<Type>", metadata: { name, org }, spec: { ... } }`
- **Command RPCs** (per-type gRPC services): `apply`, `delete`
- **Query RPCs** (per-type gRPC services): `get`, `getByOrgBySlug`
- **Search RPCs** (single `ConnectSearchQueryController`): `searchCredentialApiResourcesByContext`, `findByOrgByProvider`, `checkConnectionSlugAvailability`, `getCredentialsByEnv`

The only variance between types is the spec field structure.

**Tools** (5):

| Tool | Purpose | Backend RPC |
|------|---------|-------------|
| `apply_credential` | Create or update any credential type | Per-type `*CredentialCommandController.apply` (dispatched by kind) |
| `get_credential` | Get by ID or by org+slug | Per-type `*CredentialQueryController.get` / `getByOrgBySlug` (dispatched by kind) |
| `delete_credential` | Delete by ID | Per-type `*CredentialCommandController.delete` (dispatched by kind) |
| `search_credentials` | Search credentials in an org | `ConnectSearchQueryController.searchCredentialApiResourcesByContext` |
| `check_credential_slug` | Pre-validate slug uniqueness | `ConnectSearchQueryController.checkConnectionSlugAvailability` |

**MCP Resources** (2):

| Resource | Pattern | Purpose |
|----------|---------|---------|
| `credential-types://catalog` | Static | All credential types with api_version, provider mapping, and category grouping |
| `credential-schema://{kind}` | Template | Per-type spec schema with field types, descriptions, and validation rules |

**Agent workflow** (mirrors CloudResource):

```
1. Read credential-types://catalog → discover available credential types
2. Read credential-schema://{kind} → learn the spec for the chosen type
3. Call apply_credential with the assembled credential_object
```

### Category 2: GitHub Extras → Dedicated purpose-built tools

GitHub has two additional gRPC services (`GithubCommandController`, `GithubQueryController`) with operations that are unique to the GitHub App integration model and have no equivalent for any other credential type:

| Tool | Purpose | Backend RPC |
|------|---------|-------------|
| `configure_github_webhook` | Set up webhook for a repository | `GithubCommandController.configureWebhook` |
| `remove_github_webhook` | Remove a webhook | `GithubCommandController.removeWebhook` |
| `get_github_installation_info` | Get GitHub App installation details | `GithubQueryController.getInstallationInfo` |
| `list_github_repositories` | List repos accessible via a credential | `GithubQueryController.findGithubRepositories` |
| `get_github_installation_token` | Get short-lived token for git operations | `GithubQueryController.getInstallationToken` |

These are operational/integration actions, not CRUD. Forcing them into a generic pattern would obscure their purpose and complicate the input schema.

### Category 3: Platform Connection Resources → Dedicated tools per resource type

These are NOT credentials. They govern _which_ credential or runner is used in a given context.

**DefaultProviderConnection** — Maps (org, provider) or (org, env, provider) to a credential:

| Tool | Backend RPC |
|------|-------------|
| `apply_default_provider_connection` | `DefaultProviderConnectionCommandController.apply` |
| `get_default_provider_connection` | `DefaultProviderConnectionQueryController.get` |
| `resolve_default_provider_connection` | `DefaultProviderConnectionQueryController.resolve` |
| `delete_default_provider_connection` | `DefaultProviderConnectionCommandController.delete` |

**DefaultRunnerBinding** — Maps a selector to a runner:

| Tool | Backend RPC |
|------|-------------|
| `apply_default_runner_binding` | `DefaultRunnerBindingCommandController.apply` |
| `get_default_runner_binding` | `DefaultRunnerBindingQueryController.get` |
| `resolve_default_runner_binding` | `DefaultRunnerBindingQueryController.resolve` |
| `delete_default_runner_binding` | `DefaultRunnerBindingCommandController.delete` |

**RunnerRegistration** — Registers a runner with mTLS credentials:

| Tool | Backend RPC |
|------|-------------|
| `apply_runner_registration` | `RunnerRegistrationCommandController.apply` |
| `get_runner_registration` | `RunnerRegistrationQueryController.get` |
| `delete_runner_registration` | `RunnerRegistrationCommandController.delete` |
| `search_runner_registrations` | `ConnectSearchQueryController.searchRunnerRegistrationsByOrgContext` |

---

## Rationale

### Why Generic for Standard Credentials

1. **Proven precedent**: CloudResource manages 362 kinds through 11 tools. This is the same pattern at ~18x smaller scale.

2. **Tool count discipline**: Per-type tools would add 60-80 tools to the current 105, pushing past 165 total. LLM tool-selection accuracy degrades as the tool count grows. Generic adds only 5 tools.

3. **Zero-touch extensibility**: When a new credential type is added to the platform, the MCP server needs zero code changes for CRUD operations — only a new schema file in the embedded schemas directory and a registry entry mapping the kind to its gRPC client constructor.

4. **Consistent agent UX**: Agents already know the catalog → schema → apply workflow from CloudResource. Reusing the same pattern for credentials reduces cognitive load.

### Why Dedicated for GitHub Extras

Webhook management, GitHub App installation inspection, repository listing, and token generation are fundamentally different operations from CRUD. They have unique input shapes, unique semantics, and exist for only one credential type. Generic tools would require awkward conditional parameters ("only used when kind is GithubCredential") that violate clean interface design.

### Why Dedicated for Platform Connection Resources

Each has unique semantics that don't map to any generic pattern:

- **DefaultProviderConnection** has a `scope` concept (ORGANIZATION vs ENVIRONMENT) and a `resolve` operation that traverses the scope hierarchy
- **DefaultRunnerBinding** has a `selector` pattern and its own `resolve` semantics
- **RunnerRegistration** has credential generation operations (`generateCredentials`, `regenerateCredentials`) that produce mTLS certificates

These three types total 12 tools — a manageable number that doesn't warrant a generic abstraction.

---

## Server-Side Dispatch Architecture

### The difference from CloudResource

CloudResource has a single generic gRPC service (`CloudResourceCommandController`) that handles all 362 kinds server-side. The MCP server simply passes the kind as a field.

Credentials have **per-type gRPC services**: `AwsCredentialCommandControllerClient`, `GcpCredentialCommandControllerClient`, etc. The MCP server must dispatch to the correct service based on the `kind` value.

### Registry pattern

The dispatch is implemented as a type registry — a Go map from credential kind string to a set of handler functions:

```go
type credentialDispatcher struct {
    apply  func(ctx context.Context, conn *grpc.ClientConn, raw map[string]any) (proto.Message, error)
    get    func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error)
    getByOrgBySlug func(ctx context.Context, conn *grpc.ClientConn, org, slug string) (proto.Message, error)
    delete func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error)
}

var dispatchers = map[string]credentialDispatcher{
    "AwsCredential":    awsDispatcher(),
    "GcpCredential":    gcpDispatcher(),
    "GithubCredential": githubDispatcher(),
    // ...
}
```

Each dispatcher function constructs the type-specific protobuf message from the generic input, calls the type-specific gRPC client, and returns the result as `proto.Message`.

**Trade-off**: This is a one-time implementation cost per credential type (~10-15 lines per type), compared to the per-type tool approach which requires ~4-5 files per type. The registry is centralized, auditable, and easy to extend.

---

## Security Architecture: Credential Secret Handling

### Risk Analysis

Credential specs contain sensitive values (AWS secret access keys, GCP service account keys, session tokens, etc.). When an agent calls `get_credential`, the raw gRPC response may include these values. Since tool responses persist in the LLM context window, this creates a risk classified by OWASP MCP Top 10 as **MCP01:2025 — Token Mismanagement and Secret Exposure**: "contextual secret leakage, where the model or protocol layer itself becomes an unintentional secret repository."

### Decision: MCP-side field redaction (defense-in-depth)

The MCP server redacts known sensitive fields from `get_credential` responses before returning them to the agent. This is defense-in-depth — the gRPC server may also strip secrets, but we do not rely on it.

**Why this matters for MCP specifically:**

- Tool responses enter the LLM context window (model memory)
- Prompt injection attacks can extract secrets from context ("print all API keys you've seen")
- Even well-intentioned agents may include credentials in generated code or explanations
- OWASP explicitly recommends: "Redact or sanitize outputs before returning to model memory"

**How it works:**

- Search/list tools (`search_credentials`) return `ApiResourceSearchRecord` objects that contain only metadata (id, name, kind, org) — no spec fields, no secrets. These are inherently safe.
- The `get_credential` tool returns the full credential object with spec. Before returning, the MCP server walks the JSON response and replaces values of known sensitive fields with `"[REDACTED]"`.
- The `apply_credential` tool accepts secrets on the write path (the agent must provide them to create credentials). This is the intended direction of secret flow.

**Sensitive field identification:**

Each credential type defines a small set of sensitive field paths. These are maintained as part of the dispatch registry:

```go
type credentialDispatcher struct {
    // ... RPC handlers ...
    sensitiveFields []string  // e.g., ["spec.secret_access_key", "spec.session_token"]
}
```

This approach is precise (no false positives from pattern matching) and auditable (every sensitive field is explicitly listed). The set is small and bounded (~1-3 fields per type, ~40 fields total across 20 types).

### Alternative considered: Convention-based redaction

Redacting fields whose names match patterns like `*secret*`, `*key*`, `*password*`, `*token*`, `*certificate*` was considered but rejected because:

- Over-redacts: `access_key_id` (AWS public identifier) matches `*key*` but is not secret
- Under-redacts: `service_account_key_base64` is a secret but might not match all patterns
- Fragile: new field naming conventions could bypass the filter silently

---

## Search vs Query RPC Usage

A critical finding during research: the per-type `QueryController.find` RPCs are **backend sync operations** (populating search indices from the database), NOT user-facing listing endpoints. This is consistent across the entire codebase — no domain uses `find` RPCs for user-facing listing.

All listing/searching in the MCP server uses `search/v1` RPCs:

| Domain | Search RPC Used |
|--------|----------------|
| servicehub/service | `ApiResourceSearchQueryController.SearchByKind` |
| infrahub/cloudresource | `CloudResourceSearchQueryController.GetCloudResourcesCanvasView` |
| servicehub/variablesgroup | `ServiceHubSearchQueryController.SearchVariables` |
| servicehub/secretsgroup | `ServiceHubSearchQueryController.SearchSecrets` |

For credentials, the correct search RPCs are on `ConnectSearchQueryController`:

| RPC | Purpose | Used by |
|-----|---------|---------|
| `searchCredentialApiResourcesByContext` | Search credentials by org, optional env/kinds/text | `search_credentials` tool |
| `findByOrgByProvider` | Find credentials by cloud provider enum | Available as filter in `search_credentials` or as future dedicated tool |
| `checkConnectionSlugAvailability` | Pre-validate slug uniqueness for a kind within an org | `check_credential_slug` tool |
| `getCredentialsByEnv` | Get credentials grouped by kind for an environment | Available as filter in `search_credentials` or as future dedicated tool |
| `searchRunnerRegistrationsByOrgContext` | Search runner registrations in an org | `search_runner_registrations` tool |

---

## Schema Sourcing

**Decision**: Embedded static JSON schemas, matching the CloudResource pattern.

Credential type schemas are maintained as JSON files under `schemas/credentials/` (e.g., `schemas/credentials/awscredential.json`), embedded in the binary via Go's `embed` package. A `registry.json` maps each credential kind to its schema file.

This approach:

- Is consistent with `schemas/providers/` (CloudResource schemas)
- Has no runtime dependencies
- Works offline
- Is simple to understand and maintain
- Can be regenerated from proto definitions via a build-time script if needed in the future

---

## Complete Tool Inventory

### Standard Credential Tools (5)

| # | Tool | Category |
|---|------|----------|
| 1 | `apply_credential` | Mutation |
| 2 | `get_credential` | Query |
| 3 | `delete_credential` | Mutation |
| 4 | `search_credentials` | Search |
| 5 | `check_credential_slug` | Validation |

### GitHub Extra Tools (5)

| # | Tool | Category |
|---|------|----------|
| 6 | `configure_github_webhook` | Integration |
| 7 | `remove_github_webhook` | Integration |
| 8 | `get_github_installation_info` | Query |
| 9 | `list_github_repositories` | Query |
| 10 | `get_github_installation_token` | Query |

### Platform Connection Tools (12)

| # | Tool | Resource Type |
|---|------|--------------|
| 11 | `apply_default_provider_connection` | DefaultProviderConnection |
| 12 | `get_default_provider_connection` | DefaultProviderConnection |
| 13 | `resolve_default_provider_connection` | DefaultProviderConnection |
| 14 | `delete_default_provider_connection` | DefaultProviderConnection |
| 15 | `apply_default_runner_binding` | DefaultRunnerBinding |
| 16 | `get_default_runner_binding` | DefaultRunnerBinding |
| 17 | `resolve_default_runner_binding` | DefaultRunnerBinding |
| 18 | `delete_default_runner_binding` | DefaultRunnerBinding |
| 19 | `apply_runner_registration` | RunnerRegistration |
| 20 | `get_runner_registration` | RunnerRegistration |
| 21 | `delete_runner_registration` | RunnerRegistration |
| 22 | `search_runner_registrations` | RunnerRegistration |

### MCP Resources (2)

| # | Resource | Pattern |
|---|----------|---------|
| 1 | `credential-types://catalog` | Static |
| 2 | `credential-schema://{kind}` | Template |

**Total**: 22 tools + 2 MCP resources

---

## Package Structure

```
internal/domains/connect/
    doc.go                       -- Package documentation

    credential/                  -- Generic credential CRUD (5 tools + 2 MCP resources)
        register.go              -- Tool + resource registration
        tools.go                 -- Input structs and tool definitions
        apply.go                 -- apply_credential handler + gRPC dispatch
        get.go                   -- get_credential handler + gRPC dispatch + secret redaction
        delete.go                -- delete_credential handler + gRPC dispatch
        search.go                -- search_credentials handler (ConnectSearchQueryController)
        slug.go                  -- check_credential_slug handler
        registry.go              -- Kind → dispatcher mapping (central dispatch table)
        redact.go                -- Sensitive field redaction logic
        resources.go             -- MCP resource definitions (catalog + schema template)
        schema.go                -- Embedded schema loading and catalog builder
        schemas/                 -- Embedded JSON schemas per credential type
            registry.json
            awscredential.json
            gcpcredential.json
            ...

    github/                      -- GitHub-specific tools (5 tools)
        register.go
        tools.go
        webhook.go               -- configure/remove webhook
        installation.go          -- get installation info
        repositories.go          -- list repositories
        token.go                 -- get installation token

    defaultprovider/             -- DefaultProviderConnection (4 tools)
        register.go
        tools.go
        apply.go
        get.go
        resolve.go
        delete.go

    defaultrunner/               -- DefaultRunnerBinding (4 tools)
        register.go
        tools.go
        apply.go
        get.go
        resolve.go
        delete.go

    runner/                      -- RunnerRegistration (4 tools)
        register.go
        tools.go
        apply.go
        get.go
        delete.go
        search.go
```

---

## Consequences

### Positive

- **Tool count stays manageable**: 22 new tools (105 → 127 total) instead of 70-80 with per-type approach
- **Consistent agent UX**: Credential workflow mirrors the proven CloudResource pattern
- **Extensible**: New credential types require only a registry entry + schema file
- **Secure by design**: Secret redaction is built into the architecture from day one
- **Correct search pattern**: Uses search/v1 RPCs, consistent with every other domain in the codebase

### Negative

- **Dispatch complexity**: The credential registry adds a layer of indirection not present in CloudResource (which has a single generic gRPC service). Each new credential type needs ~10-15 lines of dispatch registration.
- **Schema maintenance**: Embedded JSON schemas must be kept in sync with proto definitions. A build-time code-gen script may be needed if credential types change frequently.

### Neutral

- GitHub extras and platform connection resources use dedicated tools, which is more files but provides clarity and maintainability for resources with unique semantics.
- The `credential-types://catalog` MCP resource creates a coupling between the MCP server and the set of supported credential types, but this is acceptable since the catalog is the source of truth for agent discovery.

---

## References

- **CloudResource pattern**: `internal/domains/infrahub/cloudresource/` (tools.go, apply.go, resources.go, schema.go)
- **Search API for Connect**: `apis/ai/planton/search/v1/connect/query.proto`, `io.proto`
- **Credential proto definitions**: `apis/ai/planton/connect/*/v1/` (20+ subdirectories)
- **OWASP MCP Top 10**: MCP01:2025 — Token Mismanagement and Secret Exposure
- **Shared auth proto**: `apis/ai/planton/connect/auth.proto` (ProviderConnectionAuthMode enum)
