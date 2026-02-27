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

**Current Task**: Phase 1C (InfraPipeline tools)
**Status**: Not started — ready to plan

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
- 🔵 Next: **Phase 1C: InfraPipeline tools** (5 tools: list, get, get latest, run, cancel)

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

## Objectives for Next Conversations

### Option A (Recommended): Phase 1C — InfraPipeline (5 tools)

The natural continuation. Completes the Phase 1 trifecta (Chart + Project + Pipeline) giving agents full composition-to-deployment observability.

| Tool | Backend RPC | Purpose |
|------|-------------|---------|
| `list_infra_pipelines` | `InfraPipelineQueryController.listByFilters` | List pipelines by project, status |
| `get_infra_pipeline` | `InfraPipelineQueryController.get` | Full pipeline status and details |
| `get_latest_infra_pipeline` | `InfraPipelineQueryController.getLastInfraPipelineByInfraProjectId` | Last pipeline for a project |
| `run_infra_pipeline` | `InfraPipelineCommandController.runInfraProjectChartSourcePipeline` | Trigger pipeline run |
| `cancel_infra_pipeline` | `InfraPipelineCommandController.cancel` | Cancel a running pipeline |

Files to create: `internal/domains/infrahub/infrapipeline/` (tools.go, list.go, get.go, latest.go, run.go, cancel.go)

### Option B: Phase 2A — Graph / Dependency Intelligence (4 tools)

The "wow factor" differentiator — impact analysis, dependency graphs, org topology.

### Option C: Phase 1C + Phase 2A combined

If scope allows, tackle both Pipeline and Graph in one session (9 tools, reaching 36 total).

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
- `internal/domains/infrahub/infraproject/` — full lifecycle: search, get, apply, delete, slug, undeploy (6 tools)
- `internal/domains/infrahub/stackjob/` — read-only query tools (3 tools)
- `internal/domains/infrahub/preset/` — search + get pair (2 tools)

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
