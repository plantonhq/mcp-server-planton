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
**Current Status**: Phase 1 COMPLETED → Ready for Phase 2

## Quick Context

Refactoring mcp-server-planton to follow stigmer/mcp-server architecture:
- Domain-driven tool structure
- Consistent apply/delete/get patterns
- Codegen pipeline from day one (proto → schema → Go input types)
- Three tools: `apply_cloud_resource`, `delete_cloud_resource`, `get_cloud_resource`

## Current Step

- ✅ **T01 Plan** — APPROVED
- ✅ **Phase 1: Clean Slate + Shared Utilities** (2026-02-26)
  - Migrated MCP SDK: `mark3labs/mcp-go v0.6.0` → `modelcontextprotocol/go-sdk v1.3.0`
  - Deleted all old domain code (55 files, ~9400 lines removed)
  - Built 12-file Stigmer-pattern foundation (auth, config, grpc, domains, server, pkg/mcpserver, main)
  - `go build ./...` and `go vet ./...` pass cleanly
- 🔵 Next: **Phase 2: Codegen Pipeline**

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

## Execution Order

### Phase 1: Clean Slate + Shared Utilities ✅
Delete existing domain code, set up Stigmer-style foundation.

### Phase 2: Codegen Pipeline
Adapt Stigmer's two-stage codegen:
- Stage 1: `proto2schema` — Proto → JSON schemas
- Stage 2: `generator --target=mcp` — JSON schemas → Go input types with `ToProto()`
- Makefile targets for codegen

### Phase 3: Implement apply_cloud_resource
First working MCP tool with generated input types.

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

## Resolved Decisions

1. **Cloud object format**: Full OpenMCF message (api_version, kind, metadata) but NOT status
2. **Tool naming**: `apply_cloud_resource` / `delete_cloud_resource` / `get_cloud_resource`
3. **Codegen**: Build from day one, no hand-written types
4. **get_cloud_resource**: Included in scope
5. **MCP SDK**: Official `modelcontextprotocol/go-sdk` (not community `mark3labs/mcp-go`)
6. **HTTP transport**: Streamable HTTP (native SDK support, replaces SSE proxy)
7. **Logging**: `slog` structured logging (replaces `log.Printf`)
