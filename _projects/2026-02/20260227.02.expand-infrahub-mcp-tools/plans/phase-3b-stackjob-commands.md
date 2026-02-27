---
name: Phase 3B StackJob Commands
overview: "Add 4 MCP tools (`rerun_stack_job`, `cancel_stack_job`, `resume_stack_job`, `check_stack_job_essentials`) to the existing stackjob domain, expanding the server from 55 to 59 tools. This completes the operational control surface for stack jobs: agents can now retry failures, cancel stuck runs, approve blocked jobs, and pre-validate deployment prerequisites."
todos:
  - id: create-rerun
    content: Create `rerun.go` -- Rerun RPC wrapper (StackJobCommandController.Rerun)
    status: completed
  - id: create-cancel
    content: Create `cancel.go` -- Cancel RPC wrapper (StackJobCommandController.Cancel)
    status: completed
  - id: create-resume
    content: Create `resume.go` -- Resume RPC wrapper (StackJobCommandController.Resume)
    status: completed
  - id: create-essentials
    content: Create `essentials.go` -- Check RPC wrapper (StackJobEssentialsQueryController.Check) with kind resolution and CloudResourceOwner construction
    status: completed
  - id: update-tools
    content: Update `tools.go` -- Add 4 input structs, 4 tool definitions, 4 handler functions
    status: completed
  - id: update-server
    content: Update `server.go` -- Register 4 new tools, update count 55->59, update tool name list
    status: completed
  - id: update-doc
    content: Update `doc.go` -- Expand stackjob subpackage description
    status: completed
  - id: verify-build
    content: "Verify: `go build ./...` + `go vet ./...` + `go test ./...`"
    status: completed
isProject: false
---

# Phase 3B: StackJob Commands / Lifecycle Control (4 tools)

## Proto Analysis Summary

All RPCs live in `ai.planton.infrahub.stackjob.v1`:

- `**StackJobCommandController**` ([command.proto](apis/ai/planton/infrahub/stackjob/v1/command.proto))
  - `rerun(StackJobId) returns (StackJob)` -- re-run a failed stack job
  - `cancel(StackJobId) returns (StackJob)` -- gracefully cancel a running job (signal-based, not immediate)
  - `resume(StackJobId) returns (StackJob)` -- approve and resume an awaiting-approval job
- `**StackJobEssentialsQueryController**` ([query.proto](apis/ai/planton/infrahub/stackjob/v1/query.proto))
  - `check(CheckStackJobEssentialsInput) returns (CheckStackJobEssentialsResponse)` -- pre-validate all 4 deployment prerequisites

## Tool Definitions

### Tool 1: `rerun_stack_job`

- **Input**: `id` (required, stack job ID)
- **RPC**: `StackJobCommandController.Rerun(StackJobId) -> StackJob`
- **Pattern**: Identical to `cancel_infra_pipeline` -- ID in, updated resource out
- **File**: `internal/domains/infrahub/stackjob/rerun.go`

### Tool 2: `cancel_stack_job`

- **Input**: `id` (required, stack job ID)
- **RPC**: `StackJobCommandController.Cancel(StackJobId) -> StackJob`
- **Pattern**: Same as above
- **File**: `internal/domains/infrahub/stackjob/cancel.go`
- **Key detail for description**: Cancellation is **graceful** -- the currently executing IaC operation completes fully; only remaining operations are skipped. Infrastructure from completed operations remains (no rollback). The proto documents this thoroughly in [command.proto lines 32-74](apis/ai/planton/infrahub/stackjob/v1/command.proto).

### Tool 3: `resume_stack_job`

- **Input**: `id` (required, stack job ID)
- **RPC**: `StackJobCommandController.Resume(StackJobId) -> StackJob`
- **Pattern**: Same as above
- **File**: `internal/domains/infrahub/stackjob/resume.go`
- **Context**: Resumes stack jobs blocked by flow control policies in `awaiting_approval` state. Combined with `cancel_stack_job`, this gives agents a complete approval surface (approve = resume, reject = cancel).

### Tool 4: `check_stack_job_essentials`

- **Input**:
  - `cloud_resource_kind` (required, PascalCase string, e.g. "AwsEksCluster")
  - `org` (required, organization ID)
  - `env` (optional, environment name)
- **RPC**: `StackJobEssentialsQueryController.Check(CheckStackJobEssentialsInput) -> CheckStackJobEssentialsResponse`
- **Request construction**: Resolve kind via existing `domains.ResolveKind()`, wrap org/env in `CloudResourceOwner{Org, Env}`
- **File**: `internal/domains/infrahub/stackjob/essentials.go`
- **Response**: 4 preflight checks (iac_module, backend_credential, flow_control, provider_credential), each with `passed` bool + `errors` list

Proto types referenced:

- `CheckStackJobEssentialsInput` from [preflight.proto](apis/ai/planton/infrahub/stackjob/v1/preflight.proto) -- `cloud_resource_kind` (CloudResourceKind enum) + `cloud_resource_owner` (CloudResourceOwner{org, env})
- `CloudResourceOwner` from [io.proto](apis/ai/planton/commons/apiresource/io.proto) -- `org` (required) + `env` (optional)

## Files to Create

- `internal/domains/infrahub/stackjob/rerun.go` -- Rerun RPC, ~25 lines (mirrors [get.go](internal/domains/infrahub/stackjob/get.go) pattern)
- `internal/domains/infrahub/stackjob/cancel.go` -- Cancel RPC, ~25 lines
- `internal/domains/infrahub/stackjob/resume.go` -- Resume RPC, ~25 lines
- `internal/domains/infrahub/stackjob/essentials.go` -- Check RPC, ~35 lines (needs kind resolution + CloudResourceOwner construction)

## Files to Modify

- [internal/domains/infrahub/stackjob/tools.go](internal/domains/infrahub/stackjob/tools.go) -- Add 4 input structs + 4 tool defs + 4 handlers. Grows from 137 to ~310 lines (consistent with infrapipeline/tools.go at 299 lines for 7 tools).
- [internal/server/server.go](internal/server/server.go) -- Add 4 `mcp.AddTool` calls after existing stackjob registrations (line 69). Update tool count 55 -> 59. Add 4 tool names to the log list.
- [internal/domains/infrahub/doc.go](internal/domains/infrahub/doc.go) -- Update stackjob description from "observability for IaC stack jobs (get, list, latest)" to include rerun, cancel, resume, essentials.

## Patterns to Follow

All 4 tools follow established conventions:

- **RPC wrapper functions**: `domains.WithConnection` + `New*Client` + RPC call + `domains.RPCError` + `domains.MarshalJSON` (see existing [get.go](internal/domains/infrahub/stackjob/get.go))
- **Tool definitions**: Input struct with `json` + `jsonschema` tags, `*Tool()` returning `*mcp.Tool`, `*Handler()` returning typed handler closure (see [infrapipeline/tools.go](internal/domains/infrahub/infrapipeline/tools.go))
- **Kind resolution**: `domains.ResolveKind()` from [kind.go](internal/domains/kind.go) for the essentials tool
- **Server registration**: `mcp.AddTool(srv, stackjob.*Tool(), stackjob.*Handler(serverAddress))` grouped with existing stackjob tools

## Deferred

- `which`* RPCs (whichIacRunner, whichIacModule, whichBackendCredential, whichFlowControlPolicy, whichProviderCredential) -- granular debugging lookups. The `check` RPC covers the combined preflight. Can add later if agents need to debug specific failures.
- `pipelineCancel` -- internal RPC for pipeline-to-stackjob cancellation chain. Not an agent-facing operation.

## Verification

- `go build ./...`
- `go vet ./...`
- `go test ./...`

