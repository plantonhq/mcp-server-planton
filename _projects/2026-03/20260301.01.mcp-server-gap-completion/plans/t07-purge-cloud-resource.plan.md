---
name: T07 purge cloud resource
overview: Add the `purge_cloud_resource` MCP tool to the CloudResource domain. This is the "nuke everything" lifecycle operation that destroys infrastructure via IaC and then deletes the record, orchestrated as a Temporal workflow on the backend. The `cleanup` RPC is excluded because it is a platform-operator-only emergency operation, consistent with the existing pattern of not exposing `platform/operator` RPCs.
todos:
  - id: purge-domain-fn
    content: Create `purge.go` with the `Purge` domain function (modeled on `delete.go`)
    status: completed
  - id: purge-tool-def
    content: Add `PurgeCloudResourceInput`, `PurgeTool()`, `PurgeHandler()` to `tools.go`; update package doc comment
    status: completed
  - id: purge-register
    content: Register `purge_cloud_resource` in `register.go`
    status: completed
  - id: verify-build
    content: Run `go build ./...` and `go vet ./...` to verify everything compiles cleanly
    status: completed
  - id: update-next-task
    content: Update `next-task.md` to mark T07 as COMPLETED with scope note
    status: completed
isProject: false
---

# T07: CloudResource Lifecycle Completion -- purge_cloud_resource

## Scope

One new tool: `purge_cloud_resource`. The `cleanup` RPC is excluded (platform-operator-only, consistent with existing exclusion of `updateOutputs`, `pipelineApply`, `pipelineDestroy`).

## Domain Analysis

### Protobuf contract

From `[command.proto](apis/ai/planton/infrahub/cloudresource/v1/command.proto)`, lines 55-65:

```protobuf
rpc purge(CloudResourceId) returns (CloudResource) {
  // auth: cloud_resource / delete (user-level, same as destroy)
}
```

- **Input**: `CloudResourceId` (a single `value` string -- the resource ID)
- **Output**: `CloudResource` (the resource that was purged)
- **Semantics**: Temporal workflow: destroy IaC stack, wait for completion, delete record

### Closest existing analog: `delete.go`

The `purge` RPC takes `CloudResourceId` (just an ID string), identical to how `delete` takes `ApiResourceDeleteInput.resource_id`. The implementation pattern in `[delete.go](internal/domains/infrahub/cloudresource/delete.go)` is the exact template:

1. `resolveResourceID(ctx, conn, id)` -- resolve slug-path to ID if needed
2. Call the command RPC with the resolved ID
3. Marshal and return the response

The only difference: `purge` passes `CloudResourceId{Value: resourceID}` instead of `ApiResourceDeleteInput{ResourceId: resourceID}`.

## Implementation

### 1. New file: `internal/domains/infrahub/cloudresource/purge.go`

Domain function following the `delete.go` pattern exactly:

- `Purge(ctx, serverAddress, id ResourceIdentifier) (string, error)`
- Uses `domains.WithConnection` -> `resolveResourceID` -> `cmdClient.Purge` -> `domains.MarshalJSON`

### 2. Additions to `internal/domains/infrahub/cloudresource/tools.go`

- `PurgeCloudResourceInput` struct -- same fields as `DeleteCloudResourceInput` (id or kind+org+env+slug via `ResourceIdentifier`)
- `PurgeTool()` -- returns `*mcp.Tool` with name `purge_cloud_resource` and a description that clearly conveys:
  - This is destroy + delete in one atomic workflow
  - It is destructive and irreversible
  - Use `get_latest_stack_job` to monitor the destroy phase
  - Contrast with `destroy_cloud_resource` (keeps record) and `delete_cloud_resource` (removes record only)
- `PurgeHandler(serverAddress)` -- validates identifier, calls `Purge`, returns `domains.TextResult`
- Update the package doc comment from "Eleven tools" to "Twelve tools" and add `purge_cloud_resource` to the list

### 3. Registration in `internal/domains/infrahub/cloudresource/register.go`

Add `mcp.AddTool(srv, PurgeTool(), PurgeHandler(serverAddress))` to the `Register` function.

### 4. Update `next-task.md`

Mark T07 as COMPLETED with a note that only `purge` was added; `cleanup` was intentionally excluded per the platform-operator exclusion pattern.

## What NOT to do

- No new shared utilities needed -- all building blocks exist (`ResourceIdentifier`, `validateIdentifier`, `resolveResourceID`, `describeIdentifier`)
- No new test files for the domain function itself (consistent with `delete.go`, `destroy.go`, `get.go` which also have no unit tests -- they'd require gRPC mocking for minimal value)
- No changes to `server.go` -- registration is handled inside the domain's `Register()` function

