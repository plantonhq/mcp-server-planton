---
name: Phase 6A Implementation
overview: Implement `list_cloud_resources` and `destroy_cloud_resource` — two new tools that complete the cloud resource lifecycle in the MCP server. This touches 4 files (2 new, 2 modified) plus test files.
todos:
  - id: resolve-kinds
    content: Add resolveKinds helper to kind.go (maps []string to []CloudResourceKind, fails fast on unknown)
    status: completed
  - id: resolve-resource
    content: Add resolveResource helper to identifier.go (resolves ResourceIdentifier to full *CloudResource proto, parallels resolveResourceID)
    status: completed
  - id: list-domain
    content: Create list.go with List domain function (calls CloudResourceSearchQueryController.GetCloudResourcesCanvasView)
    status: completed
  - id: destroy-domain
    content: Create destroy.go with Destroy domain function (resolves resource then calls CommandController.Destroy)
    status: completed
  - id: tool-definitions
    content: Add ListCloudResourcesInput, ListTool, ListHandler, DestroyCloudResourceInput, DestroyTool, DestroyHandler to tools.go
    status: completed
  - id: server-registration
    content: Register both new tools in server.go registerTools function
    status: completed
  - id: unit-tests
    content: Create list_test.go with tests for resolveKinds and org validation
    status: completed
  - id: verify-build
    content: Run go build and go test ./... to verify everything compiles and passes
    status: completed
isProject: false
---

# Phase 6A: Complete the Resource Lifecycle

## Scope

Two new tools added to the existing `cloudresource` domain:

- `**list_cloud_resources**` -- query the search index for browsable resource listing (calls a *different* gRPC service than existing tools)
- `**destroy_cloud_resource`** -- tear down real cloud infrastructure while keeping the resource record

## Files Changed


| Action       | File                                             | Purpose                                             |
| ------------ | ------------------------------------------------ | --------------------------------------------------- |
| **New**      | `internal/domains/cloudresource/list.go`         | `List` domain function + input validation           |
| **New**      | `internal/domains/cloudresource/destroy.go`      | `Destroy` domain function                           |
| **Modified** | `internal/domains/cloudresource/tools.go`        | Tool definitions + handlers for both new tools      |
| **Modified** | `internal/server/server.go`                      | Register both new tools                             |
| **New**      | `internal/domains/cloudresource/list_test.go`    | Tests for kind-list conversion and input validation |
| **New**      | `internal/domains/cloudresource/destroy_test.go` | (only if testable pure logic emerges)               |


---

## Design Decision 1: Import Alias for Search Stubs

**The surprise:** The generated Go stubs for `CloudResourceSearchQueryController` live in package `cloudresource` at:

```
github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/infrahub/cloudresource
```

This **collides** with our own `package cloudresource` (the domain package). We need an explicit import alias. Proposed alias: `cloudresourcesearch`, following the existing convention of `cloudresourcev1` for the infrahub stubs:

```go
cloudresourcesearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/infrahub/cloudresource"
```

Used as: `cloudresourcesearch.NewCloudResourceSearchQueryControllerClient(conn)`, `cloudresourcesearch.ExploreCloudResourcesRequest{...}`.

---

## Design Decision 2: `destroy` Needs Full Resource Resolution

**The surprise:** Unlike `delete` which only needs a resource ID (`ApiResourceDeleteInput{ResourceId}`), the `destroy` RPC signature is:

```
destroy(CloudResource) returns (CloudResource)
```

It requires the **full `CloudResource` proto object**, not just an ID. This means `destroy` must first **fetch** the full resource via the query controller, then pass it to the command controller.

**Approach:** Create a `resolveResource` helper in [identifier.go](internal/domains/cloudresource/identifier.go) that returns `*cloudresourcev1.CloudResource` instead of just a string ID. This parallels the existing `resolveResourceID` pattern:

```go
func resolveResource(ctx context.Context, conn *grpc.ClientConn, id ResourceIdentifier) (*cloudresourcev1.CloudResource, error)
```

This helper handles both identification paths (ID-based and slug-based) and returns the full proto. The `Destroy` function in `destroy.go` calls `resolveResource` then passes the result to `CommandController.Destroy`.

**Note:** I considered also refactoring `Get` in `get.go` to reuse `resolveResource`, but that would change an existing, working file for a marginal deduplication. I recommend we leave `get.go` untouched for now and revisit during the Hardening phase if the pattern repeats across more tools. This keeps Phase 6A's blast radius minimal.

---

## Design Decision 3: Kind-List Conversion for `list_cloud_resources`

The `list_cloud_resources` tool accepts an optional `kinds` field as a list of PascalCase strings (e.g., `["AwsVpc", "GcpCloudSqlDatabase"]`). The proto expects `[]cloudresourcekind.CloudResourceKind` enum values.

**Approach:** Create a `resolveKinds` helper (in [kind.go](internal/domains/cloudresource/kind.go)) that maps a `[]string` to `[]CloudResourceKind`, calling the existing `resolveKind` for each entry. Fails fast on the first unknown kind with a clear error message.

```go
func resolveKinds(kindStrs []string) ([]cloudresourcekind.CloudResourceKind, error)
```

This is testable pure logic, and tests will go in `list_test.go`.

---

## Tool 1: `list_cloud_resources`

### Input Schema

```go
type ListCloudResourcesInput struct {
    Org        string   `json:"org"        jsonschema:"required,Organization identifier."`
    Envs       []string `json:"envs,omitempty"       jsonschema:"Environment slugs to filter by."`
    SearchText string   `json:"search_text,omitempty" jsonschema:"Free-text search query."`
    Kinds      []string `json:"kinds,omitempty"      jsonschema:"PascalCase cloud resource kinds to filter by (e.g. AwsVpc). Read cloud-resource-kinds://catalog for valid kinds."`
}
```

### Handler Flow

1. Validate `org` is non-empty
2. If `kinds` provided, convert via `resolveKinds` (fail fast on unknown kind)
3. Call `List(ctx, serverAddress, input)` domain function
4. Return `domains.TextResult(text)`

### Domain Function (`list.go`)

```go
func List(ctx context.Context, serverAddress string, org string, envs []string, searchText string, kinds []cloudresourcekind.CloudResourceKind) (string, error) {
    return domains.WithConnection(ctx, serverAddress,
        func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
            client := cloudresourcesearch.NewCloudResourceSearchQueryControllerClient(conn)
            resp, err := client.GetCloudResourcesCanvasView(ctx, &cloudresourcesearch.ExploreCloudResourcesRequest{
                Org:        org,
                Envs:       envs,
                SearchText: searchText,
                Kinds:      kinds,
            })
            if err != nil {
                return "", domains.RPCError(err, fmt.Sprintf("cloud resources in org %q", org))
            }
            return domains.MarshalJSON(resp)
        })
}
```

### Cross-Service gRPC Verification

This is the **first tool calling a service from a different proto package** (`search/v1/infrahub/cloudresource` vs `infrahub/cloudresource/v1`). The existing `WithConnection` creates a raw `grpc.ClientConn` to `serverAddress` -- different gRPC service clients simply register different service paths on the same connection. This is standard gRPC multiplexing and should work if the Planton API gateway routes all services. We will verify this works during implementation by running the tool against the real backend.

---

## Tool 2: `destroy_cloud_resource`

### Input Schema

Same dual-path `ResourceIdentifier` pattern as `get` and `delete`:

```go
type DestroyCloudResourceInput struct {
    ID   string `json:"id,omitempty"   jsonschema:"System-assigned resource ID. Provide this alone OR provide all of kind, org, env, and slug."`
    Kind string `json:"kind,omitempty" jsonschema:"PascalCase cloud resource kind (e.g. AwsEksCluster). Required with org, env, slug when id is not provided. Read cloud-resource-kinds://catalog for valid kinds."`
    Org  string `json:"org,omitempty"  jsonschema:"Organization identifier. Required with kind, env, slug when id is not provided."`
    Env  string `json:"env,omitempty"  jsonschema:"Environment identifier. Required with kind, org, slug when id is not provided."`
    Slug string `json:"slug,omitempty" jsonschema:"Immutable unique resource slug within (org, env, kind). Required with kind, org, env when id is not provided."`
}
```

### Tool Description

Distinct from `delete_cloud_resource` -- this must be crystal clear to agents:

> "Destroy the cloud infrastructure (Terraform/Pulumi destroy) for a resource while keeping the resource record on the Planton platform. This tears down the actual cloud resources (VPCs, clusters, databases, etc.). Use delete_cloud_resource to remove the record itself. WARNING: This is a destructive operation that will destroy real cloud infrastructure. Identify the resource by 'id' alone, or by all of 'kind', 'org', 'env', and 'slug' together."

### Handler Flow

1. Build `ResourceIdentifier` from input
2. Validate via `validateIdentifier`
3. Call `Destroy(ctx, serverAddress, id)` domain function
4. Return `domains.TextResult(text)`

### Domain Function (`destroy.go`)

```go
func Destroy(ctx context.Context, serverAddress string, id ResourceIdentifier) (string, error) {
    return domains.WithConnection(ctx, serverAddress,
        func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
            cr, err := resolveResource(ctx, conn, id)
            if err != nil {
                return "", err
            }
            desc := describeIdentifier(id)
            cmdClient := cloudresourcev1.NewCloudResourceCommandControllerClient(conn)
            result, err := cmdClient.Destroy(ctx, cr)
            if err != nil {
                return "", domains.RPCError(err, desc)
            }
            return domains.MarshalJSON(result)
        })
}
```

---

## server.go Changes

Add two lines to `registerTools` and update the count:

```go
mcp.AddTool(srv, cloudresource.ListTool(), cloudresource.ListHandler(serverAddress))
mcp.AddTool(srv, cloudresource.DestroyTool(), cloudresource.DestroyHandler(serverAddress))
```

Update log: `"count", 5` and add tool names to the slice.

---

## Testing Strategy

- `**list_test.go**`: Table-driven tests for `resolveKinds` (valid kinds, unknown kind, empty list, mixed valid/invalid)
- `**list_test.go**`: Validation test for empty `org` field
- `**destroy_test.go**`: Only if `resolveResource` introduces testable branching beyond what `identifier_test.go` already covers. If `resolveResource` is a thin gRPC wrapper with no pure logic, skip the test file (consistent with existing pattern where `delete.go` has no `delete_test.go`).

---

## Implementation Order

1. Add `resolveKinds` to `kind.go`
2. Add `resolveResource` to `identifier.go`
3. Create `list.go` with `List` domain function
4. Create `destroy.go` with `Destroy` domain function
5. Add tool definitions and handlers to `tools.go`
6. Register tools in `server.go`
7. Create `list_test.go` with unit tests
8. Run `go build` and `go test ./...` to verify

