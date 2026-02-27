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

**Current Task**: Phase 1B (InfraProject tools)
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
- 🔵 Next: **Phase 1B: InfraProject tools** (6 tools: search, get, apply, delete, slug check, undeploy)

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

## Objectives for Next Conversations

### Option A (Recommended): Phase 1B — InfraProject (6 tools)

The natural continuation. InfraProject tools give agents full project lifecycle management.

| Tool | Backend RPC | Purpose |
|------|-------------|---------|
| `search_infra_projects` | `InfraHubSearchQueryController.searchInfraProjects` | Search projects by org, env, text |
| `get_infra_project` | `InfraProjectQueryController.get` / `getByOrgBySlug` | Retrieve full project |
| `apply_infra_project` | `InfraProjectCommandController.apply` | Create or update |
| `delete_infra_project` | `InfraProjectCommandController.delete` | Remove project |
| `check_infra_project_slug` | `InfraProjectQueryController.checkSlugAvailability` | Slug uniqueness check |
| `undeploy_infra_project` | `InfraProjectCommandController.undeploy` | Tear down all deployed resources |

Files to create: `internal/domains/infrahub/infraproject/` (tools.go, search.go, get.go, apply.go, delete.go, slug.go, undeploy.go)

### Option B: Phase 1C — InfraPipeline (5 tools)

Pipeline observability and control (list, get, get latest, run, cancel).

### Option C: Phase 1B + 1C combined

If scope allows, tackle both InfraProject and InfraPipeline in one session.

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
