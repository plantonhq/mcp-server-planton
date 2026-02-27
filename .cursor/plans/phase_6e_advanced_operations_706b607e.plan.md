---
name: Phase 6E Advanced Operations
overview: Implement 5 advanced operations tools (locks, rename, env var map, value references) in the cloudresource domain, expanding the MCP server from 13 to 18 tools. Two proto surprises discovered during analysis require design decisions before implementation.
todos:
  - id: resolve-surprises
    content: Resolve the two proto surprises (get_env_var_map input, resolve_value_references input) and dual-path design decision with user
    status: completed
  - id: impl-locks
    content: Implement list_cloud_resource_locks and remove_cloud_resource_locks in locks.go + tools.go
    status: completed
  - id: impl-rename
    content: Implement rename_cloud_resource in rename.go + tools.go
    status: completed
  - id: impl-envvarmap
    content: Implement get_env_var_map in envvarmap.go + tools.go (based on Surprise 1 resolution)
    status: completed
  - id: impl-references
    content: Implement resolve_value_references in references.go + tools.go (based on Surprise 2 resolution)
    status: completed
  - id: register-server
    content: Register all 5 new tools in server.go (13 to 18), update package doc comment in tools.go
    status: completed
  - id: verify-build
    content: Verify go build ./... and go test ./... pass with zero errors
    status: completed
isProject: false
---

# Phase 6E: Advanced Operations (5 tools, 13 to 18)

## Proto Surprises — Decisions Required Before Implementation

### Surprise 1: `get_env_var_map` — Input is raw YAML, not ID + manifest

The revised plan (T01_2) assumed:

- Input: `id` (required) + `manifest` (required, as map)

**Actual proto** (`GetEnvVarMapRequest` in [io.proto](apis/ai/planton/infrahub/cloudresource/v1/io.proto)):

```protobuf
message GetEnvVarMapRequest {
  string yaml_content = 1 [(buf.validate.field).required = true];
}
```

The RPC takes a single field: `yaml_content` — a raw YAML string of the cloud resource manifest. There is no `id` field. The proto comment on the query.proto RPC says: *"Authorization: Handled in request handler. Requires 'get' permission on the cloud resource resolved from the provided YAML. Cannot be configured via RPC options because the cloud resource ID must be extracted from the YAML first."*

**Impact**: The MCP tool's input is just a YAML string, not `id` + `manifest`. The server parses the YAML internally to determine the resource, check authorization, extract env vars/secrets, and resolve valueFrom references.

**Response** (`GetEnvVarMapResponse`):

- `variables`: map of env var name to resolved plain string value
- `secrets`: map of secret name to resolved plain string value (note: K8s secretRef resolution NOT included)

**Question**: The agent would need to provide raw YAML content. This could come from:

1. A file the agent is working with
2. Constructing it from a prior `get_cloud_resource` JSON response (converted to YAML)
3. A preset template

Should we expose this tool as-is (taking `yaml_content` string), or should we consider wrapping it to also accept an `id` and fetch + convert the resource internally? The raw YAML interface is what the backend provides, so wrapping adds complexity and potentially different semantics.

### Surprise 2: `resolve_value_references` — No `references` list, resolves ALL references

The revised plan (T01_2) assumed:

- Input: `cloud_resource_id` (required) + `references` (required, list of references to resolve)

**Actual proto** (`ResolveValueFromReferencesRequest` in [io.proto](apis/ai/planton/infrahub/cloudresource/v1/io.proto)):

```protobuf
message ResolveValueFromReferencesRequest {
  CloudResourceKind cloud_resource_kind = 1;
  string cloud_resource_id = 2;
}
```

There is no `references` field. The RPC takes `cloud_resource_kind` (enum) + `cloud_resource_id` (string). The server loads the full resource from the database and resolves ALL valueFrom references in its spec, returning the fully transformed YAML.

**Response** (`ResolveValueFromReferencesResponse`):

- `is_resolved`: bool — whether all references resolved successfully
- `errors`: list of error strings
- `diagnostics`: list of diagnostic messages
- `cloud_resource_yaml`: the fully transformed cloud resource as YAML

**Impact**: This is a "resolve everything for this resource" operation, not a "resolve these specific references" operation. The tool name `resolve_value_references` still fits, but the input is simpler than planned. The tool requires `kind` (for the enum), which means we use `domains.ResolveKind`.

---

## Design Decision: Dual-path (ID vs slug) for locks and rename

Within the `cloudresource` domain, all existing tools that identify a resource support the `ResourceIdentifier` dual-path: `id` alone OR `kind + org + env + slug`. This includes `get`, `delete`, and `destroy`.

The revised plan says locks and rename take only `id`. But the agent may only know a resource by its slug path. Requiring a separate `get_cloud_resource` call just to obtain the ID adds friction.

**Proposal**: Support the dual-path for `list_cloud_resource_locks`, `remove_cloud_resource_locks`, and `rename_cloud_resource`. Internally, use `resolveResourceID` to convert slug path to ID before calling the lock/rename RPCs. This adds one extra gRPC call on the slug path but is consistent with the domain convention.

For `resolve_value_references`: also support dual-path — the RPC needs both `cloud_resource_id` and `cloud_resource_kind`. If the agent provides the slug path, we already have the kind and can resolve the ID via `resolveResourceID`.

For `get_env_var_map`: dual-path does not apply — the RPC takes YAML content, not a resource identifier.

---

## Tool-by-Tool Implementation

### Tool 1: `list_cloud_resource_locks`

- **RPC**: `CloudResourceLockController.listLocks(CloudResourceId) returns (CloudResourceLockInfo)`
- **New gRPC client**: `NewCloudResourceLockControllerClient` (first use of lock controller)
- **Input**: `ResourceIdentifier` dual-path (pending decision above; falls back to `id`-only if rejected)
- **Domain function**: `ListLocks(ctx, serverAddress, id ResourceIdentifier) (string, error)`
  - Uses `resolveResourceID` for slug path
  - Creates lock controller client, calls `ListLocks`
  - Returns `CloudResourceLockInfo` as JSON
- **File**: [internal/domains/cloudresource/locks.go](internal/domains/cloudresource/locks.go) (new)

### Tool 2: `remove_cloud_resource_locks`

- **RPC**: `CloudResourceLockController.removeLocks(CloudResourceId) returns (CloudResourceLockRemovalResponse)`
- **Input**: Same as `list_cloud_resource_locks`
- **Domain function**: `RemoveLocks(ctx, serverAddress, id ResourceIdentifier) (string, error)`
- **Tool description warning**: Must include caution about state corruption risk
- **File**: shares [internal/domains/cloudresource/locks.go](internal/domains/cloudresource/locks.go) with `list_cloud_resource_locks`

### Tool 3: `rename_cloud_resource`

- **RPC**: `CloudResourceCommandController.rename(RenameCloudResourceRequest) returns (CloudResource)`
- **Proto fields**: `id` (string) + `name` (string)
- **MCP input fields**: Uses `ResourceIdentifier` dual-path + `new_name` (required). Note: proto field is `name` but MCP field is `new_name` for clarity
- **Domain function**: `Rename(ctx, serverAddress, id ResourceIdentifier, newName string) (string, error)`
  - Uses `resolveResourceID` for slug path
  - Creates command controller client, calls `Rename` with `RenameCloudResourceRequest{Id: resolvedID, Name: newName}`
- **File**: [internal/domains/cloudresource/rename.go](internal/domains/cloudresource/rename.go) (new)

### Tool 4: `get_env_var_map` (pending Surprise 1 resolution)

- **RPC**: `CloudResourceQueryController.getEnvVarMap(GetEnvVarMapRequest) returns (GetEnvVarMapResponse)`
- **Input**: `yaml_content` (required string)
- **Domain function**: `GetEnvVarMap(ctx, serverAddress, yamlContent string) (string, error)`
- **File**: [internal/domains/cloudresource/envvarmap.go](internal/domains/cloudresource/envvarmap.go) (new)

### Tool 5: `resolve_value_references` (updated per Surprise 2)

- **RPC**: `CloudResourceQueryController.resolveValueFromReferences(ResolveValueFromReferencesRequest) returns (ResolveValueFromReferencesResponse)`
- **Input**: `ResourceIdentifier` dual-path + `kind` (required PascalCase string). When using ID path, `kind` is still required (proto needs both). When using slug path, `kind` is already in the identifier.
- **Domain function**: `ResolveValueReferences(ctx, serverAddress, id ResourceIdentifier) (string, error)`
  - Slug path: kind comes from identifier, ID resolved via `resolveResourceID`
  - ID path: kind is required as separate input
- **File**: [internal/domains/cloudresource/references.go](internal/domains/cloudresource/references.go) (new)

---

## Files Changed

**New files** (4):

- `internal/domains/cloudresource/locks.go` — `ListLocks` + `RemoveLocks` domain functions
- `internal/domains/cloudresource/rename.go` — `Rename` domain function
- `internal/domains/cloudresource/envvarmap.go` — `GetEnvVarMap` domain function
- `internal/domains/cloudresource/references.go` — `ResolveValueReferences` domain function

**Modified files** (2):

- `internal/domains/cloudresource/tools.go` — 5 new input structs, tool defs, handlers; package doc comment updated (6 to 11 tools)
- `internal/server/server.go` — 5 new `mcp.AddTool` calls, import unchanged (same package), tool count 13 to 18

---

## Patterns Reused (no new patterns needed)

- `domains.WithConnection` for gRPC lifecycle
- `domains.RPCError` for error classification
- `domains.TextResult` for wrapping responses
- `domains.MarshalJSON` for protojson serialization
- `domains.ResolveKind` for PascalCase kind to proto enum
- `resolveResourceID` for slug-path to ID resolution (from [identifier.go](internal/domains/cloudresource/identifier.go))
- `validateIdentifier` for input validation
- `describeIdentifier` for error context

## Testing

- No new pure domain logic that warrants unit tests in this phase (all tools are thin RPC wrappers with existing validation helpers)
- If dual-path is adopted for locks/rename/references, the existing `resolveResourceID` and `validateIdentifier` are already tested
- Integration testing deferred to Hardening phase

