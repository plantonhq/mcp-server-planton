---
name: Phase 1B InfraProject Tools
overview: Add 6 MCP tools for InfraProject lifecycle management (search, get, apply, delete, slug check, undeploy), expanding the server from 21 to 27 tools. Follows established patterns from infrachart/cloudresource packages with domain-appropriate simplifications.
todos:
  - id: get-and-helpers
    content: Create get.go with Get(), resolveProjectID(), resolveProject() -- foundation helpers
    status: completed
  - id: search
    content: Create search.go with Search() via InfraHubSearchQueryController.SearchInfraProjects
    status: completed
  - id: slug
    content: Create slug.go with CheckSlugAvailability()
    status: completed
  - id: apply
    content: Create apply.go with Apply() using protojson.Unmarshal passthrough
    status: completed
  - id: delete
    content: Create delete.go with Delete() using resolveProjectID()
    status: completed
  - id: undeploy
    content: Create undeploy.go with Undeploy() using resolveProjectID()
    status: completed
  - id: tools
    content: Create tools.go with package doc, 6 input structs, 6 Tool/Handler pairs
    status: completed
  - id: server-registration
    content: Update server.go (6 AddTool calls, count 21->27) and doc.go
    status: completed
  - id: verify-build
    content: Run go build, go vet, go test to verify clean build
    status: completed
isProject: false
---

# Phase 1B: InfraProject Tools (6 tools, 21 -> 27)

## Domain Analysis

InfraProject is the unit of infrastructure deployment on Planton. It has two source types:

- **infra_chart**: Created from an InfraChart template with parameter overrides (template YAML + params)
- **git_repo**: Backed by a Git repository with webhook-driven pipelines

Unlike CloudResource (362 polymorphic kinds needing generated parsers), InfraProject is a **single typed proto** with a well-defined spec. This means simpler identification (no `kind` dimension) and direct proto passthrough for apply.

### Proto Surface (from planton stubs)

**Query Controller** (`InfraProjectQueryControllerClient`):

- `Get(InfraProjectId)` -> `InfraProject`
- `GetByOrgBySlug(ApiResourceByOrgBySlugRequest)` -> `InfraProject`
- `CheckSlugAvailability(InfraProjectSlugAvailabilityCheckRequest)` -> `InfraProjectSlugAvailabilityCheckResponse`
- `Find(FindApiResourcesRequest)` -> `InfraProjectList` (pagination, not used -- search is more powerful)
- `Build(InfraProject)` -> `InfraProject` (not exposed in this phase)

**Command Controller** (`InfraProjectCommandControllerClient`):

- `Apply(InfraProject)` -> `InfraProject`
- `Delete(ApiResourceDeleteInput)` -> `InfraProject`
- `Undeploy(InfraProjectId)` -> `InfraProject`
- `Purge(InfraProjectId)` -> `InfraProject` (not exposed -- see design decision below)

**Search Controller** (`InfraHubSearchQueryControllerClient`):

- `SearchInfraProjects(SearchInfraProjectsRequest)` -> `ApiResourceSearchRecordList`

---

## Design Decisions

### DD-1: Identification pattern -- simpler than CloudResource

CloudResource uses a 4-field slug path (`kind`, `org`, `env`, `slug`) because slugs are scoped to (org, env, kind). InfraProject slugs are scoped to **org only** -- so identification needs just:

- **ID path**: `id` alone
- **Slug path**: `org` + `slug`

No need for the heavyweight `ResourceIdentifier` from cloudresource. A lighter `resolveProjectID` helper suffices.

### DD-2: Apply accepts full InfraProject JSON (confirmed)

The tool accepts a JSON object matching the `InfraProject` proto shape. The handler uses `protojson.Unmarshal` to create the proto, then calls `Apply`. This:

- Supports both source types (infra_chart and git_repo)
- Makes `get_infra_project` output directly usable as `apply_infra_project` input
- Avoids leaky abstractions

### DD-3: Purge is intentionally excluded

The backend has a `Purge` RPC (undeploy all resources + delete record atomically). We deliberately exclude it because:

- It's a destructive, irreversible compound operation
- An AI agent can achieve the same result via `undeploy_infra_project` followed by `delete_infra_project`, giving the human a review checkpoint between steps
- Matches the AD-01 principle of keeping dangerous operations bounded

### DD-4: Search over Find

The query controller has both `Find` (pagination-only) and the search controller has `SearchInfraProjects` (org + env + free-text + pagination). We expose **search** because it's strictly more powerful and what an agent naturally reaches for.

---

## Tool Specifications

### 1. `search_infra_projects`

- **RPC**: `InfraHubSearchQueryController.SearchInfraProjects`
- **Input**: `org` (required), `env` (optional), `search_text` (optional), `page_num` (optional, 1-based), `page_size` (optional)
- **Response**: `ApiResourceSearchRecordList` -- lightweight records with IDs for follow-up `get_infra_project` calls
- **Pattern**: Follows [preset/search.go](internal/domains/infrahub/preset/search.go) but uses a different search method

### 2. `get_infra_project`

- **RPC**: `InfraProjectQueryController.Get` or `GetByOrgBySlug`
- **Input**: `id` (mutually exclusive with org+slug), `org`, `slug`
- **Response**: Full `InfraProject` proto as JSON (metadata, spec with source config, status with rendered YAML and DAG)
- **Pattern**: Follows [infrachart/get.go](internal/domains/infrahub/infrachart/get.go) but with dual identification

### 3. `apply_infra_project`

- **RPC**: `InfraProjectCommandController.Apply`
- **Input**: `infra_project` (JSON object matching InfraProject proto)
- **Deserialization**: `encoding/json.Marshal` the map -> `protojson.Unmarshal` into `*InfraProject`
- **Response**: The applied `InfraProject` as JSON (includes server-assigned ID, audit info)
- **Pattern**: New pattern -- no existing tool does protojson passthrough. CloudResource's apply is too different (polymorphic parsing)

### 4. `delete_infra_project`

- **RPC**: `InfraProjectCommandController.Delete`
- **Input**: `id` or `org` + `slug` (resolved to ID via query controller)
- **Response**: The deleted `InfraProject` as JSON
- **Pattern**: Follows [cloudresource/delete.go](internal/domains/infrahub/cloudresource/delete.go) -- resolve ID then call Delete with `ApiResourceDeleteInput`

### 5. `check_infra_project_slug`

- **RPC**: `InfraProjectQueryController.CheckSlugAvailability`
- **Input**: `org` (required), `slug` (required)
- **Response**: `InfraProjectSlugAvailabilityCheckResponse` (is_available boolean)
- **Note**: Simpler than CloudResource's slug check -- no kind/env dimensions. Uses `InfraProjectSlugAvailabilityCheckRequest{Org, Slug}`

### 6. `undeploy_infra_project`

- **RPC**: `InfraProjectCommandController.Undeploy`
- **Input**: `id` or `org` + `slug` (resolved to ID)
- **Response**: The `InfraProject` as JSON (status.pipeline_id will be set to the triggered pipeline)
- **Pattern**: Follows [cloudresource/destroy.go](internal/domains/infrahub/cloudresource/destroy.go) conceptually (teardown infra, keep record)

---

## File Plan

### New files (7) -- all under `internal/domains/infrahub/infraproject/`

- `**tools.go`** -- Package doc, 6 input structs, 6 `*Tool()` functions, 6 `*Handler()` functions. ~200 lines.
- `**search.go`** -- `Search()` function calling `InfraHubSearchQueryController.SearchInfraProjects`. ~45 lines.
- `**get.go**` -- `Get()` function with dual identification (ID or org+slug). Includes `resolveProject()` and `resolveProjectID()` helpers. ~60 lines.
- `**apply.go**` -- `Apply()` function with `protojson.Unmarshal` deserialization. ~40 lines.
- `**delete.go**` -- `Delete()` function resolving ID then calling command controller. ~30 lines.
- `**slug.go**` -- `CheckSlugAvailability()` function. ~25 lines.
- `**undeploy.go**` -- `Undeploy()` function resolving ID then calling command controller. ~30 lines.

### Modified files (2)

- `**[internal/server/server.go](internal/server/server.go)**` -- Add import for `infraproject` package, 6 `mcp.AddTool()` calls, update tool count from 21 to 27, add 6 tool names to the log list.
- `**[internal/domains/infrahub/doc.go](internal/domains/infrahub/doc.go)**` -- Add `infraproject` to the subpackage list.

---

## Implementation Order

The tools should be implemented in dependency order:

1. `**get.go**` first -- contains `resolveProjectID()` and `resolveProject()` helpers that `delete.go` and `undeploy.go` depend on
2. `**search.go**` -- independent, can be done in parallel with get
3. `**slug.go**` -- independent, simple
4. `**apply.go**` -- independent (no shared helpers)
5. `**delete.go**` -- depends on `resolveProjectID()` from get.go
6. `**undeploy.go**` -- depends on `resolveProjectID()` from get.go
7. `**tools.go**` -- ties everything together with input structs and handlers
8. **Server registration** -- wire into server.go + update doc.go

---

## Proto Imports Required

```
infraprojectv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infraproject/v1"
apiresource    "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
rpc            "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"
infrahubsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/infrahub"
```

---

## Verification Criteria

- `go build ./...` passes
- `go vet ./...` passes
- `go test ./...` passes
- Server starts and reports 27 registered tools
- All 6 tool names appear in the startup log

