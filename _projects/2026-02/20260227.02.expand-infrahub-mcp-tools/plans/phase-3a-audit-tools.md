---
name: Phase 3A Audit Tools
overview: Add 3 MCP tools for resource version history and change tracking under a new `internal/domains/audit/` bounded context, expanding the server from 52 to 55 tools.
todos:
  - id: audit-doc
    content: Create `internal/domains/audit/doc.go` with package documentation
    status: completed
  - id: audit-enum
    content: Create `internal/domains/audit/enum.go` with `resolveApiResourceKind()` resolver using `apiresourcekind.ApiResourceKind_value` map
    status: completed
  - id: audit-list
    content: Create `internal/domains/audit/list.go` — `List()` function calling `ListByFilters` RPC with kind resolution + 1-based-to-0-based pagination
    status: completed
  - id: audit-get
    content: Create `internal/domains/audit/get.go` — `Get()` function calling `GetByIdWithContextSize` RPC with default context_size=3
    status: completed
  - id: audit-count
    content: Create `internal/domains/audit/count.go` — `Count()` function calling `GetCount` RPC
    status: completed
  - id: audit-tools
    content: Create `internal/domains/audit/tools.go` — 3 input structs, 3 tool defs, 3 handlers following established patterns
    status: completed
  - id: server-register
    content: Update `internal/server/server.go` — add audit import + 3 mcp.AddTool calls, update count 52→55 and tool name list
    status: completed
  - id: verify-build
    content: Run `go build ./...`, `go vet ./...`, `go test ./...` to verify clean build
    status: completed
isProject: false
---

# Phase 3A: Audit / Version History Tools

## Domain Analysis

Audit is a **cross-cutting read-only query domain** — it answers "what changed, when, and what was the diff?" for any platform resource type. This distinguishes it from all prior phases where each domain package hard-codes its own `ApiResourceKind` (e.g., `infrachart/list.go` hard-codes `ApiResourceKind_infra_chart`). Audit must accept the kind from the user because it queries across resource types.

The proto surface is clean: 4 RPCs, 3 tools (we merge `Get` and `GetByIdWithContextSize` into one).

## Proto Foundation

**Service**: `ai.planton.audit.apiresourceversion.v1.ApiResourceVersionQueryController`
**Import**: `apiresourceversionv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/audit/apiresourceversion/v1"`

**RPCs**:

- `ListByFilters(ListApiResourceVersionsInput) → ApiResourceVersionList`
- `GetByIdWithContextSize(ApiResourceVersionWithContextSizeInput) → ApiResourceVersion`
- `GetCount(GetApiResourceVersionCountInput) → ApiResourceVersionCount`
- `Get(ApiResourceVersionId) → ApiResourceVersion` — superseded by `GetByIdWithContextSize`, not exposed

**Key request types**:

- `ListApiResourceVersionsInput` — `PageInfo`, `Kind` (ApiResourceKind enum), `ResourceId` (string)
- `ApiResourceVersionWithContextSizeInput` — `VersionId` (string), `ContextSize` (int32)
- `GetApiResourceVersionCountInput` — `Kind` (ApiResourceKind enum), `Id` (string)

**ApiResourceKind enum** (from `apiresourcekind` package) — commonly relevant values for audit:
`cloud_resource`, `infra_project`, `infra_chart`, `infra_pipeline`, `variable`, `secret`, `environment`, `organization`, `service`, `stack_job`

## Tools (3)

### Tool 1: `list_resource_versions`

**RPC**: `ListByFilters`
**Required params**: `resource_id` (string), `kind` (string — resolved to `ApiResourceKind` enum)
**Optional params**: `page_num` (int32, default 1), `page_size` (int32, default 20)
**Pagination**: 1-based input, converted to 0-based for proto `PageInfo.Num` (follows the newer convention established in Phases 1A+)
**Returns**: Paginated list of version entries with metadata, event type, timestamps

### Tool 2: `get_resource_version`

**RPC**: `GetByIdWithContextSize` (NOT plain `Get` — strictly more useful)
**Required params**: `version_id` (string)
**Optional params**: `context_size` (int32, default 3 — matches standard unified diff `-U3`)
**Returns**: Full version with original/new state YAML, unified diff, event type, stack job link, cloud object version details

### Tool 3: `get_resource_version_count`

**RPC**: `GetCount`
**Required params**: `resource_id` (string), `kind` (string — resolved to `ApiResourceKind` enum)
**Returns**: Integer count of versions

## Architecture Decision: `ApiResourceKind` Enum Resolver

Unlike existing domains that hard-code their kind, audit needs a dynamic resolver. Design:

- Create `internal/domains/audit/enum.go` with `resolveApiResourceKind(string) (ApiResourceKind, error)`
- Use the established `EnumType_value` map pattern from [internal/domains/infrahub/stackjob/enum.go](internal/domains/infrahub/stackjob/enum.go)
- Accept ALL valid enum values (don't filter — the backend rejects nonsensical queries)
- Exclude `unspecified` from the error message's valid-values list
- Document the most useful values in tool descriptions for agent discoverability
- Keep the resolver in the `audit` package (single consumer); extract to `domains` if a second package needs it

## File Plan

All files under `internal/domains/audit/`:

- `**doc.go`** — Package documentation. References the `ApiResourceVersionQueryController` service.
- `**enum.go`** — `resolveApiResourceKind()` using the `apiresourcekind.ApiResourceKind_value` map. Follows the `stackjob/enum.go` pattern with `joinEnumValues` (reuse approach — not the function itself, since `joinEnumValues` is unexported and in a different package). We define a local `joinEnumValues` or, better, promote a shared version to the `domains` package (since it's already duplicated in `stackjob/enum.go` and `graph/enum.go`).
- `**tools.go**` — 3 input structs, 3 tool definitions, 3 typed handlers. Follows the established pattern from [internal/domains/infrahub/stackjob/tools.go](internal/domains/infrahub/stackjob/tools.go).
- `**list.go**` — `List()` function calling `ListByFilters`. Follows pagination pattern from [internal/domains/configmanager/variable/list.go](internal/domains/configmanager/variable/list.go).
- `**get.go**` — `Get()` function calling `GetByIdWithContextSize`. Simple single-ID lookup like [internal/domains/infrahub/stackjob/get.go](internal/domains/infrahub/stackjob/get.go).
- `**count.go**` — `Count()` function calling `GetCount`. Lightweight — returns the count integer.

**Modified files**:

- `**internal/server/server.go`** — Add `audit` import + 3 `mcp.AddTool` calls. Update tool count 52 → 55 and tool name list.

## Open Question: `joinEnumValues` Duplication

The `joinEnumValues` helper (sorts enum map keys, excludes unspecified, joins with commas) is currently duplicated in:

- `internal/domains/infrahub/stackjob/enum.go`
- `internal/domains/graph/enum.go`
- Presumably also in `internal/domains/configmanager/variable/enum.go` and `secret/enum.go`

Adding yet another copy in `audit/enum.go` increases the duplication. **Options**:

- **Option A**: Copy again (consistent with existing pattern, no cross-package changes)
- **Option B**: Extract to `internal/domains/join.go` as `JoinEnumValues()` and update all consumers

I will follow **Option A** (copy) during implementation to avoid touching unrelated packages in this phase. A cleanup refactor can happen separately. Will pause and confirm with you if the duplication count feels wrong during implementation.

## Verification

- `go build ./...`
- `go vet ./...`
- `go test ./...`

