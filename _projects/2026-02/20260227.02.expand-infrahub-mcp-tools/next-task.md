# Next Task: 20260227.02.expand-infrahub-mcp-tools

## ⚠️ RULES OF ENGAGEMENT - READ FIRST ⚠️

**When this file is loaded in a new conversation, the AI MUST:**

1. **DO NOT AUTO-EXECUTE** - Never start implementing without explicit user approval
2. **GATHER CONTEXT SILENTLY** - Read project files without outputting
3. **PRESENT STATUS SUMMARY** - Show what's done, what's pending, agreed next steps
4. **SHOW OPTIONS** - List recommended and alternative actions
5. **WAIT FOR DIRECTION** - Do NOT proceed until user explicitly says "go" or chooses an option

---

## Project: 20260227.02.expand-infrahub-mcp-tools

**Description**: Add InfraChart, InfraProject, InfraPipeline, Graph, ConfigManager, Audit, StackJob commands, Deployment Component catalog tools to the MCP server, and restructure generated code under domain-scoped directories.
**Goal**: Expand the MCP server from 18 tools (cloud resource CRUD only) to ~50+ tools covering the full InfraHub composition, pipeline observability, dependency intelligence, configuration lifecycle, audit trail, and operational control surface.
**Tech Stack**: Go/gRPC/MCP
**Components**: internal/domains/infrahub/, internal/domains/graph/, internal/domains/configmanager/, internal/domains/audit/, gen/, cmd/server

## Current Status

**Current Task**: Phase 3B/3C (StackJob Commands or Catalog)
**Status**: Ready to pick next phase

**Current step:**
- ✅ **Phase 0: Restructure Generated Code** (2026-02-27)
  - Moved gen/cloudresource/ -> gen/infrahub/cloudresource/ (362 providers, 17 clouds)
  - Updated code generator, Makefile, consumer import
  - Build, tests, vet all pass clean
  - Committed: `refactor(codegen): move generated code to gen/infrahub/cloudresource/`
- ✅ **Phase 1A: InfraChart tools** (2026-02-27)
  - 3 new tools: `list_infra_charts`, `get_infra_chart`, `build_infra_chart`
  - Server expanded from 18 to 21 tools
  - Build, tests, vet all pass clean
- ✅ **Phase 1B: InfraProject tools** (2026-02-27)
  - 6 new tools: `search_infra_projects`, `get_infra_project`, `apply_infra_project`, `delete_infra_project`, `check_infra_project_slug`, `undeploy_infra_project`
  - Server expanded from 21 to 27 tools
  - Build, vet, tests all pass clean
- ✅ **Phase 1C: InfraPipeline tools** (2026-02-28)
  - 7 new tools (expanded from planned 5): `list_infra_pipelines`, `get_infra_pipeline`, `get_latest_infra_pipeline`, `run_infra_pipeline`, `cancel_infra_pipeline`, `resolve_infra_pipeline_env_gate`, `resolve_infra_pipeline_node_gate`
  - Server expanded from 27 to 34 tools
  - Build, vet, tests all pass clean
  - Phase 1 trifecta complete (Chart + Project + Pipeline)
- ✅ **Phase 2A: Graph / Dependency Intelligence** (2026-02-28)
  - 7 new tools (expanded from planned 4): `get_organization_graph`, `get_environment_graph`, `get_service_graph`, `get_cloud_resource_graph`, `get_dependencies`, `get_dependents`, `get_impact_analysis`
  - Server expanded from 34 to 41 tools
  - First domain outside infrahub bounded context (`internal/domains/graph/`)
  - Build, vet, tests all pass clean
- ✅ **Phase 2B: ConfigManager / Variables & Secrets** (2026-02-28)
  - 11 new tools (expanded from planned 5): `list_variables`, `get_variable`, `apply_variable`, `delete_variable`, `resolve_variable`, `list_secrets`, `get_secret`, `apply_secret`, `delete_secret`, `create_secret_version`, `list_secret_versions`
  - Server expanded from 41 to 52 tools
  - Second domain outside infrahub bounded context (`internal/domains/configmanager/`)
  - Three sub-packages: variable/ (5 tools), secret/ (4 tools), secretversion/ (2 tools)
  - Build, vet, tests all pass clean
  - Phase 2 pair complete (Graph + ConfigManager)
- ✅ **Phase 3A: Audit / Version History** (2026-02-28)
  - 3 new tools: `list_resource_versions`, `get_resource_version`, `get_resource_version_count`
  - Server expanded from 52 to 55 tools
  - Third domain outside infrahub bounded context (`internal/domains/audit/`)
  - Dynamic `ApiResourceKind` enum resolver (first cross-resource-type domain)
  - Build, vet, tests all pass clean
- 🔵 Next: **Phase 3B: StackJob Commands** (3 tools) or **Phase 3C: Catalog** (3 tools)

---

### ✅ COMPLETED: Phase 0 — Gen Code Restructure (2026-02-27)

**Restructured generated code from flat gen/cloudresource/ to domain-scoped gen/infrahub/cloudresource/**

**What was delivered:**

1. **Code generator update** - Default output path changed to `gen/infrahub/cloudresource/`
   - `tools/codegen/generator/main.go` - flag default + doc comment
   - `tools/codegen/generator/codegen.go` - doc comment
   - `tools/codegen/generator/registry.go` - doc comment

2. **Makefile update** - `codegen-types` target updated
   - `Makefile` - rm and output paths

3. **Consumer import update** - Single import path change
   - `internal/domains/infrahub/cloudresource/tools.go`

4. **Regeneration** - 362 provider types across 17 clouds at new path

**Key Decisions Made:**
- Delete + regenerate strategy (not git mv) since generated files have timestamps in headers
- Zero logic changes — pure path/import restructuring

**Files Changed/Created:**
- `tools/codegen/generator/main.go` - Updated default output dir + doc
- `tools/codegen/generator/codegen.go` - Updated doc comment
- `tools/codegen/generator/registry.go` - Updated doc comment
- `Makefile` - Updated codegen-types target
- `internal/domains/infrahub/cloudresource/tools.go` - Updated import path
- `gen/infrahub/cloudresource/` - 367 regenerated files (17 provider dirs + registry)

**Verification:** `go build ./...` ✅ | `go test ./...` ✅ | `make vet` ✅

---

### ✅ COMPLETED: Phase 1A — InfraChart Tools (2026-02-27)

**Added 3 MCP tools for infra chart discovery and rendering, expanding the server from 18 to 21 tools.**

**What was delivered:**

1. **`list_infra_charts`** - Paginated listing via `InfraChartQueryController.Find`
   - `internal/domains/infrahub/infrachart/list.go` - Find RPC with org/env filters
   - 1-based page numbers in tool API, converted to 0-based for proto

2. **`get_infra_chart`** - Retrieve by ID via `InfraChartQueryController.Get`
   - `internal/domains/infrahub/infrachart/get.go` - Standard get-by-ID pattern

3. **`build_infra_chart`** - Preview rendered chart via Get + Build two-step
   - `internal/domains/infrahub/infrachart/build.go` - Fetches chart, merges param overrides, builds
   - `mergeParams()` validates override names against chart's defined params

4. **Tool definitions and handlers**
   - `internal/domains/infrahub/infrachart/tools.go` - Input structs, tool defs, typed handlers

5. **Server registration**
   - `internal/server/server.go` - Import + 3 `mcp.AddTool` calls, count 18→21
   - `internal/domains/infrahub/doc.go` - Updated subpackage list

**Key Decisions Made:**
- Named `list_infra_charts` (not `search_`) because the Find RPC only supports org/env/pagination, not free-text search
- Build tool uses simplified `chart_id` + `params` map input (two RPCs) instead of raw InfraChart proto passthrough
- Hard-coded `ApiResourceKind_infra_chart` (enum value 31) for the Find request's required kind field

**Files Created:**
- `internal/domains/infrahub/infrachart/tools.go`
- `internal/domains/infrahub/infrachart/get.go`
- `internal/domains/infrahub/infrachart/list.go`
- `internal/domains/infrahub/infrachart/build.go`

**Files Modified:**
- `internal/server/server.go` - Tool registration + count
- `internal/domains/infrahub/doc.go` - Subpackage docs

**Verification:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅

---

### ✅ COMPLETED: Phase 1B — InfraProject Tools (2026-02-27)

**Added 6 MCP tools for infra project lifecycle management, expanding the server from 21 to 27 tools.**

**What was delivered:**

1. **`search_infra_projects`** - Free-text search via `InfraHubSearchQueryController.SearchInfraProjects`
   - `internal/domains/infrahub/infraproject/search.go` - Org-scoped search with env/text/pagination filters
   - 1-based page numbers in tool API, converted to 0-based for proto

2. **`get_infra_project`** - Dual identification (ID or org+slug) via `InfraProjectQueryController`
   - `internal/domains/infrahub/infraproject/get.go` - Get + GetByOrgBySlug with resolveProject/resolveProjectID helpers

3. **`apply_infra_project`** - Create/update via `InfraProjectCommandController.Apply`
   - `internal/domains/infrahub/infraproject/apply.go` - Full JSON passthrough using protojson.Unmarshal
   - Supports both infra_chart and git_repo source types

4. **`delete_infra_project`** - Remove record via `InfraProjectCommandController.Delete`
   - `internal/domains/infrahub/infraproject/delete.go` - Resolves ID then calls Delete with ApiResourceDeleteInput

5. **`check_infra_project_slug`** - Slug availability via `InfraProjectQueryController.CheckSlugAvailability`
   - `internal/domains/infrahub/infraproject/slug.go` - Org-scoped only (simpler than CloudResource)

6. **`undeploy_infra_project`** - Tear down infra via `InfraProjectCommandController.Undeploy`
   - `internal/domains/infrahub/infraproject/undeploy.go` - Keeps record, triggers undeploy pipeline

7. **Tool definitions and handlers**
   - `internal/domains/infrahub/infraproject/tools.go` - 6 input structs, 6 tool defs, 6 handlers, shared validateIdentification()

8. **Server registration**
   - `internal/server/server.go` - Import + 6 `mcp.AddTool` calls, count 21→27
   - `internal/domains/infrahub/doc.go` - Updated subpackage list

**Key Decisions Made:**
- DD-1: Simpler identification pattern (ID or org+slug) vs CloudResource's 4-field slug path — InfraProject slugs are org-scoped only
- DD-2: Apply accepts full InfraProject JSON via protojson.Unmarshal — honest passthrough, supports both source types
- DD-3: Purge RPC intentionally excluded — compound destructive operation; agents can undeploy then delete with review checkpoint
- DD-4: Search over Find — SearchInfraProjects (free-text + org + env + pagination) is strictly more powerful than Find (pagination-only)

**Files Created:**
- `internal/domains/infrahub/infraproject/tools.go`
- `internal/domains/infrahub/infraproject/search.go`
- `internal/domains/infrahub/infraproject/get.go`
- `internal/domains/infrahub/infraproject/apply.go`
- `internal/domains/infrahub/infraproject/delete.go`
- `internal/domains/infrahub/infraproject/slug.go`
- `internal/domains/infrahub/infraproject/undeploy.go`

**Files Modified:**
- `internal/server/server.go` - Tool registration + count
- `internal/domains/infrahub/doc.go` - Subpackage docs

**Verification:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅

---

### ✅ COMPLETED: Phase 1C — InfraPipeline Tools (2026-02-28)

**Added 7 MCP tools for infra pipeline observability and control, expanding the server from 27 to 34 tools. Completes the Phase 1 trifecta (Chart + Project + Pipeline).**

**What was delivered:**

1. **`list_infra_pipelines`** - Paginated listing via `InfraPipelineQueryController.ListByFilters`
   - `internal/domains/infrahub/infrapipeline/list.go` - Org-scoped with optional project ID filter
   - 1-based page numbers in tool API, converted to 0-based for proto

2. **`get_infra_pipeline`** - Retrieve by ID via `InfraPipelineQueryController.Get`
   - `internal/domains/infrahub/infrapipeline/get.go` - Standard get-by-ID pattern

3. **`get_latest_infra_pipeline`** - Most recent pipeline for a project via `InfraPipelineQueryController.GetLastInfraPipelineByInfraProjectId`
   - `internal/domains/infrahub/infrapipeline/latest.go` - Mirrors stackjob/latest.go pattern

4. **`run_infra_pipeline`** - Unified trigger for both project source types
   - `internal/domains/infrahub/infrapipeline/run.go` - Dispatches to RunInfraProjectChartSourcePipeline or RunGitCommit based on commit_sha presence

5. **`cancel_infra_pipeline`** - Cancel a running pipeline via `InfraPipelineCommandController.Cancel`
   - `internal/domains/infrahub/infrapipeline/cancel.go` - Returns updated pipeline after cancellation

6. **`resolve_infra_pipeline_env_gate`** - Approve/reject environment manual gate
   - `internal/domains/infrahub/infrapipeline/gate.go` - Maps "approve"/"reject" to proto's "yes"/"no" enum

7. **`resolve_infra_pipeline_node_gate`** - Approve/reject DAG node manual gate
   - `internal/domains/infrahub/infrapipeline/gate.go` - Shared gate resolution with resolveDecision() helper

8. **Tool definitions and handlers**
   - `internal/domains/infrahub/infrapipeline/tools.go` - 7 input structs, 7 tool defs, 7 handlers

9. **Server registration**
   - `internal/server/server.go` - Import + 7 `mcp.AddTool` calls, count 27→34
   - `internal/domains/infrahub/doc.go` - Updated subpackage list

**Key Decisions Made:**
- DD-1: Unified run tool — single `run_infra_pipeline` with optional `commit_sha` dispatches to chart-source or git-commit RPC
- DD-2: Manual gate tools included — without them, agents hit dead ends at approval gates
- DD-3: User-friendly gate decisions — tool accepts "approve"/"reject", translates to proto's "yes"/"no"
- DD-4: Streaming RPCs excluded — GetStatusStream/GetLogStream incompatible with MCP; get_infra_pipeline snapshot suffices
- DD-5: Pipeline CRUD excluded — Apply/Create/Update/Delete are internal platform operations
- DD-6: ListByFilters correction — proto only supports org + infra_project_id + pagination (no status/result filters)
- DD-7: 0-based pagination — follows infrachart/infraproject convention (stackjob inconsistency noted)

**Files Created:**
- `internal/domains/infrahub/infrapipeline/tools.go`
- `internal/domains/infrahub/infrapipeline/list.go`
- `internal/domains/infrahub/infrapipeline/get.go`
- `internal/domains/infrahub/infrapipeline/latest.go`
- `internal/domains/infrahub/infrapipeline/run.go`
- `internal/domains/infrahub/infrapipeline/cancel.go`
- `internal/domains/infrahub/infrapipeline/gate.go`

**Files Modified:**
- `internal/server/server.go` - Tool registration + count
- `internal/domains/infrahub/doc.go` - Subpackage docs

**Verification:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅

---

### ✅ COMPLETED: Phase 2A — Graph / Dependency Intelligence (2026-02-28)

**Added 7 MCP tools for dependency intelligence and impact analysis, expanding the server from 34 to 41 tools. First domain outside the infrahub bounded context.**

**What was delivered:**

1. **`get_organization_graph`** - Full resource topology via `GraphQueryController.GetOrganizationGraph`
   - `internal/domains/graph/organization.go` - Org-scoped with env/node-type/depth filters
   - Supports topological ordering for deployment order analysis

2. **`get_environment_graph`** - Environment-scoped graph via `GraphQueryController.GetEnvironmentGraph`
   - `internal/domains/graph/environment.go` - Everything deployed in a specific environment

3. **`get_service_graph`** - Service-centric subgraph via `GraphQueryController.GetServiceGraph`
   - `internal/domains/graph/service.go` - Service deployments per environment with upstream/downstream traversal

4. **`get_cloud_resource_graph`** - Resource-centric dependency view via `GraphQueryController.GetCloudResourceGraph`
   - `internal/domains/graph/cloudresource.go` - Services, credentials, and dependency neighbors

5. **`get_dependencies`** - Upstream traversal via `GraphQueryController.GetDependencies`
   - `internal/domains/graph/dependency.go` - "What does this resource depend on?" with relationship type filter

6. **`get_dependents`** - Downstream traversal via `GraphQueryController.GetDependents`
   - `internal/domains/graph/dependency.go` - "What depends on this resource?" with relationship type filter

7. **`get_impact_analysis`** - Impact assessment via `GraphQueryController.GetImpactAnalysis`
   - `internal/domains/graph/impact.go` - Direct + transitive impacts, counts, breakdown by type

8. **Enum resolvers**
   - `internal/domains/graph/enum.go` - resolveNodeTypes, resolveRelationshipTypes, resolveChangeType with user-friendly error messages

9. **Tool definitions and handlers**
   - `internal/domains/graph/tools.go` - 7 input structs, 7 tool defs, 7 typed handlers

10. **Server registration**
    - `internal/server/server.go` - Import + 7 `mcp.AddTool` calls, count 34→41

**Key Decisions Made:**
- DD-1: Expanded from planned 4 tools to 7 — `getEnvironmentGraph`, `getServiceGraph`, `getDependents` discovered during proto analysis, all high-value read-only queries
- DD-2: New bounded context — `internal/domains/graph/` (not under `infrahub/`) mirrors proto package `ai.planton.graph.v1`
- DD-3: Shared `DependencyInput` struct for `get_dependencies` and `get_dependents` — identical parameter shapes, same response type
- DD-4: Enum handling follows `stackjob/enum.go` pattern — proto `_value` maps with `joinEnumValues` for error messages

**Files Created:**
- `internal/domains/graph/doc.go`
- `internal/domains/graph/enum.go`
- `internal/domains/graph/tools.go`
- `internal/domains/graph/organization.go`
- `internal/domains/graph/environment.go`
- `internal/domains/graph/service.go`
- `internal/domains/graph/cloudresource.go`
- `internal/domains/graph/dependency.go`
- `internal/domains/graph/impact.go`

**Files Modified:**
- `internal/server/server.go` - Tool registration + count

**Verification:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅

---

### ✅ COMPLETED: Phase 2B — ConfigManager / Variables & Secrets (2026-02-28)

**Added 11 MCP tools for configuration lifecycle management across 3 sub-domains, expanding the server from 41 to 52 tools. Second domain outside infrahub bounded context. Completes Phase 2 pair (Graph + ConfigManager).**

**What was delivered:**

1. **Variable tools (5)** — `internal/domains/configmanager/variable/`
   - `list_variables` — Paginated listing via `VariableQueryController.Find` with `ApiResourceKind_variable` (130)
   - `get_variable` — Dual identification (ID or org+scope+slug) via Get/GetByOrgByScopeBySlug
   - `apply_variable` — Create/update with explicit params (name, org, scope, env, description, value)
   - `delete_variable` — Resolve ID then Delete via `ApiResourceDeleteInput`
   - `resolve_variable` — Quick value lookup returning plain string via `VariableQueryController.Resolve`
   - `enum.go` — `resolveScope` for `VariableSpec_Scope` (organization, environment)

2. **Secret tools (4)** — `internal/domains/configmanager/secret/`
   - `list_secrets` — Paginated listing via `SecretQueryController.Find` with `ApiResourceKind_secret` (38)
   - `get_secret` — Dual identification (ID or org+scope+slug) — metadata only, no values
   - `apply_secret` — Create/update with explicit params (name, org, scope, env, description, backend)
   - `delete_secret` — Destructive: removes secret AND all versions permanently
   - `enum.go` — `resolveScope` for `SecretSpec_Scope`

3. **SecretVersion tools (2)** — `internal/domains/configmanager/secretversion/`
   - `create_secret_version` — Store encrypted key-value pairs via `SecretVersionCommandController.Create`
   - `list_secret_versions` — List version metadata via `SecretVersionQueryController.ListBySecret` (no data field)

4. **Domain documentation** — `internal/domains/configmanager/doc.go`

5. **Server registration** — `internal/server/server.go` — 3 imports + 11 `mcp.AddTool` calls, count 41→52

**Key Decisions Made:**
- DD-1: Expanded from planned 5 to 11 tools — proto analysis revealed 6 gRPC services with ~20 RPCs; 11 tools cover complete lifecycle
- DD-2: Write-only secret values — `create_secret_version` included, `get_secret_version`/`get_latest_secret_version` excluded (security boundary per AD-01 logic)
- DD-3: Explicit parameters for apply tools — Variable and Secret have simple schemas; better agent UX than JSON passthrough
- DD-4: Exclude `refresh_variable` — specialized operation, low initial demand
- DD-5: `delete_secret` includes destructive WARNING — same pattern as `destroy_cloud_resource`

**Files Created:**
- `internal/domains/configmanager/doc.go`
- `internal/domains/configmanager/variable/tools.go`
- `internal/domains/configmanager/variable/enum.go`
- `internal/domains/configmanager/variable/list.go`
- `internal/domains/configmanager/variable/get.go`
- `internal/domains/configmanager/variable/apply.go`
- `internal/domains/configmanager/variable/delete.go`
- `internal/domains/configmanager/variable/resolve.go`
- `internal/domains/configmanager/secret/tools.go`
- `internal/domains/configmanager/secret/enum.go`
- `internal/domains/configmanager/secret/list.go`
- `internal/domains/configmanager/secret/get.go`
- `internal/domains/configmanager/secret/apply.go`
- `internal/domains/configmanager/secret/delete.go`
- `internal/domains/configmanager/secretversion/tools.go`
- `internal/domains/configmanager/secretversion/create.go`
- `internal/domains/configmanager/secretversion/list.go`

**Files Modified:**
- `internal/server/server.go` - Tool registration + count 41→52

**Verification:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅

---

### ✅ COMPLETED: Phase 3A — Audit / Version History (2026-02-28)

**Added 3 MCP tools for resource version history and change tracking, expanding the server from 52 to 55 tools. Third domain outside the infrahub bounded context.**

**What was delivered:**

1. **`list_resource_versions`** — Paginated version history via `ApiResourceVersionQueryController.ListByFilters`
   - `internal/domains/audit/list.go` — Kind + resource ID scoping with 1-based-to-0-based pagination
   - Requires both `resource_id` and `kind` to scope the query

2. **`get_resource_version`** — Full version with diff via `ApiResourceVersionQueryController.GetByIdWithContextSize`
   - `internal/domains/audit/get.go` — Returns YAML states, unified diff, event type, stack job link
   - `context_size` parameter controls diff context (default 3, like `git diff -U3`)
   - Uses `GetByIdWithContextSize` RPC (supersedes plain `Get` — strictly more useful)

3. **`get_resource_version_count`** — Lightweight count via `ApiResourceVersionQueryController.GetCount`
   - `internal/domains/audit/count.go` — Returns version count for quick "has anything changed?" checks

4. **Dynamic `ApiResourceKind` enum resolver**
   - `internal/domains/audit/enum.go` — `resolveApiResourceKind()` accepts any valid kind string
   - First cross-resource-type domain — unlike infrachart/variable/secret that hard-code their kind

5. **Tool definitions and handlers**
   - `internal/domains/audit/tools.go` — 3 input structs, 3 tool defs, 3 typed handlers

6. **Server registration**
   - `internal/server/server.go` — Import + 3 `mcp.AddTool` calls, count 52→55

**Key Decisions Made:**
- Merge `Get` and `GetByIdWithContextSize` into single `get_resource_version` tool — `GetByIdWithContextSize` is strictly more useful
- Default `context_size` to 3 — matches standard unified diff convention (`-U3`)
- Both `resource_id` and `kind` are required for list and count — backend needs both to scope the query
- `joinEnumValues` duplicated locally (Option A from plan) — avoids touching unrelated packages; cleanup refactor deferred

**Files Created:**
- `internal/domains/audit/doc.go`
- `internal/domains/audit/enum.go`
- `internal/domains/audit/tools.go`
- `internal/domains/audit/list.go`
- `internal/domains/audit/get.go`
- `internal/domains/audit/count.go`

**Files Modified:**
- `internal/server/server.go` — Tool registration + count 52→55

**Verification:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅

---

## Objectives for Next Conversations

### Option A (Recommended): Phase 3B — StackJob Commands / Lifecycle Control (3 tools)

Retry, cancel, and pre-validate stack jobs. Would bring the server to 58 tools.

### Option B: Phase 3C — Deployment Component & IaC Module Catalog (3 tools)

Browse cloud resource types and IaC modules. Would bring the server to 58 tools.

---

## Essential Files to Review

### 1. Latest Checkpoint (if exists)
```
_projects/2026-02/20260227.02.expand-infrahub-mcp-tools/checkpoints/
```

### 2. Master Plan
```
_projects/2026-02/20260227.02.expand-infrahub-mcp-tools/tasks/T01_0_plan.md
```

### 3. Plans
```
_projects/2026-02/20260227.02.expand-infrahub-mcp-tools/plans/
```

### 4. Design Decisions
```
_projects/2026-02/20260227.02.expand-infrahub-mcp-tools/design-decisions/
```
- AD-01: Exclude credential management
- AD-02: Restructure gen code by domain (implemented)

### 5. Patterns to Follow
Existing domain implementations to use as reference:
- `internal/domains/infrahub/cloudresource/` — full CRUD + search (11 tools)
- `internal/domains/infrahub/infrachart/` — list + get + two-step build (3 tools)
- `internal/domains/infrahub/infrapipeline/` — pipeline observability + control + gate resolution (7 tools)
- `internal/domains/infrahub/infraproject/` — full lifecycle: search, get, apply, delete, slug, undeploy (6 tools)
- `internal/domains/infrahub/stackjob/` — read-only query tools (3 tools)
- `internal/domains/infrahub/preset/` — search + get pair (2 tools)
- `internal/domains/graph/` — dependency intelligence + impact analysis (7 tools, first non-infrahub bounded context)
- `internal/domains/configmanager/variable/` — variable CRUD + resolve (5 tools, scope-aware identification)
- `internal/domains/configmanager/secret/` — secret metadata CRUD (4 tools, scope-aware identification)
- `internal/domains/configmanager/secretversion/` — version create + list (2 tools, write-only security boundary)
- `internal/domains/audit/` — resource version history + change tracking (3 tools, third non-infrahub bounded context, dynamic ApiResourceKind resolver)

---

## Resume Checklist

When starting a new session:

1. [ ] Read this file for status and agreed next steps
2. [ ] Check `plans/README.md` for active plans
3. [ ] Review `design-decisions/` for architectural context
4. [ ] Review `coding-guidelines/`, `wrong-assumptions/`, `dont-dos/`
5. [ ] Present status summary and wait for direction

---

*This file provides direct paths to all project resources for quick context loading.*
