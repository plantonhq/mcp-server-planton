---
name: Tier 3 VariablesGroup SecretsGroup
overview: Implement 16 MCP tools across two symmetric ServiceHub entities (VariablesGroup and SecretsGroup), covering group-level CRUD, entry-level mutations, value resolution, search, and config key transformation -- all following the established patterns from Tier 1/2.
todos:
  - id: variablesgroup-tools
    content: Create tools.go with 8 input structs, 8 Tool/Handler pairs, and validateIdentification for VariablesGroup package
    status: completed
  - id: variablesgroup-get
    content: Implement get.go with Get, resolveGroup, resolveGroupID, describeGroup helpers
    status: completed
  - id: variablesgroup-search
    content: Implement search.go via ServiceHubSearchQueryController.searchVariables
    status: completed
  - id: variablesgroup-apply
    content: Implement apply.go via VariablesGroupCommandController.Apply with protojson unmarshal
    status: completed
  - id: variablesgroup-delete
    content: Implement delete.go via VariablesGroupCommandController.Delete with dual-path ID resolution
    status: completed
  - id: variablesgroup-upsert-entry
    content: Implement upsert_entry.go via VariablesGroupCommandController.UpsertEntry with entry JSON→proto conversion
    status: completed
  - id: variablesgroup-delete-entry
    content: Implement delete_entry.go via VariablesGroupCommandController.DeleteEntry
    status: completed
  - id: variablesgroup-get-value
    content: Implement get_value.go via VariablesGroupQueryController.GetValue with StringValue unwrapping
    status: completed
  - id: variablesgroup-transform
    content: Implement transform.go via VariablesGroupQueryController.Transform
    status: completed
  - id: variablesgroup-register
    content: Create register.go wiring all 8 VariablesGroup tools
    status: completed
  - id: secretsgroup-tools
    content: Create tools.go with 8 input structs, 8 Tool/Handler pairs for SecretsGroup (mirror variablesgroup)
    status: completed
  - id: secretsgroup-impl
    content: Implement all 8 SecretsGroup operation files (get, search, apply, delete, upsert_entry, delete_entry, get_value, transform)
    status: completed
  - id: secretsgroup-register
    content: Create register.go wiring all 8 SecretsGroup tools
    status: completed
  - id: server-wiring
    content: Add servicehubvariablesgroup and servicehubsecretsgroup imports and Register calls to server.go
    status: completed
  - id: build-verify
    content: Run go build, go vet, go test to verify clean compilation
    status: completed
isProject: false
---

# Tier 3: VariablesGroup + SecretsGroup MCP Tools (16 tools)

## Scope

8 tools per entity, mirrored across VariablesGroup and SecretsGroup:

- **Group-level CRUD**: `get`, `apply`, `delete`
- **Entry-level mutations**: `upsert_entry`, `delete_entry`
- **Value resolution**: `get_value`
- **Search**: entry-level full-text search across all groups in an org
- **Transform**: batch-resolve `$variables-group/` or `$secrets-group/` references

## Tool Catalogue

### VariablesGroup (8 tools)


| Tool Name                | RPC                      | Controller                        | Notes                                |
| ------------------------ | ------------------------ | --------------------------------- | ------------------------------------ |
| `search_variables`       | `searchVariables`        | `ServiceHubSearchQueryController` | Entry-level search across all groups |
| `get_variables_group`    | `get` / `getByOrgBySlug` | `VariablesGroupQueryController`   | By ID or org+slug                    |
| `apply_variables_group`  | `apply`                  | `VariablesGroupCommandController` | Idempotent create-or-update          |
| `delete_variables_group` | `delete`                 | `VariablesGroupCommandController` | Remove group                         |
| `upsert_variable`        | `upsertEntry`            | `VariablesGroupCommandController` | Add/update single entry              |
| `delete_variable`        | `deleteEntry`            | `VariablesGroupCommandController` | Remove single entry                  |
| `get_variable_value`     | `getValue`               | `VariablesGroupQueryController`   | Resolve by org+group_name+entry_name |
| `transform_variables`    | `transform`              | `VariablesGroupQueryController`   | Batch-resolve references             |


### SecretsGroup (8 tools)


| Tool Name              | RPC                      | Controller                        | Notes                                                       |
| ---------------------- | ------------------------ | --------------------------------- | ----------------------------------------------------------- |
| `search_secrets`       | `searchSecrets`          | `ServiceHubSearchQueryController` | Entry-level search across all groups                        |
| `get_secrets_group`    | `get` / `getByOrgBySlug` | `SecretsGroupQueryController`     | By ID or org+slug                                           |
| `apply_secrets_group`  | `apply`                  | `SecretsGroupCommandController`   | Idempotent create-or-update                                 |
| `delete_secrets_group` | `delete`                 | `SecretsGroupCommandController`   | Remove group                                                |
| `upsert_secret`        | `upsertEntry`            | `SecretsGroupCommandController`   | Add/update single entry                                     |
| `delete_secret`        | `deleteEntry`            | `SecretsGroupCommandController`   | Remove single entry                                         |
| `get_secret_value`     | `getValue`               | `SecretsGroupQueryController`     | Resolve by org+group_name+entry_name; **plaintext warning** |
| `transform_secrets`    | `transform`              | `SecretsGroupQueryController`     | Batch-resolve references                                    |


## Design Decisions

### DD-T3-1: Entry-level search via dedicated RPC

Unlike Service (which uses the generic `ApiResourceSearchQueryController.searchByKind`), VariablesGroup and SecretsGroup have dedicated search RPCs on `ServiceHubSearchQueryController` that search individual **entries** across all groups in an org. This returns `VariableEntrySearchRecord` / `SecretEntrySearchRecord` with group context (group_name, group_id) embedded in each result.

- Input: `SearchConfigEntriesRequest` (org, search_text, page_info) -- same message type for both
- Shared import: `search/v1/servicehub` stubs

### DD-T3-2: Dual-path identification for group-level tools

Consistent with the Service tools pattern (`[tools.go](internal/domains/servicehub/service/tools.go)` `validateIdentification`):

- **get/delete/apply**: Accept `id` alone OR `org`+`slug`
- Helpers: `resolveGroup()`, `resolveGroupID()`, `describeGroup()` -- same pattern as `resolveService()`, `resolveServiceID()`, `describeService()` in the service package

### DD-T3-3: Entry mutation tools accept group_id OR org+group_slug

The proto RPCs `upsertEntry` and `deleteEntry` require `group_id`. But for agent ergonomics, these tools will also accept `org`+`group_slug` and internally resolve to a group_id via `getByOrgBySlug` (one extra round-trip). This matches how Service tools handle disconnect/webhook/branches.

### DD-T3-4: Upsert entry input format

`upsert_variable` / `upsert_secret` accept the entry as a **nested JSON object** under the `entry` field:

```json
{
  "group_id": "vg-abc123",
  "entry": {
    "name": "DATABASE_HOST",
    "description": "Primary database hostname",
    "value": "postgres.internal.example.com"
  }
}
```

The handler marshals the entry map to JSON, then unmarshals via `protojson` into the typed proto -- consistent with how `apply_service` handles the full Service resource. This handles both literal values and `value_from` references without special-casing.

### DD-T3-5: StringValue unwrapping for get_value tools

`getValue` returns `google.protobuf.StringValue`. The handler unwraps it and returns the plain string as text content (not JSON-wrapped). If the StringValue is nil/empty, return a clear "no value found" message.

### DD-T3-6: Security warning on get_secret_value

The tool description for `get_secret_value` must include:

> "WARNING: This returns the secret value in PLAINTEXT. Only use when the user explicitly requests to see a secret value. Never log or display secret values unless specifically asked."

### DD-T3-7: No shared abstraction between packages

VariablesGroup and SecretsGroup are separate bounded contexts with their own proto types (`VariablesGroup` vs `SecretsGroup`, `VariablesGroupEntry` vs `SecretsGroupEntry`, etc.). Even though structurally similar, they get independent packages with no shared generics. Reasons:

- Different proto import paths and generated types
- Different ubiquitous language (variables vs secrets)
- Different security semantics (secrets require sensitivity warnings)
- Premature abstraction creates coupling where none is needed

## Package Structure

```
internal/domains/servicehub/
├── variablesgroup/
│   ├── register.go        # Register(srv, serverAddress) — 8 AddTool calls
│   ├── tools.go           # 8 input structs, 8 Tool/Handler pairs, validateIdentification
│   ├── search.go          # search_variables → ServiceHubSearchQueryController.searchVariables
│   ├── get.go             # get_variables_group + resolveGroup, resolveGroupID, describeGroup
│   ├── apply.go           # apply_variables_group → VariablesGroupCommandController.Apply
│   ├── delete.go          # delete_variables_group → VariablesGroupCommandController.Delete
│   ├── upsert_entry.go    # upsert_variable → VariablesGroupCommandController.UpsertEntry
│   ├── delete_entry.go    # delete_variable → VariablesGroupCommandController.DeleteEntry
│   ├── get_value.go       # get_variable_value → VariablesGroupQueryController.GetValue
│   └── transform.go       # transform_variables → VariablesGroupQueryController.Transform
└── secretsgroup/
    ├── register.go        # Register(srv, serverAddress) — 8 AddTool calls
    ├── tools.go           # 8 input structs, 8 Tool/Handler pairs, validateIdentification
    ├── search.go          # search_secrets → ServiceHubSearchQueryController.searchSecrets
    ├── get.go             # get_secrets_group + resolveGroup, resolveGroupID, describeGroup
    ├── apply.go           # apply_secrets_group → SecretsGroupCommandController.Apply
    ├── delete.go          # delete_secrets_group → SecretsGroupCommandController.Delete
    ├── upsert_entry.go    # upsert_secret → SecretsGroupCommandController.UpsertEntry
    ├── delete_entry.go    # delete_secret → SecretsGroupCommandController.DeleteEntry
    ├── get_value.go       # get_secret_value → SecretsGroupQueryController.GetValue
    └── transform.go       # transform_secrets → SecretsGroupQueryController.Transform
```

**Server wiring** (`[internal/server/server.go](internal/server/server.go)`):

- Add imports: `servicehubvariablesgroup` and `servicehubsecretsgroup`
- Add 2 Register calls in `registerTools()`

## Key Proto Import Paths (Go stubs)

- `variablesgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/variablesgroup/v1"`
- `secretsgroupv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/secretsgroup/v1"`
- `servicehubsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/servicehub"`
- `apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"`
- `rpc "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/rpc"`

## Implementation Order

Execute sequentially within each entity, then wire:

1. **VariablesGroup**: `tools.go` (all structs + definitions) → `get.go` (foundation + resolve helpers) → `search.go` → `apply.go` → `delete.go` → `upsert_entry.go` → `delete_entry.go` → `get_value.go` → `transform.go` → `register.go`
2. **SecretsGroup**: Mirror the same sequence
3. **Server wiring**: Add both to `server.go`
4. **Verification**: `go build ./...`, `go vet ./...`

## Excluded RPCs (consistent with T01 plan)

- `create` / `update` — Covered by `apply` (idempotent)
- `find` — Platform-operator-only admin query

## Total Deliverable

- **20 new Go files** (10 per package)
- **1 modified file** (`internal/server/server.go`)
- **16 MCP tools** (8 + 8)

