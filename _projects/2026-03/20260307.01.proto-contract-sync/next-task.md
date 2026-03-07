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
**Last Updated**: 2026-03-08 00:17

**Current step:**
- ✅ **Phase 1: Fix the Build** — T01 Credential-to-Connection Migration (2026-03-07)
  - Build is green, all 20 connection types callable, schemas migrated, redaction removed
- 🔵 Next: **Phase 2: Enrich Existing Connect Tools** (T02.1–T02.5) or choose from other objectives

---

## Objectives for Next Conversations

### Option A (Recommended): Phase 2 — Enrich Existing Connect Tools
Quick wins — expose new gRPC methods that already exist but aren't wired as MCP tools:
- T02.1: `connect/defaultprovider` — GetOrgDefault, GetEnvDefault, DeleteOrgDefault, DeleteEnvDefault
- T02.2: `connect/runner` — GenerateCredentials, RegenerateCredentials
- T02.3: `connect/defaultrunner` — GetBySelector, DeleteBySelector
- T02.4: `connect/providerauth` — Sync, DeleteBySemanticKey, Find
- T02.5: Provider-specific controllers (AWS CloudFormation, GCP OAuth, Azure OAuth)

### Option B: Phase 3 — New Resources in Existing Domains
Implement MCP tools for entirely new resources:
- T03.1: `configmanager/secretbackend`
- T03.2: `configmanager/variablegroup`
- T03.3: `iam/serviceaccount`
- T03.4: `infrahub/iacprovisionermapping`

### Option C: Phase 4 — Evaluate New Domains
Survey 7 new domains (agentfleet, search, integration, runner, billing, copilot, reporting) and decide which need MCP tool coverage.

### Option D: Phase 5 — Enrich Existing Stable Domains
Low priority enrichments to configmanager, iam, etc.

---

## Quick Commands

After loading context:
- "Continue with Phase 2" - Start enriching connect tools
- "Show project status" - Get overview of progress
- "Create checkpoint" - Save current progress
- "Review plan" - See full T01 plan details

---

*This file provides direct paths to all project resources for quick context loading.*
