# Next Task: 20260307.01.proto-contract-sync

## ⚠️ RULES OF ENGAGEMENT - READ FIRST ⚠️

**When this file is loaded in a new conversation, the AI MUST:**

1. **DO NOT AUTO-EXECUTE** - Never start implementing without explicit user approval
2. **GATHER CONTEXT SILENTLY** - Read project files without outputting
3. **PRESENT STATUS SUMMARY** - Show what's done, what's pending, agreed next steps
4. **SHOW OPTIONS** - List recommended and alternative actions
5. **WAIT FOR DIRECTION** - Do NOT proceed until user explicitly says "go" or chooses an option

---

## Project: 20260307.01.proto-contract-sync

**Description**: Migrate MCP server tools to match restructured protobuf contracts. The connect domain renamed all *Credential types to *ProviderConnection/*Connection, new resources were added across configmanager, IAM, and infrahub, and 7 entirely new domains appeared in the proto definitions.

**Goal**: Get the MCP server build green by migrating broken connect/credential and connect/github imports, then add tool implementations for new resources (secretbackend, variablegroup, serviceaccount, iacprovisionermapping) and evaluate new domains (agentfleet, billing, copilot, search, reporting, integration, runner) for tool coverage.

**Tech Stack**: Go/gRPC/MCP

**Components**: internal/domains/connect/connection, internal/domains/connect/github, internal/domains/configmanager, internal/domains/iam, internal/domains/infrahub, schemas/connections, gen/go

---

### ✅ COMPLETED: Phase 1 — Fix the Build (2026-03-07)

**Migrated connect/credential to connect/connection and updated all import paths to use gen/go proto stubs.**

**What was delivered:**

1. **connect/credential → connect/connection package migration** — Renamed package, rewrote all 11 Go files with new proto types, import paths, gRPC controller names, and tool/resource definitions
   - `registry.go` — 20 connection dispatchers with updated proto imports and gRPC client constructors
   - `tools.go` — MCP tools renamed: `apply_credential` → `apply_connection`, etc.
   - `resources.go` — MCP resources renamed: `credential-types://catalog` → `connection-types://catalog`, `credential-schema://{kind}` → `connection-schema://{kind}`
   - `search.go` — 20-entry `connectionKindToAPIResourceKind` map with new enum values
   - `get.go` — Removed all redaction logic (design decision: secret slugs are not sensitive)
   - `redact.go` — Deleted entirely

2. **schemas/credentials → schemas/connections migration** — Renamed directory, renamed 19 JSON schema files, created new `cloudflareworkerscriptsr2connection.json`, rewrote `registry.json`

3. **schemas/embed.go** — Updated embed path to `connections`, renamed `CredentialFS` → `ConnectionFS`

4. **connect/github package** — Updated import paths and type references (`GithubCredential*` → `GithubConnection*`)

5. **server.go wiring** — Updated import alias and function calls to use `connectconnection` package

6. **All domain import paths** — 150+ files across audit, configmanager, graph, iam, infrahub, resourcemanager, servicehub updated from `planton/apis/stubs` to `gen/go` imports

7. **Build verification** — `go build ./...` and `go vet ./...` pass cleanly

**Key Decisions Made:**
- Rename tool names from "credential" to "connection" (Option A — match proto definitions)
- Remove redaction entirely — secret slugs in `ConnectionFieldSecretRef` are not sensitive data, no redaction needed
- Add `CloudflareWorkerScriptsR2Connection` as the 20th connection type (new type with no old equivalent)

**Files Changed/Created:**
- `internal/domains/connect/connection/` — 11 new Go files (apply, delete, doc, get, register, registry, resources, schema, search, slug, tools)
- `internal/domains/connect/credential/` — Deleted entirely (12 files including redact.go)
- `schemas/connections/` — 21 JSON files (20 schemas + registry.json)
- `schemas/credentials/` — Deleted entirely
- `schemas/embed.go` — Updated embed paths and variable names
- `internal/domains/connect/github/tools.go`, `doc.go` — Updated imports/types
- `internal/domains/connect/doc.go` — Updated package documentation
- `internal/server/server.go` — Updated wiring
- `go.mod`, `go.sum` — Updated dependencies
- `buf.gen.go.yaml` — Added buf code generation config
- `tools/buf-generate-go.sh` — Added buf generate script
- 150+ domain files — Import path updates from stubs to gen/go

---

## Essential Files to Review

### 1. Latest Checkpoint (if exists)
Check for the most recent checkpoint file:
```
_projects/2026-03/20260307.01.proto-contract-sync/checkpoints/
```

### 2. Current Task
Review the current task status and plan:
```
_projects/2026-03/20260307.01.proto-contract-sync/tasks/
```

### 3. Plans
Review implementation plans and their status:
```
_projects/2026-03/20260307.01.proto-contract-sync/plans/
```

### 4. Design Decisions
```
_projects/2026-03/20260307.01.proto-contract-sync/design-decisions/
```

---

## Knowledge Folders

- **Coding Guidelines**: `_projects/2026-03/20260307.01.proto-contract-sync/coding-guidelines/`
- **Wrong Assumptions**: `_projects/2026-03/20260307.01.proto-contract-sync/wrong-assumptions/`
- **Don't Dos**: `_projects/2026-03/20260307.01.proto-contract-sync/dont-dos/`

---

## Current Status

**Created**: 2026-03-07 19:45
**Last Updated**: 2026-03-08

**Current step:**
- ✅ **Phase 1: Fix the Build** — T01 Credential-to-Connection Migration (2026-03-07)
  - Build is green, all 20 connection types callable, schemas migrated, redaction removed
- ✅ **Phase 2: Enrich Existing Connect Tools** — T02.1–T02.4 (2026-03-08)
  - 9 new tools + 1 enhanced tool + 1 bug fix across 4 connect sub-packages
  - T02.5 (provider-specific controllers) deferred — needs design decision on OAuth scope
- ✅ **Phase 3: New Resources in Existing Domains** — T03.1–T03.4 (2026-03-08)
  - 23 new tools across 4 new packages, 1 security boundary (SecretBackend redaction)
- 🔵 Next: **T02.5** or choose from Phase 4/5

---

### ✅ COMPLETED: Phase 2 — Enrich Existing Connect Tools (2026-03-08)

**What was delivered:**

1. **Bug fix: defaultprovider ResolveHandler** — Provider field was collected from user but never passed to gRPC call (sent UNSPECIFIED). Now correctly resolves and passes the provider enum.

2. **T02.1: defaultprovider — 4 new tools** (4→8 tools)
   - `get_org_default_provider_connection` — explicit org-level lookup (no fallback)
   - `get_env_default_provider_connection` — explicit env-level lookup (no fallback)
   - `delete_org_default_provider_connection` — delete org-level default by org+provider
   - `delete_env_default_provider_connection` — delete env-level default by org+provider+env

3. **T02.2: runner — 2 new tools** (4→6 tools)
   - `generate_runner_credentials` — generate initial auth credentials (sensitive output warning)
   - `regenerate_runner_credentials` — rotate/regenerate auth credentials (sensitive output warning)

4. **T02.4: providerauth — 1 new tool + 1 enhanced tool** (3→4 tools)
   - `sync_provider_connection_authorization` — reconcile authorization state by semantic key
   - Enhanced `delete_provider_connection_authorization` — now accepts ID or semantic key (org+provider+connection), mirroring how Get already works

5. **T02.3: defaultrunner — 2 new tools + ApiResourceKind resolver** (4→6 tools)
   - `get_default_runner_binding_by_selector` — lookup by resource selector (kind+ID)
   - `delete_default_runner_binding_by_selector` — delete by resource selector (kind+ID)
   - Added `ResolveApiResourceKind` to shared domains package

6. **Build verification** — `go build ./...` and `go vet ./...` pass cleanly

**Key Design Decisions:**
- Separate tools for org-level vs env-level operations (not overloaded into existing Get/Delete)
- Enhanced providerauth Delete to support semantic key (instead of a new 50-char tool name)
- `Find` methods explicitly skipped — proto docs say "restricted to platform operators only"
- `Create`/`Update` methods skipped — already covered by `Apply`
- T02.5 deferred — OAuth callback handlers are browser redirect endpoints, not agent-callable

**Files Changed:**
- `internal/domains/connect/defaultprovider/` — tools.go (bug fix + 4 new tools), register.go, doc.go
- `internal/domains/connect/runner/` — tools.go (2 new tools), register.go, doc.go
- `internal/domains/connect/providerauth/` — tools.go (1 new tool + enhanced delete), register.go, doc.go
- `internal/domains/connect/defaultrunner/` — tools.go (2 new tools), register.go, doc.go
- `internal/domains/kind.go` — Added `ResolveApiResourceKind` and `apiResourceKindResolver`

---

### ✅ COMPLETED: Phase 3 — New Resources in Existing Domains (2026-03-08)

**Implemented 23 MCP tools across 4 new packages for new proto resources, with dedicated security boundary for sensitive credential fields.**

**What was delivered:**

1. **T03.4: `infrahub/iacprovisionermapping`** — 3 tools
   - `apply_iac_provisioner_mapping` — Envelope Apply (idempotent)
   - `get_iac_provisioner_mapping` — Get by ID
   - `delete_iac_provisioner_mapping` — Delete by ID

2. **T03.1: `configmanager/secretbackend`** — 4 tools + redaction
   - `apply_secret_backend` — Envelope Apply with sensitive field redaction
   - `get_secret_backend` — Dual-path (ID or org+slug) with redaction
   - `list_secret_backends` — ListByOrg with per-entry redaction
   - `delete_secret_backend` — Dual-path resolve-then-delete with redaction
   - `redact.go` — Defense-in-depth redaction of 10 sensitive fields across 8 config blocks

3. **T03.3: `iam/serviceaccount`** — 8 tools
   - `create_service_account` — Explicit params (org, display_name, description)
   - `get_service_account` — Get by ID
   - `update_service_account` — Read-modify-write pattern
   - `delete_service_account` — Delete with cascade warning
   - `list_service_accounts` — FindByOrg
   - `create_service_account_key` — Sensitive output warning (key shown once)
   - `revoke_service_account_key` — By service_account_id + api_key_id
   - `list_service_account_keys` — By service account ID

4. **T03.2: `configmanager/variablegroup`** — 8 tools
   - `apply_variable_group` — Envelope Apply
   - `get_variable_group` — Dual-path (ID or org+scope+slug)
   - `delete_variable_group` — Dual-path resolve-then-delete
   - `upsert_variable_group_entry` — Explicit params
   - `delete_variable_group_entry` — Explicit params
   - `refresh_variable_group_entry` — Explicit params
   - `refresh_all_variable_group_entries` — By group ID
   - `resolve_variable_group_entry` — Quick value lookup (returns plain string)

**Key Architectural Decisions:**
- AD-03: SecretBackend Apply uses Envelope style (6 mutually exclusive config blocks)
- AD-04: SecretBackend responses MUST redact sensitive credential fields (tokens, keys, secrets)
- AD-05: ServiceAccount — skip GetByIdentityAccountId (internal plumbing, not agent-facing)
- AD-06: VariableGroup — Envelope for Apply, Explicit for entry operations (discoverability)
- AD-07: Skip Find RPCs across all 4 resources (operator-only, better alternatives exist)
- AD-08: Skip Create/Update where Apply exists (ServiceAccount is exception — no Apply in proto)

**Files Created:**
- `internal/domains/infrahub/iacprovisionermapping/` — doc.go, register.go, tools.go
- `internal/domains/configmanager/secretbackend/` — doc.go, register.go, tools.go, get.go, apply.go, delete.go, list.go, redact.go
- `internal/domains/iam/serviceaccount/` — doc.go, register.go, tools.go, create.go, get.go, update.go, delete.go, list.go, keys.go
- `internal/domains/configmanager/variablegroup/` — doc.go, register.go, tools.go, enum.go, get.go, apply.go, delete.go, entry.go, resolve.go

**Files Modified:**
- `internal/server/server.go` — 4 new imports + Register calls
- `internal/domains/configmanager/doc.go` — Updated to 6 sub-packages / 23 tools
- `internal/domains/iam/doc.go` — Updated to 6 sub-packages / added serviceaccount
- `internal/domains/infrahub/doc.go` — Added iacprovisionermapping to sub-package list

---

## Objectives for Next Conversations

### Option A: T02.5 — Provider-Specific Controllers (Pending Decision)
Wire CloudFormation setup + OAuth initiation tools:
- `initiate_aws_cloudformation_setup` + `get_aws_cloudformation_setup_status` (2 tools, new package)
- `initiate_gcp_oauth` (1 tool, new package)
- `initiate_azure_oauth` (1 tool, new package)
- Skip: OAuth callback handlers (browser redirect endpoints)

### Option B: Phase 4 — Evaluate New Domains
Survey 7 new domains (agentfleet, search, integration, runner, billing, copilot, reporting) and decide which need MCP tool coverage.

### Option C: Phase 5 — Enrich Existing Stable Domains
Low priority enrichments to configmanager, iam, etc.

---

## Quick Commands

After loading context:
- "Continue with T02.5" - Wire provider-specific controllers
- "Start Phase 3" - Implement new resources
- "Show project status" - Get overview of progress
- "Create checkpoint" - Save current progress

---

*This file provides direct paths to all project resources for quick context loading.*
