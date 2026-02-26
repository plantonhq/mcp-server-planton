# Next Task: Refactor mcp-server-planton (Stigmer Patterns)

## ⚠️ RULES OF ENGAGEMENT - READ FIRST ⚠️

**When this file is loaded in a new conversation, the AI MUST:**

1. **DO NOT AUTO-EXECUTE** - Never start implementing without explicit user approval
2. **GATHER CONTEXT SILENTLY** - Read project files without outputting
3. **PRESENT STATUS SUMMARY** - Show what's done, what's pending, agreed next steps
4. **SHOW OPTIONS** - List recommended and alternative actions
5. **WAIT FOR DIRECTION** - Do NOT proceed until user explicitly says "go" or chooses an option

---

**Project**: `_projects/2026-02/20260226.01.refactor-mcp-server-stigmer-patterns/`
**Current Status**: Phase 3 COMPLETED → Ready for Phase 4 (delete_cloud_resource + get_cloud_resource)

## Quick Context

Refactoring mcp-server-planton to follow stigmer/mcp-server architecture:
- Domain-driven tool structure
- Consistent apply/delete/get patterns
- Codegen pipeline from day one (proto → schema → Go input types)
- Three tools: `apply_cloud_resource`, `delete_cloud_resource`, `get_cloud_resource`
- MCP resource templates for per-kind schema discovery

## Current Step

- ✅ **T01 Plan** — APPROVED
- ✅ **Phase 1: Clean Slate + Shared Utilities** (2026-02-26)
  - Migrated MCP SDK: `mark3labs/mcp-go v0.6.0` → `modelcontextprotocol/go-sdk v1.3.0`
  - Deleted all old domain code (55 files, ~9400 lines removed)
  - Built 12-file Stigmer-pattern foundation (auth, config, grpc, domains, server, pkg/mcpserver, main)
  - `go build ./...` and `go vet ./...` pass cleanly
- ✅ **Phase 2 planning** — Design decisions finalized (2026-02-26)
  - See "Resolved Decisions" section below for full list
- ✅ **Phase 2 Stage 1: proto2schema** (2026-02-26)
  - Built proto2schema codegen tool (5 source files, ~760 lines)
  - Parses OpenMCF provider .proto files via local SCM_ROOT convention
  - Generated JSON schemas for 362 providers across 17 cloud platforms
  - StringValueOrRef simplified to string with referenceKind metadata (Option C)
  - Extracts buf.validate rules and OpenMCF custom options via protowire
  - `make codegen-schemas` target added, zero parse errors
- ✅ **Phase 2 Stage 2: schema2go generator** (2026-02-26)
  - Built schema2go codegen tool (3 source files, ~850 lines)
  - Generates typed Go input structs from JSON schemas for all 362 providers
  - snake_case JSON tags (matches PlantON backend `preservingProtoFieldNames()`)
  - Per-provider Parse{Kind}() functions with validate/applyDefaults/toMap
  - Shared type deduplication via types_gen.go per cloud package
  - Central registry with ParseFunc dispatch by kind
  - Hand-written internal/parse/helpers.go for shared utilities
  - `make codegen-types` and `make codegen` targets added
  - 367 generated Go files compile and vet cleanly
- ✅ **Phase 3: apply_cloud_resource + MCP Resource Templates** (2026-02-26)
  - Promoted JSON schemas to top-level `schemas/` package with `go:embed` for runtime access
  - Built `internal/domains/cloudresource/` domain package (6 files): tool, handler, apply gRPC, kind mapping, metadata extraction, schema lookup, resource template
  - Registered `apply_cloud_resource` tool + `cloud-resource-schema://{kind}` resource template in server
  - Added `plantonhq/planton/apis` and `plantonhq/openmcf` dependencies
  - Added shared `ResourceResult` helper to `internal/domains/toolresult.go`
- 🔵 Next: **Phase 4: Implement delete_cloud_resource + get_cloud_resource**

---

### ✅ COMPLETED: Phase 1 — Clean Slate + Shared Utilities (2026-02-26)

**Rebuilt mcp-server-planton foundation from the ground up following Stigmer MCP server patterns.**

**What was delivered:**

1. **MCP SDK migration** — Swapped `mark3labs/mcp-go v0.6.0` (community) for `modelcontextprotocol/go-sdk v1.3.0` (official). Enables typed tool handlers, proper context propagation, native Streamable HTTP.

2. **Context-based auth** (`internal/auth/credentials.go`) — Clean `WithAPIKey`/`APIKey`/`GetAPIKey` context pattern + `TokenAuth` gRPC credentials. Eliminated the global mutex API key store hack.

3. **Config** (`internal/config/config.go`) — Stigmer-pattern env-based config preserving all existing `PLANTON_*` env vars. Added `PLANTON_MCP_LOG_FORMAT` and `PLANTON_MCP_LOG_LEVEL` for structured logging.

4. **gRPC client factory** (`internal/grpc/client.go`) — Centralized `NewConnection` with TLS/:443 convention and optional PerRPCCredentials.

5. **Domain shared utilities** (`internal/domains/`) — `WithConnection` lifecycle helper, `MarshalJSON` protojson, `RPCError` gRPC error classification, `TextResult`/`CallFetch` tool result helpers.

6. **MCP server** (`internal/server/`) — Server init with tool registration placeholder, STDIO transport with context auth injection, Streamable HTTP transport with auth middleware and health probe.

7. **Public embedding API** (`pkg/mcpserver/`) — `Config`/`DefaultConfig`/`Run` for embedding the MCP server in other Go programs.

8. **Entry point** (`cmd/mcp-server-planton/main.go`) — Simplified CLI with subcommand override (stdio/http/both).

**Key Decisions Made:**
- Official MCP SDK (`modelcontextprotocol/go-sdk`) over community SDK — enables typed handlers and context propagation
- Preserved all existing `PLANTON_*` env vars for backward compatibility
- HTTP transport moved from SSE proxy hack to native Streamable HTTP (breaking change for existing HTTP clients)
- Migrated logging from `log.Printf` to `slog` (structured logging)

**Files Changed/Created:**
- `go.mod` — Dependency swap (mcp-go → go-sdk, removed buf.build deps temporarily)
- `cmd/mcp-server-planton/main.go` — Rewritten
- `internal/auth/credentials.go` — New
- `internal/config/config.go` — Rewritten
- `internal/grpc/client.go` — New
- `internal/domains/conn.go` — New
- `internal/domains/marshal.go` — New
- `internal/domains/rpcerr.go` — New
- `internal/domains/toolresult.go` — New
- `internal/server/server.go` — New
- `internal/server/http.go` — New
- `pkg/mcpserver/config.go` — New
- `pkg/mcpserver/run.go` — New
- Deleted: `internal/common/`, `internal/domains/` (old), `internal/mcp/`, `IMPLEMENTATION_SUMMARY*.md`

---

### ✅ COMPLETED: Phase 2 Stage 1 — proto2schema Codegen Tool (2026-02-26)

**Built the proto2schema codegen tool that parses OpenMCF provider .proto files and generates JSON schemas for code generation and MCP resource template discovery.**

**What was delivered:**

1. **proto2schema tool** (`tools/codegen/proto2schema/`) — 5-file Go CLI tool adapted from Stigmer's codegen pipeline. Parses all OpenMCF provider protos using `jhump/protoreflect`, extracts spec fields, nested types, validation rules, and custom OpenMCF options.

2. **362 provider schemas** — Generated JSON schemas for all providers across 17 cloud platforms (AWS, GCP, Azure, Kubernetes, DigitalOcean, Civo, Cloudflare, Confluent, Auth0, OpenFGA, Snowflake, Atlas, AliCloud, HetznerCloud, OCI, OpenStack, Scaleway). Zero parse errors.

3. **Provider registry** (`tools/codegen/schemas/providers/registry.json`) — Kind-to-schema-path index for all 362 providers. Used by Stage 2 generator and MCP resource template handlers.

4. **Shared metadata schema** (`tools/codegen/schemas/shared/metadata.json`) — CloudResourceMetadata fields shared across all providers, with nested CloudResourceRelationship type.

5. **Makefile target** — `make codegen-schemas` runs the full pipeline.

**Key Decisions Made:**
- StringValueOrRef → simplified to `string` with `referenceKind`/`referenceFieldPath` metadata (Option C). Respects bounded context boundary between specification and provisioning layers.
- OpenMCF custom options (`default_kind`, `default_kind_field_path`, `default`, `recommended_default`) extracted via protowire from unknown fields.
- Proto file resolution via `SCM_ROOT` convention (`$HOME/scm/github.com/{org}/{repo}/`) with `--openmcf-apis-dir` CLI override.
- Split into 5 focused files (vs Stigmer's single file) for maintainability at this scale (362 providers vs Stigmer's ~15).

**Files Created:**
- `tools/codegen/proto2schema/main.go` — CLI entry point, provider scanning, buf cache detection
- `tools/codegen/proto2schema/schema.go` — Schema type definitions
- `tools/codegen/proto2schema/parser.go` — Proto parsing, field extraction, validation
- `tools/codegen/proto2schema/options.go` — OpenMCF custom option extraction via protowire
- `tools/codegen/proto2schema/registry.go` — Registry and file writing
- `tools/codegen/schemas/` — 362 provider schemas + registry + shared metadata

**Files Modified:**
- `go.mod` / `go.sum` — Added `jhump/protoreflect`, `buf.build/gen/go/bufbuild/protovalidate`
- `Makefile` — Added `codegen-schemas` target

---

### ✅ COMPLETED: Phase 2 Stage 2 — schema2go Generator (2026-02-26)

**Built the schema2go codegen generator that transforms JSON schemas into typed Go input structs with validation, defaults, map conversion, and a central kind-to-parser registry.**

**What was delivered:**

1. **schema2go generator** (`tools/codegen/generator/`) — 3-file Go CLI tool. Loads provider registry and JSON schemas, generates typed Go input structs with `validate()`, `applyDefaults()`, `toMap()` methods, and top-level `Parse{Kind}()` functions per provider.

2. **367 generated Go files** — 362 per-provider input types, 5 shared `types_gen.go` files for deduplicated nested types, 1 `registry_gen.go` central dispatch. All organized under `gen/cloudresource/{cloud}/` (17 cloud packages).

3. **Central registry** (`gen/cloudresource/registry_gen.go`) — `ParseFunc` type, `GetParser(kind)` lookup, `KnownKinds()` enumeration. Imports all 17 cloud packages.

4. **Shared parse helpers** (`internal/parse/helpers.go`) — Hand-written utilities (`ValidateHeader`, `ExtractSpecMap`, `RebuildCloudObject`) shared by all generated Parse functions. Prevents circular dependencies.

5. **Makefile targets** — `make codegen-types` (Stage 2 only), `make codegen` (full pipeline: schemas + types).

**Key Decisions Made:**
- **snake_case JSON tags** — PlantON backend uses `JsonFormat.printer().preservingProtoFieldNames()` and MongoDB stores with snake_case keys. Verified via Java backend `CloudResourceMapper` and `ValueFromToValueResolver` source code.
- **toMap() instead of ToProto()** — Generated types convert to `map[string]any` (for `structpb.Struct`) rather than concrete proto messages, since `cloud_object` uses `google.protobuf.Struct`.
- **Multi-package structure** — One Go package per cloud provider under `gen/cloudresource/` for clean namespacing at scale (362 providers).
- **Shared type deduplication** — Common nested types (e.g., `ContainerInput`, `ProbeInput`) generated once per cloud package in `types_gen.go`.
- **Generate all 362 providers** — Marginal cost is minimal; ensures comprehensive coverage from day one.

**Files Created:**
- `tools/codegen/generator/main.go` — CLI entry point, schema loading, orchestration
- `tools/codegen/generator/codegen.go` — Core struct/method/parse-function generation
- `tools/codegen/generator/registry.go` — Registry file generation
- `internal/parse/helpers.go` — Hand-written shared utilities
- `gen/cloudresource/` — 367 generated `.go` files across 17 cloud packages

**Files Modified:**
- `Makefile` — Added `codegen-types` and `codegen` targets

---

### ✅ COMPLETED: Phase 3 — apply_cloud_resource + MCP Resource Templates (2026-02-26)

**Implemented the first working MCP tool with typed validation via generated parsers, and MCP resource templates for per-kind schema discovery.**

**What was delivered:**

1. **Schema promotion** — Moved 362 JSON schemas from `tools/codegen/schemas/` to top-level `schemas/` package. Created `schemas/embed.go` with `//go:embed` directive. Updated codegen tool defaults and Makefile targets. Clean separation: `tools/codegen/` is build tooling, `schemas/` is shared domain data.

2. **`apply_cloud_resource` MCP tool** (`internal/domains/cloudresource/tools.go`) — Accepts opaque `cloud_object` map keeping tool schema small (no 50k+ token provider explosion). Handler pipeline: extract kind → get parser from registry → validate + normalize spec → build CloudResource proto → gRPC Apply → return JSON response.

3. **Domain functions** (`internal/domains/cloudresource/apply.go`) — `Apply()` calls `CloudResourceCommandController.Apply` via `domains.WithConnection`, `buildCloudResource()` assembles the full proto with api_version, kind, metadata, and spec.cloud_object.

4. **Kind mapping** (`internal/domains/cloudresource/kind.go`) — `resolveKind()` maps PascalCase kind strings to `CloudResourceKind` enum values from openmcf proto stubs.

5. **Metadata extraction** (`internal/domains/cloudresource/metadata.go`) — `extractMetadata()` maps cloud_object["metadata"] to `ApiResourceMetadata` proto. Required: name, org, env. Optional: slug, id, labels, annotations, tags, version.message.

6. **MCP resource templates** (`internal/domains/cloudresource/resources.go`, `schema.go`) — `cloud-resource-schema://{kind}` URI template backed by embedded JSON schemas. Registry-based lookup with `sync.Once` caching. Agents discover per-kind schemas before calling apply.

7. **Server registration** (`internal/server/server.go`) — New `registerResources()` function alongside existing `registerTools()`. Both tool and resource template registered at startup.

8. **Shared resource result helper** (`internal/domains/toolresult.go`) — Added `ResourceResult()` for constructing `ReadResourceResult` responses.

**Key Decisions Made:**
- JSON schemas promoted from `tools/codegen/schemas/` to `schemas/` — respects bounded context boundary (build tooling vs runtime data)
- Raw JSON schemas served via resource templates (not Go struct reflection) — schemas contain richer validation rules, descriptions, and metadata than generated Go types
- `cloud-resource-schema://` custom URI scheme — standard URL parsing, kind as host component
- Registry cached with `sync.Once` — loaded once from embedded FS, no repeated I/O

**Files Created:**
- `schemas/embed.go` — `go:embed` package for runtime schema access
- `internal/domains/cloudresource/tools.go` — Tool definition + typed handler
- `internal/domains/cloudresource/apply.go` — gRPC Apply + proto assembly
- `internal/domains/cloudresource/kind.go` — CloudResourceKind enum resolution
- `internal/domains/cloudresource/metadata.go` — ApiResourceMetadata extraction
- `internal/domains/cloudresource/resources.go` — MCP resource template definition + handler
- `internal/domains/cloudresource/schema.go` — Embedded FS schema lookup + URI parsing

**Files Modified:**
- `internal/server/server.go` — Added tool + resource template registration
- `internal/domains/toolresult.go` — Added `ResourceResult()` helper
- `Makefile` — Updated codegen targets for new schema location
- `tools/codegen/proto2schema/main.go` — Updated default `--output-dir`
- `tools/codegen/generator/main.go` — Updated default `--schemas-dir`
- `go.mod` / `go.sum` — Added `plantonhq/planton/apis`, `plantonhq/openmcf`
- `schemas/` — 362 JSON schemas + registry + shared metadata (moved from `tools/codegen/schemas/`)

---

## Execution Order

### Phase 1: Clean Slate + Shared Utilities ✅
Delete existing domain code, set up Stigmer-style foundation.

### Phase 2: Codegen Pipeline
Adapt Stigmer's two-stage codegen for OpenMCF provider specs:
- Stage 1: `proto2schema` — Parse OpenMCF provider .proto files → JSON schemas ✅
- Stage 2: `generator` — JSON schemas → Go input types with `toMap()` for each provider kind ✅
- Central kind-to-parser registry for runtime dispatch ✅
- Makefile targets: `codegen-schemas` (Stage 1) ✅, `codegen-types` (Stage 2) ✅, `codegen` (full pipeline) ✅

### Phase 3: Implement apply_cloud_resource + MCP Resource Templates
- First working MCP tool with generated input types
- `cloud_object` stays opaque in tool schema (typed validation happens inside handler)
- MCP resource templates expose per-kind typed schemas for client discovery
- No separate schema lookup tool — agents use MCP resources

### Phase 4: Implement delete_cloud_resource + get_cloud_resource
Complete the tool set.

### Phase 5: Testing + Documentation

## Key References

- **Stigmer MCP server** (reference): `@stigmer/mcp-server/`
- **Stigmer codegen**: `@stigmer/tools/codegen/`
- **Planton cloud resource protos**: `@planton/apis/ai/planton/infrahub/cloudresource/v1/`
- **OpenMCF provider specs**: `@openmcf/apis/org/openmcf/provider/`
- **Design decisions**: `_projects/2026-02/20260226.01.refactor-mcp-server-stigmer-patterns/design-decisions/`
- **Phase 1 plan**: `_projects/2026-02/20260226.01.refactor-mcp-server-stigmer-patterns/plans/phase-1-foundation.plan.md`
- **Proto2schema plan**: `_projects/2026-02/20260226.01.refactor-mcp-server-stigmer-patterns/plans/proto2schema-codegen-tool.plan.md`
- **Generated schemas**: `schemas/` (promoted from `tools/codegen/schemas/` in Phase 3)
- **Cloud resource domain**: `internal/domains/cloudresource/`
- **Phase 3 plan**: `_projects/2026-02/20260226.01.refactor-mcp-server-stigmer-patterns/plans/phase-3-apply-cloud-resource.plan.md`

## Resolved Decisions

1. **Cloud object format**: Full OpenMCF message (api_version, kind, metadata) but NOT status
2. **Tool naming**: `apply_cloud_resource` / `delete_cloud_resource` / `get_cloud_resource`
3. **Codegen**: Build from day one, no hand-written types
4. **get_cloud_resource**: Included in scope
5. **MCP SDK**: Official `modelcontextprotocol/go-sdk` (not community `mark3labs/mcp-go`)
6. **HTTP transport**: Streamable HTTP (native SDK support, replaces SSE proxy)
7. **Logging**: `slog` structured logging (replaces `log.Printf`)
8. **Typed provider codegen (Option B)**: Generate typed Go input structs for OpenMCF provider specs
   starting with a subset (~10-20 most common providers), expand to all ~150 later.
   Same discriminated-union pattern as Stigmer workflow task configs
   (kind + Struct with `discriminated_by`).
9. **Tool schema stays small**: Do NOT expand 150+ typed provider fields into the
   `apply_cloud_resource` tool input schema — that would be 50,000-100,000+ tokens,
   overwhelming MCP clients. The tool keeps `cloud_object` as opaque `map[string]any`.
   Typed validation happens inside the handler using generated input structs.
10. **Schema discovery via MCP resource templates**: Expose per-kind typed schemas
    as MCP resource templates (e.g., `cloud-resource-schema://{kind}`). Clients
    fetch the schema for the specific kind they need before calling apply.
    No separate `get_cloud_resource_schema` tool — agents use MCP resources.
11. **Dependency: Stigmer agent runner MCP resources support**: The Stigmer agent
    runner currently only uses MCP tools (via `langchain_mcp_adapters`), not MCP
    resources. A separate project in the stigmer repo will add MCP resources
    support so agents can auto-discover schemas.
    See: `stigmer/_projects/2026-02/20260226.02.agent-runner-mcp-resources/`
