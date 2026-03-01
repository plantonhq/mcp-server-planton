---
name: T02 Credential Architecture Decision
overview: "Produce an Architecture Decision Record (ADR) for the Connect domain that covers three categories: standard credential CRUD, GitHub-specific extras, and platform connection resources (DefaultProviderConnection, DefaultRunnerBinding, RunnerRegistration)."
todos:
  - id: research-complete
    content: Deep exploration of codebase patterns, Connect proto APIs, and existing MCP resources (DONE)
    status: completed
  - id: discuss-open-points
    content: "Discuss the 3 open points with the user: schema sourcing, list_credentials scope, and security boundary"
    status: completed
  - id: write-adr
    content: Write the ADR document at design-decisions/DD01-connect-domain-tool-architecture.md covering all three categories
    status: completed
isProject: false
---

# T02: Connect Domain Tool Architecture Decision

## Domain Analysis (Architect Review)

### What the codebase tells us

The MCP server currently has **105 tools** across 22 domains and **2 MCP resources**. The most relevant precedent is **CloudResource**, which solved the exact same problem at a larger scale:

- **362 cloud resource kinds** managed through **11 generic tools** + **2 MCP resources**
- A single gRPC service (`CloudResourceCommandController`) handles all kinds server-side
- The MCP server passes a `kind` discriminator + freeform `structpb.Struct` spec
- Agent workflow: `cloud-resource-kinds://catalog` -> `cloud-resource-schema://{kind}` -> `apply_cloud_resource`

Key files:

- [internal/domains/infrahub/cloudresource/tools.go](internal/domains/infrahub/cloudresource/tools.go) -- generic input with `cloud_object map[string]any`
- [internal/domains/infrahub/cloudresource/apply.go](internal/domains/infrahub/cloudresource/apply.go) -- single gRPC service dispatch
- [internal/domains/infrahub/cloudresource/resources.go](internal/domains/infrahub/cloudresource/resources.go) -- MCP resource definitions
- [internal/domains/infrahub/cloudresource/schema.go](internal/domains/infrahub/cloudresource/schema.go) -- embedded schema catalog

### The Connect domain API surface

Exploration of `/Users/suresh/scm/github.com/plantonhq/planton/apis/ai/planton/connect/` reveals **three distinct categories** of resources, each with fundamentally different characteristics:

---

## Category 1: Standard Credentials (20 types)

**Types**: AwsCredential, GcpCredential, AzureCredential, KubernetesClusterCredential, DigitalOceanCredential, CivoCredential, CloudflareCredential, Auth0Credential, ConfluentCredential, MongodbAtlasCredential, OpenFgaCredential, SnowflakeCredential, GithubCredential, GitlabCredential, DockerCredential, MavenCredential, NpmCredential, TerraformBackendCredential, PulumiBackendCredential (+ CloudflareWorkerScriptSR2Bucket)

**API uniformity**: Every type has the exact same RPC surface:

- `CommandController`: `apply`, `create`, `update`, `delete`
- `QueryController`: `get`, `getByOrgBySlug`, `find`

**Variance**: Only in spec fields. All share the same OpenMCF envelope:

```
{ api_version: "connect.planton.ai/v1", kind: "<Type>", metadata: { name, org }, spec: { ... } }
```

**Recommendation: Generic tools (Option B)**

Rationale:

- Identical to the CloudResource problem, just smaller scale (20 vs 362 kinds)
- Per-type tools would add **60-80 tools** (105 -> 165-185 total) -- LLM tool selection degrades with more tools
- Generic adds **4-5 tools** (105 -> 109-110 total)
- Agent UX mirrors the proven CloudResource workflow
- Adding new credential types requires zero MCP server tool code changes

**Proposed tools:**

- `apply_credential` -- Create or update any credential type (accepts `credential_object` in OpenMCF format)
- `get_credential` -- Get by ID, or by org + kind + slug
- `delete_credential` -- Delete by ID
- `list_credentials` -- List/search credentials in an org, optionally filtered by kind

**Proposed MCP resources:**

- `credential-types://catalog` -- Static catalog of all credential types with api_version, grouped by category (cloud providers, SaaS, SCM, package registries, IaC backends)
- `credential-schema://{kind}` -- Per-type spec schema (same pattern as `cloud-resource-schema://{kind}`)

### Surprise: Server-side dispatch difference

Unlike CloudResource (single generic gRPC service), credentials have **per-type gRPC services** (`AwsCredentialCommandControllerClient`, `GcpCredentialCommandControllerClient`, etc.). The MCP server will need a **dispatch registry** that maps `kind` -> gRPC client constructor.

This is a well-understood pattern (type switch or map-based registry) and a one-time implementation cost. Example structure:

```go
type credentialHandler struct {
    apply  func(ctx context.Context, conn *grpc.ClientConn, obj *structpb.Struct) (proto.Message, error)
    get    func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error)
    delete func(ctx context.Context, conn *grpc.ClientConn, id string) (proto.Message, error)
    // ...
}
var registry = map[string]credentialHandler{ ... }
```

---

## Category 2: GitHub Extras

GitHub has **two additional gRPC services** beyond standard CRUD that no other credential type has:

- `GithubCommandController`: `configureWebhook`, `removeWebhook`
- `GithubQueryController`: `getInstallationInfo`, `findGithubRepositories`, `getInstallationToken`

**Recommendation: Dedicated purpose-built tools**

These are operational/integration actions, not CRUD. They don't fit any generic pattern and shouldn't be forced into one.

**Proposed tools:**

- `configure_github_webhook` -- Set up webhook for a GitHub repository
- `remove_github_webhook` -- Remove a webhook
- `get_github_installation_info` -- Get GitHub App installation details
- `list_github_repositories` -- List repos accessible via a GitHub credential
- `get_github_installation_token` -- Get a short-lived token for git operations

**Package**: `internal/domains/connect/github/` (separate from generic credential package)

---

## Category 3: Platform Connection Resources

Three resource types that are **not credentials** but **govern which credential/runner is used**:

**DefaultProviderConnection** -- Maps (org, provider) or (org, env, provider) to a credential slug:

- Command: `apply`, `delete`, `deleteOrgDefault`, `deleteEnvDefault`
- Query: `get`, `resolve`, `getOrgDefault`, `getEnvDefault`
- Has `scope` concept (ORGANIZATION vs ENVIRONMENT)

**DefaultRunnerBinding** -- Maps a selector to a runner:

- Command: `apply`, `delete`, `deleteBySelector`
- Query: `get`, `getBySelector`, `resolve`

**RunnerRegistration** -- Registers a runner with mTLS credentials:

- Command: `apply`, `delete`, `generateCredentials`, `regenerateCredentials`
- Query: `get`, `getBySelectorBySlug`

**Recommendation: Dedicated tools per resource type**

Each has unique semantics (resolve, scope, selector, credential generation) that don't map to any generic pattern. Low count (3 types, ~8-10 tools total) makes dedicated tools practical.

**Proposed tools:**

- DefaultProviderConnection: `apply_default_provider_connection`, `get_default_provider_connection`, `resolve_default_provider_connection`, `delete_default_provider_connection`
- DefaultRunnerBinding: `apply_default_runner_binding`, `get_default_runner_binding`, `resolve_default_runner_binding`, `delete_default_runner_binding`
- RunnerRegistration: `apply_runner_registration`, `get_runner_registration`, `delete_runner_registration`

**Package**: `internal/domains/connect/defaultprovider/`, `internal/domains/connect/defaultrunner/`, `internal/domains/connect/runner/`

---

## Summary: Tool and Resource Impact


| Category             | Approach           | New Tools | New MCP Resources |
| -------------------- | ------------------ | --------- | ----------------- |
| Standard Credentials | Generic            | 4         | 2                 |
| GitHub Extras        | Dedicated          | 5         | 0                 |
| Platform Connections | Dedicated per type | 10        | 0                 |
| **Total**            |                    | **19**    | **2**             |


vs. original T01 estimate of 25-30: lower because the generic approach eliminates ~40-60 per-type tools.

---

## Package Structure

```
internal/domains/connect/
  credential/          -- Generic credential CRUD (4 tools, 2 MCP resources)
    tools.go           -- Input types and tool definitions
    register.go        -- Tool + resource registration
    apply.go           -- Dispatch to per-type gRPC services
    get.go
    delete.go
    list.go
    resources.go       -- credential-types://catalog + credential-schema://{kind}
    schema.go          -- Embedded schemas and catalog builder
    registry.go        -- Kind -> gRPC handler dispatch table
    schemas/           -- Embedded JSON schemas per credential type
  github/              -- GitHub-specific non-CRUD tools (5 tools)
    tools.go
    register.go
    webhook.go
    installation.go
    repositories.go
  defaultprovider/     -- DefaultProviderConnection (4 tools)
    tools.go
    register.go
    apply.go
    get.go
    delete.go
  defaultrunner/       -- DefaultRunnerBinding (3 tools)
    tools.go
    register.go
    ...
  runner/              -- RunnerRegistration (3 tools)
    tools.go
    register.go
    ...
```

---

## Open Discussion Points

1. **Credential schema sourcing**: Should we embed static JSON schemas (like CloudResource does from `schemas/providers/`) or generate them from proto reflection? Embedded is simpler and consistent with the existing pattern. Proto reflection is more maintainable but adds a dependency.
2. `**list_credentials` scope**: Should `list_credentials` call the per-type `find` RPCs internally (dispatching N queries for N types and merging), or should there be a server-side cross-type search endpoint? If neither exists, we may need to note this as a backend requirement or limit `list_credentials` to require a `kind` filter.
3. **Security boundary**: Credential secrets are write-only. The `get_credential` tool response should strip secret fields (or the server already does this). Need to verify server behavior.

## Deliverable

A design decision document at:

```
_projects/2026-03/20260301.01.mcp-server-gap-completion/design-decisions/DD01-connect-domain-tool-architecture.md
```

Contents: Context, Decision, Rationale, Package Structure, Agent Workflow, Implementation Guidance for T05, and Consequences.