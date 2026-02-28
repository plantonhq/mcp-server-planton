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
**Current Task**: Tier 3 — VariablesGroup + SecretsGroup (12 tools)
**Status**: Tier 1 and Tier 2 completed, ready for Tier 3

**Current step:**
- ✅ Completed T01 planning (architecture and tool catalogue for all 35 tools)
- ✅ Completed Tier 1 — Service tools (7 tools) (2026-02-28)
  - search_services, get_service, apply_service, delete_service
  - disconnect_service_git_repo, configure_service_webhook, list_service_branches
  - Wired into server.go, clean build verified
- ✅ **Completed Tier 2 — Pipeline tools (9 tools)** (2026-02-28)
  - list_pipelines, get_pipeline, get_last_pipeline, run_pipeline, rerun_pipeline
  - cancel_pipeline, resolve_pipeline_gate, list_pipeline_files, update_pipeline_file
  - 3 design decisions confirmed: DD-T2-1 (branch required), DD-T2-2 (bytes-to-string decode), DD-T2-3 (single gate tool)
  - Custom marshaling for pipeline files (bytes→UTF-8), success message for run (Empty response)
  - Wired into server.go, clean build verified
- Next: **Tier 3 — VariablesGroup + SecretsGroup** (12 tools)

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

### ✅ COMPLETED: Tier 2 — Pipeline Tools (2026-02-28)

**Added 9 MCP tools for ServiceHub Pipeline observability, lifecycle control, gate resolution, and repository pipeline file management.**

**What was delivered:**

1. **New package `internal/domains/servicehub/pipeline/`** — 11 Go files
   - `register.go` — Register function wiring all 9 tools
   - `tools.go` — 9 input structs, 9 Tool/Handler pairs
   - `get.go` — Get pipeline by ID via PipelineQueryController.Get
   - `list.go` — List pipelines with org/service/envs filters via PipelineQueryController.ListByFilters
   - `latest.go` — Most recent pipeline for a service via PipelineQueryController.GetLastPipelineByServiceId
   - `run.go` — Trigger pipeline via PipelineCommandController.RunGitCommit (branch required, commit_sha optional)
   - `rerun.go` — Re-run pipeline via PipelineCommandController.Rerun
   - `cancel.go` — Cancel running pipeline via PipelineCommandController.Cancel
   - `gate.go` — Resolve manual gate via PipelineCommandController.ResolveManualGate + resolveDecision helper
   - `files.go` — List Tekton pipeline files with custom bytes→UTF-8 marshaling
   - `update_file.go` — Update pipeline file with string→bytes encoding + optimistic locking

2. **Server wiring** — `internal/server/server.go` updated with `servicehubpipeline.Register`

**Key Decisions Made:**
- DD-T2-1: `run_pipeline` requires `branch` (proto required field), `commit_sha` optional; RPC returns Empty so tool returns success message directing agent to `get_last_pipeline`
- DD-T2-2: Custom marshaling for pipeline files — decode `bytes` content to plain UTF-8 string for agent readability (deviates from standard `domains.MarshalJSON` pattern)
- DD-T2-3: Single `resolve_pipeline_gate` tool with `deployment_task_name` (unlike InfraPipeline's two gate tools)

**Files Created:**
- `internal/domains/servicehub/pipeline/register.go` — New
- `internal/domains/servicehub/pipeline/tools.go` — New
- `internal/domains/servicehub/pipeline/get.go` — New
- `internal/domains/servicehub/pipeline/list.go` — New
- `internal/domains/servicehub/pipeline/latest.go` — New
- `internal/domains/servicehub/pipeline/run.go` — New
- `internal/domains/servicehub/pipeline/rerun.go` — New
- `internal/domains/servicehub/pipeline/cancel.go` — New
- `internal/domains/servicehub/pipeline/gate.go` — New
- `internal/domains/servicehub/pipeline/files.go` — New
- `internal/domains/servicehub/pipeline/update_file.go` — New

**Files Modified:**
- `internal/server/server.go` — Added import + Register call (2 lines)

**Verification:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅

---

## Objectives for Next Conversations

### Option A (Recommended): Tier 3 — VariablesGroup + SecretsGroup (12 tools)
Configuration management. Two entities with symmetric API surface (get, apply, delete, upsert_entry, delete_entry, get_value each).

### Option B: Tier 4+5 — DnsDomain + TektonPipeline + TektonTask (7 tools)
Quick wins. Simple CRUD entities with 2-3 tools each.

## Quick Commands

After loading context:
- "Continue with Tier 3" - Start VariablesGroup + SecretsGroup tools implementation
- "Show project status" - Get overview of progress
- "Create checkpoint" - Save current progress
- "Review guidelines" - Check established patterns

---

*This file provides direct paths to all project resources for quick context loading.*
