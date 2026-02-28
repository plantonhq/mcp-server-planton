---
name: Tier 1 Service MCP Tools
overview: Implement 7 MCP tools for the ServiceHub Service entity in a new `internal/domains/servicehub/service/` package, following the established infrahub patterns exactly, and wire them into the server.
todos:
  - id: create-package
    content: Create `internal/domains/servicehub/service/` package directory
    status: completed
  - id: impl-get
    content: "Implement get_service: `get.go` (Get, resolveService, resolveServiceID, describeService) + tool/handler in `tools.go`"
    status: completed
  - id: impl-search
    content: "Implement search_services: `search.go` (Search using generic ApiResourceSearchQueryController.searchByKind) + tool/handler in `tools.go`"
    status: completed
  - id: impl-apply
    content: "Implement apply_service: `apply.go` (Apply with protojson unmarshal) + tool/handler in `tools.go`"
    status: completed
  - id: impl-delete
    content: "Implement delete_service: `delete.go` (Delete using resolveServiceID + ApiResourceDeleteInput) + tool/handler in `tools.go`"
    status: completed
  - id: impl-disconnect-webhook
    content: "Implement disconnect_service_git_repo and configure_service_webhook: `disconnect.go` + `webhook.go` + tools/handlers in `tools.go`"
    status: completed
  - id: impl-branches
    content: "Implement list_service_branches: `branches.go` + tool/handler in `tools.go`"
    status: completed
  - id: impl-register
    content: Create `register.go` with Register function adding all 7 tools
    status: completed
  - id: wire-server
    content: "Wire into `internal/server/server.go`: add import and Register call"
    status: completed
  - id: verify-build
    content: "Verify the project compiles cleanly: `go build ./...`"
    status: completed
isProject: false
---

# Tier 1: ServiceHub Service MCP Tools

## Discovery Summary

After thorough exploration of the codebase, I have mapped every pattern used by the existing infrahub tools (`infraproject`, `preset`, `cloudresource`, etc.) and the ServiceHub Service proto API surface. Two surprises emerged:

**Surprise 1 (Resolved)**: No dedicated `searchServices` RPC exists. We will use the generic `ApiResourceSearchQueryController.searchByKind` with `api_resource_kind = service` (enum value 16). This is a new pattern for this codebase -- no existing MCP tool uses it yet.

**Surprise 2 (Decision Point in Plan)**: The T01 plan mentioned reusing `cloudresource.GetParser(kind)` to validate `cloud_object` inside `deployment_targets` when `deployment_config_source = inline`. I am recommending we **skip this** for now. The `apply_infra_project` pattern simply unmarshals the raw JSON into the proto and lets the backend validate. Client-side validation of `cloud_object` would couple the servicehub package to infrahub/cloudresource and duplicate backend business logic. We follow the same thin-client approach as infraproject.

---

## Tool Catalogue (7 tools)

All 7 tools follow the established pattern:

- **Tool function** returning `*mcp.Tool` (Name + Description)
- **Handler function** capturing `serverAddress`, returning typed handler with input struct
- **Operation function** using `domains.WithConnection`, `domains.MarshalJSON`, `domains.RPCError`


| #   | Tool Name                     | RPC Controller                     | RPC Method               | Input Identification                             |
| --- | ----------------------------- | ---------------------------------- | ------------------------ | ------------------------------------------------ |
| 1   | `search_services`             | `ApiResourceSearchQueryController` | `searchByKind`           | org (required), search_text, page_num, page_size |
| 2   | `get_service`                 | `ServiceQueryController`           | `get` / `getByOrgBySlug` | id OR org+slug                                   |
| 3   | `apply_service`               | `ServiceCommandController`         | `apply`                  | Full Service JSON object                         |
| 4   | `delete_service`              | `ServiceCommandController`         | `delete`                 | id OR org+slug                                   |
| 5   | `disconnect_service_git_repo` | `ServiceCommandController`         | `disconnectGitRepo`      | id OR org+slug                                   |
| 6   | `configure_service_webhook`   | `ServiceCommandController`         | `configureWebhook`       | id OR org+slug                                   |
| 7   | `list_service_branches`       | `ServiceQueryController`           | `listBranches`           | id OR org+slug                                   |


---

## Package Structure

```
internal/domains/servicehub/
  service/
    register.go       # Register(srv, serverAddress) -- adds all 7 tools
    tools.go          # Input structs, Tool/Handler functions, validateIdentification
    search.go         # Search() operation using ApiResourceSearchQueryController
    get.go            # Get(), resolveService(), resolveServiceID(), describeService()
    apply.go          # Apply() operation
    delete.go         # Delete() operation
    disconnect.go     # DisconnectGitRepo() operation
    webhook.go        # ConfigureWebhook() operation
    branches.go       # ListBranches() operation
```

---

## Key Implementation Details

### `tools.go` -- Central Definitions

Contains all 7 input structs, tool definitions, handlers, and the shared `validateIdentification` helper. Follows [infraproject/tools.go](internal/domains/infrahub/infraproject/tools.go) pattern exactly.

**Identification pattern** (reused by get, delete, disconnect, webhook, branches): ID alone OR org+slug. Identical to infraproject's `validateIdentification`.

**Input structs and their fields:**

- **SearchServicesInput**: `org` (required), `search_text`, `page_num`, `page_size`
- **GetServiceInput**: `id`, `org`, `slug` (mutually exclusive paths)
- **ApplyServiceInput**: `service map[string]any` (required, full Service JSON)
- **DeleteServiceInput**: `id`, `org`, `slug`
- **DisconnectServiceGitRepoInput**: `id`, `org`, `slug`
- **ConfigureServiceWebhookInput**: `id`, `org`, `slug`
- **ListServiceBranchesInput**: `id`, `org`, `slug`

### `search.go` -- New Pattern

Uses `ApiResourceSearchQueryController.SearchByKind` instead of a domain-specific search controller. This is the first MCP tool to use this generic endpoint.

Key import: `apiresourcesearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/apiresource"`

The `api_resource_kind` field is set to `apiresourcekind.ApiResourceKind_service` (enum value 16).

### `get.go` -- Resolve Pattern

Follows [infraproject/get.go](internal/domains/infrahub/infraproject/get.go) exactly:

- `resolveService()` returns the full `*servicev1.Service` proto
- `resolveServiceID()` returns just the ID string (used by delete, disconnect, webhook, branches)
- `describeService()` returns human-readable resource description for error messages

Key imports:

- `servicev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/service/v1"`
- `apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"`

### `apply.go` -- Thin Client

Follows [infraproject/apply.go](internal/domains/infrahub/infraproject/apply.go): marshal raw map to JSON bytes, unmarshal into `servicev1.Service` proto via `protojson.Unmarshal`, send via `ServiceCommandController.Apply`. No client-side validation of `cloud_object` in `deployment_targets`.

### `delete.go`

Follows [infraproject/delete.go](internal/domains/infrahub/infraproject/delete.go): resolve ID (directly or via org+slug), call `ServiceCommandController.Delete` with `ApiResourceDeleteInput`.

### `disconnect.go` and `webhook.go`

Both call `ServiceCommandController` RPCs that take `ServiceId{Value: id}` and return `Service`. Pattern: resolve ID, call RPC, marshal response. These are simpler than delete since the RPCs just take `ServiceId`.

### `branches.go`

Calls `ServiceQueryController.ListBranches` which takes `ServiceId` and returns `StringList`. Pattern: resolve ID, call RPC, marshal response.

---

## Server Wiring

Add one import and one `Register` call in [internal/server/server.go](internal/server/server.go):

```go
import (
    servicehubservice "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/service"
)

func registerTools(srv *mcp.Server, serverAddress string) {
    // ... existing registrations ...
    servicehubservice.Register(srv, serverAddress)
    // ...
}
```

The import alias `servicehubservice` avoids collision with the `service` keyword and follows the pattern of disambiguating domain entity packages.

---

## Tool Descriptions (Critical for AI Agent UX)

Each tool description must be precise, actionable, and include:

- What the tool does
- When to use it vs. alternatives
- Key warnings (e.g., delete does NOT tear down deployed resources)
- Cross-references to related tools

I will draft these carefully for each tool during implementation.

---

## Implementation Order

Within Tier 1, I recommend implementing in this order (each step is independently testable):

1. **get_service** first -- establishes the resolve pattern and gRPC client wiring
2. **search_services** -- new pattern (generic search), good to validate early
3. **apply_service** -- idempotent create/update
4. **delete_service** -- reuses resolveServiceID from get
5. **disconnect_service_git_repo** and **configure_service_webhook** -- both simple ServiceId RPCs
6. **list_service_branches** -- simple ServiceId RPC with StringList response
7. **Server wiring** -- add Register call to server.go

Steps 1-6 can be done file-by-file and each produces a compilable package. Step 7 wires everything together.