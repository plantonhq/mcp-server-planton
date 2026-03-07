# Task T01: Proto Contract Sync — Full Migration Plan

**Created**: 2026-03-07
**Status**: PENDING REVIEW
**Type**: Migration

> This plan requires your review before execution.

## Summary

The proto definitions have diverged significantly from the MCP server tool implementations. The build is broken due to the connect domain's credential-to-connection rename. Beyond the build fix, new resources and entire domains have appeared. This plan lays out a phased approach to get everything aligned.

## Current State

- **Build status**: BROKEN — 2 files (`credential/registry.go`, `github/tools.go`) import from the deleted `plantonhq/planton/apis/stubs` module
- **Proto inventory**: 764 `.pb.go` files across 18 domains (all in `gen/go/`, all untracked)
- **Existing MCP tools**: 10 domain packages, ~80+ tools total
- **Stable domains**: audit, graph, resourcemanager, servicehub, infrahub (existing sub-packages), iam (existing sub-packages), configmanager (existing sub-packages) — all compile fine

---

## Phase 1: Fix the Build (Critical Path)

**Goal**: `go build ./...` passes with zero errors.

### T01.1 — Migrate `connect/credential` package

The generic credential dispatcher handles 19 credential types through a registry pattern. Every type has been renamed:

| Old (`*credential`) | New (`*providerconnection` / `*connection`) |
|---|---|
| `AwsCredential` | `AwsProviderConnection` |
| `GcpCredential` | `GcpProviderConnection` |
| `AzureCredential` | `AzureProviderConnection` |
| `GithubCredential` | `GithubConnection` |
| `GitlabCredential` | `GitlabConnection` |
| `KubernetesClusterCredential` | `KubernetesProviderConnection` |
| `DockerCredential` | `ContainerRegistryConnection` |
| `MongodbAtlasCredential` | `AtlasProviderConnection` |
| `CloudflareCredential` | `CloudflareProviderConnection` |
| `CivoCredential` | `CivoProviderConnection` |
| `ConfluentCredential` | `ConfluentProviderConnection` |
| `DigitalOceanCredential` | `DigitalOceanProviderConnection` |
| `Auth0Credential` | `Auth0ProviderConnection` |
| `OpenFgaCredential` | `OpenFgaProviderConnection` |
| `SnowflakeCredential` | `SnowflakeProviderConnection` |
| `MavenCredential` | `MavenConnection` |
| `NpmCredential` | `NpmConnection` |
| `PulumiBackendCredential` | `PulumiBackendConnection` |
| `TerraformBackendCredential` | `TerraformBackendConnection` |

**Work items:**
- [ ] Update `registry.go` — change all 19 import paths from `planton/apis/stubs/.../` to `gen/go/.../`
- [ ] Update type names in the dispatcher map (credential -> connection types)
- [ ] Update `apply.go`, `get.go`, `delete.go` — adapt to new message types
- [ ] Update `search.go` and `slug.go` — verify `ConnectSearchQueryController` contract
- [ ] Review `redact.go` — new contracts use secret slug references instead of plaintext; redaction logic may be simplified or removed
- [ ] Update `schemas/credentials/` embedded JSON schemas to match new type names
- [ ] Add `CloudflareWorkerScriptsR2Connection` — brand new type with no old equivalent
- [ ] Update MCP tool names/descriptions if "credential" terminology should change to "connection"

### T01.2 — Migrate `connect/github` package

- [ ] Change import from `planton/apis/stubs/.../githubcredential/v1` to `gen/go/.../githubconnection/v1`
- [ ] Update message type references (`GithubCredential` -> `GithubConnection`, etc.)
- [ ] Verify webhook and query method signatures haven't changed

### T01.3 — Verify build

- [ ] Run `go build ./...` — must pass
- [ ] Run `go vet ./...` — must pass

---

## Phase 2: Enrich Existing Connect Tools (Quick Wins)

These packages already compile but have new gRPC methods not yet exposed as MCP tools.

### T02.1 — `connect/defaultprovider`

New methods: `GetOrgDefault`, `GetEnvDefault`, `DeleteOrgDefault`, `DeleteEnvDefault`

### T02.2 — `connect/runner`

New methods: `GenerateCredentials`, `RegenerateCredentials`

### T02.3 — `connect/defaultrunner`

New methods: `GetBySelector`, `DeleteBySelector`

### T02.4 — `connect/providerauth`

New methods: `Sync`, `DeleteBySemanticKey`, `Find`

### T02.5 — New provider-specific controllers

- AWS: `CloudFormationSetupController` — `InitiateCloudFormationSetup`, `GetCloudFormationSetupStatus`
- GCP: `OAuthController` — `InitiateOAuth`, `HandleOAuthCallback`
- Azure: `AzureOAuthController`

---

## Phase 3: New Resources in Existing Domains

### T03.1 — `configmanager/secretbackend`

New resource — manages secret backend configuration (platform, openbao, AWS Secrets Manager, HashiCorp Vault, GCP Secret Manager, Azure Key Vault).

Tools needed: `get_secret_backend`, `apply_secret_backend`, `delete_secret_backend`, `list_secret_backends`

### T03.2 — `configmanager/variablegroup`

New resource — groups variables with scope (organization, environment) and entry-level operations.

Tools needed: `get_variable_group`, `apply_variable_group`, `delete_variable_group`, `list_variable_groups`, `resolve_variable_group`, plus entry-level ops (`upsert_entry`, `delete_entry`, `refresh_entry`, `refresh_all`)

### T03.3 — `iam/serviceaccount`

New resource — programmatic identity with key management.

Tools needed: `create_service_account`, `get_service_account`, `update_service_account`, `delete_service_account`, `list_service_accounts`, `create_service_account_key`, `revoke_service_account_key`, `list_service_account_keys`

### T03.4 — `infrahub/iacprovisionermapping`

New resource — maps IAC provisioners.

Tools needed: `get_iac_provisioner_mapping`, `apply_iac_provisioner_mapping`, `delete_iac_provisioner_mapping`, `find_iac_provisioner_mappings`

---

## Phase 4: Evaluate New Domains

These 7 domains are entirely new in the protos with zero MCP tool coverage. Each needs evaluation for whether MCP tools add value.

| Domain | Resources | Initial Assessment |
|---|---|---|
| **agentfleet** | agent, agenttestsuite, execution, mcpserver, session, skill | Medium priority — meta-management of AI agents and MCP servers |
| **search** | quickaction, apiresource, connect, iam, infrahub, servicehub, resourcemanager | Medium priority — cross-domain search is highly useful for LLM agents |
| **integration** | gcp, git, kubernetes/cost, kubernetes/kubernetesobject, tekton, vcs | Medium priority — operational integrations |
| **runner** | cloudapis/provider (AWS/Azure/GCP/K8s) | Medium priority — runner-side cloud API queries |
| **billing** | billingaccount | Low priority — billing management via MCP seems niche |
| **copilot** | copilotagent, copilotchat | Low priority — copilot-specific, may overlap with agentfleet |
| **reporting** | iac/v1 | Low priority — read-only reporting queries |

---

## Phase 5: Enrich Existing Stable Domains (Optional)

Low priority — existing tools work, but new methods are available:

- `configmanager/variable` — `Refresh` method, `Source` field
- `configmanager/secretversion` — `Get`, `Delete`, `GetLatestBySecret`
- `iam/identity` — `GetByGithubUsername`, `GetActorInfo`
- `iam/apikey` — `Get`, `Update`
- `iam/policy` (v2) — `GrantOwnership`, `GrantMemberAccess`, `ListAuthorizedPrincipals`

---

## Execution Order

1. **Phase 1** (T01.1-T01.3) — Fix the build. This is the critical blocker.
2. **Phase 2** (T02.1-T02.5) — Quick enrichments to already-working connect tools.
3. **Phase 3** (T03.1-T03.4) — New resource implementations following established patterns.
4. **Phase 4** — Domain evaluation and selective implementation.
5. **Phase 5** — Optional enrichments as time permits.

## Design Decision Needed

**Should we rename the MCP tool surface from "credential" to "connection"?**

The backend has renamed everything, but existing MCP users may have muscle memory around `apply_credential`, `get_credential`, etc. Options:
- A) Rename tools to match protos (`apply_aws_provider_connection`, `get_credential` -> `get_connection`)
- B) Keep old tool names but rewire internals (simpler migration, but naming drift)
- C) Support both with aliases during a transition period

This decision affects Phase 1 scope.

---

## Success Criteria

1. `go build ./...` succeeds (Phase 1)
2. All 19+1 connection types callable through MCP tools (Phase 1)
3. New resource tools for secretbackend, variablegroup, serviceaccount, iacprovisionermapping (Phase 3)
4. New domains evaluated with documented decisions (Phase 4)
5. Credential schemas/resources updated to match new contracts (Phase 1)
