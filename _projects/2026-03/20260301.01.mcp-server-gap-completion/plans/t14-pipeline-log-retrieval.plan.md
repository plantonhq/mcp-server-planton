---
name: T14 Log Investigation
overview: Add pipeline log retrieval tools that internally collect streaming RPC responses and return them as batch results. Two new tools (get_pipeline_logs, get_infra_pipeline_logs) plus a shared generic stream drain utility.
todos:
  - id: t14-shared-utility
    content: Add DrainStream generic utility to internal/domains/stream.go and StreamCollectTimeout to internal/grpc/client.go
    status: completed
  - id: t14-pipeline-logs
    content: Add get_pipeline_logs tool to ServiceHub Pipeline domain (logs.go, tools.go, register.go)
    status: completed
  - id: t14-infrapipeline-logs
    content: Add get_infra_pipeline_logs tool to InfraPipeline domain (logs.go, tools.go, register.go)
    status: completed
  - id: t14-verify
    content: Verify build compiles and lint passes
    status: completed
  - id: t14-close
    content: Document T14 in next-task.md session history
    status: completed
isProject: false
---

# T14: Pipeline Log Retrieval Tools

## Investigation Findings

### Finding 1: Status Polling Already Works — No New Tools Needed

The existing unary `get` RPCs return **rich, structured status** that is sufficient for an agent to monitor progress via polling. This is identical to how AI agents already monitor terminal commands (run, wait, check, repeat).

**StackJob** — `get_stack_job` returns:

- Overall `status` (queued / running / completed / awaiting_approval)
- Overall `result` (succeeded / failed / cancelled / skipped)
- `errors` (repeated string -- Pulumi/Terraform diagnostic errors)
- Per-IaC-operation state (init, refresh, update_preview, update, destroy_preview, destroy) -- each with its own status, result, errors, timestamps
- IaC operation `snapshot` -- resource change diffs, diagnostic messages, summary (create/update/delete counts), stack outputs

**InfraPipeline** — `get_infra_pipeline` returns:

- Overall `progress_status` and `progress_result`
- `status_reason`
- Build stage -- per-task status, result, error via `TektonPipelineDag`
- Deployment stage -- per-environment status, result, error, diagnostic_message
- Per-node DAG execution -- status, result, status_reason, stack_job_id

**ServiceHub Pipeline** — `get_pipeline` returns:

- Same structure as InfraPipeline: `PipelineStatus` with `progress_status`, `progress_result`
- Build stage -- per-task status, result, error, diagnostic_message, `pod_id`, `image_build_failure_analysis`
- Deployment stage -- per-task status, result, error, diagnostic_message, `stack_job_id`, `errors`

**Existing tool descriptions already guide agents to poll:**

- `rerun_stack_job`: "Use get_stack_job to monitor progress"
- `run_infra_pipeline`: "Use get_infra_pipeline or get_latest_infra_pipeline to monitor progress"

**Conclusion: The agent polling pattern is already fully supported for status and progress monitoring.** No new tools are needed for this use case.

---

### Finding 2: Status vs Logs — Two Distinct Information Layers

A critical distinction: **the unary `get` RPCs return structured status, NOT raw logs.**

**What `getLogStream` returns** (streaming only):

- `TektonTaskLogEntry`: `owner`, `task_name`, `log_message` -- raw stdout/stderr lines from Tekton pipeline task pods
- This content is NOT stored in `Pipeline.status` or `InfraPipeline.status`

**What `get` (unary) returns in the status** -- per build/deployment task:

- `error` (string) -- failure cause summary
- `diagnostic_message` (string) -- reason for current state
- `pod_id` (string) -- ID of the pod that ran the task (could theoretically be used to fetch logs separately, but no unary RPC exists for that)
- For StackJob: `snapshot.diagnostic_messages`, `snapshot.prelude_messages`, `snapshot.resource_diffs` -- structured IaC engine events, not raw CLI output

**Every log-retrieval RPC across the entire Planton API is server-streaming:**

- `InfraPipelineQueryController.getLogStream` -> `stream TektonTaskLogEntry`
- `PipelineQueryController.getLogStream` -> `stream TektonTaskLogEntry`
- `KubernetesPodLogsQueryController.streamPodLogs` (runner, cloudops) -> streaming
- `KubernetesObjectQueryController.streamPodLogs` (integration) -> streaming

There are **zero** unary RPCs that return log content anywhere in the API.

---

### Finding 3: Practical Assessment — Do Agents Need Raw Logs?

**What agents already get WITHOUT raw logs:**

- Exact error messages from Pulumi/Terraform diagnostics (`errors`, `diagnostic_messages`)
- Per-operation/per-task failure reasons (`error`, `diagnostic_message` fields)
- AI-generated error resolution recommendations (`get_error_resolution_recommendation`)
- IaC resource state, diffs, and change summaries (via `snapshot`)
- Per-environment deployment errors
- Build failure analysis (`image_build_failure_analysis` in ServiceHub Pipeline)

**What raw logs would add:**

- Full Terraform/Pulumi CLI stdout (useful for deep debugging beyond structured diagnostics)
- Tekton task stdout/stderr (build output, test results, container image build logs)
- Chronological event timeline

**Honest assessment:** For "did my deployment succeed, and if not, why?" -- the structured status data is more useful than raw logs. For "show me the full build output" or "what did the test runner print?" -- raw logs are the only answer.

---

## Decision: Stream-Collect-Return Pattern (Confirmed)

The server replays logs from the beginning. We will create MCP tools that internally call the streaming `getLogStream` RPC, collect all entries until EOF (completed job) or timeout (running job), and return the batch as a single text response. From the MCP protocol perspective, the tool is unary request-response.

**Not for StackJob** -- StackJob has no `getLogStream` RPC. Its structured progress events (errors, diagnostics, IaC snapshots) are already available via `get_stack_job`.

---

## Implementation Plan

### Architecture

```
Agent calls get_pipeline_logs(pipeline_id)
    |
    v
MCP Tool Handler (unary request-response)
    |
    v
domains.WithConnection (30s outer timeout)
    |
    v
client.GetLogStream(streamCtx, PipelineId)  <-- 15s stream timeout
    |
    v
DrainStream: Recv() loop until EOF or timeout
    |
    v
Format as text, return to agent
```

**Two behaviors based on job state:**

- **Completed/failed job**: Server replays all historical logs, sends EOF. Tool returns complete logs instantly (seconds).
- **Running job**: Server replays history, then emits new logs in real-time. Tool collects for up to 15s, then returns snapshot. Agent calls again later for updates.

### Scope: 2 New Tools


| Tool                      | Domain              | RPC                                         | Input               |
| ------------------------- | ------------------- | ------------------------------------------- | ------------------- |
| `get_pipeline_logs`       | ServiceHub Pipeline | `PipelineQueryController.GetLogStream`      | `pipeline_id`       |
| `get_infra_pipeline_logs` | InfraPipeline       | `InfraPipelineQueryController.GetLogStream` | `infra_pipeline_id` |


### Shared Utility: Generic Stream Drain

A generic `DrainStream[T]` function in `internal/domains/stream.go` that works with any `grpc.ServerStreamingClient[T]`:

```go
func DrainStream[T any](
    stream grpc.ServerStreamingClient[T],
    maxEntries int,
    format func(*T) string,
) (text string, count int, err error)
```

- Reads until `io.EOF` (stream complete) or context error (timeout)
- Caps at `maxEntries` (default 1000) to prevent massive responses
- Each caller provides its own `format` function -- keeps proto imports out of the shared package
- Returns partial results on timeout (not an error -- expected for running jobs)

### Timeout Strategy

- `StreamCollectTimeout = 15 * time.Second` constant in [internal/grpc/client.go](internal/grpc/client.go)
- Created as a derived context inside the `WithConnection` callback -- no changes to the shared `WithConnection` helper
- For completed jobs: stream closes before timeout (fast)
- For running jobs: timeout fires after 15s, tool returns collected logs

### Output Format

Plain text, grouped by task name (not JSON -- logs are inherently textual, fewer tokens for the agent):

```
Collected 42 log entries.

[build-step] Downloading dependencies...
[build-step] npm install completed
[test-step] Running tests...
[test-step] FAIL: test_auth.go:42 - expected 200, got 401
```

If truncated at max entries:

```
--- output truncated at 1000 entries ---
```

### Files to Create


| File                                              | Purpose                          |
| ------------------------------------------------- | -------------------------------- |
| `internal/domains/stream.go`                      | Generic `DrainStream[T]` utility |
| `internal/domains/servicehub/pipeline/logs.go`    | `GetLogs` domain function        |
| `internal/domains/infrahub/infrapipeline/logs.go` | `GetLogs` domain function        |


### Files to Modify


| File                                                                                                       | Change                                                                                        |
| ---------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| [internal/grpc/client.go](internal/grpc/client.go)                                                         | Add `StreamCollectTimeout` constant                                                           |
| [internal/domains/servicehub/pipeline/tools.go](internal/domains/servicehub/pipeline/tools.go)             | Add `GetLogsInput`, `GetLogsTool()`, `GetLogsHandler()` -- update package doc (9 -> 10 tools) |
| [internal/domains/servicehub/pipeline/register.go](internal/domains/servicehub/pipeline/register.go)       | Register `GetLogsTool`                                                                        |
| [internal/domains/infrahub/infrapipeline/tools.go](internal/domains/infrahub/infrapipeline/tools.go)       | Add `GetLogsInput`, `GetLogsTool()`, `GetLogsHandler()` -- update package doc (8 -> 9 tools)  |
| [internal/domains/infrahub/infrapipeline/register.go](internal/domains/infrahub/infrapipeline/register.go) | Register `GetLogsTool`                                                                        |


### Key Reference Code

**Domain function pattern** (from [internal/domains/servicehub/pipeline/get.go](internal/domains/servicehub/pipeline/get.go)):

```go
func Get(ctx context.Context, serverAddress, pipelineID string) (string, error) {
    return domains.WithConnection(ctx, serverAddress,
        func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
            client := pipelinev1.NewPipelineQueryControllerClient(conn)
            resp, err := client.Get(ctx, &pipelinev1.PipelineId{Value: pipelineID})
            if err != nil {
                return "", domains.RPCError(err, fmt.Sprintf("pipeline %q", pipelineID))
            }
            return domains.MarshalJSON(resp)
        })
}
```

**gRPC stub signature** (from `query_grpc.pb.go`):

```go
GetLogStream(ctx context.Context, in *PipelineId, opts ...grpc.CallOption) (grpc.ServerStreamingClient[tekton.TektonTaskLogEntry], error)
```

**TektonTaskLogEntry** (from `log.pb.go`):

```go
type TektonTaskLogEntry struct {
    Owner      string  // owner of the log entry
    TaskName   string  // name of the pipeline task
    LogMessage string  // log message
}
```

### Tool Description Guidance

Tool descriptions should communicate:

1. Returns logs collected up to the current point in time
2. For completed/failed pipelines: returns all logs
3. For running pipelines: returns a snapshot -- call again for updates
4. Capped at 1000 entries per call
5. Use `get_pipeline` / `get_infra_pipeline` for status and error summaries -- logs are for raw output

