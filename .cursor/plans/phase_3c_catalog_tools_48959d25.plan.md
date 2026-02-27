---
name: Phase 3C Catalog Tools
overview: Add Deployment Component and IaC Module catalog tools to the MCP server. Proto analysis reveals a richer API surface than originally planned, warranting 4 tools (not 3), and an opportunity to reduce accumulated tech debt in the shared enum helpers.
todos:
  - id: shared-helpers
    content: "Add shared enum helpers to internal/domains/: enum.go (JoinEnumValues), provider.go (ResolveProvider, ResolveProvisioner). Update 3 existing consumers of local joinEnumValues to use the shared version."
    status: completed
  - id: deploymentcomponent-pkg
    content: "Create internal/domains/infrahub/deploymentcomponent/ with tools.go, search.go, get.go (2 tools: search_deployment_components, get_deployment_component)"
    status: completed
  - id: iacmodule-pkg
    content: "Create internal/domains/infrahub/iacmodule/ with tools.go, search.go, get.go (2 tools: search_iac_modules, get_iac_module)"
    status: completed
  - id: server-registration
    content: Register 4 new tools in server.go, update doc.go, update tool count to 63
    status: completed
  - id: verify-build
    content: Run go build, go vet, go test to verify clean build
    status: completed
isProject: false
---

# Phase 3C: Deployment Component and IaC Module Catalog

## Proto Analysis -- Surprises vs. Master Plan

The master plan specified **3 tools** against 2 gRPC services. Proto analysis reveals **3 gRPC services** with **7 relevant RPCs**, and a natural 4-tool design.

### Surprise 1: `DeploymentComponentQueryController` exists

The master plan only used `InfraHubSearchQueryController.SearchDeploymentComponentsByFilter` for deployment components. But a dedicated `DeploymentComponentQueryController` exists with three RPCs:

- `Get(ApiResourceId)` -- lookup by ID
- `GetByCloudResourceKind(CloudResourceKindRequest)` -- lookup by kind string (e.g., "AwsEksCluster")
- `Find(FindApiResourcesRequest)` -- paginated find

**Recommendation**: Add a `get_deployment_component` tool with dual identification (by ID or by kind). This completes the search-then-drill-down flow. The `GetByCloudResourceKind` RPC is especially valuable because agents already know kind strings from the `cloud-resource-kinds://catalog` resource.

### Surprise 2: IaC modules have both org-context and official search RPCs

The `InfraHubSearchQueryController` has two IaC module search RPCs:

- `SearchIacModulesByOrgContext` -- org-scoped, includes official + org modules, rich filters
- `SearchOfficialIacModules` -- public, no org required

This is the **exact same pattern** as the existing preset domain (`SearchCloudObjectPresetsByOrgContext` / `SearchOfficialCloudObjectPresets`). The `search_iac_modules` tool should dispatch to the appropriate RPC based on whether `org` is provided, exactly as [preset/search.go](internal/domains/infrahub/preset/search.go) does.

### Surprise 3: `FindDeploymentComponentIacModulesByOrgContext` RPC

A bridge RPC exists that finds IaC modules for a specific deployment component kind within an org. However, `SearchIacModulesByOrgContext` already accepts a `cloud_resource_kind` filter that achieves the same result with richer filtering. **Recommendation**: Do NOT create a separate tool for this -- the `search_iac_modules` tool with `kind` parameter subsumes it.

## Proposed Tools (4, up from planned 3)


| Tool                           | RPC(s)                                                       | Package                |
| ------------------------------ | ------------------------------------------------------------ | ---------------------- |
| `search_deployment_components` | `SearchDeploymentComponentsByFilter`                         | `deploymentcomponent/` |
| `get_deployment_component`     | `Get` + `GetByCloudResourceKind`                             | `deploymentcomponent/` |
| `search_iac_modules`           | `SearchIacModulesByOrgContext` or `SearchOfficialIacModules` | `iacmodule/`           |
| `get_iac_module`               | `IacModuleQueryController.Get`                               | `iacmodule/`           |


Server expands from **59 to 63 tools**.

## Shared Enum Helpers -- Debt Reduction

Both new packages need `resolveProvider` (CloudResourceProvider) and `resolveProvisioner` (IacProvisioner) helpers. Rather than duplicating per-package, add them to the shared `internal/domains/` package alongside the existing `ResolveKind` in [kind.go](internal/domains/kind.go).

Additionally, `joinEnumValues` is currently duplicated identically in 3 packages (`audit/enum.go`, `graph/enum.go`, `stackjob/enum.go`). Since we are adding new enum resolvers to the shared package anyway, lift `joinEnumValues` to `internal/domains/enum.go` and update the 3 existing consumers. This prevents a 4th copy and reduces accumulated tech debt.

## Package Structure

### `internal/domains/infrahub/deploymentcomponent/`

- `**tools.go`** -- 2 input structs, 2 tool definitions, 2 typed handlers
  - `SearchDeploymentComponentsInput`: search_text, page_num, page_size, provider (optional, single string like "aws")
  - `GetDeploymentComponentInput`: id (mutually exclusive with kind), kind (PascalCase CloudResourceKind string)
- `**search.go`** -- `SearchDeploymentComponentsByFilter` RPC call
  - Public endpoint, no org required
  - Optional provider filter resolved via `domains.ResolveProvider`
  - Pagination: 1-based tool API, 0-based proto (established convention)
- `**get.go**` -- Dispatches to `Get` or `GetByCloudResourceKind` based on input
  - Dual identification: by ID or by kind string
  - Kind resolution via `domains.ResolveKind`

### `internal/domains/infrahub/iacmodule/`

- `**tools.go**` -- 2 input structs, 2 tool definitions, 2 typed handlers
  - `SearchIacModulesInput`: search_text, page_num, page_size, org (optional), kind (optional), provisioner (optional), provider (optional)
  - `GetIacModuleInput`: id (required)
- `**search.go**` -- Dispatches to org-context or official RPC based on `org` presence
  - Follows [preset/search.go](internal/domains/infrahub/preset/search.go) pattern exactly
  - When org provided: `SearchIacModulesByOrgContext` with `IsIncludeOfficial: true`, `IsIncludeOrganizationModules: true`
  - When org absent: `SearchOfficialIacModules`
- `**get.go**` -- `IacModuleQueryController.Get` RPC call (standard get-by-ID pattern)

### `internal/domains/` (shared helpers)

- `**enum.go**` (new) -- `JoinEnumValues(m map[string]int32, exclude string) string` (exported, lifted from 3 existing duplicates)
- `**provider.go**` (new) -- `ResolveProvider(s string) (CloudResourceProvider, error)` and `ResolveProvisioner(s string) (IacProvisioner, error)`

### Existing file updates

- [internal/server/server.go](internal/server/server.go) -- 2 new imports, 4 `mcp.AddTool` calls, count 59->63
- [internal/domains/infrahub/doc.go](internal/domains/infrahub/doc.go) -- Add deploymentcomponent and iacmodule to subpackage list
- `internal/domains/audit/enum.go` -- Replace local `joinEnumValues` with `domains.JoinEnumValues`
- `internal/domains/graph/enum.go` -- Replace local `joinEnumValues` with `domains.JoinEnumValues`
- `internal/domains/infrahub/stackjob/enum.go` -- Replace local `joinEnumValues` with `domains.JoinEnumValues`

## Key Patterns to Follow

- **Preset search pattern** for IaC module search (org-context vs. official RPC dispatch)
- **InfraChart get pattern** for simple get-by-ID (`get_iac_module`)
- **InfraProject get pattern** for dual-identification get (`get_deployment_component`)
- **Shared `domains.ResolveKind`** for kind string resolution
- **1-based page numbers** in tool API, converted to 0-based for proto (established convention)

## Verification

After implementation: `go build ./...`, `go vet ./...`, `go test ./...` must all pass clean.