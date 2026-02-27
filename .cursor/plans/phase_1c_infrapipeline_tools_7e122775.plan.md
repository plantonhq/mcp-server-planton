---
name: Phase 1C InfraPipeline Tools
overview: Add 7 MCP tools for InfraPipeline observability and control to the infrahub domain, expanding the server from 27 to 34 tools. This completes the Phase 1 trifecta (Chart + Project + Pipeline).
todos:
  - id: query-tools
    content: Create list.go, get.go, latest.go -- the three read-only query functions
    status: completed
  - id: run-cancel
    content: Create run.go and cancel.go -- the two command functions
    status: completed
  - id: gate
    content: Create gate.go -- both gate resolution functions + resolveDecision() enum helper
    status: completed
  - id: tools-file
    content: Create tools.go -- 7 input structs, 7 tool definitions, 7 handlers
    status: completed
  - id: registration
    content: Update server.go (7 AddTool calls, count 27->34) and doc.go (add infrapipeline to subpackage list)
    status: completed
  - id: verify
    content: Run go build, go vet, go test -- all must pass clean
    status: completed
isProject: false
---

# Phase 1C: InfraPipeline Tools

## Context

The server currently has 27 tools. This phase adds 7 tools (not the originally planned 5) covering pipeline listing, retrieval, execution, cancellation, and manual gate resolution.

**Proto package**: `ai.planton.infrahub.infrapipeline.v1`
**Go import**: `github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrapipeline/v1`

## Design Decisions Made During Planning

- **DD-1: Unified run tool** -- A single `run_infra_pipeline` with optional `commit_sha` dispatches to either `RunInfraProjectChartSourcePipeline` (chart-sourced) or `RunGitCommit` (git-sourced) based on whether `commit_sha` is present.
- **DD-2: Manual gate tools included** -- Both `ResolveEnvironmentManualGate` and `ResolveNodeManualGate` are included. Without them, an agent hits a dead end when a pipeline pauses at an approval gate.
- **DD-3: User-friendly gate decisions** -- The tool accepts `"approve"` / `"reject"` and translates to the proto's `"yes"` / `"no"` enum values, which are poor UX for an agent.
- **DD-4: Streaming RPCs excluded** -- `GetStatusStream` and `GetLogStream` are server-streaming and incompatible with MCP. The `get_infra_pipeline` snapshot covers the status use case.
- **DD-5: Pipeline CRUD excluded** -- `Apply`, `Create`, `Update`, `Delete` are internal platform operations. Agents run/cancel pipelines, they don't manage pipeline records.
- **DD-6: ListByFilters correction** -- The proto only supports `org` + `infra_project_id` + pagination (no status/result filters). The plan's original claim of "filter by status" was inaccurate.
- **DD-7: 0-based pagination** -- Following the `infrachart`/`infraproject` convention (1-based tool API, convert to 0-based for PageInfo). Note: `stackjob` does NOT do this conversion, which appears to be a pre-existing inconsistency.

## Tools (7 total)

- `**list_infra_pipelines`** -- List pipelines by org with optional project filter. RPC: `QueryController.ListByFilters`
- `**get_infra_pipeline`** -- Retrieve full pipeline by ID. RPC: `QueryController.Get`
- `**get_latest_infra_pipeline**` -- Get the most recent pipeline for a project. RPC: `QueryController.GetLastInfraPipelineByInfraProjectId`
- `**run_infra_pipeline**` -- Trigger a pipeline run. Dispatches to `RunInfraProjectChartSourcePipeline` or `RunGitCommit` based on `commit_sha` presence. Returns the new pipeline ID.
- `**cancel_infra_pipeline**` -- Cancel a running pipeline. RPC: `CommandController.Cancel`. Returns the updated pipeline.
- `**resolve_infra_pipeline_env_gate**` -- Approve or reject a manual gate for an entire environment. RPC: `CommandController.ResolveEnvironmentManualGate`. Returns confirmation text (RPC returns `Empty`).
- `**resolve_infra_pipeline_node_gate**` -- Approve or reject a manual gate for a specific DAG node. RPC: `CommandController.ResolveNodeManualGate`. Returns confirmation text.

## Files to Create

All under `internal/domains/infrahub/infrapipeline/`:

- **[tools.go](internal/domains/infrahub/infrapipeline/tools.go)** -- 7 input structs, 7 tool definitions, 7 typed handlers. Follows the pattern in [infraproject/tools.go](internal/domains/infrahub/infraproject/tools.go).
- **[list.go](internal/domains/infrahub/infrapipeline/list.go)** -- `List()` function. Builds `ListInfraPipelinesByFiltersInput` with org + optional project ID + pagination. Follows [infraproject/search.go](internal/domains/infrahub/infraproject/search.go) pagination pattern (1-based to 0-based).
- **[get.go](internal/domains/infrahub/infrapipeline/get.go)** -- `Get()` function. Simple get-by-ID. Follows [stackjob/get.go](internal/domains/infrahub/stackjob/get.go).
- **[latest.go](internal/domains/infrahub/infrapipeline/latest.go)** -- `GetLatest()` function. Get by InfraProjectId. Mirrors [stackjob/latest.go](internal/domains/infrahub/stackjob/latest.go) but uses `InfraProjectId` instead of `CloudResourceId`.
- **[run.go](internal/domains/infrahub/infrapipeline/run.go)** -- `Run()` function. Accepts project ID + optional commit SHA. Dispatches to the correct CommandController RPC. Returns marshaled `InfraPipelineId`.
- **[cancel.go](internal/domains/infrahub/infrapipeline/cancel.go)** -- `Cancel()` function. Takes pipeline ID, calls `CommandController.Cancel`. Returns marshaled `InfraPipeline`.
- **[gate.go](internal/domains/infrahub/infrapipeline/gate.go)** -- `ResolveEnvGate()` and `ResolveNodeGate()` functions + `resolveDecision()` helper that maps `"approve"`/`"reject"` to the `WorkflowStepManualGateDecision` enum. Returns human-readable confirmation strings (since both RPCs return `Empty`).

## Files to Modify

- **[internal/server/server.go](internal/server/server.go)** -- Add `infrapipeline` import, 7 `mcp.AddTool` calls, update tool count from 27 to 34, add 7 tool names to the log slice.
- **[internal/domains/infrahub/doc.go](internal/domains/infrahub/doc.go)** -- Add `infrapipeline` to the subpackage list.

## Key Patterns to Follow

Each file follows established conventions from the codebase:

- **Connection**: `domains.WithConnection(ctx, serverAddress, func(ctx, conn) (string, error))`
- **Client creation**: `infrapipelinev1.NewInfraPipelineQueryControllerClient(conn)` / `...CommandControllerClient(conn)`
- **Error handling**: `domains.RPCError(err, "descriptive context")`
- **Result marshaling**: `domains.MarshalJSON(resp)` for proto responses, `domains.TextResult(text)` for handler returns
- **Input validation**: Explicit checks in handlers before calling domain functions

## Execution Order

1. Create `list.go`, `get.go`, `latest.go` (read-only query tools -- lowest risk, establish the package)
2. Create `run.go`, `cancel.go` (command tools)
3. Create `gate.go` (gate resolution + enum helper)
4. Create `tools.go` (all definitions and handlers, wiring everything together)
5. Modify `server.go` and `doc.go` (registration)
6. Verify: `go build ./...`, `go vet ./...`, `go test ./...`

