---
name: T06 StackJob AI Tools
overview: "Add 5 new AI-native and diagnostic tools to the StackJob domain: error resolution recommendations, IaC resource inspection (by stack job and by API resource), cross-environment service deployment status, and safe stack input debugging. All are read-only Query RPCs extending the existing stackjob package."
todos:
  - id: iac-resources
    content: Implement find_iac_resources_by_stack_job and find_iac_resources_by_api_resource (iac_resources.go + tools.go + register.go)
    status: completed
  - id: stack-input
    content: Implement get_stack_job_input (stack_input.go + tools.go + register.go)
    status: completed
  - id: service-stack-jobs
    content: Implement find_service_stack_jobs_by_env (service_stack_jobs.go + tools.go + register.go)
    status: completed
  - id: error-recommendation
    content: Implement get_error_resolution_recommendation (error_recommendation.go + tools.go + register.go)
    status: completed
  - id: finalize
    content: Update package doc, final go build, lint check, verify all 12 tools registered
    status: completed
isProject: false
---

# T06: StackJob AI-Native Tools

## Scope

5 new tools, all backed by `StackJobQueryController` Query RPCs (read-only). The existing 7 tools remain untouched. Total after: **12 tools**.

**Dropped from original plan:** `get_last_stack_job_by_cloud_resource` — already exists as `get_latest_stack_job` ([latest.go](internal/domains/infrahub/stackjob/latest.go)).

**Added beyond original plan:** `get_stack_job_input` — safe stack input for debugging (user approved).

## Tool Specifications

### 1. `get_error_resolution_recommendation` — AI error analysis (highest-ROI)

- **RPC:** `StackJobQueryController.GetErrorResolutionRecommendation`
- **Input:** `stack_job_id` (required), `error_message` (required)
- **Response:** `google.protobuf.StringValue` — AI-generated recommendation text
- **Auth:** `is_skip_authorization = true` (any authenticated user)
- **Note:** Returns a plain string, not JSON. Use `resp.Value` directly, wrapped in `domains.TextResult`.

### 2. `find_iac_resources_by_stack_job` — IaC state for a specific job

- **RPC:** `StackJobQueryController.FindIacResourcesByStackJobId`
- **Input:** `stack_job_id` (required)
- **Response:** `IacResources{entries: []IacResource}` where each has `address`, `resource_type`, `provider`, `logical_name`, `resource_external_id`

### 3. `find_iac_resources_by_api_resource` — IaC state for any API resource

- **RPC:** `StackJobQueryController.FindIacResourcesByApiResourceId`
- **Input:** `api_resource_id` (required)
- **Response:** Same `IacResources` type as above
- **Proto note:** There is a `todo` in the proto about accepting `ApiResourceSelector` instead; for now we follow the current contract.

### 4. `find_service_stack_jobs_by_env` — Cross-environment deployment overview

- **RPC:** `StackJobQueryController.FindServiceStackJobsByEnv`
- **Input:** `service_id` (required)
- **Response:** `ServiceEnvStackJobs{entries: map<string, StackJob>}` — env name to latest stack job
- **Cross-domain import:** Takes `servicev1.ServiceId` — first servicehub import in the stackjob package, but this mirrors the proto API design.

### 5. `get_stack_job_input` — Safe stack input debugging

- **RPC:** `StackJobQueryController.GetCloudObjectStackInput`
- **Input:** `stack_job_id` (required)
- **Response:** `cloudresourcev1.CloudObjectStackInput{value: google.protobuf.Struct}` — the exact inputs fed to Pulumi/Terraform with credentials redacted
- **Note:** Response is a `Struct` (arbitrary JSON). Can be large. Tool description should set expectations.

## File Plan

All work is within `internal/domains/infrahub/stackjob/`.

### New Files (4 domain function files)

- `**error_recommendation.go`** — `GetErrorRecommendation(ctx, serverAddress, stackJobID, errorMessage) (string, error)`
- `**iac_resources.go`** — Two functions:
  - `FindIacResourcesByStackJob(ctx, serverAddress, stackJobID) (string, error)`
  - `FindIacResourcesByApiResource(ctx, serverAddress, apiResourceID) (string, error)`
- `**service_stack_jobs.go`** — `FindServiceStackJobsByEnv(ctx, serverAddress, serviceID) (string, error)`
- `**stack_input.go**` — `GetStackInput(ctx, serverAddress, stackJobID) (string, error)`

### Modified Files (2 existing files)

- **[tools.go](internal/domains/infrahub/stackjob/tools.go)** — Add 5 input structs, 5 `*Tool()` functions, 5 `*Handler()` functions. Update package doc comment from "Seven tools" to "Twelve tools" and add the 5 new tool names to the list.
- **[register.go](internal/domains/infrahub/stackjob/register.go)** — Add 5 `mcp.AddTool` calls.

### Go Import Paths (verified from existing code)

- `stackjobv1` — `github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1`
- `cloudresourcev1` — `github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1`
- `apiresource` — `github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource`
- `servicev1` — `github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/service/v1`

## Implementation Order

Build from simplest to most complex, compiling after each:

1. **IaC resources** (tools #2 and #3) — simplest pattern, both return `IacResources`, share a file
2. **Stack input** (tool #5) — single ID input, simple response marshaling
3. **Service stack jobs** (tool #4) — introduces cross-domain `servicev1` import
4. **Error recommendation** (tool #1) — unique response type (`StringValue`), wrap `.Value` as text

After each tool pair: `go build ./...` and lint check.

## Patterns to Follow (from existing code)

Every domain function follows the same shape established in [get.go](internal/domains/infrahub/stackjob/get.go):

```go
func DomainFunc(ctx context.Context, serverAddress string, args...) (string, error) {
    return domains.WithConnection(ctx, serverAddress,
        func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
            client := stackjobv1.NewStackJobQueryControllerClient(conn)
            resp, err := client.RPC(ctx, &RequestType{...})
            if err != nil {
                return "", domains.RPCError(err, fmt.Sprintf("description %q", id))
            }
            return domains.MarshalJSON(resp)
        })
}
```

Exception: `get_error_resolution_recommendation` returns `StringValue`, so use `resp.Value` directly via `domains.TextResult` instead of `MarshalJSON`.

## RPCs Explicitly Excluded

- `getProgressEventStream` / `getStackJobStatusStream` — streaming, not supported by MCP transport
- `streamStackJobsByOrg` — platform operator + streaming
- `find` — platform operator
- `getCloudResourceStackExecuteInput` — platform operator (contains credentials)
- `pipelineCancel` — internal pipeline use only

## Risks and Watchpoints

- `**getErrorResolutionRecommendation` latency:** This calls an AI backend (ChatGPT per proto comments). May be slow. The 30s `DefaultRPCTimeout` should be sufficient, but worth monitoring.
- `**CloudObjectStackInput` response size:** The `Struct` field can be large for complex cloud resources. `MarshalJSON` will serialize it fully. No truncation needed for now; MCP handles large text content.
- **Proto `todo` on `findIacResourcesByApiResourceId`:** Proto notes it should accept `ApiResourceSelector` instead of `ApiResourceId`. We follow the current contract. If the proto changes, this tool's input will need updating.

