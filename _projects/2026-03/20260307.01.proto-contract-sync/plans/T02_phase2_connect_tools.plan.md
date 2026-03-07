---
name: Phase 2 Connect Tools
overview: Enrich existing connect domain MCP tools by wiring unwired gRPC methods across 4 sub-packages, fixing a confirmed bug in the existing Resolve handler, and making a design decision on provider-specific OAuth/CloudFormation controllers.
todos:
  - id: fix-resolve-bug
    content: "Fix confirmed bug: defaultprovider ResolveHandler not passing Provider field to gRPC call"
    status: completed
  - id: t02-1-defaultprovider
    content: "T02.1: Add 4 new tools to defaultprovider (get_org_default, get_env_default, delete_org_default, delete_env_default)"
    status: completed
  - id: t02-2-runner
    content: "T02.2: Add 2 new tools to runner (generate_runner_credentials, regenerate_runner_credentials) with security warnings"
    status: completed
  - id: t02-4-providerauth
    content: "T02.4: Add sync_provider_connection_authorization tool + enhance existing delete to support semantic key"
    status: completed
  - id: t02-3-defaultrunner
    content: "T02.3: Add 2 new tools to defaultrunner (get/delete by selector) + ApiResourceKind resolver"
    status: completed
  - id: t02-5-decision
    content: "T02.5: Get user decision on provider-specific controllers (AWS CloudFormation, GCP/Azure OAuth initiation)"
    status: pending
  - id: build-verify
    content: "Build verification: go build ./... and go vet ./..."
    status: completed
  - id: update-docs
    content: Update doc.go files and project checkpoint
    status: completed
isProject: false
---

# Phase 2: Enrich Existing Connect Tools

## Confirmed Bug (Fix First)

The existing `resolve_default_provider_connection` handler in [defaultprovider/tools.go](internal/domains/connect/defaultprovider/tools.go) **does not pass the `Provider` field** to the gRPC call, even though:

- The proto `ResolveDefaultProviderConnectionRequest` has a required `Provider` field (buf-validated, non-zero)
- The MCP input struct marks `provider` as required and collects it from the user

The handler currently sends `Provider=0` (UNSPECIFIED) to the server, which will either fail validation or return wrong results. Fix this before adding new tools.

```go
// Current (broken):
resp, err := client.Resolve(ctx, &defaultproviderconnectionv1.ResolveDefaultProviderConnectionRequest{
    Org:         input.Org,
    Environment: input.Environment,
})

// Fixed:
providerEnum, resolveErr := domains.ResolveProvider(input.Provider)
if resolveErr != nil {
    return nil, nil, resolveErr
}
resp, err := client.Resolve(ctx, &defaultproviderconnectionv1.ResolveDefaultProviderConnectionRequest{
    Org:         input.Org,
    Provider:    providerEnum,
    Environment: input.Environment,
})
```

---

## T02.1 -- defaultprovider: 4 New Tools

**Package:** [internal/domains/connect/defaultprovider/](internal/domains/connect/defaultprovider/)

All new request types use `Org` + `Provider` (+ optional `Env`), matching the existing `ResolveInput` pattern. Each uses `domains.ResolveProvider()` to map the user's string to the proto enum.


| New Tool                                 | gRPC Method                | Request Type                                  | Semantics                                            |
| ---------------------------------------- | -------------------------- | --------------------------------------------- | ---------------------------------------------------- |
| `get_org_default_provider_connection`    | Query.`GetOrgDefault`      | `GetOrgDefaultRequest{Org, Provider}`         | Get specifically the org-level default (no fallback) |
| `get_env_default_provider_connection`    | Query.`GetEnvDefault`      | `GetEnvDefaultRequest{Org, Provider, Env}`    | Get specifically the env-level default (no fallback) |
| `delete_org_default_provider_connection` | Command.`DeleteOrgDefault` | `DeleteOrgDefaultRequest{Org, Provider}`      | Delete org-level default by org+provider             |
| `delete_env_default_provider_connection` | Command.`DeleteEnvDefault` | `DeleteEnvDefaultRequest{Org, Provider, Env}` | Delete env-level default by org+provider+env         |


**Why separate tools instead of overloading existing Get/Delete:** The existing `get_default_provider_connection` takes an ID. The existing `resolve_default_provider_connection` applies fallback logic (env -> org). The new tools are **explicit level lookups** -- "show me exactly what's set at the org level" vs "show me exactly what's set at the env level" -- distinct semantics that deserve distinct tools.

**Files to change:**

- `defaultprovider/tools.go` -- add 4 tool+handler pairs, fix Resolve bug
- `defaultprovider/register.go` -- register 4 new tools
- `defaultprovider/doc.go` -- update tool listing (4 -> 8 tools + 1 bug fix)

---

## T02.2 -- runner: 2 New Tools

**Package:** [internal/domains/connect/runner/](internal/domains/connect/runner/)


| New Tool                        | gRPC Method                     | Request Type    | Semantics                                      |
| ------------------------------- | ------------------------------- | --------------- | ---------------------------------------------- |
| `generate_runner_credentials`   | Command.`GenerateCredentials`   | `ApiResourceId` | Generate initial auth credentials for a runner |
| `regenerate_runner_credentials` | Command.`RegenerateCredentials` | `ApiResourceId` | Rotate/regenerate runner auth credentials      |


Both return `RunnerCredentials` containing: org, runner ID, channel identifier, API key, planton API endpoint, tunnel endpoint, CA certificate, agent certificate, and agent private key.

**Security concern:** These responses contain private keys and certificates. Tool descriptions must include prominent warnings that the output contains sensitive cryptographic material. The existing codebase pattern does not redact (Phase 1 decision: secret slugs are not sensitive), but private keys ARE sensitive. I will add explicit warnings in the tool description text but will NOT add redaction logic -- the caller (LLM agent) must handle the response appropriately.

**Files to change:**

- `runner/tools.go` -- add 2 tool+handler pairs
- `runner/register.go` -- register 2 new tools
- `runner/doc.go` -- update tool listing (4 -> 6 tools)

---

## T02.3 -- defaultrunner: 2 New Tools

**Package:** [internal/domains/connect/defaultrunner/](internal/domains/connect/defaultrunner/)


| New Tool                                    | gRPC Method                | Request Type                    | Semantics                                     |
| ------------------------------------------- | -------------------------- | ------------------------------- | --------------------------------------------- |
| `get_default_runner_binding_by_selector`    | Query.`GetBySelector`      | `ApiResourceSelector{Kind, Id}` | Get binding by resource kind + resource ID    |
| `delete_default_runner_binding_by_selector` | Command.`DeleteBySelector` | `ApiResourceSelector{Kind, Id}` | Delete binding by resource kind + resource ID |


`ApiResourceSelector` takes an `ApiResourceKind` enum + an ID string. This is an indirect lookup pattern -- useful when you have a reference to a binding through another resource's selector field and want to resolve or delete it without knowing the binding ID.

**Note:** These tools require mapping user-supplied kind strings to `ApiResourceKind` enum values. We will need a new enum resolver in `internal/domains/` (similar to the existing `ResolveProvider` pattern) or accept the raw enum string.

**Files to change:**

- `defaultrunner/tools.go` -- add 2 tool+handler pairs
- `defaultrunner/register.go` -- register 2 new tools
- `defaultrunner/doc.go` -- update listing (4 -> 6 tools)
- Potentially `internal/domains/kind.go` or new file -- add `ApiResourceKind` resolver

---

## T02.4 -- providerauth: 1 New Tool + 1 Enhancement

**Package:** [internal/domains/connect/providerauth/](internal/domains/connect/providerauth/)

**New tool:**


| New Tool                                 | gRPC Method    | Request Type                                                            | Semantics                                      |
| ---------------------------------------- | -------------- | ----------------------------------------------------------------------- | ---------------------------------------------- |
| `sync_provider_connection_authorization` | Command.`Sync` | `ProviderConnectionAuthorizationSyncRequest{Org, Provider, Connection}` | Reconcile authorization state for a connection |


**Enhancement to existing tool:**

The existing `delete_provider_connection_authorization` only accepts an ID. The proto surface also exposes `DeleteBySemanticKey(DeleteBySemanticKeyRequest{Org, Provider, Connection})`. Rather than adding a new tool with a 50-character name, **enhance the existing delete tool** to accept either an ID or a semantic key (org + provider + connection) -- mirroring how `get_provider_connection_authorization` already works.

**Files to change:**

- `providerauth/tools.go` -- add Sync tool, enhance Delete to support semantic key
- `providerauth/register.go` -- register 1 new tool (Sync)
- `providerauth/doc.go` -- update listing

---

## T02.5 -- Provider-Specific Controllers: Design Decision

**My recommendation (challenge to the original plan):**

The original plan lists T02.5 as "AWS CloudFormation, GCP OAuth, Azure OAuth". After analyzing the proto contracts, I believe we should **split this into two categories**:

**Wire as MCP tools (3 tools):**

- `initiate_aws_cloudformation_setup` -- returns a pre-filled CloudFormation URL for the user to visit
- `get_aws_cloudformation_setup_status` -- polls setup completion (useful: agent can check status)
- `initiate_gcp_oauth` -- returns a Google authorization URL for the user to visit

**Probably wire (1 tool):**

- `initiate_azure_oauth` -- returns a Microsoft authorization URL

**Do NOT wire as MCP tools (2 methods):**

- `HandleOAuthCallback` (GCP) -- browser redirect endpoint; receives auth code + CSRF state from Google's redirect. An LLM agent will never be the OAuth callback target.
- `HandleOAuthCallback` (Azure) -- same reasoning.

The OAuth callbacks are server-side redirect handlers, not operations an agent would invoke. Wiring them would create misleading tools that encourage misuse.

**Package structure for the tools we do wire:**

- `connect/awssetup/` (doc.go, register.go, tools.go) -- 2 tools
- `connect/gcpoauth/` (doc.go, register.go, tools.go) -- 1 tool
- `connect/azureoauth/` (doc.go, register.go, tools.go) -- 1 tool
- Wire in [server.go](internal/server/server.go) -- 3 new import+Register calls

Small packages, but each has its own controller, types, and lifecycle -- clean bounded contexts.

---

## What We Explicitly Skip

- `**Find` methods** across all packages -- proto docs say `FindApiResourcesRequest` is "restricted to platform operators only" and "not available to end customers (they use the search API instead)". Wiring `Find` as MCP tools would violate the intended access model.
- `**Create`/`Update` methods** -- already covered by `Apply` (idempotent create-or-update).
- `**GetByOrgBySlug`** on providerauth -- the existing Get tool already supports semantic key lookup.
- `**GetBySelectorBySlug**` on runner -- niche internal reference pattern, low MCP value for now.

---

## Execution Order

1. Fix Resolve bug (T02.1 prerequisite, touches same file)
2. T02.1 -- defaultprovider (4 new tools, biggest batch)
3. T02.2 -- runner (2 new tools, security-sensitive)
4. T02.4 -- providerauth (1 new tool + 1 enhancement)
5. T02.3 -- defaultrunner (2 new tools, needs new enum resolver)
6. T02.5 -- provider-specific (pending your decision on scope)
7. Build verification: `go build ./...` and `go vet ./...`

---

## Total Impact

- **New tools:** 9 (+ up to 4 more if T02.5 is approved)
- **Enhanced tools:** 1 (providerauth delete)
- **Bug fixes:** 1 (Resolve Provider field)
- **New packages:** 0-3 (depending on T02.5 decision)
- **Files changed:** ~12-20
- **No changes to server.go** for T02.1-T02.4 (packages already registered)

