---
name: T03 Org Full CRUD
overview: Add 4 new MCP tools (get, create, update, delete) to the Organization domain, completing full CRUD lifecycle. Follows existing domain patterns with separate create/update tools matching the proto auth model.
todos:
  - id: get
    content: Create get.go with Get domain function calling OrganizationQueryController.Get
    status: completed
  - id: create
    content: Create create.go with Create domain function constructing Organization proto and calling OrganizationCommandController.Create
    status: completed
  - id: update
    content: Create update.go with UpdateFields struct and Update domain function doing read-modify-write via Get then OrganizationCommandController.Update
    status: completed
  - id: delete
    content: Create delete.go with Delete domain function calling OrganizationCommandController.Delete
    status: completed
  - id: tools
    content: Add 4 input structs, tool definitions, and handlers to tools.go; update package doc
    status: completed
  - id: register
    content: Add 4 mcp.AddTool calls to register.go
    status: completed
  - id: verify
    content: Run go build to verify compilation; run lints on all changed files
    status: completed
isProject: false
---

# T03: Organization Full CRUD

## Current State

The Organization domain at `internal/domains/resourcemanager/organization/` has 3 files and 1 tool (`list_organizations`). The proto API exposes 7 command RPCs and 4 query RPCs, but only `findOrganizations` is surfaced today.

## What We Are Adding

4 new tools, bringing the Organization domain from 1 to 5 tools:


| Tool | RPC | Input | Purpose |
| ---- | --- | ----- | ------- |


Scratch that -- no tables in plans. Here are the tools:

- `**get_organization**` -- calls `OrganizationQueryController.Get(OrganizationId)`. Input: `org_id` (required). Returns full Organization including metadata, spec, status.
- `**create_organization**` -- calls `OrganizationCommandController.Create(Organization)`. Input: `slug` (required), `name`, `description`, `contact_email` (optional). Constructs the full Organization proto internally (sets api_version, kind, metadata, spec). The `create` RPC skips auth -- any authenticated user can create an org.
- `**update_organization**` -- calls `OrganizationQueryController.Get` then `OrganizationCommandController.Update`. Input: `org_id` (required), `name`, `description`, `contact_email`, `logo_url` (all optional). Does a read-modify-write within a single `WithConnection` callback so both RPCs share one gRPC connection. Only non-empty fields are applied to the fetched org.
- `**delete_organization**` -- calls `OrganizationCommandController.Delete(OrganizationId)`. Input: `org_id` (required). Requires org delete permission (enforced server-side).

## RPCs Intentionally Excluded

- `repairFgaTuples` -- platform operator permission; consistent with existing policy (MCP tools = user-level RPCs only)
- `toggleGettingStartedTasks` / `toggleGettingStartedTask` -- UI onboarding concern, not relevant to MCP workflows
- `find` -- platform operator paginated search; `findOrganizations` (already exposed as `list_organizations`) covers the user-facing use case
- `apply` -- we use separate `create` + `update` (per design decision above)
- `checkSlugAvailability` -- useful but not core CRUD; can be added as a follow-up if needed

## Design Decision: Separate create + update (not apply)

Rationale documented in the plan and aligned with user:

- `create` skips auth, `update` requires org update permission -- different auth models
- Very different input shapes: create needs slug, update needs org_id + partial fields
- Clearer LLM intent signals ("create a new org" vs "modify this org")
- Organization is a single type, not polymorphic like CloudResource

## Key Pattern: Read-Modify-Write for Update

The `update` RPC expects a full `Organization` proto. The tool does:

1. GET current org by ID (query controller)
2. Merge caller-provided fields (only non-empty strings overwrite)
3. UPDATE with the merged org (command controller)

Both calls happen within one `domains.WithConnection` callback, sharing a single gRPC connection and timeout. This is the same pattern as CloudResource's multi-step operations (e.g., slug-path resolution then delete).

**Trade-off acknowledged**: empty string cannot be used to *clear* a field. This is acceptable for v1 -- the `omitempty` JSON tag means "not provided" and "empty" are indistinguishable. If field-clearing becomes a real need, we can add explicit `clear_`* booleans later.

## File Structure

All files in `internal/domains/resourcemanager/organization/`:

- `register.go` -- **MODIFY**: add 4 new `mcp.AddTool` calls
- `tools.go` -- **MODIFY**: add 4 new input structs, tool definitions, and handlers; update package doc comment (1 -> 5 tools)
- `list.go` -- **UNCHANGED**
- `get.go` -- **NEW**: `Get(ctx, serverAddress, orgID) (string, error)`
- `create.go` -- **NEW**: `Create(ctx, serverAddress, slug, name, description, contactEmail) (string, error)`
- `update.go` -- **NEW**: `Update(ctx, serverAddress, orgID, fields UpdateFields) (string, error)` with `UpdateFields` struct
- `delete.go` -- **NEW**: `Delete(ctx, serverAddress, orgID) (string, error)`

## Key Implementation Details

**Proto imports needed** (from existing stubs):

- `organizationv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/organization/v1"` (already imported in `list.go`)
- `apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"` (for `ApiResourceMetadata` in create)

**Difference from CloudResource delete**: Organization delete uses `OrganizationId{Value: orgID}`, not `ApiResourceDeleteInput{ResourceId: ...}`. This matches the proto definition in `command.proto`.

**Create proto construction**:

```go
&organizationv1.Organization{
    ApiVersion: "resource-manager.planton.ai/v1",
    Kind:       "Organization",
    Metadata:   &apiresource.ApiResourceMetadata{Slug: slug, Name: name},
    Spec:       &organizationv1.OrganizationSpec{Description: desc, ContactEmail: email},
}
```

The const values for `api_version` and `kind` come from the proto validation rules in `api.proto`.

## What Will NOT Change

- `internal/server/server.go` -- Organization is already registered; no changes needed
- `internal/domains/` shared utilities -- no changes needed
- `list.go` -- existing list tool unchanged

