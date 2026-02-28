# Next Task: 20260228.01.servicehub-mcp-tools

## RULES OF ENGAGEMENT - READ FIRST

**When this file is loaded in a new conversation, the AI MUST:**

1. **DO NOT AUTO-EXECUTE** - Never start implementing without explicit user approval
2. **GATHER CONTEXT SILENTLY** - Read project files without outputting
3. **PRESENT STATUS SUMMARY** - Show what's done, what's pending, agreed next steps
4. **SHOW OPTIONS** - List recommended and alternative actions
5. **WAIT FOR DIRECTION** - Do NOT proceed until user explicitly says "go" or chooses an option

---

## Quick Resume Instructions

Drop this file into your conversation to quickly resume work on this project.

## Project: 20260228.01.servicehub-mcp-tools

**Description**: Add MCP tools for the ServiceHub domain — Service, Pipeline, VariablesGroup, SecretsGroup, DnsDomain, TektonPipeline, and TektonTask API resources.
**Goal**: Implement 35 MCP tools across 7 ServiceHub bounded contexts (Service, Pipeline, VariablesGroup, SecretsGroup, DnsDomain, TektonPipeline, TektonTask), following the existing infrahub tool patterns.
**Tech Stack**: Go/gRPC/MCP
**Components**: internal/domains/servicehub/, internal/server/server.go

## Essential Files to Review

### 1. Latest Checkpoint (if exists)
Check for the most recent checkpoint file:
```
_projects/2026-02/20260228.01.servicehub-mcp-tools/checkpoints/
```

### 2. Current Task
Review the current task status and plan:
```
_projects/2026-02/20260228.01.servicehub-mcp-tools/tasks/
```

### 3. Plans
Review implementation plans and their status:
```
_projects/2026-02/20260228.01.servicehub-mcp-tools/plans/
```

### 4. Project Documentation
- **README**: `_projects/2026-02/20260228.01.servicehub-mcp-tools/README.md`

## Knowledge Folders to Check

### Design Decisions
```
_projects/2026-02/20260228.01.servicehub-mcp-tools/design-decisions/
```
Review architectural and strategic choices made for this project.

### Coding Guidelines
```
_projects/2026-02/20260228.01.servicehub-mcp-tools/coding-guidelines/
```
Check project-specific patterns and conventions established.

### Wrong Assumptions
```
_projects/2026-02/20260228.01.servicehub-mcp-tools/wrong-assumptions/
```
Review misconceptions discovered to avoid repeating them.

### Don't Dos
```
_projects/2026-02/20260228.01.servicehub-mcp-tools/dont-dos/
```
Check anti-patterns and failed approaches to avoid.

## Resume Checklist

When starting a new session:

1. [ ] Read the latest checkpoint (if any) from `_projects/2026-02/20260228.01.servicehub-mcp-tools/checkpoints/`
2. [ ] Check current task status in `_projects/2026-02/20260228.01.servicehub-mcp-tools/tasks/`
3. [ ] Review plans in `_projects/2026-02/20260228.01.servicehub-mcp-tools/plans/`
4. [ ] Review any new design decisions in `_projects/2026-02/20260228.01.servicehub-mcp-tools/design-decisions/`
5. [ ] Check coding guidelines in `_projects/2026-02/20260228.01.servicehub-mcp-tools/coding-guidelines/`
6. [ ] Review lessons learned in `_projects/2026-02/20260228.01.servicehub-mcp-tools/wrong-assumptions/` and `_projects/2026-02/20260228.01.servicehub-mcp-tools/dont-dos/`
7. [ ] Continue with the next task or complete the current one

## Current Status

**Created**: 2026-02-28 18:12
**Current Task**: Tier 2 — Pipeline tools (9 tools)
**Status**: Tier 1 completed, ready for Tier 2

**Current step:**
- Completed T01 planning (architecture and tool catalogue for all 35 tools)
- Completed Tier 1 — Service tools (7 tools) (2026-02-28)
  - search_services, get_service, apply_service, delete_service
  - disconnect_service_git_repo, configure_service_webhook, list_service_branches
  - Wired into server.go, clean build verified
- Next: **Tier 2 — Pipeline tools** (9 tools: list, get, get_last, run, rerun, cancel, gate, files, update_file)

### Completed: Tier 1 — Service Tools (2026-02-28)

**Implemented 7 MCP tools for the ServiceHub Service entity.**

**What was delivered:**

1. **New package `internal/domains/servicehub/service/`** — 8 Go files
   - `register.go` — Register function wiring all 7 tools
   - `tools.go` — Input structs, tool definitions, handlers, validateIdentification
   - `search.go` — Search via generic ApiResourceSearchQueryController.searchByKind
   - `get.go` — Get, resolveService, resolveServiceID, describeService
   - `apply.go` — Apply via protojson unmarshal + ServiceCommandController.Apply
   - `delete.go` — Delete via resolveServiceID + ApiResourceDeleteInput
   - `disconnect.go` — DisconnectGitRepo via ServiceCommandController
   - `webhook.go` — ConfigureWebhook via ServiceCommandController
   - `branches.go` — ListBranches via ServiceQueryController

2. **Server wiring** — `internal/server/server.go` updated with `servicehubservice.Register`

**Key Decisions Made:**
- Used generic `ApiResourceSearchQueryController.searchByKind` for search (no dedicated Service search RPC exists)
- Skipped client-side `cloud_object` validation in apply — follows thin-client pattern, lets backend validate
- Import alias `servicehubservice` to avoid collision with `service` keyword

**Files Changed/Created:**
- `internal/domains/servicehub/service/register.go` — New
- `internal/domains/servicehub/service/tools.go` — New
- `internal/domains/servicehub/service/search.go` — New
- `internal/domains/servicehub/service/get.go` — New
- `internal/domains/servicehub/service/apply.go` — New
- `internal/domains/servicehub/service/delete.go` — New
- `internal/domains/servicehub/service/disconnect.go` — New
- `internal/domains/servicehub/service/webhook.go` — New
- `internal/domains/servicehub/service/branches.go` — New
- `internal/server/server.go` — Modified (added import + Register call)

---

## Objectives for Next Conversations

### Option A (Recommended): Tier 2 — Pipeline Tools (9 tools)
Highest operational value. Implements list_pipelines, get_pipeline, get_last_pipeline, run_pipeline, rerun_pipeline, cancel_pipeline, resolve_pipeline_gate, list_pipeline_files, update_pipeline_file.

### Option B: Tier 3 — VariablesGroup + SecretsGroup (12 tools)
Configuration management. Two entities with symmetric API surface (get, apply, delete, upsert_entry, delete_entry, get_value each).

### Option C: Tier 4+5 — DnsDomain + TektonPipeline + TektonTask (7 tools)
Quick wins. Simple CRUD entities with 2-3 tools each.

## Quick Commands

After loading context:
- "Continue with Tier 2" - Start Pipeline tools implementation
- "Show project status" - Get overview of progress
- "Create checkpoint" - Save current progress
- "Review guidelines" - Check established patterns

---

*This file provides direct paths to all project resources for quick context loading.*
