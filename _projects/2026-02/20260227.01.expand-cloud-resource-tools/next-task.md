# Next Task: 20260227.01.expand-cloud-resource-tools

## ⚠️ RULES OF ENGAGEMENT - READ FIRST ⚠️

**When this file is loaded in a new conversation, the AI MUST:**

1. **DO NOT AUTO-EXECUTE** - Never start implementing without explicit user approval
2. **GATHER CONTEXT SILENTLY** - Read project files without outputting
3. **PRESENT STATUS SUMMARY** - Show what's done, what's pending, agreed next steps
4. **SHOW OPTIONS** - List recommended and alternative actions
5. **WAIT FOR DIRECTION** - Do NOT proceed until user explicitly says "go" or chooses an option

---

## Quick Resume Instructions

Drop this file into your conversation to quickly resume work on this project.

## Project: 20260227.01.expand-cloud-resource-tools

**Description**: Expand the MCP server's cloud resource tool surface from the current 3 tools (apply, get, delete) to 18 tools covering the full lifecycle: listing/search, infrastructure destroy, stack job observability, org/env discovery, slug validation, presets, locks, rename, env var maps, and cross-resource reference resolution.
**Goal**: Give AI agents full autonomous capability over cloud resource lifecycle — from discovering their operating context (orgs, environments) through CRUD operations to observing provisioning outcomes (stack jobs) and managing operational concerns (locks, presets, references).
**Tech Stack**: Go/gRPC/MCP
**Components**: internal/domains/cloudresource/, internal/domains/stackjob/, internal/domains/organization/, internal/domains/environment/, internal/domains/preset/, internal/server/server.go

## Essential Files to Review

### 1. Latest Checkpoint (if exists)
Check for the most recent checkpoint file:
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/checkpoints/
```

### 2. Current Task
Review the current task status and plan:
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/tasks/
```

### 3. Project Documentation
- **README**: `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/README.md`

## Knowledge Folders to Check

### Design Decisions
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/design-decisions/
```
Review architectural and strategic choices made for this project.

### Coding Guidelines
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/coding-guidelines/
```
Check project-specific patterns and conventions established.

### Wrong Assumptions
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/wrong-assumptions/
```
Review misconceptions discovered to avoid repeating them.

### Don't Dos
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/dont-dos/
```
Check anti-patterns and failed approaches to avoid.

## Resume Checklist

When starting a new session:

1. [ ] Read the latest checkpoint (if any) from `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/checkpoints/`
2. [ ] Check current task status in `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/tasks/`
3. [ ] Review any new design decisions in `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/design-decisions/`
4. [ ] Check coding guidelines in `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/coding-guidelines/`
5. [ ] Review lessons learned in `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/wrong-assumptions/` and `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-02/20260227.01.expand-cloud-resource-tools/dont-dos/`
6. [ ] Continue with the next task or complete the current one

## Current Status

**Created**: 2026-02-27
**Current Task**: Hardening (docs complete, H4 refactor pending)
**Status**: All 5 feature phases (6A–6E) and documentation hardening complete. 18 tools implemented and documented.
**Last Session**: 2026-02-27

### Session Progress (2026-02-27, Session 6)

**Hardening: Documentation and Tool Description Review — DONE**

Completed the documentation hardening pass (H2 + H5 + surprise H3):

**What was delivered:**

1. **Tool description normalization** — "Planton platform" → "Planton Cloud" across 9 descriptions, added cross-tool references for agent workflow, added org discovery hints to `org` jsonschema tags
2. **README.md update** — Intro paragraph expanded, tool table from 3 to 18 tools across 4 grouped sections
3. **docs/tools.md rewrite** — 229 → ~470 lines, 15 new tool sections, Resource Identification Pattern section, expanded Agent Cheat Sheet with 6 decision guides
4. **docs/development.md update** — Project structure and test files table updated with all new domain packages

**Surprise discovery:** `docs/tools.md` was the most important update but wasn't in the original Hardening plan. Added to scope.

**Files modified:**
- `internal/domains/cloudresource/tools.go` — description string normalization (+29 lines)
- `internal/domains/stackjob/tools.go` — description string normalization (+4 lines)
- `README.md` — intro paragraph + tool table expansion (+42 lines)
- `docs/tools.md` — full rewrite from 3 to 18 tools (+456 lines)
- `docs/development.md` — project structure + test table (+14 lines)

**Verification:** `go build ./...` and `go test ./...` both pass. Zero linter errors.

---

### Session Progress (2026-02-27, Session 5)

**Phase 6E: Advanced Operations — DONE**

Implemented 5 new tools (`list_cloud_resource_locks`, `remove_cloud_resource_locks`, `rename_cloud_resource`, `get_env_var_map`, `resolve_value_references`) in the existing `cloudresource/` domain package, expanding the MCP server from 13 to 18 tools. This completes the full planned tool surface.

**Proto surprises discovered and resolved:**
- `get_env_var_map`: plan assumed `id` + `manifest` (as map), but actual proto takes `yaml_content` (raw YAML string). No `id` field — the server extracts identity from the YAML for authorization. Implemented as-is, matching the proto.
- `resolve_value_references`: plan assumed `cloud_resource_id` + `references` (list), but actual proto takes `cloud_resource_kind` (enum) + `cloud_resource_id` (string). No `references` list — the server resolves ALL valueFrom references automatically. Implemented with `kind` always required.

**Design decisions made:**
- Dual-path (ID or kind/org/env/slug) for locks, remove locks, and rename — consistent with all other cloudresource tools that identify a resource
- Dual-path for `resolve_value_references` with `kind` always required — custom handler validation since `validateIdentifier` treats kind as slug-path-only, but this tool's RPC always needs kind
- `get_env_var_map` takes raw `yaml_content` string only — no dual-path, no wrapper, matches proto directly
- MCP input field `new_name` (not `name`) for rename tool — clearer agent-facing semantics vs proto's generic `name`

**Files created:**
- `internal/domains/cloudresource/locks.go` — `ListLocks` and `RemoveLocks` domain functions using `CloudResourceLockControllerClient`
- `internal/domains/cloudresource/rename.go` — `Rename` domain function using `CloudResourceCommandControllerClient`
- `internal/domains/cloudresource/envvarmap.go` — `GetEnvVarMap` domain function using `CloudResourceQueryControllerClient`
- `internal/domains/cloudresource/references.go` — `ResolveValueReferences` domain function using `CloudResourceQueryControllerClient`

**Files modified:**
- `internal/domains/cloudresource/tools.go` — 5 new input structs, tool defs, handlers; package doc comment updated (6 → 11 tools)
- `internal/server/server.go` — 5 new `mcp.AddTool` calls, tool count updated (13 → 18)

**Verification:** `go build ./...` and `go test ./...` both pass. Zero linter errors.

### Session Progress (2026-02-27, Session 4)

**Phase 6D: Agent Quality-of-Life — DONE**

Implemented 3 new tools (`check_slug_availability`, `search_cloud_object_presets`, `get_cloud_object_preset`) and a new domain package (`preset/`), expanding the MCP server from 10 to 13 tools. Also performed a cross-cutting refactor to extract `resolveKind` to the shared `domains` package.

**Files created:**
- `internal/domains/kind.go` — Shared `ResolveKind` function extracted from duplicated copies across domains
- `internal/domains/kind_test.go` — Tests for `ResolveKind` (known, unknown, empty)
- `internal/domains/cloudresource/slug.go` — `CheckSlugAvailability` domain function calling `CloudResourceQueryController.CheckSlugAvailability`
- `internal/domains/preset/tools.go` — Package doc comment (proto provenance for both `InfraHubSearchQueryController` and `CloudObjectPresetQueryController`), input structs, tool defs, handlers for search + get
- `internal/domains/preset/search.go` — `Search` domain function with org/official branching logic, kind resolution
- `internal/domains/preset/get.go` — `Get` domain function calling `CloudObjectPresetQueryController.Get` via `ApiResourceId`

**Files modified:**
- `internal/domains/cloudresource/tools.go` — Added `CheckSlugAvailabilityInput`, `CheckSlugAvailabilityTool()`, `CheckSlugAvailabilityHandler()`, updated package doc comment (5 → 6 tools)
- `internal/domains/cloudresource/kind.go` — Removed local `resolveKind`, updated `resolveKinds` to use `domains.ResolveKind`
- `internal/domains/cloudresource/kind_test.go` — Removed `TestResolveKind_*` tests (moved to `domains/kind_test.go`)
- `internal/domains/cloudresource/apply.go` — `resolveKind` → `domains.ResolveKind`
- `internal/domains/cloudresource/get.go` — `resolveKind` → `domains.ResolveKind`
- `internal/domains/cloudresource/identifier.go` — `resolveKind` → `domains.ResolveKind` (2 call sites)
- `internal/domains/stackjob/enum.go` — Removed local `resolveKind` and unused `cloudresourcekind` import
- `internal/domains/stackjob/enum_test.go` — Removed `TestResolveKind_*` tests (moved to `domains/kind_test.go`)
- `internal/domains/stackjob/list.go` — `resolveKind` → `domains.ResolveKind`
- `internal/server/server.go` — Imported `preset` package, registered 3 new tools (count 10 → 13)

**Design decisions made:**
- Extracted `resolveKind` to shared `domains.ResolveKind` — 3-way duplication threshold crossed; `resolveKind` depends only on the `cloudresourcekind` proto package (leaf dependency), not on any domain; moving to `domains` doesn't create cross-domain coupling
- Preset search branching: when `org` is provided, both `IsIncludeOfficial` and `IsIncludeOrganizationPresets` set to `true` — agent always gets the most useful result set without needing to understand internal flags
- `providers` filter not exposed — `kind` already implies provider; deferred to Hardening if needed
- Pagination not exposed on preset search — preset datasets are small (1-5 per kind); no `PageInfo` sent
- Proto surprise documented: `CloudResourceSlugAvailabilityCheckResponse` only has `is_available: bool`, not the `existing_resource_id` claimed in T01_2 revised plan

**Verification:** `go build ./...` and `go test ./...` both pass. Zero linter errors.

### Session Progress (2026-02-27, Session 3)

**Phase 6C: Context Discovery — DONE**

Implemented 2 new tools (`list_organizations`, `list_environments`) in two new domain packages (`organization/`, `environment/`), expanding the MCP server from 8 to 10 tools.

**Files created:**
- `internal/domains/organization/tools.go` — Package doc comment (proto provenance), empty `ListOrganizationsInput`, `ListTool()`, `ListHandler()`
- `internal/domains/organization/list.go` — `List` domain function calling `OrganizationQueryController.FindOrganizations`
- `internal/domains/environment/tools.go` — Package doc comment (proto provenance), `ListEnvironmentsInput` with required `org`, `ListTool()`, `ListHandler()` with validation
- `internal/domains/environment/list.go` — `List` domain function calling `EnvironmentQueryController.FindAuthorized`

**Files modified:**
- `internal/server/server.go` — Imported `organization` and `environment` packages, registered 2 new tools (count 8 → 10)

**Design decisions made:**
- Flat domain package structure maintained — evaluated mirroring proto hierarchy (`infrahub/`, `resourcemanager/`) but decided against it: only 5-6 packages at full expansion, Go idiom favors flat, MCP server is an adapter layer not the platform, no shared code within groupings
- Provenance documented via package doc comments mapping each package to its proto service origin
- `FindOrganizations` chosen over `find` (platform-operator-only, wrong layer)
- `FindAuthorized` chosen over `findByOrg` (FGA-filtered, permission-aware)
- No unit tests for these tools — they are thin RPC wrappers with no domain logic (no enum resolution, no identifier validation, no kind conversion)

**Verification:** `go build ./...` and `go test ./...` both pass. Zero linter errors.

### Session Progress (2026-02-27, Session 2)

**Phase 6B: Stack Job Observability — DONE**

Implemented 3 new tools (`get_stack_job`, `get_latest_stack_job`, `list_stack_jobs`) in a new `internal/domains/stackjob/` domain package, expanding the MCP server from 5 to 8 tools.

**Files created:**
- `internal/domains/stackjob/enum.go` — Four enum resolvers: `resolveOperationType`, `resolveExecutionStatus`, `resolveExecutionResult`, `resolveKind`, plus shared `joinEnumValues` helper
- `internal/domains/stackjob/enum_test.go` — 11 unit tests covering all resolvers (valid values, unknown values, empty strings)
- `internal/domains/stackjob/get.go` — `Get` domain function calling `StackJobQueryController.Get`
- `internal/domains/stackjob/latest.go` — `GetLatest` domain function calling `StackJobQueryController.GetLastStackJobByCloudResourceId`
- `internal/domains/stackjob/list.go` — `List` domain function calling `StackJobQueryController.ListByFilters` with enum resolution and pagination defaults
- `internal/domains/stackjob/tools.go` — Tool definitions, input structs, and handlers for all 3 tools

**Files modified:**
- `internal/server/server.go` — Imported `stackjob` package, registered 3 new tools (count 5 → 8)

**Design decisions made:**
- Added `get_stack_job` (by ID) as a third tool beyond the original 2-tool plan — workflow analysis revealed dead-ends for polling, user-provided IDs, and cross-turn references
- Renamed `get_stack_job_status` to `get_latest_stack_job` — the original name implied a status fragment but returns a full `StackJob`, identical to `get_stack_job`. The distinction is the lookup key (resource ID vs job ID), not the return shape
- Corrected `org` from optional to required on `list_stack_jobs` — server-side `buf.validate` constraint discovered during proto analysis
- Enum resolvers kept domain-local in `stackjob/enum.go` — no cross-domain coupling; `resolveKind` duplicated (trivial one-liner) rather than importing from `cloudresource`
- Pagination exposed with defaults (`page_num=1`, `page_size=20`) — stack jobs can be numerous unlike canvas view resources

**Verification:** `go build ./...` and `go test ./...` both pass. Zero linter errors.

### Session Progress (2026-02-27, Session 1)

**Phase 6A: Complete the Resource Lifecycle — DONE**

Implemented 2 new tools (`list_cloud_resources`, `destroy_cloud_resource`), expanding the MCP server from 3 to 5 tools.

**Files created:**
- `internal/domains/cloudresource/list.go` — `List` domain function calling `CloudResourceSearchQueryController.GetCloudResourcesCanvasView`
- `internal/domains/cloudresource/destroy.go` — `Destroy` domain function (resolves full resource, then calls `CommandController.Destroy`)
- `internal/domains/cloudresource/list_test.go` — 7 tests for `resolveKinds` (nil, empty, single, multiple, unknown, mixed)

**Files modified:**
- `internal/domains/cloudresource/kind.go` — added `resolveKinds` batch helper
- `internal/domains/cloudresource/identifier.go` — added `resolveResource` (returns full `*CloudResource` proto, parallels `resolveResourceID`)
- `internal/domains/cloudresource/tools.go` — added tool definitions, handlers, input structs for both new tools
- `internal/server/server.go` — registered both new tools (count 3 → 5)

**Design decisions made:**
- Import alias `cloudresourcesearch` for search stubs (package name collision with domain package)
- `resolveResource` helper created for `destroy` (RPC requires full `CloudResource`, not just ID)
- `get.go` left untouched to minimize blast radius (refactor deferred to Hardening)
- `resolveKinds` with fail-fast semantics for kind-list conversion

**Verification:** `go build ./...` and `go test ./...` both pass. Zero linter errors.

### Context
- Full API surface analysis completed across `plantonhq/planton/apis` (8 gRPC services)
- 15 new tools identified (18 total), organized into 5 phases (6A–6E + Hardening)
- All 5 open questions resolved (see `tasks/T01_1_feedback.md`)
- Revised plan: `tasks/T01_2_revised_plan.md`
- Key decisions:
  - `list_cloud_resources` uses `CloudResourceSearchQueryController.getCloudResourcesCanvasView` (not `find` or `streamByOrg` — both are wrong layer)
  - No MCP-level destroy confirmation (agent responsibility)
  - Full stack job responses (no secrets present)
  - Presets split into search + get (adds 1 tool)
  - Phase order confirmed: 6A → 6B → 6C → 6D → 6E
  - Total tool count adjusted from 17 to 18 (added `get_stack_job`)
  - `get_stack_job_status` renamed to `get_latest_stack_job` (naming honesty)

## Current Step

- ✅ Phase 6A: Complete the Resource Lifecycle (2 tools, 3 → 5) — 2026-02-27
- ✅ Phase 6B: Stack Job Observability (3 tools, 5 → 8) — 2026-02-27
- ✅ Phase 6C: Context Discovery (2 tools, 8 → 10) — 2026-02-27
- ✅ Phase 6D: Agent Quality-of-Life (3 tools, 10 → 13) — 2026-02-27
  - Also: extracted `resolveKind` to shared `domains` package
- ✅ Phase 6E: Advanced Operations (5 tools, 13 → 18) — 2026-02-27
  - Proto surprises: `get_env_var_map` takes YAML not ID+manifest; `resolve_value_references` resolves ALL refs not specific ones
- ✅ Hardening: Documentation and Tool Description Review — 2026-02-27
  - Normalized tool descriptions, expanded README/tools.md/development.md
  - Surprise: docs/tools.md was the most important update (not in original plan)
- 🔵 Next: **Remaining Hardening** — H4 (`get.go` refactor) or project close-out

## Next Steps

1. **H4: `get.go` refactor** (optional, low priority)
   - Refactor `get.go` to use `resolveResource` instead of its current inline logic (deferred since Phase 6A)
   - Code-quality improvement only — no user-facing behavior change
2. **Project close-out** — if H4 is skipped, the project is complete

## Quick Commands

After loading context:
- "Start Hardening" - Begin the hardening phase
- "Show project status" - Get overview of progress
- "Create checkpoint" - Save current progress
- "Review guidelines" - Check established patterns

---

*This file provides direct paths to all project resources for quick context loading.*
