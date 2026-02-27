---
name: Phase 5 Hardening
overview: Clean up dead code, write comprehensive unit tests for all pure domain logic, and rewrite the README and docs to reflect the new 3-tool architecture.
todos:
  - id: dead-code-cleanup
    content: Remove FetchFunc and CallFetch from toolresult.go, verify go build passes
    status: completed
  - id: test-identifier
    content: Unit tests for validateIdentifier and describeIdentifier (identifier_test.go)
    status: completed
  - id: test-kind
    content: Unit tests for extractKindFromCloudObject and resolveKind (kind_test.go)
    status: completed
  - id: test-metadata
    content: Unit tests for extractMetadata, toStringMap, toStringSlice (metadata_test.go)
    status: completed
  - id: test-schema
    content: Unit tests for parseSchemaURI + embedded FS integration tests (schema_test.go)
    status: completed
  - id: test-apply
    content: Unit tests for buildCloudResource and describeResource (apply_test.go)
    status: completed
  - id: test-rpcerr
    content: Unit tests for classifyCode (rpcerr_test.go)
    status: completed
  - id: test-parse-helpers
    content: Unit tests for ValidateHeader, ExtractSpecMap, RebuildCloudObject (helpers_test.go)
    status: completed
  - id: test-config
    content: Unit tests for Config.Validate and ParseLogLevel (config_test.go)
    status: completed
  - id: delete-stale-docs
    content: Delete all 5 stale files under docs/
    status: completed
  - id: rewrite-readme
    content: Rewrite README.md for the new 3-tool architecture with MCP resources
    status: completed
  - id: recreate-docs
    content: Recreate docs/configuration.md and docs/development.md with current content
    status: completed
isProject: false
---

# Phase 5: Hardening -- Dead Code, Unit Tests, Documentation

## Architectural Decisions

- **Test pure functions, not gRPC wiring.** The tool handlers (Apply, Get, Delete) are thin gRPC proxies. The real risk is in validation, metadata extraction, identifier logic, URI parsing, and config validation. All of these are pure functions that can be tested without mocks.
- **No mock gRPC server.** Adds complexity without proportional confidence. Defer to integration testing against a real backend (manual or CI).
- **No codegen determinism tests yet.** The codegen output is validated by `go build` and `go vet`. Determinism tests add value only when the codegen stabilizes and multiple contributors touch it.
- **README rewrite is a separate stage from tests.** Both are critical, but should be reviewed independently.

---

## Stage 1: Dead Code Cleanup

Remove unused code that was superseded by Phase 4's `ResourceIdentifier` pattern.

**File:** [internal/domains/toolresult.go](internal/domains/toolresult.go)

- Delete `FetchFunc` type (line 9-11)
- Delete `CallFetch` function (lines 22-29)
- Keep `TextResult` and `ResourceResult` (both actively used)

**Verification:** `go build ./...` must pass after removal.

---

## Stage 2: Unit Tests for Pure Domain Logic

Create test files alongside the source files (standard Go convention). Every test file targets pure functions with no network/filesystem dependencies (except embedded FS).

### 2a. `internal/domains/cloudresource/identifier_test.go`

Test `validateIdentifier` and `describeIdentifier`:

- ID-only path (valid)
- Slug-only path with all 4 fields (valid)
- Both ID and slug fields present (error)
- Partial slug fields with specific missing-field error messages
- All fields empty (error)
- `describeIdentifier` output for both paths

### 2b. `internal/domains/cloudresource/kind_test.go`

Test `extractKindFromCloudObject` and `resolveKind`:

- Valid kind string extraction
- Missing "kind" key
- Non-string kind value
- Empty string kind
- `resolveKind` with known kind (e.g. "AwsVpc")
- `resolveKind` with unknown kind

### 2c. `internal/domains/cloudresource/metadata_test.go`

Test `extractMetadata`, `toStringMap`, `toStringSlice`:

- Valid metadata with all required fields
- Missing metadata key
- Non-object metadata value
- Missing each required field (name, org, env) individually
- Optional fields present (slug, id, labels, annotations, tags, version)
- Optional fields absent (still succeeds)
- Labels/annotations with non-string values (silently skipped)

### 2d. `internal/domains/cloudresource/schema_test.go`

Test `parseSchemaURI` and integration with embedded FS:

- Valid URI: `cloud-resource-schema://AwsVpc`
- Malformed URI
- Wrong scheme
- Missing kind (empty host)
- `loadRegistry` returns consistent data
- `loadProviderSchema` for a known kind returns valid JSON
- `loadProviderSchema` for unknown kind returns error
- `buildKindCatalog` produces valid JSON with expected structure

### 2e. `internal/domains/cloudresource/apply_test.go`

Test `buildCloudResource` (pure proto assembly, no gRPC):

- Valid inputs produce correct CloudResource proto fields
- Invalid kind string propagates error from `resolveKind`
- Missing metadata propagates error from `extractMetadata`
- `describeResource` output with and without metadata

### 2f. `internal/domains/rpcerr_test.go`

Test `classifyCode`:

- Each gRPC code (NotFound, PermissionDenied, Unauthenticated, Unavailable, DeadlineExceeded, InvalidArgument) maps to expected message
- Unknown code falls through to default

### 2g. `internal/parse/helpers_test.go`

Test `ValidateHeader`, `ExtractSpecMap`, `RebuildCloudObject`:

- Correct api_version and kind pass
- Wrong api_version returns error
- Wrong kind returns error
- Missing spec returns error
- Non-object spec returns error
- `RebuildCloudObject` preserves top-level fields and replaces spec

### 2h. `internal/config/config_test.go`

Test `Validate` and `ParseLogLevel`:

- Valid config passes
- Invalid transport rejected
- Empty server address rejected
- Invalid log format rejected
- Missing API key with stdio transport rejected
- Missing API key with http transport passes
- Each log level string parses correctly
- Invalid log level rejected

---

## Stage 3: README and Documentation Rewrite

### 3a. Delete stale docs

Remove all 5 files under `docs/` -- they document tools and domains that no longer exist:

- `docs/service-hub-tools.md` (Service Hub domain deleted)
- `docs/installation.md` (references old binary, old tools)
- `docs/http-transport.md` (likely partially valid but needs full rewrite)
- `docs/development.md` (references old dependencies, old Go version info)
- `docs/configuration.md` (partially valid but needs rewrite)

### 3b. Rewrite README.md

The new README should cover:

- Project overview with the new 3-tool architecture
- Available tools: `apply_cloud_resource`, `get_cloud_resource`, `delete_cloud_resource`
- MCP resources: `cloud-resource-kinds://catalog`, `cloud-resource-schema://{kind}`
- 3-step workflow: catalog, schema, apply
- Integration guides (Cursor, Claude Desktop, LangGraph) -- updated for new tool names
- Configuration (env vars -- same as current, these are correct)
- Codegen pipeline overview (`make codegen`)
- Development guide (build, test, format, lint)

### 3c. Recreate docs/ with focused guides

- `docs/configuration.md` -- rewrite with current env vars (mostly correct, remove stale tool references)
- `docs/development.md` -- rewrite for current Go version, codegen pipeline, test commands

---

## What Is NOT In Scope

- Mock gRPC integration tests (thin proxy, low ROI)
- Codegen determinism tests (validated by `go build`)
- End-to-end MCP protocol tests (depends on mock backend)
- CI/CD pipeline changes (separate concern)

