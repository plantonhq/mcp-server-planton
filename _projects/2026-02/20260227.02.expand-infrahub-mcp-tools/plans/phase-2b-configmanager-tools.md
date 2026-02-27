---
name: Phase 2B ConfigManager Tools
overview: Add 11 MCP tools across three configmanager sub-domains (variable, secret, secretversion) to the Planton MCP server, expanding from 41 to 52 tools. This is the second domain outside infrahub, completing the Phase 2 pair (Graph + ConfigManager).
todos:
  - id: variable-package
    content: "Implement Variable package (5 tools): enum.go, list.go, get.go, apply.go, delete.go, resolve.go, tools.go"
    status: completed
  - id: secret-package
    content: "Implement Secret package (4 tools): enum.go, list.go, get.go, apply.go, delete.go, tools.go"
    status: completed
  - id: secretversion-package
    content: "Implement SecretVersion package (2 tools): create.go, list.go, tools.go"
    status: completed
  - id: server-registration
    content: Wire 11 new tools in server.go, add imports, update count 41->52
    status: completed
  - id: domain-doc
    content: Create configmanager/doc.go with domain-level documentation
    status: completed
  - id: verification
    content: "Full verification: go build, go vet, go test -- all must pass clean"
    status: completed
isProject: false
---

# Phase 2B: ConfigManager / Variables and Secrets

## Domain Analysis (Architect Role)

The ConfigManager domain has three sub-resources with distinct security profiles:

- **Variable** -- configuration values that are readable and writable (plaintext)
- **Secret** -- metadata containers for encrypted values (metadata is safe to expose)
- **SecretVersion** -- the actual encrypted payloads (write-only from agent perspective; reading decrypted data crosses a security boundary per AD-01 logic)

Variables and secrets share a **scope** dimension (`organization` or `environment`) that determines their uniqueness key: `(org, scope, slug)`. This is different from InfraProject's simpler `(org, slug)` pattern and requires scope-aware enum resolution.

## Design Decisions (Approved)

- **DD-1: Expanded tool surface (11 tools, up from planned 5)** -- proto analysis reveals 6 gRPC services with ~20 RPCs; 11 tools cover the complete lifecycle (list/get/apply/delete for variables and secrets, create/list for secret versions)
- **DD-2: Write-only secret values** -- `create_secret_version` included (agents write values), `get_secret_version`/`get_latest_secret_version` excluded (reading decrypted data is a security boundary). Same principle as AD-01 for credentials
- **DD-3: Explicit parameters for apply tools** -- Variable and Secret have simple, stable schemas; explicit params (name, org, scope, value) provide better agent UX and validation vs JSON passthrough. `VariableSpec.source` (ValueFromRef) excluded initially; can be added later
- **DD-4: Exclude refresh_variable** -- `VariableCommandController.Refresh` re-reads from source reference; specialized operation, low initial demand
- **DD-5: delete_secret includes WARNING** -- `SecretCommandController.Delete` destroys the secret AND all its versions permanently; tool description must carry a strong warning (same pattern as `destroy_cloud_resource`)

## Proto References (Verified)

- `ApiResourceKind_variable` (130), `ApiResourceKind_secret` (38) -- confirmed in apiresourcekind enum
- `VariableSpec_Scope` / `SecretSpec_Scope` -- `organization` (1), `environment` (2)
- Import paths:
  - `github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/variable/v1`
  - `github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secret/v1`
  - `github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secretversion/v1`

## Tool Catalog (11 tools)

### Variable (5 tools)


| Tool               | RPC                               | Parameters                                                         | Notes                                 |
| ------------------ | --------------------------------- | ------------------------------------------------------------------ | ------------------------------------- |
| `list_variables`   | `VariableQueryController.Find`    | org, env (optional), page_num, page_size                           | Kind=`ApiResourceKind_variable` (130) |
| `get_variable`     | `.Get` / `.GetByOrgByScopeBySlug` | id OR (org, scope, slug)                                           | Dual identification like infraproject |
| `apply_variable`   | `VariableCommandController.Apply` | name, org, scope, env (when scope=environment), description, value | Construct proto internally            |
| `delete_variable`  | `.Delete`                         | id OR (org, scope, slug)                                           | Resolves to ID then calls Delete      |
| `resolve_variable` | `VariableQueryController.Resolve` | org, scope, slug                                                   | Returns plain string value only       |


### Secret (4 tools)


| Tool            | RPC                               | Parameters                                                                      | Notes                              |
| --------------- | --------------------------------- | ------------------------------------------------------------------------------- | ---------------------------------- |
| `list_secrets`  | `SecretQueryController.Find`      | org, env (optional), page_num, page_size                                        | Kind=`ApiResourceKind_secret` (38) |
| `get_secret`    | `.Get` / `.GetByOrgByScopeBySlug` | id OR (org, scope, slug)                                                        | Metadata only, no secret data      |
| `apply_secret`  | `SecretCommandController.Apply`   | name, org, scope, env (when scope=environment), description, backend (optional) | Construct proto internally         |
| `delete_secret` | `.Delete`                         | id OR (org, scope, slug)                                                        | WARNING: deletes all versions too  |


### SecretVersion (2 tools)


| Tool                    | RPC                                         | Parameters                          | Notes                                  |
| ----------------------- | ------------------------------------------- | ----------------------------------- | -------------------------------------- |
| `create_secret_version` | `SecretVersionCommandController.Create`     | secret_id, data (map[string]string) | Agent writes encrypted key-value pairs |
| `list_secret_versions`  | `SecretVersionQueryController.ListBySecret` | secret_id                           | Metadata only (no `data` field)        |


## Package Structure

```
internal/domains/configmanager/
  doc.go                        -- domain-level docs listing all 3 sub-domains
  variable/
    tools.go                    -- 5 input structs + 5 tool defs + 5 handlers
    list.go                     -- Find RPC (pattern: infrachart/list.go)
    get.go                      -- Get + GetByOrgByScopeBySlug (pattern: infraproject/get.go)
    apply.go                    -- Apply with explicit param -> proto construction
    delete.go                   -- Resolve ID then Delete
    resolve.go                  -- Resolve RPC returning plain string
    enum.go                     -- resolveScope(string) -> VariableSpec_Scope
  secret/
    tools.go                    -- 4 input structs + 4 tool defs + 4 handlers
    list.go                     -- Find RPC
    get.go                      -- Get + GetByOrgByScopeBySlug
    apply.go                    -- Apply with explicit param -> proto construction
    delete.go                   -- Resolve ID then Delete
    enum.go                     -- resolveScope(string) -> SecretSpec_Scope
  secretversion/
    tools.go                    -- 2 input structs + 2 tool defs + 2 handlers
    create.go                   -- Create RPC
    list.go                     -- ListBySecret RPC
```

Modified: [internal/server/server.go](internal/server/server.go) -- add imports + 11 `mcp.AddTool` calls, update count 41 -> 52

## Key Patterns to Follow

- **Dual identification** (get/delete): follows [infraproject/get.go](internal/domains/infrahub/infraproject/get.go) with `resolveVariable`/`resolveSecret` + `resolveVariableID`/`resolveSecretID` helpers; scope adds a third dimension to the slug path
- **Paginated listing**: follows [infrachart/list.go](internal/domains/infrahub/infrachart/list.go) with `FindApiResourcesRequest`, 1-based page API -> 0-based proto
- **Enum resolution**: follows [graph/enum.go](internal/domains/graph/enum.go) with `resolveScope` + `joinEnumValues` for error messages
- **Connection handling**: all RPC calls go through `domains.WithConnection(ctx, serverAddress, fn)` from [internal/domains/conn.go](internal/domains/conn.go)
- **Error handling**: all gRPC errors go through `domains.RPCError(err, resourceDesc)` from [internal/domains/rpcerr.go](internal/domains/rpcerr.go)
- **Serialization**: all responses go through `domains.MarshalJSON(resp)` from [internal/domains/marshal.go](internal/domains/marshal.go)

## Implementation Sequence

Each step is independently verifiable (`go build ./... && go vet ./...`):

1. **Variable package** -- establishes the configmanager pattern; 5 tools, 7 files
2. **Secret package** -- follows variable pattern closely; 4 tools, 6 files
3. **SecretVersion package** -- simplest sub-domain; 2 tools, 3 files
4. **Server registration** -- wire all 11 tools, update count
5. **Domain doc.go** -- top-level configmanager package docs
6. **Full verification** -- `go build ./...`, `go vet ./...`, `go test ./...`

## Pause Points (Will Collaborate)

- If the `Find` RPC for variables/secrets behaves differently from infrachart's `Find` (e.g. doesn't need `Kind` field, or requires scope as a filter)
- If the `Delete` RPCs require additional input beyond `ApiResourceDeleteInput`
- If the scope enum interacts with the `env` metadata field in unexpected ways (e.g. org-scoped variables ignoring env vs requiring it empty)
- If `resolve_variable` returns errors for variables that have no value set yet

