---
name: Phase 1A InfraChart tools
overview: Add 3 MCP tools for the InfraChart domain (list, get, build) following established patterns from preset/ and stackjob/, with a simplified build input that fetches the chart internally before rendering.
todos:
  - id: resolve-enum
    content: Inspect apiresourcekind.ApiResourceKind enum to find the infra_chart constant name
    status: completed
  - id: create-tools-go
    content: Create infrachart/tools.go with tool definitions, input structs, and handlers for all 3 tools
    status: completed
  - id: create-get-go
    content: Create infrachart/get.go with Get() implementation
    status: completed
  - id: create-list-go
    content: Create infrachart/list.go with List() implementation including pagination and kind constant
    status: completed
  - id: create-build-go
    content: Create infrachart/build.go with two-step Build() implementation (get + merge params + build)
    status: completed
  - id: register-tools
    content: Update server.go to import infrachart and register the 3 new tools
    status: completed
  - id: update-doc
    content: Update infrahub/doc.go subpackage list
    status: completed
  - id: verify-build
    content: Run go build, go vet, go test to verify everything compiles and passes
    status: completed
isProject: false
---

# Phase 1A: InfraChart MCP Tools

## Decisions made

- **Tool naming**: `list_infra_charts` (not `search_`) because the backing `Find` RPC only supports org/env/pagination filters, not free-text search. Consistent with `list_stack_jobs`, `list_cloud_resources`.
- **Build input**: Simplified `chart_id` + `params` map. Internally fetches the chart via `Get`, merges param overrides, then calls `Build`. Two RPCs but trivial input for AI agents.
- `**FindApiResourcesRequest` operator restriction**: The proto comments note this request type is "restricted to platform operators only." MCP users are developers/operators, so this is acceptable. Worth noting but not blocking.

## Tools to implement


| Tool                | Backing RPC                      | Input                             | Output                                 |
| ------------------- | -------------------------------- | --------------------------------- | -------------------------------------- |
| `list_infra_charts` | `InfraChartQueryController.Find` | org?, env?, page_num?, page_size? | Paginated `InfraChartList` JSON        |
| `get_infra_chart`   | `InfraChartQueryController.Get`  | id (required)                     | Full `InfraChart` JSON                 |
| `build_infra_chart` | `Get` then `Build`               | chart_id (required), params?      | `InfraChart` JSON with rendered output |


## Files to create

All under `internal/domains/infrahub/infrachart/`:

### 1. `tools.go` — Tool definitions, input structs, handlers

Pattern: follow [preset/tools.go](internal/domains/infrahub/preset/tools.go) structure.

- Package doc comment listing all 3 tools
- `ListInfraChartsInput` struct: `org`, `env` (both optional), `page_num`, `page_size` (both optional, 1-based defaults)
- `ListTool()` + `ListHandler(serverAddress)`
- `GetInfraChartInput` struct: `id` (required)
- `GetTool()` + `GetHandler(serverAddress)`
- `BuildInfraChartInput` struct: `chart_id` (required), `params` (optional `map[string]any` for param name-to-value overrides)
- `BuildTool()` + `BuildHandler(serverAddress)`

All handlers follow the established pattern: validate required fields, call domain function, return via `domains.TextResult()`.

### 2. `get.go` — Get implementation

Pattern: follow [preset/get.go](internal/domains/infrahub/preset/get.go) (7 lines of logic).

```go
func Get(ctx context.Context, serverAddress, chartID string) (string, error)
```

- `domains.WithConnection` wrapping `InfraChartQueryControllerClient.Get`
- Input: `apiresource.ApiResourceId{Value: chartID}`
- `domains.RPCError` for error translation
- `domains.MarshalJSON` for response serialization

### 3. `list.go` — List/Find implementation

Pattern: follow [stackjob/list.go](internal/domains/infrahub/stackjob/list.go) for pagination defaults.

```go
type ListInput struct { Org, Env string; PageNum, PageSize int32 }
func List(ctx context.Context, serverAddress string, input ListInput) (string, error)
```

- Construct `apiresource.FindApiResourcesRequest` with:
  - `Page`: convert 1-based input to 0-based proto convention
  - `Kind`: hard-coded to the InfraChart `ApiResourceKind` enum value (not user-supplied)
  - `Org`, `Env`: pass through from input
- Use `domains.WithConnection` + `InfraChartQueryControllerClient.Find`

**Open item to resolve during implementation**: exact `ApiResourceKind` enum value name for infra charts (likely `infra_chart` or similar). Will inspect `apiresourcekind` package.

### 4. `build.go` — Build implementation (two-step)

New pattern not seen elsewhere in the codebase. Must be implemented carefully.

```go
type BuildInput struct { ChartID string; Params map[string]any }
func Build(ctx context.Context, serverAddress string, input BuildInput) (string, error)
```

Implementation steps:

1. Open gRPC connection via `domains.WithConnection`
2. Call `InfraChartQueryControllerClient.Get` with the chart ID
3. If `Params` map is non-empty, iterate over the chart's `Spec.Params` slice and merge overrides using `structpb.NewValue()` for type conversion
4. Call `InfraChartQueryControllerClient.Build` with the modified `InfraChart`
5. Return `domains.MarshalJSON(resp)`

Important: both Get and Build must use the **same connection** within one `WithConnection` call.

## Files to modify

### 5. [internal/server/server.go](internal/server/server.go) — Register new tools

Add import for `infrachart` package. Add 3 `mcp.AddTool` calls in `registerTools()`. Update the tool count from 18 to 21 and add tool names to the log slice.

### 6. [internal/domains/infrahub/doc.go](internal/domains/infrahub/doc.go) — Update subpackage list

Add `infrachart` to the subpackage documentation list.

## Verification

- `go build ./...` must pass
- `go vet ./...` must pass
- `go test ./...` must pass (no new tests in this phase since we can't mock the gRPC backend, consistent with existing pattern where only enum resolution has unit tests)

## Proto imports reference

```
infrachartv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrachart/v1"
apiresource  "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
rpc          "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
```

## Risks and open items

- `**ApiResourceKind` enum value**: Need to inspect the `apiresourcekind` package to find the correct constant for infra charts. Will resolve at implementation start.
- `**FindApiResourcesRequest.Kind` is marked required** via buf/validate: must hard-code the correct value, cannot leave it zero.
- `**structpb.NewValue` type conversion in build.go**: Need to handle type mismatches gracefully (e.g., agent passes string "true" for a bool param). Will add validation with clear error messages.

