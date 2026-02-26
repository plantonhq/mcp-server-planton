# Phase 4: Get and Delete Cloud Resource Tools

**Date**: February 27, 2026

## Summary

Implemented `get_cloud_resource` and `delete_cloud_resource` MCP tools, completing the cloud resource tool set (apply/get/delete). Both tools support dual-path resource identification: by system-assigned ID or by the composite key (kind, org, env, slug), matching the Planton backend's gRPC API surface exactly.

## Problem Statement

The MCP server only exposed `apply_cloud_resource` (Phase 3), leaving agents unable to retrieve or remove existing resources. The full CRUD lifecycle requires all three operations.

### Pain Points

- Agents could create/update cloud resources but not inspect or delete them
- No programmatic way for agents to look up a resource by its human-readable identifiers
- The tool set was incomplete, blocking downstream workflows that depend on get/delete

## Solution

Added two new MCP tools with a shared `ResourceIdentifier` abstraction that supports two mutually exclusive identification paths:

1. **ID path**: single `id` field, routes directly to `QueryController.Get` or `CommandController.Delete`
2. **Slug path**: four fields (`kind`, `org`, `env`, `slug`), routes to `QueryController.GetByOrgByEnvByKindBySlug` for get, or resolves to ID then deletes for delete

The design was driven by the discovery that slug uniqueness is scoped to (org, env, kind) in the Planton backend, requiring all four fields when using the slug path. This was confirmed during planning and aligns with the gRPC `CloudResourceByOrgByEnvByKindBySlugRequest` message definition.

## Implementation Details

### New Files (3)

- **`internal/domains/cloudresource/identifier.go`** (~110 lines) — `ResourceIdentifier` struct, `validateIdentifier()` ensuring exactly one path is fully specified (with clear error messages for partial inputs), `describeIdentifier()` for human-readable error context, and `resolveResourceID()` for the delete slug-path's two-step resolution pattern.

- **`internal/domains/cloudresource/get.go`** (~50 lines) — `Get()` function with two distinct code paths within a single `domains.WithConnection` call. ID path calls `QueryController.Get(CloudResourceId)`, slug path resolves the PascalCase kind string to the proto enum via `resolveKind()` then calls `QueryController.GetByOrgByEnvByKindBySlug()`. No unnecessary resolution step — both RPCs return the full `CloudResource`.

- **`internal/domains/cloudresource/delete.go`** (~35 lines) — `Delete()` function following the Stigmer two-step pattern: resolve identifier to resource ID (via `resolveResourceID` which handles both paths), then call `CommandController.Delete(ApiResourceDeleteInput)`. Both calls share a single gRPC connection.

### Modified Files (2)

- **`internal/domains/cloudresource/tools.go`** (+97 lines) — Added `GetCloudResourceInput`, `GetTool()`, `GetHandler()`, `DeleteCloudResourceInput`, `DeleteTool()`, `DeleteHandler()`. Input structs have 5 optional fields with `jsonschema` annotations guiding agents on the dual-path identification pattern. Validation happens at the handler boundary.

- **`internal/server/server.go`** (+4/-6 lines) — Registered both new tools, replaced Phase 4 placeholder comments with actual registrations, updated tool count log to 3 with named tool list.

### Design Decisions

- **Slug path requires `kind`**: Confirmed during planning that slug is unique within (org, env, kind), not globally. The backend's `getByOrgByEnvByKindBySlug` RPC requires all four fields.
- **Delete kept simple**: Skipped `version_message` and `force` fields from `ApiResourceDeleteInput` — can be added later if needed.
- **Error handling**: Kind validation errors pass through directly (already user-friendly); gRPC errors go through `domains.RPCError()` for classification. `resolveResourceID` owns its error formatting to avoid double-wrapping.
- **Validation at the boundary**: Handlers validate the `ResourceIdentifier` before calling `Get()`/`Delete()`, which assume valid inputs.

## Benefits

- **Complete tool set**: Agents can now apply, get, and delete cloud resources — full lifecycle coverage
- **Flexible identification**: Agents can use whichever identifier they have (ID from a previous response, or the human-readable composite key they used during apply)
- **Clear error guidance**: Partial slug-path inputs get specific messages listing which fields are missing
- **Consistent patterns**: Follows established Phase 3 conventions (handler → validation → gRPC function → response) and Stigmer reference architecture

## Impact

- MCP server now exposes 3 tools + 2 resources (catalog + schema template)
- Agents can execute full cloud resource management workflows end-to-end
- Phase 4 completes the core tool set; only Phase 5 (Testing + Documentation) remains

## Related Work

- Phase 3: apply_cloud_resource + MCP Resource Templates (prerequisite)
- Kind Catalog Resource (discovery workflow used by all tools)
- Stigmer MCP server delete/get patterns (reference architecture)

---

**Status**: Production Ready
**Timeline**: ~1 hour (planning + implementation)
