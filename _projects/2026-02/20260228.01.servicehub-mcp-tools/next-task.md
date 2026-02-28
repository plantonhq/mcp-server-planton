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
**Current Task**: PROJECT COMPLETE
**Status**: All 5 tiers completed (37 tools across 7 bounded contexts). Ready for documentation update.

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
- ✅ **Completed Tier 3 — VariablesGroup + SecretsGroup (16 tools)** (2026-02-28)
  - 8 VariablesGroup tools: search_variables, get_variables_group, apply_variables_group, delete_variables_group, upsert_variable, delete_variable, get_variable_value, transform_variables
  - 8 SecretsGroup tools: search_secrets, get_secrets_group, apply_secrets_group, delete_secrets_group, upsert_secret, delete_secret, get_secret_value, transform_secrets
  - 4 tools added beyond original plan: search_variables, search_secrets, transform_variables, transform_secrets (discovered dedicated RPCs during proto review)
  - 7 design decisions confirmed: DD-T3-1 through DD-T3-7
  - Wired into server.go, clean build verified
- ✅ **Completed Tier 4+5 — DnsDomain + TektonPipeline + TektonTask (9 tools)** (2026-03-01)
  - 3 DnsDomain tools: get_dns_domain, apply_dns_domain, delete_dns_domain
  - 3 TektonPipeline tools: get_tekton_pipeline, apply_tekton_pipeline, delete_tekton_pipeline
  - 3 TektonTask tools: get_tekton_task, apply_tekton_task, delete_tekton_task
  - 2 tools added beyond original plan: delete_tekton_pipeline, delete_tekton_task (added for complete CRUD surface)
  - 4 design decisions confirmed: DD-T4-1 (slug normalization), DD-T4-2 (entity-based Tekton delete), DD-T4-3 (no search tools), DD-T4-4 (separate bounded contexts)
  - Wired into server.go, clean build verified
- ✅ **PROJECT COMPLETE** — All 37 tools across 7 ServiceHub bounded contexts implemented

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

### ✅ COMPLETED: Tier 3 — VariablesGroup + SecretsGroup (16 tools) (2026-02-28)

**Added 16 MCP tools for ServiceHub configuration management — VariablesGroup and SecretsGroup entities with group-level CRUD, entry-level mutations, value resolution, search, and config key transformation.**

**What was delivered:**

1. **New package `internal/domains/servicehub/variablesgroup/`** — 10 Go files
   - `register.go` — Register function wiring all 8 tools
   - `tools.go` — 8 input structs, 8 Tool/Handler pairs, validateIdentification, validateGroupIdentification
   - `search.go` — Entry-level search via ServiceHubSearchQueryController.SearchVariables
   - `get.go` — Get + resolveGroup, resolveGroupID, describeGroup helpers
   - `apply.go` — Apply via protojson unmarshal + VariablesGroupCommandController.Apply
   - `delete.go` — Delete via resolveGroupID + ApiResourceDeleteInput
   - `upsert_entry.go` — UpsertEntry via protojson entry conversion + dual-path group identification
   - `delete_entry.go` — DeleteEntry via dual-path group identification
   - `get_value.go` — GetValue with StringValue unwrapping
   - `transform.go` — Transform batch reference resolution

2. **New package `internal/domains/servicehub/secretsgroup/`** — 10 Go files (mirrors variablesgroup with secrets-specific naming and security warnings)

3. **Server wiring** — `internal/server/server.go` updated with `servicehubvariablesgroup.Register` and `servicehubsecretsgroup.Register`

**Key Decisions Made:**
- DD-T3-1: Entry-level search via dedicated `ServiceHubSearchQueryController` RPCs (not generic searchByKind)
- DD-T3-2: Dual-path identification (ID or org+slug) on all group-level tools
- DD-T3-3: Entry mutation tools accept group_id OR org+group_slug (resolves internally)
- DD-T3-4: Entry input as nested JSON object with protojson unmarshalling
- DD-T3-5: StringValue unwrapping for get_value tools (plain text, not JSON-wrapped)
- DD-T3-6: `get_secret_value` and `transform_secrets` include plaintext security warnings
- DD-T3-7: No shared abstraction between packages — separate bounded contexts

**Files Created:**
- `internal/domains/servicehub/variablesgroup/` — 10 new files (register, tools, search, get, apply, delete, upsert_entry, delete_entry, get_value, transform)
- `internal/domains/servicehub/secretsgroup/` — 10 new files (same structure)

**Files Modified:**
- `internal/server/server.go` — Added 2 imports + 2 Register calls

**Verification:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅

### ✅ COMPLETED: Tier 4+5 — DnsDomain + TektonPipeline + TektonTask (9 tools) (2026-03-01)

**Added 9 MCP tools for DNS domain management, Tekton pipeline templates, and Tekton task templates — completing the ServiceHub MCP tools project.**

**What was delivered:**

1. **New package `internal/domains/servicehub/dnsdomain/`** — 5 Go files
   - `register.go` — Register function wiring all 3 tools
   - `tools.go` — 3 input structs, 3 Tool/Handler pairs, validateIdentification
   - `get.go` — Get, resolveDomain, resolveDomainID, describeDomain (DnsDomainId + ApiResourceByOrgBySlugRequest)
   - `apply.go` — Apply via protojson unmarshal + DnsDomainCommandController.Apply
   - `delete.go` — Delete via resolveDomainID + ApiResourceDeleteInput

2. **New package `internal/domains/servicehub/tektonpipeline/`** — 5 Go files
   - `register.go` — Register function wiring all 3 tools
   - `tools.go` — 3 input structs, 3 Tool/Handler pairs, validateIdentification
   - `get.go` — Get, resolvePipeline, describePipeline (ApiResourceId + GetByOrgAndNameInput with slug→Name mapping)
   - `apply.go` — Apply via protojson unmarshal + TektonPipelineCommandController.Apply
   - `delete.go` — Delete via get-then-delete pattern (RPC takes full entity)

3. **New package `internal/domains/servicehub/tektontask/`** — 5 Go files (mirrors tektonpipeline structure)

4. **Server wiring** — `internal/server/server.go` updated with `servicehubdnsdomain.Register`, `servicehubtektonpipeline.Register`, `servicehubtektontask.Register`

**Key Decisions Made:**
- DD-T4-1: Normalized all dual-path lookup to `slug` for consistency with existing tools (TektonPipeline/TektonTask `name` field presented as `slug` in MCP input)
- DD-T4-2: Added delete tools for TektonPipeline/TektonTask (originally excluded; uses get-then-delete pattern since RPC takes full entity)
- DD-T4-3: No search tools (no search RPCs exist for these entities)
- DD-T4-4: Separate packages for TektonPipeline and TektonTask despite near-identical structure (per DD-T3-7 precedent)

**Files Created:**
- `internal/domains/servicehub/dnsdomain/` — 5 new files (register, tools, get, apply, delete)
- `internal/domains/servicehub/tektonpipeline/` — 5 new files (register, tools, get, apply, delete)
- `internal/domains/servicehub/tektontask/` — 5 new files (register, tools, get, apply, delete)

**Files Modified:**
- `internal/server/server.go` — Added 3 imports + 3 Register calls

**Verification:** `go build ./...` ✅ | `go vet ./...` ✅ | `go test ./...` ✅

---

## Project Complete

All 37 ServiceHub MCP tools have been implemented across 7 bounded contexts:

| Tier | Entity | Tools | Status |
|------|--------|-------|--------|
| 1 | Service | 7 | ✅ Complete |
| 2 | Pipeline | 9 | ✅ Complete |
| 3a | VariablesGroup | 8 | ✅ Complete |
| 3b | SecretsGroup | 8 | ✅ Complete |
| 4 | DnsDomain | 3 | ✅ Complete |
| 5a | TektonPipeline | 3 | ✅ Complete |
| 5b | TektonTask | 3 | ✅ Complete |
| | **Total** | **37** (originally planned 35 + 2 discovered tools) | |

## Objectives for Next Conversations

### Option A (Recommended): Update project documentation and tools reference
Update README.md with the new tool count (63 → 100 tools), update any docs to reflect the full ServiceHub tool surface.

### Option B: Wrap up project and archive
Close out the project, move to archived status, finalize all documentation.

## Quick Commands

After loading context:
- "Update documentation" - Update README.md and tools reference
- "Show project status" - Get overview of progress
- "Archive project" - Mark project as complete and archive

---

*This file provides direct paths to all project resources for quick context loading.*
