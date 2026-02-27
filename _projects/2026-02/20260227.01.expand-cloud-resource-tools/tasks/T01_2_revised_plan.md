# Task T01: Revised Implementation Plan — Expand Cloud Resource Tools

**Created**: 2026-02-27
**Revised**: 2026-02-27
**Status**: PENDING APPROVAL
**Type**: Feature Development

## Objective

Expand the MCP server from 3 tools to 17 tools, covering the full cloud resource lifecycle: discovery, CRUD, provisioning observability, and operational management.

## Changes from Original Plan

1. **Dropped `find` AND `streamByOrg`** for `list_cloud_resources` — `find` is platform-operator only (meant for re-indexing the search index). `streamByOrg` is a raw database stream not intended for user-facing listing. Instead, using `CloudResourceSearchQueryController.getCloudResourcesCanvasView` from the search domain — this queries the search index with server-side filtering (envs, kinds, text search) and returns lightweight search records.
2. **No MCP-level destroy confirmation** — agent is responsible for confirming intent with user.
3. **Full stack job response** — no truncation, no secrets present.
4. **Split preset tool into search + get** — `search_cloud_object_presets` returns metadata, new `get_cloud_object_preset` returns full content. Total tools: 14 new (17 total).
5. **Phase ordering confirmed** as 6A → 6B → 6C → 6D → 6E → Hardening.

---

## Tool Inventory (14 new + 3 existing = 17 total)

### Existing Tools (3)
| # | Tool | Domain |
|---|------|--------|
| 1 | `apply_cloud_resource` | `cloudresource/` |
| 2 | `get_cloud_resource` | `cloudresource/` |
| 3 | `delete_cloud_resource` | `cloudresource/` |

### New Tools (14)
| # | Tool | Phase | Domain | Backing RPC |
|---|------|-------|--------|-------------|
| 4 | `list_cloud_resources` | 6A | `cloudresource/` | `CloudResourceSearchQueryController.getCloudResourcesCanvasView` |
| 5 | `destroy_cloud_resource` | 6A | `cloudresource/` | `CommandController.destroy` |
| 6 | `get_stack_job_status` | 6B | `stackjob/` (new) | `StackJobQueryController.getLastStackJobByCloudResourceId` |
| 7 | `list_stack_jobs` | 6B | `stackjob/` (new) | `StackJobQueryController.listByFilters` |
| 8 | `list_organizations` | 6C | `organization/` (new) | `OrganizationQueryController.findOrganizations` |
| 9 | `list_environments` | 6C | `environment/` (new) | `EnvironmentQueryController.findAuthorized` |
| 10 | `check_slug_availability` | 6D | `cloudresource/` | `QueryController.checkSlugAvailability` |
| 11 | `search_cloud_object_presets` | 6D | `preset/` (new) | `InfraHubSearchQueryController.searchOfficialCloudObjectPresets` + `searchCloudObjectPresetsByOrgContext` |
| 12 | `get_cloud_object_preset` | 6D | `preset/` (new) | `CloudObjectPresetQueryController.get` |
| 13 | `list_cloud_resource_locks` | 6E | `cloudresource/` | `LockController.listLocks` |
| 14 | `remove_cloud_resource_locks` | 6E | `cloudresource/` | `LockController.removeLocks` |
| 15 | `rename_cloud_resource` | 6E | `cloudresource/` | `CommandController.rename` |
| 16 | `get_env_var_map` | 6E | `cloudresource/` | `QueryController.getEnvVarMap` |
| 17 | `resolve_value_references` | 6E | `cloudresource/` | `QueryController.resolveValueFromReferences` |

---

## Phase 6A: Complete the Resource Lifecycle (2 tools)

### `list_cloud_resources`

**Backing RPC**: `CloudResourceSearchQueryController.getCloudResourcesCanvasView(ExploreCloudResourcesRequest) returns (ExploreCloudResourcesCanvasViewResponse)`
- Service: `ai.planton.search.v1.infrahub.cloudresource.CloudResourceSearchQueryController`
- Proto: `ai/planton/search/v1/infrahub/cloudresource/query.proto`
- Auth: `get` permission on `organization`

**Input**:
- `org` (required) — organization ID
- `envs` (optional) — list of environment slugs to filter by
- `search_text` (optional) — free-text search query
- `kinds` (optional) — list of cloud resource kinds to filter by

**Output**: The response is `ExploreCloudResourcesCanvasViewResponse` which contains a list of `CanvasEnvironment` objects. Each `CanvasEnvironment` has:
- `env_id`, `env_slug`
- `resource_kind_mapping` — a `map<string, ApiResourceSearchRecords>` grouping search records by kind

Each `ApiResourceSearchRecord` is a lightweight search record containing: `id`, `name`, `kind`, `org`, `env`, `slug`, `tags`, `description`, `created_by`, `created_at`.

**Implementation**:
- Maps MCP tool input to `ExploreCloudResourcesRequest` proto
- Calls `getCloudResourcesCanvasView` via gRPC
- Serializes the response as JSON

**Design notes**:
- This queries the search index (not the database directly), so results are fast and lightweight
- All filtering (envs, kinds, text) is handled server-side
- No pagination needed — canvas view is designed for browsable result sets
- `lookupCloudResource` is also available on the same service for single-resource lookup by org/env/kind/name, but `get_cloud_resource` already covers that use case

**New files**: `internal/domains/cloudresource/list.go`

### `destroy_cloud_resource`

**Backing RPC**: `CloudResourceCommandController.destroy(CloudResource) returns (CloudResource)`
- Auth: `delete` permission on `cloud_resource`
- Proto comment: "meant to be called from CLI or web-app, not from pipelines"

**Input**: Same `ResourceIdentifier` dual-path as existing delete/get tools:
- `id` (cloud resource ID), OR
- `kind` + `org` + `env` + `slug`

**Implementation**:
- Resolves to full `CloudResource` via `get` (same pattern as `delete`)
- Calls `CommandController.destroy(CloudResource)`
- Returns the resource post-destroy as JSON

**Tool description**: Must clearly distinguish from `delete_cloud_resource`:
> "Tears down the real cloud infrastructure (Terraform/Pulumi destroy) while keeping the resource record. Use `delete_cloud_resource` to remove the record from the system. WARNING: This is a destructive operation that will destroy real cloud infrastructure."

**No MCP-level confirmation** — the agent/user is responsible for confirming intent before calling.

**New files**: `internal/domains/cloudresource/destroy.go`

---

## Phase 6B: Stack Job Observability (2 tools, new domain)

**New domain package**: `internal/domains/stackjob/`

### `get_stack_job_status`

**Backing RPC**: `StackJobQueryController.getLastStackJobByCloudResourceId(CloudResourceId) returns (StackJob)`
- Auth: `get` permission on `cloud_resource`

**Input**:
- `cloud_resource_id` (required)

**Output**: Full `StackJob` object as JSON — includes operation type, progress status, result, timestamps, error messages, IaC resource counts. No truncation (no secrets in this response).

**Use case**: Primary tool agents call after `apply` or `destroy` to check "did it work?"

**New files**: `internal/domains/stackjob/tools.go`, `internal/domains/stackjob/get.go`

### `list_stack_jobs`

**Backing RPC**: `StackJobQueryController.listByFilters(ListStackJobsByFiltersQueryInput) returns (StackJobList)`

**Input** (all optional filters):
- `org`, `env`, `kind`, `cloud_resource_id`
- `operation_type`, `status`, `result`

**Output**: Full list of matching stack jobs as JSON.

**Use case**: "Show me all failed deployments in production"

**New files**: `internal/domains/stackjob/list.go`

---

## Phase 6C: Context Discovery (2 tools, new domains)

### `list_organizations`

**Backing RPC**: `OrganizationQueryController.findOrganizations(CustomEmpty) returns (Organizations)`
- Auth: No authorization required — returns only orgs the authenticated user belongs to

**Input**: None (uses authenticated user context)

**Output**: List of organizations (id, name, slug)

**Use case**: Often the first tool an agent calls to establish operating context.

**New domain**: `internal/domains/organization/`
**New files**: `internal/domains/organization/tools.go`, `internal/domains/organization/list.go`

### `list_environments`

**Backing RPC**: `EnvironmentQueryController.findAuthorized(OrganizationId) returns (Environments)`
- Auth: `get` permission on `organization`
- Returns ONLY environments where the calling user has at least `get` permission (FGA-filtered)

**Input**:
- `org` (required) — organization ID

**Output**: List of authorized environments (id, name, slug)

**New domain**: `internal/domains/environment/`
**New files**: `internal/domains/environment/tools.go`, `internal/domains/environment/list.go`

---

## Phase 6D: Agent Quality-of-Life (3 tools)

### `check_slug_availability`

**Backing RPC**: `CloudResourceQueryController.checkSlugAvailability(CloudResourceSlugAvailabilityCheckRequest) returns (CloudResourceSlugAvailabilityCheckResponse)`
- Auth: `get` permission on `organization`

**Input**:
- `org` (required), `env` (required), `kind` (required), `slug` (required)

**Output**: Availability status (available/taken), and if taken, the existing resource's ID.

**New files**: `internal/domains/cloudresource/slug.go`

### `search_cloud_object_presets`

**Backing RPCs**:
- `InfraHubSearchQueryController.searchOfficialCloudObjectPresets(SearchOfficialCloudObjectPresetsInput) returns (ApiResourceSearchRecordList)` — public endpoint
- `InfraHubSearchQueryController.searchCloudObjectPresetsByOrgContext(SearchCloudObjectPresetsByOrgContextInput) returns (ApiResourceSearchRecordList)` — requires `get` on `organization`

**Input**:
- `kind` (optional) — filter by cloud resource kind
- `org` (optional) — if provided, searches org-scoped presets in addition to official ones
- `query` (optional) — text search

**Output**: Lightweight search records (metadata: id, name, description, kind, rank). No full YAML content.

**Implementation**: If `org` is provided, call `searchCloudObjectPresetsByOrgContext` (with `is_include_official = true`). If no `org`, call `searchOfficialCloudObjectPresets`.

**New domain**: `internal/domains/preset/`
**New files**: `internal/domains/preset/tools.go`, `internal/domains/preset/search.go`

### `get_cloud_object_preset`

**Backing RPC**: `CloudObjectPresetQueryController.get(ApiResourceId) returns (CloudObjectPreset)`

**Input**:
- `id` (required) — the preset ID (from search results)

**Output**: Full `CloudObjectPreset` including spec with YAML content.

**Use case**: Agent searches presets to find one, then gets the full content to use as a template for `apply_cloud_resource`.

**New files**: `internal/domains/preset/get.go`

---

## Phase 6E: Advanced Operations (5 tools)

### `list_cloud_resource_locks`

**Backing RPC**: `CloudResourceLockController.listLocks(CloudResourceId) returns (CloudResourceLockInfo)`
- Auth: `update` permission on `cloud_resource`

**Input**: `id` (required — cloud resource ID)

**Output**: Lock info including `is_locked`, current lock holder (workflow ID, acquired timestamp, TTL remaining), and queue entries.

**New files**: `internal/domains/cloudresource/locks.go`

### `remove_cloud_resource_locks`

**Backing RPC**: `CloudResourceLockController.removeLocks(CloudResourceId) returns (CloudResourceLockRemovalResponse)`
- Auth: `update` permission on `cloud_resource`

**Input**: `id` (required)

**Output**: Removal result (lock_removed, queue_entries_removed, message).

**Tool description warning**:
> "Use with caution — removing locks on a resource with an active stack job may cause state corruption."

Shares `locks.go` with `list_cloud_resource_locks`.

### `rename_cloud_resource`

**Backing RPC**: `CloudResourceCommandController.rename(RenameCloudResourceRequest) returns (CloudResource)`
- Auth: `update` permission on `cloud_resource`

**Input**: `id` (required), `new_name` (required)

**Output**: The renamed `CloudResource`.

**New files**: `internal/domains/cloudresource/rename.go`

### `get_env_var_map`

**Backing RPC**: `CloudResourceQueryController.getEnvVarMap(GetEnvVarMapRequest) returns (GetEnvVarMapResponse)`
- Auth: Handled in handler — requires `get` on the cloud resource resolved from YAML

**Input**: `id` (required), `manifest` (required — the cloud object manifest as map)

**Output**: Environment variable map the resource exposes.

**New files**: `internal/domains/cloudresource/envvarmap.go`

### `resolve_value_references`

**Backing RPC**: `CloudResourceQueryController.resolveValueFromReferences(ResolveValueFromReferencesRequest) returns (ResolveValueFromReferencesResponse)`
- Auth: `get` permission on `cloud_resource`

**Input**: `cloud_resource_id` (required), `references` (required — list of references to resolve)

**Output**: Resolved values.

**New files**: `internal/domains/cloudresource/references.go`

---

## Cross-Cutting Concerns

### Patterns to Reuse (from existing codebase)
- `ResourceIdentifier` dual-path pattern for tools accepting ID or slug composite
- `domains.WithConnection` for gRPC lifecycle management
- `domains.RPCError` for gRPC error classification
- `domains.TextResult` for wrapping responses as MCP tool results
- `domains.MarshalJSON` for protojson serialization
- Tool handler as thin proxy: validate input at boundary, delegate to domain function
- Tool registration via `mcp.AddTool(srv, *Tool(), *Handler(serverAddress))` in `server.go`

### New Patterns Needed
- **Optional filter mapping**: Several tools accept optional filters that map to proto filter types. Need a clean pattern for optional filter → proto field conversion.
- **Multi-service gRPC connections**: New tools call different gRPC services (`CloudResourceSearchQueryController`, `StackJobQueryController`, `OrganizationQueryController`, `EnvironmentQueryController`, `CloudResourceLockController`, `InfraHubSearchQueryController`, `CloudObjectPresetQueryController`). The existing `WithConnection` pattern creates a `grpc.ClientConn` to `serverAddress` — this should work if all services are behind the same API gateway. Verify during Phase 6A/6B.

### Testing Strategy
- Unit tests for all pure domain logic (input validation, filter construction, response formatting)
- No mock gRPC server (consistent with existing codebase)
- Test files co-located: `*_test.go`

### Documentation
- Update `README.md` tool table (3 → 17 tools)
- Update `docs/development.md` with new domain packages

---

## Implementation Order

| Phase | Tools | New Packages | New Files | Modify |
|-------|-------|-------------|-----------|--------|
| 6A | `list_cloud_resources`, `destroy_cloud_resource` | — | `list.go`, `destroy.go` | `tools.go`, `server.go` |
| 6B | `get_stack_job_status`, `list_stack_jobs` | `stackjob/` | `tools.go`, `get.go`, `list.go` | `server.go` |
| 6C | `list_organizations`, `list_environments` | `organization/`, `environment/` | 4 files (2× `tools.go`, `list.go`) | `server.go` |
| 6D | `check_slug_availability`, `search_cloud_object_presets`, `get_cloud_object_preset` | `preset/` | `slug.go`, `tools.go`, `search.go`, `get.go` | `tools.go`, `server.go` |
| 6E | locks, rename, envvarmap, references | — | `locks.go`, `rename.go`, `envvarmap.go`, `references.go` | `tools.go`, `server.go` |
| Hardening | Unit tests, README, docs | — | ~10 test files | `README.md`, docs |

---

## gRPC Service → Proto Package Mapping

For reference during implementation:

| gRPC Service | Proto Package | Proto File |
|-------------|---------------|------------|
| `CloudResourceQueryController` | `ai.planton.infrahub.cloudresource.v1` | `query.proto` |
| `CloudResourceCommandController` | `ai.planton.infrahub.cloudresource.v1` | `command.proto` |
| `CloudResourceLockController` | `ai.planton.infrahub.cloudresource.v1` | `lock.proto` |
| `CloudResourceSearchQueryController` | `ai.planton.search.v1.infrahub.cloudresource` | `query.proto` |
| `StackJobQueryController` | `ai.planton.infrahub.stackjob.v1` | `query.proto` |
| `OrganizationQueryController` | `ai.planton.resourcemanager.organization.v1` | `query.proto` |
| `EnvironmentQueryController` | `ai.planton.resourcemanager.environment.v1` | `query.proto` |
| `InfraHubSearchQueryController` | `ai.planton.search.v1.infrahub` | `query.proto` |
| `CloudObjectPresetQueryController` | `ai.planton.infrahub.cloudobjectpreset.v1` | `query.proto` |

---

## Success Criteria

- [ ] Plan approved by developer
- [ ] All 14 new tools implemented, registered, and functional
- [ ] Unit tests for all pure domain logic
- [ ] README updated with full 17-tool table
- [ ] Documentation updated with new domain packages
