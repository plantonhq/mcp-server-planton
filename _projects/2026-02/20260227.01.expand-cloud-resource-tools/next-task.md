# Next Task: 20260227.01.expand-cloud-resource-tools

## Quick Resume Instructions

Drop this file into your conversation to quickly resume work on this project.

## Project: 20260227.01.expand-cloud-resource-tools

**Description**: Expand the MCP server's cloud resource tool surface from the current 3 tools (apply, get, delete) to 18 tools covering the full lifecycle: listing/search, infrastructure destroy, stack job observability, org/env discovery, slug validation, presets, locks, rename, env var maps, and cross-resource reference resolution.
**Goal**: Give AI agents full autonomous capability over cloud resource lifecycle — from discovering their operating context (orgs, environments) through CRUD operations to observing provisioning outcomes (stack jobs) and managing operational concerns (locks, presets, references).
**Tech Stack**: Go/gRPC/MCP
**Components**: internal/domains/cloudresource/, internal/domains/stackjob/, internal/domains/organization/ (planned), internal/domains/environment/ (planned), internal/domains/preset/ (planned), internal/server/server.go

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
**Current Task**: Phase 6C — Context Discovery
**Status**: Phase 6B complete, ready to begin Phase 6C
**Last Session**: 2026-02-27

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

## Next Steps

1. **Phase 6C: Context Discovery** (2 tools, new `organization/` and `environment/` domains)
   - `list_organizations` — `OrganizationQueryController.findOrganizations`
   - `list_environments` — `EnvironmentQueryController.findAuthorized`
2. Phase 6D: Agent Quality-of-Life (`check_slug_availability`, `search_cloud_object_presets`, `get_cloud_object_preset`)
3. Phase 6E: Advanced Operations (locks, rename, envvarmap, references)
4. Hardening: Unit tests, README update, docs, potential `get.go` refactor

## Quick Commands

After loading context:
- "Start Phase 6C" - Begin context discovery tools
- "Show project status" - Get overview of progress
- "Create checkpoint" - Save current progress
- "Review guidelines" - Check established patterns

---

*This file provides direct paths to all project resources for quick context loading.*
