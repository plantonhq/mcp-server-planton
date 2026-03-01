---
name: T04 Environment Full CRUD
overview: Add 4 new MCP tools (get, create, update, delete) to the existing Environment domain package, completing the full CRUD lifecycle. The `get_environment` tool supports dual-resolution (by ID or by org+slug).
todos:
  - id: get-env
    content: "Implement get_environment: get.go (dual-resolution domain functions) + tool definition/handler in tools.go + registration in register.go"
    status: completed
  - id: create-env
    content: "Implement create_environment: create.go (proto assembly + CommandController.Create) + tool definition/handler in tools.go + registration in register.go"
    status: completed
  - id: update-env
    content: "Implement update_environment: update.go (UpdateFields + read-modify-write) + tool definition/handler in tools.go + registration in register.go"
    status: completed
  - id: delete-env
    content: "Implement delete_environment: delete.go (CommandController.Delete) + tool definition/handler in tools.go + registration in register.go"
    status: completed
  - id: verify
    content: Update package doc in tools.go, verify compilation, check linter errors, update next-task.md session history
    status: completed
isProject: false
---

# T04: Environment Full CRUD

## Current State

The Environment domain at `[internal/domains/resourcemanager/environment/](internal/domains/resourcemanager/environment/)` has **1 tool**:

- `list_environments` — `FindAuthorized` RPC, returns only environments the caller can access within an org

## What We're Adding

4 new tools to complete the CRUD lifecycle, following the Organization pattern established in T03 with one deliberate deviation (dual-resolution get).

## Tools

### 1. `get_environment` (dual-resolution)

- **By ID**: Calls `EnvironmentQueryController.Get(EnvironmentId)`
- **By org+slug**: Calls `EnvironmentQueryController.GetByOrgBySlug(ByOrgBySlugRequest)`
- Input: `env_id` (optional) OR `org` + `slug` (optional pair)
- Validation: exactly one resolution path must be provided (env_id, or both org and slug)
- **Why dual**: Environments are naturally referenced by slug within an org. "Get the staging environment in acme" should work without an intermediate `list_environments` call.

### 2. `create_environment`

- Calls `EnvironmentCommandController.Create(Environment)`
- Input: `org` (required), `slug` (required), `name` (optional), `description` (optional)
- Assembles full Environment proto: `api_version: "resource-manager.planton.ai/v1"`, `kind: "Environment"`, `metadata.org`, `metadata.slug`, `metadata.name`, `spec.description`
- Auth: requires `environment_create` permission on the org

### 3. `update_environment`

- Read-modify-write via `EnvironmentQueryController.Get` then `EnvironmentCommandController.Update`
- Input: `env_id` (required), `name` (optional), `description` (optional)
- Both RPCs share a single gRPC connection (same pattern as Organization update)
- Empty-string-means-no-change semantics (same as Organization)
- **Note**: EnvironmentSpec only has `description`. Metadata has `name`. These are the only two updatable fields.

### 4. `delete_environment`

- Calls `EnvironmentCommandController.Delete(EnvironmentId)`
- Input: `env_id` (required)
- Tool description must warn about cascading cleanup: deleting an environment triggers cleanup of all stack-modules, microservices, secrets, and clusters deployed to it

## RPCs Excluded (and why)


| RPC                     | Reason                                                                                                                      |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| `apply`                 | We expose separate `create` + `update` for clearer LLM intent and distinct auth models (same rationale as Organization T03) |
| `find`                  | Requires `platform/operator` — consistent exclusion pattern                                                                 |
| `findByOrg`             | Returns ALL envs in org; `list_environments` already uses the superior `findAuthorized` which respects per-env permissions  |
| `checkSlugAvailability` | Deferred — not core CRUD, same as Organization                                                                              |


## Files to Create

All new files go in `[internal/domains/resourcemanager/environment/](internal/domains/resourcemanager/environment/)`:

- `**get.go`** — `Get(ctx, serverAddress, envID string)` and `GetByOrgBySlug(ctx, serverAddress, org, slug string)` domain functions. The handler in tools.go dispatches to one or the other based on input.
- `**create.go**` — `Create(ctx, serverAddress, org, slug, name, description string)` — assembles Environment proto, calls CommandController.Create
- `**update.go**` — `UpdateFields` struct + `Update(ctx, serverAddress, envID string, fields UpdateFields)` — read-modify-write pattern
- `**delete.go**` — `Delete(ctx, serverAddress, envID string)` — calls CommandController.Delete

## Files to Modify

- `**tools.go**` — Add input structs, tool definitions, and handlers for all 4 new tools. Update package doc from "One tool" to "Five tools".
- `**register.go**` — Add 4 new `mcp.AddTool` registrations

## Key Import Paths

- Environment stubs: `github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/environment/v1`
- Organization stubs (for `ByOrgBySlugRequest`): `github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1`
- Shared helpers: `github.com/plantonhq/mcp-server-planton/internal/domains`

## Design Deviation from Organization (documented)

The `get_environment` tool accepts **either** `env_id` **or** `(org, slug)`, whereas `get_organization` accepts only `org_id`. This is a deliberate choice because:

- Environments are child resources naturally referenced by slug within an org context
- The proto API explicitly provides `getByOrgBySlug` for this use case
- It eliminates a round-trip (`list_environments` then `get_environment`) for the most natural user query pattern

## Execution Approach

Implement one tool at a time in this order: get, create, update, delete. After all 4 are implemented, verify compilation and check for linter errors.