---
name: Tier 4+5 MCP Tools
overview: Implement 9 MCP tools across 3 ServiceHub bounded contexts (DnsDomain, TektonPipeline, TektonTask) following established patterns, completing the 35-tool ServiceHub project.
todos:
  - id: dnsdomain
    content: "DnsDomain package: 5 files (register, tools, get, apply, delete) -- 3 tools following Service pattern with DnsDomainId, ApiResourceByOrgBySlugRequest, ApiResourceDeleteInput"
    status: completed
  - id: tektonpipeline
    content: "TektonPipeline package: 5 files (register, tools, get, apply, delete) -- 3 tools with ApiResourceId, GetByOrgAndNameInput (slug->Name), entity-based delete"
    status: completed
  - id: tektontask
    content: "TektonTask package: 5 files (register, tools, get, apply, delete) -- 3 tools mirroring TektonPipeline structure"
    status: completed
  - id: server-wiring
    content: Wire all 3 packages into server.go with servicehubdnsdomain/servicehubtektonpipeline/servicehubtektontask aliases
    status: completed
  - id: verify-build
    content: Run go build/vet/test to verify clean compilation
    status: completed
isProject: false
---

# Tier 4+5: DnsDomain + TektonPipeline + TektonTask MCP Tools

## Scope

9 MCP tools across 3 entities, completing the ServiceHub domain:

- **DnsDomain**: `get_dns_domain`, `apply_dns_domain`, `delete_dns_domain`
- **TektonPipeline**: `get_tekton_pipeline`, `apply_tekton_pipeline`, `delete_tekton_pipeline`
- **TektonTask**: `get_tekton_task`, `apply_tekton_task`, `delete_tekton_task`

## Design Decisions

- **DD-T4-1 (slug normalization)**: All three entities present the dual-path lookup field as `slug` in MCP tool inputs, consistent with existing Service/VariablesGroup/SecretsGroup tools. For TektonPipeline/TektonTask, the `slug` value is passed to the proto's `GetByOrgAndNameInput.Name` field (the server converts name to slug internally, so slug values work directly).
- **DD-T4-2 (Tekton delete)**: Add delete tools for TektonPipeline and TektonTask for a complete CRUD surface (originally excluded in T01 plan). These RPCs take the full entity as input, requiring a **get-then-delete** pattern (unlike DnsDomain which uses `ApiResourceDeleteInput{resource_id}`).
- **DD-T4-3 (no search tools)**: None of these three entities have search RPCs in the `ServiceHubSearchQueryController`. No search tools will be added.
- **DD-T4-4 (separate bounded contexts)**: Per DD-T3-7 precedent, TektonPipeline and TektonTask remain separate packages despite near-identical structure. No shared abstraction.

## Proto API Surface (verified)

### DnsDomain

- Query: `get(DnsDomainId)`, `getByOrgBySlug(ApiResourceByOrgBySlugRequest)`
- Command: `apply(DnsDomain) -> DnsDomain`, `delete(ApiResourceDeleteInput) -> DnsDomain`
- Spec: `domain_name` (required, regex-validated), `description`

### TektonPipeline

- Query: `get(ApiResourceId)`, `getByOrgAndName(GetByOrgAndNameInput{org, name})`
- Command: `apply(TektonPipeline) -> TektonPipeline`, `delete(TektonPipeline) -> TektonPipeline`
- Spec: `selector` (required), `description`, `git_repo`, `yaml_content`, `overview_markdown`

### TektonTask

- Query: `get(ApiResourceId)`, `getByOrgAndName(GetByOrgAndNameInput{org, name})`
- Command: `apply(TektonTask) -> TektonTask`, `delete(TektonTask) -> TektonTask`
- Spec: `selector` (required), `description`, `git_repo`, `yaml_content`, `overview_markdown`

## File Plan (15 new files + 1 modified)

### Package: `internal/domains/servicehub/dnsdomain/` (5 files)

Pattern: mirrors [internal/domains/servicehub/service/](internal/domains/servicehub/service/) exactly.

- `**register.go`** -- Register function wiring 3 tools via `mcp.AddTool`
- `**tools.go**` -- `GetDnsDomainInput`, `ApplyDnsDomainInput`, `DeleteDnsDomainInput` structs + Tool/Handler pairs + `validateIdentification`
  - Get input: `id` OR `org`+`slug` (dual-path)
  - Apply input: `dns_domain` (full JSON object)
  - Delete input: `id` OR `org`+`slug` (dual-path)
- `**get.go**` -- `Get`, `resolveDomain` (by ID via `DnsDomainId` or org+slug via `ApiResourceByOrgBySlugRequest`), `resolveDomainID`, `describeDomain`
- `**apply.go**` -- `Apply` via `json.Marshal` -> `protojson.Unmarshal` -> `DnsDomainCommandController.Apply`
- `**delete.go**` -- `Delete` via `resolveDomainID` -> `DnsDomainCommandController.Delete(ApiResourceDeleteInput)`

Go stub imports:

```
dnsdomain "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/dnsdomain/v1"
apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
```

### Package: `internal/domains/servicehub/tektonpipeline/` (5 files)

- `**register.go**` -- Register function wiring 3 tools
- `**tools.go**` -- `GetTektonPipelineInput`, `ApplyTektonPipelineInput`, `DeleteTektonPipelineInput` + Tool/Handler pairs + `validateIdentification`
  - Get input: `id` OR `org`+`slug` (dual-path, slug passed as `Name` to proto)
  - Apply input: `tekton_pipeline` (full JSON object)
  - Delete input: `id` OR `org`+`slug` (dual-path)
- `**get.go**` -- `Get`, `resolvePipeline` (by ID via `ApiResourceId` or org+slug via `GetByOrgAndNameInput{Org, Name: slug}`), `describePipeline`
- `**apply.go**` -- `Apply` via `json.Marshal` -> `protojson.Unmarshal` -> `TektonPipelineCommandController.Apply`
- `**delete.go**` -- `Delete` via `resolvePipeline` (get full entity) -> `TektonPipelineCommandController.Delete(fullEntity)` (get-then-delete pattern)

Go stub imports:

```
tektonpipelinev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/tektonpipeline/v1"
apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
```

### Package: `internal/domains/servicehub/tektontask/` (5 files)

Structurally identical to tektonpipeline, with entity-specific naming.

- `**register.go**` -- Register function wiring 3 tools
- `**tools.go**` -- `GetTektonTaskInput`, `ApplyTektonTaskInput`, `DeleteTektonTaskInput` + Tool/Handler pairs + `validateIdentification`
- `**get.go**` -- `Get`, `resolveTask`, `describeTask`
- `**apply.go**` -- `Apply`
- `**delete.go**` -- `Delete` via get-then-delete pattern

Go stub imports:

```
tektontaskv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/tektontask/v1"
apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
```

### Server wiring: `internal/server/server.go`

Add 3 import aliases and 3 `Register` calls:

```go
servicehubdnsdomain "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/dnsdomain"
servicehubtektonpipeline "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/tektonpipeline"
servicehubtektontask "github.com/plantonhq/mcp-server-planton/internal/domains/servicehub/tektontask"
```

## Implementation Order

Execute sequentially, verifying build after each entity:

1. **DnsDomain** (3 tools) -- most straightforward, follows Service pattern exactly
2. **TektonPipeline** (3 tools) -- introduces `GetByOrgAndNameInput` and entity-based delete patterns
3. **TektonTask** (3 tools) -- mirrors TektonPipeline
4. **Server wiring** -- add all 3 Register calls + verify `go build ./...`

## Verification

After all 3 entities are implemented:

- `go build ./...` -- clean compilation
- `go vet ./...` -- no issues
- `go test ./...` -- all tests pass

