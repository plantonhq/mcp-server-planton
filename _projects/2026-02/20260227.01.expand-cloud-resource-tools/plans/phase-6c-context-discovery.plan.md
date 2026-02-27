---
name: Phase 6C Context Discovery
overview: Add two context-discovery tools (list_organizations, list_environments) in two new domain packages, keeping the flat package structure rather than mirroring the proto hierarchy.
todos:
  - id: create-organization-pkg
    content: Create internal/domains/organization/ package with list.go (domain function) and tools.go (tool definition + handler)
    status: completed
  - id: create-environment-pkg
    content: Create internal/domains/environment/ package with list.go (domain function) and tools.go (tool definition + handler)
    status: completed
  - id: register-tools
    content: Update internal/server/server.go to import and register both new tools (count 8 -> 10)
    status: completed
  - id: verify-build-test
    content: Run go build ./... and go test ./... to verify everything compiles and passes
    status: completed
isProject: false
---

# Phase 6C: Context Discovery

## Architecture Decision: Flat Domain Packages

The proto definitions are organized hierarchically:

- `infrahub/` contains `cloudresource/`, `stackjob/`, `cloudobjectpreset/`
- `resourcemanager/` contains `organization/`, `environment/`, `promotionpolicy/`

**Recommendation: Keep `internal/domains/` flat.** Here is the reasoning:

**1. Total package count is small.** At full expansion (through Phase 6E), we will have ~5 domain packages: `cloudresource`, `stackjob`, `organization`, `environment`, `preset`. Five flat entries are trivially navigable.

**2. Go idiom favors flat packages.** Deep nesting is an anti-pattern in Go. The standard library uses `net/http`, not `protocols/network/http/server`. Even the generated Go stubs use flat package names (`organizationv1`, `environmentv1`) rather than nesting.

**3. The MCP server is an adapter layer, not the platform.** The proto hierarchy is the authoritative structure for the platform's domain model. Duplicating it in a thin adapter creates coupling: if proto paths are reorganized, the MCP server's internal layout would need to change too, for no functional reason.

**4. No shared code within groupings.** There is nothing that `cloudresource` and `stackjob` share with each other that they don't also share with `organization`. All shared utilities live in `internal/domains/` (conn.go, marshal.go, rpcerr.go, toolresult.go) and serve every domain equally. A `resourcemanager/` or `infrahub/` intermediate package would be empty — a directory with no `.go` files, which is meaningless in Go.

**5. Phase 6E confirms this.** All Phase 6E tools (locks, rename, envvarmap, references) call `CloudResourceCommandController` or `CloudResourceQueryController` RPCs. They belong in the existing `cloudresource` package. No new domain packages are needed for 6E. Phase 6D adds `preset/` (1 new package). Final count: 6 packages. Still flat-friendly.

**Provenance documentation:** Each package's doc comment will state which proto service(s) it wraps, giving discoverability without structural overhead. Example:

```go
// Package organization provides the MCP tools for the Organization domain,
// backed by the OrganizationQueryController RPCs
// (ai.planton.resourcemanager.organization.v1) on the Planton backend.
package organization
```

### Projected final layout

```
internal/domains/
    conn.go, marshal.go, rpcerr.go, toolresult.go   (shared)
    cloudresource/   (infrahub — 5 tools now, grows in 6E)
    stackjob/        (infrahub — 3 tools)
    organization/    (resourcemanager — Phase 6C)
    environment/     (resourcemanager — Phase 6C)
    preset/          (infrahub — Phase 6D, later)
```

---

## Phase 6C Implementation

### Tool 1: `list_organizations`

**RPC:** `OrganizationQueryController.FindOrganizations(CustomEmpty) -> Organizations`

- No input parameters — returns all orgs the authenticated caller is a member of
- No authorization annotation — membership-scoped implicitly
- Response: `Organizations { Entries []*Organization }` — each with id, name, slug, spec, status

**Files to create in `internal/domains/organization/`:**

- `**list.go`** — Domain function `List(ctx, serverAddress) (string, error)` calling `FindOrganizations` via `domains.WithConnection`. Marshals the response with `domains.MarshalJSON`.
- `**tools.go`** — Package doc comment (proto provenance), empty `ListOrganizationsInput` struct, `ListTool() *mcp.Tool` definition, `ListHandler(serverAddress)` returning the handler. Since there are no input parameters, the handler simply calls `List()` directly — no validation needed.

**Import path for gRPC client:**

```go
organizationv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1"
```

**Client call:**

```go
client := organizationv1.NewOrganizationQueryControllerClient(conn)
resp, err := client.FindOrganizations(ctx, &protobuf.CustomEmpty{})
```

### Tool 2: `list_environments`

**RPC:** `EnvironmentQueryController.FindAuthorized(OrganizationId) -> Environments`

- Input: `org` (required) — organization identifier
- Auth: requires `get` permission on the organization
- Returns only environments where the caller has at least `get` permission (FGA-filtered)
- Response: `Environments { Entries []*Environment }` — each with id, name, slug, spec, status

**Files to create in `internal/domains/environment/`:**

- `**list.go`** — Domain function `List(ctx, serverAddress, org string) (string, error)` calling `FindAuthorized` via `domains.WithConnection`. Constructs `&organizationv1.OrganizationId{Value: org}` as the input.
- `**tools.go`** — Package doc comment (proto provenance), `ListEnvironmentsInput` struct with required `org` field, `ListTool() *mcp.Tool` definition, `ListHandler(serverAddress)` with `org` validation.

**Import paths:**

```go
environmentv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/environment/v1"
organizationv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1"
```

**Client call:**

```go
client := environmentv1.NewEnvironmentQueryControllerClient(conn)
resp, err := client.FindAuthorized(ctx, &organizationv1.OrganizationId{Value: org})
```

### Server Registration

**File:** [internal/server/server.go](internal/server/server.go)

- Import `organization` and `environment` packages
- Register both tools in `registerTools()`
- Update tool count from 8 to 10 in the log statement

### Unit Tests

Both tools are thin RPC wrappers with no domain logic (no enum resolution, no identifier validation, no kind conversion). There is nothing to unit-test beyond what would effectively be testing the gRPC framework itself. Tests will be added when logic is introduced (e.g., if we add response formatting or filtering).

### Verification

- `go build ./...` must pass
- `go test ./...` must pass
- `ReadLints` on all new/modified files

---

## Phase 6E Confirmation (No New Packages)

All Phase 6E tools are operations on existing CloudResource RPCs:

- `lock_cloud_resource` / `unlock_cloud_resource` — CloudResourceCommandController
- `rename_cloud_resource` — CloudResourceCommandController
- `get_env_var_map` — CloudResourceQueryController
- `resolve_cloud_resource_references` — CloudResourceQueryController

These will be new files within the existing `internal/domains/cloudresource/` package. No new domain packages needed.

## Phase 6D Quick Note (Future)

Phase 6D adds `search_cloud_object_presets` and `get_cloud_object_preset`. These use `CloudObjectPresetQueryController` (infrahub proto path). This will be a new `internal/domains/preset/` package — still flat, consistent with the pattern.