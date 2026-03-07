---
name: Phase 3 New Resources
overview: Implement MCP tool packages for 4 new proto resources (secretbackend, variablegroup, serviceaccount, iacprovisionermapping) totaling 23 tools across 4 packages, with clear architectural decisions on apply style, security boundaries, and API surface.
todos:
  - id: t03-4
    content: "T03.4: Implement infrahub/iacprovisionermapping (3 tools: apply, get, delete)"
    status: completed
  - id: t03-1
    content: "T03.1: Implement configmanager/secretbackend (4 tools + redaction for sensitive backend config fields)"
    status: completed
  - id: t03-3
    content: "T03.3: Implement iam/serviceaccount (8 tools: CRUD + key management, read-modify-write update)"
    status: completed
  - id: t03-2
    content: "T03.2: Implement configmanager/variablegroup (8 tools: envelope apply + explicit entry operations + resolve)"
    status: completed
isProject: false
---

# Phase 3: New Resources in Existing Domains

## Architectural Decisions (Final)

### AD-03: SecretBackend Apply uses Envelope style

The spec has 6 mutually exclusive backend config blocks + encryption config (~30 fields total). Explicit params would produce a bloated tool signature where 80% of fields are irrelevant for any given backend type. This is a configuration object, not a simple entity — agents will work from a spec template.

Envelope pattern (same as `defaultrunner`): accept `backend_object` as `map[string]any`, unmarshal via protojson.

### AD-04: SecretBackend responses MUST redact sensitive fields

This is not the same situation as Phase 1. Connection slugs are references — they point to secrets but aren't secrets themselves. SecretBackend configs contain **actual credentials**: OpenBAO tokens, AWS secret access keys, Azure client secrets, GCP service account key JSON. The proto marks these with the `is_sensitive` field option.

Implementation: After Get/Apply returns the proto, walk the spec config blocks and replace sensitive field values with `"[REDACTED]"`. This is a targeted redaction — only the fields the proto authors marked as sensitive. A `redact.go` file with a single `RedactSecretBackend(*SecretBackend)` function.

This is a security boundary the platform must enforce. An agent should never have raw cloud credentials in its context window.

### AD-05: ServiceAccount — skip GetByIdentityAccountId

This is internal plumbing. Agents identify service accounts by name or ID, never by the backing machine identity (which is auto-created and never user-supplied per the proto docs). If someone needs reverse lookup later, adding one tool is trivial.

### AD-06: VariableGroup — Envelope for Apply, Explicit for entry operations

The main Apply takes the full VariableGroup object (entries is a repeated nested message with optional ValueFromRef sources). Envelope is the right fit.

But the entry-level operations (UpsertEntry, DeleteEntry, RefreshEntry, RefreshAll) have flat, simple signatures (group_id + entry_name or group_id + entry). These **must** use explicit params — they're the high-frequency agent operations and discoverability matters. An agent managing config entries one at a time should see exactly what fields are needed.

### AD-07: Skip Find RPCs across all 4 resources

Find RPCs are paginated search endpoints that accept `FindApiResourcesRequest`. In Phase 2, we established that Find is "restricted to platform operators." More importantly, every resource here has a better alternative: `ListByOrg` (secretbackend), semantic key lookup (variablegroup via org+scope+slug), or `FindByOrg` (serviceaccount). IacProvisionerMapping has no Find at all.

### AD-08: Skip Create/Update where Apply exists

SecretBackend, VariableGroup, and IacProvisionerMapping all have Apply. We don't expose Create and Update separately — Apply is idempotent and covers both. ServiceAccount is the exception (no Apply in proto), so it gets Create + Update.

---

## Execution Order

T03.4 (3 tools) then T03.1 (4 tools) then T03.3 (8 tools) then T03.2 (8 tools). Simplest first, most complex last.

---

## T03.4: infrahub/iacprovisionermapping (3 tools)

**New package:** `internal/domains/infrahub/iacprovisionermapping/`

**Proto stubs:** `gen/go/ai/planton/infrahub/iacprovisionermapping/v1/`

- Command: Apply, Delete (by `ApiResourceId`)
- Query: Get (by `ApiResourceId`)
- Spec: `selector` (ApiResourceSelector) + `provisioner` (IacProvisioner enum)

**Tools:**

- `apply_iac_provisioner_mapping` — Envelope Apply
- `get_iac_provisioner_mapping` — Get by ID
- `delete_iac_provisioner_mapping` — Delete by ID

**Files to create:**

- `doc.go` — package doc listing 3 tools
- `register.go` — Register function with 3 AddTool calls
- `tools.go` — Input structs, Tool/Handler functions

**Files to modify:**

- [internal/server/server.go](internal/server/server.go) — add import + Register call
- [internal/domains/infrahub/doc.go](internal/domains/infrahub/doc.go) — update sub-package list

---

## T03.1: configmanager/secretbackend (4 tools + redaction)

**New package:** `internal/domains/configmanager/secretbackend/`

**Proto stubs:** `gen/go/ai/planton/configmanager/secretbackend/v1/`

- Command: Apply, Delete (by `ApiResourceDeleteInput`)
- Query: Get (by `SecretBackendId`), GetByOrgBySlug (by `ApiResourceByOrgBySlugRequest`), ListByOrg (by `OrganizationId`)
- ID: `SecretBackendId{Value}`

**Tools:**

- `apply_secret_backend` — Envelope Apply
- `get_secret_backend` — Get by ID or by org+slug (same dual-path pattern as `get_secret`)
- `list_secret_backends` — List by org
- `delete_secret_backend` — Delete by ID or by org+slug (resolve first, then delete)

**Files to create:**

- `doc.go`, `register.go`, `tools.go`
- `get.go` — dual-path resolution (ID vs org+slug)
- `apply.go` — envelope unmarshal + redact response
- `delete.go` — resolve + delete via ApiResourceDeleteInput
- `list.go` — ListByOrg call
- `redact.go` — `RedactSecretBackend()` that replaces sensitive fields with `"[REDACTED]"`

**Files to modify:**

- [internal/server/server.go](internal/server/server.go) — add import + Register call
- [internal/domains/configmanager/doc.go](internal/domains/configmanager/doc.go) — update sub-package list and tool count

**Redaction targets** (fields marked `is_sensitive` in proto):

- `OpenBaoConfig.Token`
- `AwsSecretsManagerConfig.AccessKeyId`, `.SecretAccessKey`
- `HashicorpVaultConfig.Token`
- `GcpSecretManagerConfig.ServiceAccountKeyJson`
- `AzureKeyVaultConfig.ClientSecret`
- `AwsKmsKeyConfig.AccessKeyId`, `.SecretAccessKey`
- `GcpKmsKeyConfig.ServiceAccountKeyJson`
- `AzureKeyVaultKeyConfig.ClientSecret`

---

## T03.3: iam/serviceaccount (8 tools)

**New package:** `internal/domains/iam/serviceaccount/`

**Proto stubs:** `gen/go/ai/planton/iam/serviceaccount/v1/`

- Command: Create, Update, Delete (by `ServiceAccountId`), CreateKey, RevokeKey
- Query: Get (by `ServiceAccountId`), FindByOrg (by `ApiResourceId`), ListKeys (by `ServiceAccountId`)
- Spec: `display_name`, `description`, `identity_account_id` (auto, never user-supplied)

**Tools:**

- `create_service_account` — Explicit params (org, display_name, description)
- `get_service_account` — Get by ID
- `update_service_account` — Read-modify-write (same pattern as `update_team`)
- `delete_service_account` — Delete by ID
- `list_service_accounts` — FindByOrg (by org ID)
- `create_service_account_key` — By service account ID. **Sensitive output warning** in description (key value shown once, same pattern as runner credential tools)
- `revoke_service_account_key` — By service account ID + API key ID
- `list_service_account_keys` — By service account ID

**Files to create:**

- `doc.go`, `register.go`, `tools.go`
- `create.go`, `get.go`, `update.go`, `delete.go`, `list.go`
- `keys.go` — CreateKey, RevokeKey, ListKeys implementations

**Files to modify:**

- [internal/server/server.go](internal/server/server.go)
- [internal/domains/iam/doc.go](internal/domains/iam/doc.go)

---

## T03.2: configmanager/variablegroup (8 tools)

**New package:** `internal/domains/configmanager/variablegroup/`

**Proto stubs:** `gen/go/ai/planton/configmanager/variablegroup/v1/`

- Command: Apply, Delete (by `ApiResourceDeleteInput`), UpsertEntry, DeleteEntry, RefreshEntry, RefreshAll
- Query: Get (by `VariableGroupId`), GetByOrgByScopeBySlug, Resolve
- Spec: `scope` (org/env enum), `description`, `entries[]` (name, description, value, source)

**Tools:**

- `apply_variable_group` — Envelope Apply (full group with entries)
- `get_variable_group` — Get by ID or by org+scope+slug (dual-path, same pattern as `get_secret`)
- `delete_variable_group` — Delete by ID or org+scope+slug
- `upsert_variable_group_entry` — Explicit: group_id, entry_name, value, description, source (optional)
- `delete_variable_group_entry` — Explicit: group_id, entry_name
- `refresh_variable_group_entry` — Explicit: group_id, entry_name
- `refresh_all_variable_group_entries` — By group ID
- `resolve_variable_group_entry` — Explicit: org, scope, slug, entry_name (returns plain string value)

**Files to create:**

- `doc.go`, `register.go`, `tools.go`
- `enum.go` — `scopeResolver` for `VariableGroupSpec_Scope` (same pattern as `variable/enum.go`)
- `get.go` — dual-path resolution
- `apply.go` — envelope unmarshal
- `delete.go` — resolve + delete
- `entry.go` — UpsertEntry, DeleteEntry, RefreshEntry, RefreshAll handlers
- `resolve.go` — Resolve handler (returns `StringValue`)

**Files to modify:**

- [internal/server/server.go](internal/server/server.go)
- [internal/domains/configmanager/doc.go](internal/domains/configmanager/doc.go)

---

## Cross-Cutting Concerns

- **Build verification** after each sub-task: `go build ./...` and `go vet ./...`
- **Parent doc.go updates**: configmanager and iam parent packages get updated sub-package lists and tool counts
- **server.go wiring**: Each new package gets an import alias and Register call following established naming conventions (e.g., `configmanagersecretbackend`, `iamserviceaccount`, etc.)
- **No new shared utilities needed**: All required helpers already exist in `internal/domains/` (WithConnection, MarshalJSON, RPCError, TextResult, NewEnumResolver)

## Totals

- 4 new packages
- 23 new tools
- ~30 new files + ~4 modified files
- 1 new security boundary (SecretBackend redaction)

