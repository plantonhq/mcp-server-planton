# Phase 6B: Stack Job Observability Tools

**Date**: February 27, 2026

## Summary

Added three stack job observability tools to the MCP server — `get_stack_job`, `get_latest_stack_job`, and `list_stack_jobs` — in a new `internal/domains/stackjob/` domain package. This gives AI agents full visibility into provisioning outcomes, enabling them to verify whether apply/destroy operations succeeded, poll running jobs, audit deployment history, and diagnose failures. The MCP server expands from 5 to 8 registered tools.

## Problem Statement

After Phase 6A, agents could create, read, update, list, delete, and destroy cloud resources — but had no way to observe what happened next. Provisioning is asynchronous: `apply_cloud_resource` triggers a stack job (Terraform/Pulumi execution) that runs in the background. Without stack job visibility, agents had to blindly trust that operations succeeded.

### Pain Points

- Agents could not verify whether an apply or destroy completed successfully
- No way to detect or diagnose provisioning failures
- No visibility into running jobs for progress monitoring
- No audit trail for deployment history across resources or environments

## Solution

Introduced a new `stackjob` domain package with three complementary access patterns:

| Pattern | Tool | Input | Use Case |
|---|---|---|---|
| Direct lookup | `get_stack_job` | stack job ID | Revisit, poll, user-provided ID |
| Contextual | `get_latest_stack_job` | cloud resource ID | Post-apply/destroy verification |
| Discovery | `list_stack_jobs` | org + filters | Audit, troubleshooting, history |

## Implementation Details

### New domain package: `internal/domains/stackjob/`

Six new files following established patterns from the `cloudresource` domain:

- **`enum.go`** — Four enum resolvers (`resolveOperationType`, `resolveExecutionStatus`, `resolveExecutionResult`, `resolveKind`) with a shared `joinEnumValues` helper that generates user-friendly error messages listing valid values
- **`enum_test.go`** — 11 unit tests covering valid values, unknown values, and empty strings for all resolvers
- **`get.go`** — `Get` domain function backed by `StackJobQueryController.Get(StackJobId)`
- **`latest.go`** — `GetLatest` domain function backed by `StackJobQueryController.GetLastStackJobByCloudResourceId(CloudResourceId)`
- **`list.go`** — `List` domain function backed by `StackJobQueryController.ListByFilters` with enum string resolution, pagination defaults (page 1, size 20), and seven optional filter parameters
- **`tools.go`** — Tool definitions, input structs with jsonschema annotations, and handlers for all three tools

### Modified: `internal/server/server.go`

Registered three new tools, updated count from 5 to 8, updated the structured log.

### Design decisions

- **3 tools instead of original 2**: Workflow analysis revealed that without a direct by-ID lookup, agents hit dead-ends when polling running jobs, handling user-provided job IDs, or referencing jobs across conversation turns
- **Renamed `get_stack_job_status` to `get_latest_stack_job`**: Both tools return a full `StackJob` — the distinction is the lookup key (job ID vs resource ID), not the return shape. The new name is honest about what it does
- **`org` corrected to required on `list_stack_jobs`**: Server-side `buf.validate` constraint discovered during proto analysis
- **No cross-domain coupling**: `stackjob/` does not import `cloudresource/`. The `resolveKind` one-liner is duplicated to keep domains self-contained

## Benefits

- Agents can now autonomously verify provisioning outcomes after apply/destroy
- Full audit trail: "Show me all failed deployments in production" is a single tool call
- Progress monitoring: agents can poll a specific running job by ID
- Structured filtering by org, env, kind, operation type, status, and result
- Paginated results for resources with extensive job history

## Impact

- **AI agents**: Can now close the feedback loop on provisioning — the primary gap after Phase 6A
- **MCP server**: 5 → 8 tools, new `stackjob` domain package
- **Test coverage**: 11 new unit tests, all existing tests still pass
- **Project total**: Adjusted from 17 to 18 tools (Phase 6B contributes 3 instead of 2)

## Related Work

- Phase 6A (same project) — list and destroy cloud resource tools
- Phase 6C (next) — organization and environment discovery tools
- Revised plan: `_projects/2026-02/20260227.01.expand-cloud-resource-tools/tasks/T01_2_revised_plan.md`

---

**Status**: Production Ready
**Timeline**: Single session
