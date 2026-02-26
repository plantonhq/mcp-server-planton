---
name: Phase 1 Foundation
overview: "Phase 1 rebuilds the mcp-server-planton foundation by migrating from the community MCP SDK (mark3labs/mcp-go) to the official SDK (modelcontextprotocol/go-sdk), deleting all existing domain code, and laying down the Stigmer-pattern skeleton: auth, config, gRPC, shared domain utilities, server, public embedding API, and entry point."
todos:
  - id: swap-sdk
    content: "Update go.mod: replace mark3labs/mcp-go with modelcontextprotocol/go-sdk, verify Go version compatibility"
    status: completed
  - id: delete-old
    content: Delete internal/common/, internal/domains/, internal/mcp/, and old summary docs
    status: completed
  - id: auth
    content: Create internal/auth/credentials.go — context-based API key + TokenAuth (no global store)
    status: completed
  - id: config
    content: Rewrite internal/config/config.go — Stigmer pattern with Planton env vars preserved, add structured logging config
    status: completed
  - id: grpc
    content: Create internal/grpc/client.go — centralized NewConnection factory
    status: completed
  - id: domain-utils
    content: "Create internal/domains/ shared utilities: conn.go, marshal.go, rpcerr.go, toolresult.go"
    status: completed
  - id: server
    content: Create internal/server/server.go + http.go — MCP server init, STDIO, Streamable HTTP, auth middleware
    status: completed
  - id: pkg-mcpserver
    content: Create pkg/mcpserver/config.go + run.go — public embedding API
    status: completed
  - id: main
    content: Rewrite cmd/mcp-server-planton/main.go — simplified entry point using pkg/mcpserver
    status: completed
  - id: verify-build
    content: Run go build, go vet to confirm zero-tool server compiles and boots cleanly
    status: completed
isProject: false
---

# Phase 1: Clean Slate + Shared Utilities

## Surprises Discovered During Research

Three findings that were not explicitly called out in the T01 plan and need acknowledgement before we begin:

### 1. MCP SDK Migration (Critical)

The T01 plan describes adopting "Stigmer-style patterns" but does not explicitly call out the **library swap** that makes those patterns possible:

- **Current**: `github.com/mark3labs/mcp-go v0.6.0` (community SDK)
- **Target**: `github.com/modelcontextprotocol/go-sdk` (official SDK, Stigmer uses v1.3.0)

This is not cosmetic. The official SDK has a fundamentally different API:

- **Typed tool handlers**: `func(ctx, *CallToolRequest, *TypedInput) (*CallToolResult, any, error)` vs the old `func(map[string]interface{}) (*CallToolResult, error)`
- **Context propagation**: Context flows through to tool handlers natively, eliminating the need for the global API key store hack
- **Streamable HTTP**: Native `mcp.NewStreamableHTTPHandler()` replaces the SSE proxy workaround

Every file we write in Phase 1 depends on this SDK choice.

### 2. HTTP Transport is a Breaking Change

The current server uses SSE transport with a reverse-proxy workaround (internal port 18080 proxied through external port 8080, with URL rewriting). The official SDK provides native Streamable HTTP, which is the MCP spec's recommended transport. This is a clean replacement, but any existing HTTP clients will need to adapt.

### 3. Global API Key Store Hack Disappears

The current [internal/common/auth/credentials.go](internal/common/auth/credentials.go) has a `globalAPIKeyStore` with mutex-protected read/write — a workaround because `mark3labs/mcp-go` does not pass context to tool handlers. The official SDK eliminates this entirely. This is a significant quality improvement.

---

## Environment Variable Backward Compatibility

The new config will preserve all existing env var names with `PLANTON_` prefix and add two new ones for structured logging:


| Env Var                         | Status  | Notes                                |
| ------------------------------- | ------- | ------------------------------------ |
| `PLANTON_API_KEY`               | Keep    | API key for STDIO/both modes         |
| `PLANTON_APIS_GRPC_ENDPOINT`    | Keep    | Direct endpoint override             |
| `PLANTON_CLOUD_ENVIRONMENT`     | Keep    | Preset selection (live/test/local)   |
| `PLANTON_MCP_TRANSPORT`         | Keep    | stdio/http/both                      |
| `PLANTON_MCP_HTTP_PORT`         | Keep    | Default 8080                         |
| `PLANTON_MCP_HTTP_AUTH_ENABLED` | Keep    | Default true                         |
| `PLANTON_MCP_LOG_FORMAT`        | **New** | text/json (default text)             |
| `PLANTON_MCP_LOG_LEVEL`         | **New** | debug/info/warn/error (default info) |


The environment preset system (`PLANTON_CLOUD_ENVIRONMENT` -> live/test/local endpoints) is kept because it provides genuine convenience. This is Planton-specific — Stigmer doesn't need it because Stigmer only has one backend.

---

## Execution Steps

### Step 1: Update go.mod — Swap MCP SDK

Replace `github.com/mark3labs/mcp-go v0.6.0` with `github.com/modelcontextprotocol/go-sdk`. Run `go get` to resolve the version and update `go.sum`. Verify Go version compatibility (current: `go 1.24.0`, may need bump depending on SDK requirements).

Keep all existing `buf.build/gen/go/blintora/apis/...` and `buf.build/gen/go/project-planton/apis/...` dependencies — those are the Planton proto stubs and don't change.

### Step 2: Delete Old Code

Remove these directories and files entirely:

- `internal/common/` (auth hack + error helpers — replaced by new auth and rpcerr)
- `internal/domains/` (all domain implementations — rebuilt in Phase 3/4)
- `internal/mcp/` (old server.go + http_server.go — replaced by internal/server/)
- `IMPLEMENTATION_SUMMARY.md` (documents old architecture)
- `IMPLEMENTATION_SUMMARY_PIPELINE_LOGS.md` (documents old architecture)

### Step 3: Create `internal/auth/credentials.go`

Adapted from Stigmer's [mcp-server/internal/auth/credentials.go](../stigmer/stigmer/mcp-server/internal/auth/credentials.go):

- `WithAPIKey(ctx, key) context.Context` — store key in context
- `APIKey(ctx) string` — retrieve key (returns empty if absent)
- `GetAPIKey(ctx) (string, error)` — strict retrieval (errors if absent)
- `TokenAuth` struct implementing `grpc.PerRPCCredentials`
- `NewTokenAuth(token) TokenAuth`
- **No global store** — context-only

### Step 4: Create `internal/config/config.go`

Following Stigmer's pattern but preserving Planton-specific env vars and the environment preset system:

- `Config` struct with: `ServerAddress`, `APIKey`, `Transport`, `HTTPPort`, `HTTPAuthEnabled`, `LogFormat`, `LogLevel`
- `LoadFromEnv() (*Config, error)` — reads `PLANTON_*` env vars
- `Validate() error` — checks invariants
- Endpoint resolution: `PLANTON_APIS_GRPC_ENDPOINT` (override) > `PLANTON_CLOUD_ENVIRONMENT` (preset) > default `api.live.planton.ai:443`
- Structured logging support: `LogFormat` (text/json) and `LogLevel` (slog.Level)
- Transport, LogFormat as typed string constants

### Step 5: Create `internal/grpc/client.go`

Adapted from Stigmer's [mcp-server/internal/grpc/client.go](../stigmer/stigmer/mcp-server/internal/grpc/client.go):

- `DefaultRPCTimeout = 30 * time.Second`
- `NewConnection(endpoint, apiKey string) (*grpc.ClientConn, error)` — TLS for :443, insecure otherwise; attaches `auth.NewTokenAuth` as `PerRPCCredentials` when apiKey is non-empty

### Step 6: Create `internal/domains/` shared utilities

Four files, adapted from Stigmer's domain utilities:

- `**conn.go`**: `WithConnection(ctx, serverAddress, fn) (string, error)` — creates authenticated gRPC connection, applies timeout, ensures cleanup
- `**marshal.go`**: `MarshalJSON(msg proto.Message) (string, error)` — protojson with consistent options (multiline, proto names, no unpopulated)
- `**rpcerr.go**`: `RPCError(err, resourceDesc) error` — translates gRPC status codes to user-friendly messages, logs original at WARN via slog
- `**toolresult.go**`: `TextResult(text)`, `CallFetch(fn, ctx, serverAddr, ...)` — MCP result construction helpers using official SDK types

### Step 7: Create `internal/server/server.go`

Adapted from Stigmer's [mcp-server/internal/server/server.go](../stigmer/stigmer/mcp-server/internal/server/server.go):

- `Server` struct wrapping `*mcp.Server` and `*config.Config`
- `New(cfg) *Server` — creates MCP server, calls `registerTools` (empty for now — placeholder until Phase 3/4)
- `registerTools(srv, serverAddress)` — will be populated in Phase 3/4
- `ServeStdio(ctx) error` — injects API key into context, runs STDIO transport
- Build-time version injection via `var buildVersion string`

### Step 8: Create `internal/server/http.go`

Adapted from Stigmer's [mcp-server/internal/server/http.go](../stigmer/stigmer/mcp-server/internal/server/http.go):

- `ServeHTTP(ctx) error` — Streamable HTTP via `mcp.NewStreamableHTTPHandler()`, graceful shutdown on context cancellation
- `authMiddleware(next) http.Handler` — extracts Bearer token, injects via `auth.WithAPIKey`
- `healthHandler` — simple 200 OK for liveness probes
- `requestLogger` middleware with short request IDs via `crypto/rand`

### Step 9: Create `pkg/mcpserver/config.go` and `pkg/mcpserver/run.go`

Public embedding API following Stigmer's [mcp-server/pkg/mcpserver/](../stigmer/stigmer/mcp-server/pkg/mcpserver/) pattern:

- `**config.go**`: Public `Config` struct (plain Go types, no internal imports for callers), `DefaultConfig() (*Config, error)`, internal conversion helpers
- `**run.go**`: `Run(ctx, *Config) error` — initializes logger, creates server, dispatches to transport, `serveBoth` for dual mode, `isNormalShutdown` for clean exit on EOF/EPIPE

### Step 10: Rewrite `cmd/mcp-server-planton/main.go`

Adapted from Stigmer's [mcp-server/cmd/mcp-server-stigmer/main.go](../stigmer/stigmer/mcp-server/cmd/mcp-server-stigmer/main.go):

- `mcpserver.DefaultConfig()` for config loading
- Positional subcommand override (`stdio`, `http`, `both`)
- `signal.NotifyContext` for clean shutdown
- `mcpserver.Run(ctx, cfg)` — delegates to pkg/mcpserver

### Step 11: Verify Build

Run `go build ./...` and `go vet ./...` to confirm everything compiles cleanly with zero tools registered. The server should start, accept connections, and report an empty tool list.

---

## Files Created/Modified/Deleted Summary

**Deleted:**

- `internal/common/` (entire directory)
- `internal/domains/` (entire directory)
- `internal/mcp/` (entire directory)
- `IMPLEMENTATION_SUMMARY.md`
- `IMPLEMENTATION_SUMMARY_PIPELINE_LOGS.md`

**Created (new):**

- `internal/auth/credentials.go`
- `internal/grpc/client.go`
- `internal/domains/conn.go`
- `internal/domains/marshal.go`
- `internal/domains/rpcerr.go`
- `internal/domains/toolresult.go`
- `internal/server/server.go`
- `internal/server/http.go`
- `pkg/mcpserver/config.go`
- `pkg/mcpserver/run.go`

**Rewritten:**

- `internal/config/config.go`
- `cmd/mcp-server-planton/main.go`
- `go.mod` (dependency swap)

**Untouched:**

- `Makefile`, `Dockerfile`, `.goreleaser.yaml`, `.github/`, `docs/`, `README.md`, `_kustomize/`, `_changelog/`, `LICENSE`, etc.

