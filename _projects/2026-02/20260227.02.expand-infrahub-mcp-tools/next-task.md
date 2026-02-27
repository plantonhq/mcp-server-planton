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

**Current Task**: Phase 1A (InfraChart tools)
**Status**: Not started — ready to plan

**Current step:**
- ✅ **Phase 0: Restructure Generated Code** (2026-02-27)
  - Moved gen/cloudresource/ -> gen/infrahub/cloudresource/ (362 providers, 17 clouds)
  - Updated code generator, Makefile, consumer import
  - Build, tests, vet all pass clean
  - Committed: `refactor(codegen): move generated code to gen/infrahub/cloudresource/`
- 🔵 Next: **Phase 1A: InfraChart tools** (3 tools: search, get, build)

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

## Objectives for Next Conversations

### Option A (Recommended): Phase 1A — InfraChart (3 tools)

The highest-impact next step. InfraChart tools let agents discover, inspect, and preview infrastructure chart templates.

| Tool | Backend RPC | Purpose |
|------|-------------|---------|
| `search_infra_charts` | `InfraChartQueryController.find` | Browse available chart templates |
| `get_infra_chart` | `InfraChartQueryController.get` | Retrieve full chart details |
| `build_infra_chart` | `InfraChartQueryController.build` | Preview rendered output |

Files to create: `internal/domains/infrahub/infrachart/` (tools.go, search.go, get.go, build.go)

### Option B: Phase 1B — InfraProject (6 tools)

If InfraChart is blocked or lower priority, InfraProject provides full project lifecycle management.

### Option C: Phase 1A + 1B combined

If scope allows, tackle both InfraChart and InfraProject in one session.

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
