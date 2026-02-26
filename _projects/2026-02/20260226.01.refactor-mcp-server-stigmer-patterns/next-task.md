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
**Current Status**: ALL PHASES COMPLETED (Phases 1-5)

## Quick Context

Refactoring mcp-server-planton to follow stigmer/mcp-server architecture:
- Domain-driven tool structure
- Consistent apply/delete/get patterns
- Codegen pipeline from day one (proto → schema → Go input types)
- Three tools: `apply_cloud_resource`, `delete_cloud_resource`, `get_cloud_resource`
- MCP resource templates for per-kind schema discovery
- Static MCP resource for cloud resource kind catalog (agent discovery)

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
- ✅ **Kind Catalog Resource** (2026-02-26)
  - Added static `cloud-resource-kinds://catalog` MCP resource serving 362 kinds grouped by 17 providers
  - Thorough analysis confirmed PascalCase is canonical for kind values across all system layers (proto enum, Go stubs, codegen registry, JSON schemas, StringValueOrRef options)
  - Updated tool descriptions with 3-step discovery workflow: catalog → schema → apply
  - Updated error messages to direct agents to catalog resource for kind discovery
- ✅ **Phase 4: Implement delete_cloud_resource + get_cloud_resource** (2026-02-27)
  - Shared `ResourceIdentifier` type with dual-path validation (ID or kind+org+env+slug)
  - `get_cloud_resource` tool: ID path → `QueryController.Get`, slug path → `QueryController.GetByOrgByEnvByKindBySlug`
  - `delete_cloud_resource` tool: Stigmer two-step pattern (resolve to ID via query, then `CommandController.Delete`)
  - Both tools registered in server, tool count updated to 3
  - Design decision: slug uniqueness scoped to (org, env, kind) — confirmed by gRPC API requiring all four fields
  - `go build ./...` and `go vet ./...` pass cleanly
- ✅ **Phase 5: Hardening — Dead Code, Unit Tests, Documentation** (2026-02-27)
  - Removed dead code: `FetchFunc`/`CallFetch` from `toolresult.go` (superseded by `ResourceIdentifier` pattern)
  - 8 test files, 50 test cases covering all pure domain logic (zero test files → comprehensive coverage)
  - README.md rewritten from scratch for the 3-tool architecture
  - Deleted 5 stale docs files (62KB of misleading content); recreated `docs/configuration.md` and `docs/development.md`
  - `go build ./...`, `go vet ./...`, and all 50 tests pass cleanly

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

### ✅ COMPLETED: Kind Catalog Resource (2026-02-26)

**Added a static MCP resource enabling agents to discover all supported cloud resource kinds before using the schema template or calling tools.**

**What was delivered:**

1. **Kind case analysis** — Thorough cross-system audit confirmed PascalCase is the canonical form for kind values across all layers: proto enum (`AwsVpc = 216`), generated Go maps (`CloudResourceKind_value["AwsVpc"]`), codegen registry, JSON schema registry, and StringValueOrRef `default_kind` options. All 362 production kinds in the codegen registry are a perfect subset of the proto enum. No case conversion needed at any boundary.

2. **`cloud-resource-kinds://catalog` static resource** (`internal/domains/cloudresource/resources.go`) — `KindCatalogResource()` and `KindCatalogHandler()` serve a grouped JSON catalog of all 362 kinds across 17 cloud providers. Each provider entry includes `api_version` and sorted kind list.

3. **Catalog data builder** (`internal/domains/cloudresource/schema.go`) — `buildKindCatalog()` transforms the embedded registry into the grouped catalog JSON, cached with `sync.Once`. Reuses the existing `loadRegistry()` cache.

4. **Updated tool descriptions** — `apply_cloud_resource` tool description now includes a 3-step workflow: catalog → schema → apply. Error messages point agents to `cloud-resource-kinds://catalog` for kind discovery.

**Key Decisions Made:**
- PascalCase confirmed as canonical for kind values — no conversion needed anywhere
- Static `cloud-resource-kinds://catalog` URI scheme (not parameterized) since there's only one catalog
- Catalog grouped by cloud provider with `api_version` per group — agents can narrow by provider and get `api_version` without extra lookup
- Catalog size is ~12.7KB covering 362 kinds across 17 providers — well within MCP resource limits

**Files Changed:**
- `internal/domains/cloudresource/schema.go` — Added `buildKindCatalog()`, `kindCatalog`, `catalogProviderEntry` types
- `internal/domains/cloudresource/resources.go` — Added `KindCatalogResource()`, `KindCatalogHandler()`
- `internal/server/server.go` — Registered static catalog resource via `srv.AddResource()`
- `internal/domains/cloudresource/tools.go` — Updated tool descriptions with 3-step workflow
- `internal/domains/cloudresource/kind.go` — Updated error message to reference catalog

---

### ✅ COMPLETED: Phase 4 — delete_cloud_resource + get_cloud_resource (2026-02-27)

**Implemented `get_cloud_resource` and `delete_cloud_resource` MCP tools with dual-path resource identification, completing the core tool set.**

**What was delivered:**

1. **Shared `ResourceIdentifier`** (`internal/domains/cloudresource/identifier.go`) — Dual-path type supporting ID (single field) or composite key (kind, org, env, slug). `validateIdentifier()` ensures exactly one path is fully specified with clear error messages for partial inputs. `describeIdentifier()` for human-readable error context. `resolveResourceID()` for delete's two-step slug-to-ID resolution.

2. **`get_cloud_resource` tool** (`internal/domains/cloudresource/get.go`, `tools.go`) — ID path calls `QueryController.Get(CloudResourceId)`, slug path resolves PascalCase kind to proto enum then calls `QueryController.GetByOrgByEnvByKindBySlug()`. No unnecessary resolution step — both RPCs return the full `CloudResource`.

3. **`delete_cloud_resource` tool** (`internal/domains/cloudresource/delete.go`, `tools.go`) — Follows Stigmer two-step pattern: `resolveResourceID()` handles both paths within a single gRPC connection, then `CommandController.Delete(ApiResourceDeleteInput)` with the resolved ID.

4. **Tool registration** (`internal/server/server.go`) — Both tools registered, tool count updated to 3 with named list in structured log.

**Key Decisions Made:**
- Slug uniqueness scoped to (org, env, kind) — confirmed by gRPC API requiring all four fields in `CloudResourceByOrgByEnvByKindBySlugRequest`
- Delete kept simple: skipped `version_message` and `force` fields (can add later)
- Error handling: kind validation errors pass through directly; gRPC errors classified via `domains.RPCError()`. `resolveResourceID` owns its error formatting.
- Validation at handler boundary; business logic assumes valid inputs.
- `FetchFunc`/`CallFetch` in `toolresult.go` not used by cloudresource (flagged as potential dead code if no other domains are added)

**Files Created:**
- `internal/domains/cloudresource/identifier.go` — ResourceIdentifier, validation, resolution
- `internal/domains/cloudresource/get.go` — Get function (dual-path query)
- `internal/domains/cloudresource/delete.go` — Delete function (resolve + command)

**Files Modified:**
- `internal/domains/cloudresource/tools.go` — Added GetTool/GetHandler, DeleteTool/DeleteHandler, input structs
- `internal/server/server.go` — Registered both tools, updated log

---

### ✅ COMPLETED: Phase 5 — Hardening: Dead Code, Unit Tests, Documentation (2026-02-27)

**Cleaned up dead code, wrote comprehensive unit tests for all pure domain logic, and rewrote the README and docs to reflect the new 3-tool architecture.**

**What was delivered:**

1. **Dead code removal** — Deleted `FetchFunc` type and `CallFetch` function from `internal/domains/toolresult.go`. These implemented an `org + slug` lookup pattern superseded by Phase 4's `ResourceIdentifier` dual-path pattern and were never called.

2. **8 unit test files, 50 test cases** — Comprehensive coverage of all pure domain logic:
   - `identifier_test.go` (8 tests): `validateIdentifier` branching paths, `describeIdentifier` output
   - `kind_test.go` (6 tests): kind extraction, enum resolution, error cases
   - `metadata_test.go` (10 tests): required/optional field extraction, type assertion edge cases
   - `schema_test.go` (10 tests): URI parsing, embedded FS integration, catalog JSON structure
   - `apply_test.go` (6 tests): CloudResource proto assembly, `describeResource` fallback
   - `rpcerr_test.go` (7 tests): gRPC status code to user-facing message mapping
   - `helpers_test.go` (10 tests): `ValidateHeader`, `ExtractSpecMap`, `RebuildCloudObject`
   - `config_test.go` (10 tests): config validation rules, log level parsing

3. **README.md rewrite** — Complete rewrite for the new 3-tool architecture: tools table, MCP resources table, 3-step agent workflow, integration guides (Cursor/Claude Desktop/LangGraph/Docker), configuration table, architecture overview, codegen pipeline.

4. **Stale docs cleanup** — Deleted 5 files (62KB) under `docs/` that documented tools and domains deleted in Phase 1 (Service Hub, Connect/Credentials, old tool names). Recreated `docs/configuration.md` (current env var reference) and `docs/development.md` (build, test, codegen pipeline, project structure).

**Key Decisions Made:**
- Test pure functions, not gRPC wiring — tool handlers are thin proxies; risk is in validation/transformation logic
- No mock gRPC server — adds complexity without proportional confidence; defer to integration testing
- No codegen determinism tests — output validated by `go build` and `go vet`; add when codegen stabilizes

**Files Created:**
- `internal/config/config_test.go`
- `internal/domains/rpcerr_test.go`
- `internal/domains/cloudresource/identifier_test.go`
- `internal/domains/cloudresource/kind_test.go`
- `internal/domains/cloudresource/metadata_test.go`
- `internal/domains/cloudresource/schema_test.go`
- `internal/domains/cloudresource/apply_test.go`
- `internal/parse/helpers_test.go`
- `docs/configuration.md` — Rewritten
- `docs/development.md` — Rewritten

**Files Modified:**
- `internal/domains/toolresult.go` — Removed dead `FetchFunc`/`CallFetch`
- `README.md` — Complete rewrite

**Files Deleted:**
- `docs/service-hub-tools.md`
- `docs/installation.md`
- `docs/http-transport.md`
- `docs/development.md` (old)
- `docs/configuration.md` (old)

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

### Phase 4: Implement delete_cloud_resource + get_cloud_resource ✅
Complete the tool set with dual-path resource identification (ID or kind+org+env+slug).

### Phase 5: Hardening — Dead Code, Unit Tests, Documentation ✅
50 unit tests, dead code cleanup, README + docs rewrite.

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
- **Kind catalog plan**: `_projects/2026-02/20260226.01.refactor-mcp-server-stigmer-patterns/plans/cloud-resource-kinds-catalog.plan.md`
- **Phase 4 plan**: `.cursor/plans/phase_4_delete_get_tools_b9314bd3.plan.md`

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
