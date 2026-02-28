---
name: Tier 2 Pipeline MCP Tools
overview: Implement 9 MCP tools for the ServiceHub Pipeline domain in `internal/domains/servicehub/pipeline/`, covering pipeline observability (list, get, get_last), lifecycle control (run, rerun, cancel), gate resolution, and repository pipeline file management (list_files, update_file).
todos:
  - id: get
    content: "Implement get_pipeline: get.go (Get via PipelineQueryController.Get) + tool/handler in tools.go"
    status: completed
  - id: list
    content: "Implement list_pipelines: list.go (List via PipelineQueryController.ListByFilters with org/service_id/envs/pagination) + tool/handler in tools.go"
    status: completed
  - id: latest
    content: "Implement get_last_pipeline: latest.go (GetLatest via PipelineQueryController.GetLastPipelineByServiceId) + tool/handler in tools.go"
    status: completed
  - id: run
    content: "Implement run_pipeline: run.go (Run via PipelineCommandController.RunGitCommit; requires service_id+branch, optional commit_sha; returns success message) + tool/handler in tools.go"
    status: completed
  - id: rerun
    content: "Implement rerun_pipeline: rerun.go (Rerun via PipelineCommandController.Rerun) + tool/handler in tools.go"
    status: completed
  - id: cancel
    content: "Implement cancel_pipeline: cancel.go (Cancel via PipelineCommandController.Cancel) + tool/handler in tools.go"
    status: completed
  - id: gate
    content: "Implement resolve_pipeline_gate: gate.go (ResolveGate via PipelineCommandController.ResolveManualGate + resolveDecision helper) + tool/handler in tools.go"
    status: completed
  - id: files
    content: "Implement list_pipeline_files: files.go (ListFiles via PipelineQueryController.ListServiceRepoPipelineFiles; custom marshaling to decode bytes content to UTF-8 string) + tool/handler in tools.go"
    status: completed
  - id: update-file
    content: "Implement update_pipeline_file: update_file.go (UpdateFile via PipelineCommandController.UpdateServiceRepoPipelineFile; string-to-bytes content encoding + optimistic locking) + tool/handler in tools.go"
    status: completed
  - id: register-wire
    content: Create register.go and wire into internal/server/server.go with servicehubpipeline import alias
    status: completed
  - id: verify
    content: "Verify: go build ./... && go vet ./... && go test ./..."
    status: completed
isProject: false
---

# Tier 2: ServiceHub Pipeline MCP Tools

## Scope

9 tools in a new package `internal/domains/servicehub/pipeline/`, wired into [internal/server/server.go](internal/server/server.go) via `servicehubpipeline.Register`.

## Proto API Surface (from analysis)

**PipelineQueryController** (3 unary RPCs used):

- `get(PipelineId) -> Pipeline`
- `listByFilters(ListPipelinesByFiltersInput) -> PipelineList`
- `getLastPipelineByServiceId(ServiceId) -> Pipeline`
- `listServiceRepoPipelineFiles(ListServiceRepoPipelineFilesInput) -> ServiceRepoPipelineFileList`

**PipelineCommandController** (5 unary RPCs used):

- `runGitCommit(RunGitCommitPipelineRequest) -> Empty`
- `rerun(PipelineId) -> Pipeline`
- `cancel(PipelineId) -> Pipeline`
- `resolveManualGate(ResolvePipelineManualGateRequest) -> Empty`
- `updateServiceRepoPipelineFile(UpdateServiceRepoPipelineFileInput) -> UpdateServiceRepoPipelineFileResponse`

**Excluded** (same rationale as T01 plan):

- Streaming RPCs (`getStatusStream`, `getLogStream`, `streamPipelinesByOrg`)
- System CRUD (`apply`, `create`, `update`, `delete` on Pipeline itself)

## Tool Catalogue


| #   | Tool Name               | RPC                             | Key Inputs                                                                                               |
| --- | ----------------------- | ------------------------------- | -------------------------------------------------------------------------------------------------------- |
| 1   | `list_pipelines`        | `listByFilters`                 | org (req), service_id, envs[], page_num, page_size                                                       |
| 2   | `get_pipeline`          | `get`                           | id (req)                                                                                                 |
| 3   | `get_last_pipeline`     | `getLastPipelineByServiceId`    | service_id (req)                                                                                         |
| 4   | `run_pipeline`          | `runGitCommit`                  | service_id (req), branch (req), commit_sha (opt)                                                         |
| 5   | `rerun_pipeline`        | `rerun`                         | id (req)                                                                                                 |
| 6   | `cancel_pipeline`       | `cancel`                        | id (req)                                                                                                 |
| 7   | `resolve_pipeline_gate` | `resolveManualGate`             | pipeline_id (req), deployment_task_name (req), decision (req)                                            |
| 8   | `list_pipeline_files`   | `listServiceRepoPipelineFiles`  | service_id (req), branch (opt)                                                                           |
| 9   | `update_pipeline_file`  | `updateServiceRepoPipelineFile` | service_id (req), path (req), content (req), expected_base_sha (opt), commit_message (opt), branch (opt) |


## Design Decisions (confirmed with user)

**DD-T2-1: `run_pipeline` requires branch, commit_sha optional.**
Proto field `branch` is `(buf.validate.field).required = true`. The `runGitCommit` RPC returns `google.protobuf.Empty` -- no pipeline ID is returned. The tool returns a success message directing the agent to use `get_last_pipeline` to check the result.

**DD-T2-2: Decode pipeline file content from bytes to string.**
`ServiceRepoPipelineFile.content` is `bytes` (base64 in protojson). For `list_pipeline_files`, we custom-marshal the response, decoding `content` bytes to plain UTF-8 string so agents see actual YAML/pipeline content. For `update_pipeline_file`, agents provide content as a plain string, which we encode to `[]byte` before sending. This deviates from the standard `domains.MarshalJSON(resp)` pattern but is the right call for agent UX.

**DD-T2-3: Single gate tool (not two like InfraPipeline).**
ServiceHub Pipeline has a single `resolveManualGate` RPC with a `deployment_task_name` field. InfraPipeline has two separate RPCs (env gate + node gate). One tool is sufficient here.

## Key Differences from InfraPipeline

- `**run_pipeline`**: requires `branch` (required) + optional `commit_sha`. InfraPipeline dispatches between chart-source and git-commit based on `commit_sha` presence.
- **Gate resolution**: single tool with `deployment_task_name` vs. InfraPipeline's two tools (env gate + node gate).
- **Pipeline files**: two entirely new tools (`list_pipeline_files`, `update_pipeline_file`) with no InfraPipeline analog.
- `**runGitCommit` returns Empty**: run tool returns a success message, not a pipeline object.

## Package Structure

```
internal/domains/servicehub/pipeline/
  register.go       # Register(srv, serverAddress) -- adds all 9 tools
  tools.go          # Input structs, Tool/Handler functions
  list.go           # List() via PipelineQueryController.ListByFilters
  get.go            # Get() via PipelineQueryController.Get
  latest.go         # GetLatest() via PipelineQueryController.GetLastPipelineByServiceId
  run.go            # Run() via PipelineCommandController.RunGitCommit
  rerun.go          # Rerun() via PipelineCommandController.Rerun
  cancel.go         # Cancel() via PipelineCommandController.Cancel
  gate.go           # ResolveGate() via PipelineCommandController.ResolveManualGate
  files.go          # ListFiles() via PipelineQueryController.ListServiceRepoPipelineFiles
  update_file.go    # UpdateFile() via PipelineCommandController.UpdateServiceRepoPipelineFile
```

## Patterns to Follow

- **Tool/Handler pattern**: Follows [internal/domains/infrahub/infrapipeline/tools.go](internal/domains/infrahub/infrapipeline/tools.go) -- input structs with jsonschema tags, `*Tool()` returning `*mcp.Tool`, `*Handler()` returning typed handler.
- **Register pattern**: Follows [internal/domains/servicehub/service/register.go](internal/domains/servicehub/service/register.go) -- `Register(srv, serverAddress)` calling `mcp.AddTool` for each tool.
- **Operation pattern**: Follows existing `domains.WithConnection` + `domains.MarshalJSON` + `domains.RPCError` for all standard tools. Custom marshaling only for `list_pipeline_files`.
- **Gate resolution**: Follows [internal/domains/infrahub/infrapipeline/gate.go](internal/domains/infrahub/infrapipeline/gate.go) -- `resolveDecision` maps "approve"/"reject" to proto enum.
- **Import alias**: `pipelinev1` for the proto stubs, `servicehubpipeline` for the package import in `server.go`.

## Server Wiring

In [internal/server/server.go](internal/server/server.go):

- Add import: `servicehubpipeline "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/pipeline"`
- Add: `servicehubpipeline.Register(srv, serverAddress)` after the existing `servicehubservice.Register` call.

## Implementation Order

Each step produces a compilable package. Steps 1-9 are individual operation files, step 10 is registration, step 11 is server wiring.

1. `get.go` -- simplest; establishes gRPC client pattern for this package
2. `list.go` -- pagination pattern; `ListPipelinesByFiltersInput` with `envs` array filter
3. `latest.go` -- simple; uses `ServiceId` from service proto
4. `run.go` -- `RunGitCommitPipelineRequest` with required branch; returns success message (Empty response)
5. `rerun.go` -- simple `PipelineId` in, `Pipeline` out
6. `cancel.go` -- simple `PipelineId` in, `Pipeline` out
7. `gate.go` -- `resolveDecision` helper + `ResolveManualGate` RPC
8. `files.go` -- custom marshaling to decode bytes content to string
9. `update_file.go` -- string-to-bytes content encoding + optimistic locking
10. `tools.go` + `register.go` -- all input structs, tool defs, handlers, registration
11. Server wiring in `server.go`

Note: In practice, `tools.go` will be written alongside the operation files since input structs and handlers are co-dependent. The order above is logical dependency order.