# Task T01: Implementation Plan — Expand Cloud Resource Tools

**Created**: 2026-02-27
**Status**: PENDING REVIEW
**Type**: Feature Development

⚠️ **This plan requires your review before execution**

## Objective

Expand the MCP server from 3 tools to 16 tools, covering the full cloud resource lifecycle: discovery, CRUD, provisioning observability, and operational management.

## Analysis Summary

A thorough analysis of the Planton gRPC API surface (`plantonhq/planton/apis`) revealed that the current tool set (apply, get, delete) covers basic CRUD but leaves agents **blind** in several critical dimensions:

1. **No list/search** — agents cannot discover what resources exist
2. **No provisioning observability** — agents cannot see if apply/destroy succeeded (StackJobs are async)
3. **delete ≠ destroy** — `delete` removes the DB record; `destroy` tears down infrastructure
4. **No context discovery** — agents cannot discover orgs or environments
5. **No pre-validation** — no slug availability check before apply
6. **No presets** — agents must construct cloud objects from scratch
7. **No lock awareness** — locked resources produce cryptic errors
8. **No rename, env var maps, or cross-resource reference resolution**

## New Tools (13 total, organized by implementation phase)

### Phase 6A: Complete the Resource Lifecycle (2 tools)

| Tool | Backing API | Domain Package |
|------|-------------|----------------|
| `list_cloud_resources` | `CloudResourceQueryController.find(FindApiResourcesRequest)` | `cloudresource/` (existing) |
| `destroy_cloud_resource` | `CloudResourceCommandController.destroy(CloudResource)` | `cloudresource/` (existing) |

**`list_cloud_resources`**
- Input: `org` (required), `env` (optional), `kind` (optional), `page_number` (optional), `page_size` (optional)
- Calls `find(FindApiResourcesRequest)` with pagination
- Returns paginated list of cloud resources (id, kind, name, slug, org, env, status summary)
- **Design question**: `find` is marked "platform operator only" in the proto. Need to verify if user-level auth works. If not, fall back to `streamByOrg` (but streaming RPCs don't map well to MCP tool responses — may need to collect and return).

**`destroy_cloud_resource`**
- Input: Same `ResourceIdentifier` dual-path as delete (id OR kind+org+env+slug)
- Resolves to full CloudResource (via get), then calls `CommandController.destroy(CloudResource)`
- Returns the resource post-destroy
- Tool description must clearly distinguish from `delete_cloud_resource`: "Tears down the real cloud infrastructure (Terraform destroy / Pulumi destroy) while keeping the resource record. Use `delete_cloud_resource` to remove the record from the system."

### Phase 6B: Stack Job Observability (2 tools, new domain)

| Tool | Backing API | Domain Package |
|------|-------------|----------------|
| `get_stack_job_status` | `StackJobQueryController.getLastStackJobByCloudResourceId(CloudResourceId)` | `stackjob/` (NEW) |
| `list_stack_jobs` | `StackJobQueryController.listByFilters(ListStackJobsByFiltersQueryInput)` | `stackjob/` (NEW) |

**`get_stack_job_status`**
- Input: `cloud_resource_id` (required)
- Returns the latest stack job for the given cloud resource, including:
  - Operation type (update, destroy, refresh, etc.)
  - Progress status (running, completed, failed)
  - Result (success, failure, cancelled)
  - Start/end time
  - Error messages (if any)
  - IaC resources created/updated/deleted
- This is the primary tool agents use after `apply` or `destroy` to check "did it work?"

**`list_stack_jobs`**
- Input: `org` (optional), `env` (optional), `kind` (optional), `cloud_resource_id` (optional), `operation_type` (optional), `status` (optional), `result` (optional)
- Calls `listByFilters` with provided filters
- Returns list of stack jobs matching criteria
- Enables queries like "show me all failed deployments in production"

**New domain package**: `internal/domains/stackjob/`
- `tools.go` — tool definitions + handlers
- `get.go` — get last stack job by cloud resource ID
- `list.go` — list by filters

### Phase 6C: Context Discovery (2 tools, new domains)

| Tool | Backing API | Domain Package |
|------|-------------|----------------|
| `list_organizations` | `OrganizationQueryController.findOrganizations(Empty)` | `organization/` (NEW) |
| `list_environments` | `EnvironmentQueryController.findAuthorized(OrganizationId)` | `environment/` (NEW) |

**`list_organizations`**
- Input: none (uses authenticated user's context)
- Returns list of organizations the user belongs to (id, name, slug)
- This is often the first tool an agent calls to establish context

**`list_environments`**
- Input: `org` (required)
- Returns list of environments the user can access within the org (id, name, slug)
- Calls `findAuthorized` (not `findByOrg`) to respect authorization boundaries

**New domain packages**: `internal/domains/organization/`, `internal/domains/environment/`
- Each follows the same pattern: `tools.go` + `list.go`

### Phase 6D: Agent Quality-of-Life (2 tools)

| Tool | Backing API | Domain Package |
|------|-------------|----------------|
| `check_slug_availability` | `CloudResourceQueryController.checkSlugAvailability(...)` | `cloudresource/` (existing) |
| `search_cloud_object_presets` | `InfraHubSearchQueryController.searchOfficialCloudObjectPresets(...)` + `searchCloudObjectPresetsByOrgContext(...)` | `preset/` (NEW) |

**`check_slug_availability`**
- Input: `org` (required), `env` (required), `kind` (required), `slug` (required)
- Returns availability status (available/taken) and, if taken, the existing resource's ID
- Prevents wasted apply attempts

**`search_cloud_object_presets`**
- Input: `kind` (optional), `org` (optional — if provided, includes org-level presets), `query` (optional — text search)
- Searches official presets and optionally org-scoped presets
- Returns preset list: name, description, kind, rank, yaml_content (the actual cloud object YAML)
- Dramatically reduces complexity of constructing cloud objects from scratch

**New domain package**: `internal/domains/preset/`
- `tools.go` — tool definition + handler
- `search.go` — search logic

### Phase 6E: Advanced Operations (5 tools)

| Tool | Backing API | Domain Package |
|------|-------------|----------------|
| `list_cloud_resource_locks` | `CloudResourceLockController.listLocks(CloudResourceId)` | `cloudresource/` (existing) |
| `remove_cloud_resource_locks` | `CloudResourceLockController.removeLocks(CloudResourceId)` | `cloudresource/` (existing) |
| `rename_cloud_resource` | `CloudResourceCommandController.rename(RenameCloudResourceRequest)` | `cloudresource/` (existing) |
| `get_env_var_map` | `CloudResourceQueryController.getEnvVarMap(GetEnvVarMapRequest)` | `cloudresource/` (existing) |
| `resolve_value_references` | `CloudResourceQueryController.resolveValueFromReferences(...)` | `cloudresource/` (existing) |

**`list_cloud_resource_locks`**
- Input: `id` (required — cloud resource ID)
- Returns lock info: lock type, holder, timestamp
- Enables agents to understand why a resource can't be modified

**`remove_cloud_resource_locks`**
- Input: `id` (required)
- Removes all locks on the resource
- Returns success/failure with lock removal details
- Tool description should warn: "Use with caution — removing locks on a resource with an active stack job may cause state corruption"

**`rename_cloud_resource`**
- Input: `id` (required), `new_name` (required)
- Calls `CommandController.rename(RenameCloudResourceRequest)`
- Returns the renamed resource

**`get_env_var_map`**
- Input: `id` (required), `manifest` (required — the cloud object manifest as map)
- Returns the environment variable map the resource exposes
- Useful for agents wiring services together

**`resolve_value_references`**
- Input: `references` (required — list of ValueFromRef objects to resolve)
- Resolves cross-resource `StringValueOrRef` references
- Returns resolved values
- Enables agents to validate inter-resource dependencies

## Cross-Cutting Concerns

### Patterns to Reuse (from Phase 1-5)
- `ResourceIdentifier` dual-path pattern for tools that accept ID or slug composite
- `domains.WithConnection` for gRPC lifecycle
- `domains.RPCError` for gRPC error classification
- `domains.TextResult` for tool responses
- `domains.MarshalJSON` for protojson serialization
- Tool handler as thin proxy: validate at boundary, delegate to domain function

### New Patterns Needed
- **Pagination handling**: `FindApiResourcesRequest` uses cursor/offset pagination. Need a shared pagination input type and response formatting.
- **Filter mapping**: Several tools accept optional filters that map to proto filter types. Need a clean pattern for optional filter → proto conversion.
- **Multi-service connections**: Stack job tools call `StackJobQueryController`, not `CloudResourceCommandController/QueryController`. May need different service addresses or the same Planton API gateway handles routing.

### Testing Strategy
- Unit tests for all pure domain logic (input validation, filter construction, response formatting)
- No mock gRPC server (same decision as Phase 5)
- Test files co-located with source: `*_test.go`

### Documentation
- Update README.md tool table (3 → 16 tools)
- Update MCP resource descriptions if tool discovery workflow changes
- Update `docs/development.md` with new domain packages

## Implementation Order

| Phase | Tools | New Files | Estimated Effort |
|-------|-------|-----------|-----------------|
| 6A | list_cloud_resources, destroy_cloud_resource | 2 new + modify tools.go, server.go | Medium |
| 6B | get_stack_job_status, list_stack_jobs | 3 new files (new domain) + server.go | Medium |
| 6C | list_organizations, list_environments | 4 new files (2 new domains) + server.go | Low |
| 6D | check_slug_availability, search_cloud_object_presets | 3 new files (1 new domain) + modify tools.go, server.go | Medium |
| 6E | locks, rename, envvarmap, references | 4 new files + modify tools.go, server.go | Medium |
| Hardening | Unit tests, README, docs | ~8 test files, README, docs | Medium |

## Open Questions for Review

1. **`find` API auth level**: The `CloudResourceQueryController.find` RPC is annotated as platform operator only. Should we verify this works with user-level API keys? If not, `streamByOrg` is the fallback (but requires collecting the stream into a list, and lacks pagination/filtering).

2. **Destroy confirmation**: Should `destroy_cloud_resource` require any confirmation mechanism, or is the agent/user responsible for confirming intent before calling the tool?

3. **Stack job response size**: Stack jobs can contain large IaC resource lists. Should `get_stack_job_status` return a summary or the full stack job?

4. **Preset YAML content**: `search_cloud_object_presets` returns `yaml_content` which could be large. Should we return full YAML or just metadata (name, description, kind) with a separate `get_cloud_object_preset` tool for the full content?

5. **Phase ordering**: Is the proposed phase order (6A → 6B → 6C → 6D → 6E) acceptable, or would you prefer a different sequencing?

## Success Criteria for T01

- [ ] Plan reviewed and approved by developer
- [ ] Open questions resolved
- [ ] Design decisions documented (with permission)
- [ ] Ready to begin Phase 6A implementation

## Review Process

**What happens next**:
1. **You review this plan** — consider the approach, tool designs, and open questions
2. **Provide feedback** — concerns, changes, answers to open questions
3. **I'll revise the plan** — create T01_2_revised_plan.md incorporating feedback
4. **You approve** — explicit approval to begin implementation
5. **Execution begins** — Phase 6A tracked in T01_3_execution.md
