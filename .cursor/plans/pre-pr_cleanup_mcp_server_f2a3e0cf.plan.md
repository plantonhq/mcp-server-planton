---
name: Pre-PR Cleanup MCP Server
overview: "Clean up mcp-server-planton before creating the PR: fix the org name mismatch (plantoncloud -> plantonhq), remove stale .github files, harden Dockerfile/goreleaser/Makefile to match stigmer patterns, and rewrite README.md following the document writer role."
todos:
  - id: fix-org-name
    content: Fix plantoncloud -> plantonhq in go.mod, hand-written Go files, codegen generator default, then re-run make codegen-types to regenerate gen/
    status: completed
  - id: github-cleanup
    content: Delete ISSUE_TEMPLATE/, pull_request_template.md, CODEOWNERS. Update Go version in ci.yml and release.yml.
    status: completed
  - id: dockerfile
    content: "Rewrite Dockerfile: alpine runtime, non-root user, Go 1.25, ldflags, remove unnecessary copies"
    status: completed
  - id: goreleaser
    content: "Clean up .goreleaser.yaml: fix org, v2 syntax, simplify templates, remove trailing whitespace"
    status: completed
  - id: makefile
    content: "Align Makefile with stigmer: version ldflags, vet excluding gen/, fix GHCR image, add race/timeout to tests"
    status: completed
  - id: readme
    content: "Rewrite README.md following document writer role: architecture diagram, key concepts, consolidated client config, proper tool docs, Docker networking"
    status: completed
  - id: docs-fixup
    content: Fix plantoncloud -> plantonhq in docs/configuration.md and docs/development.md
    status: completed
  - id: verify
    content: Run go mod tidy, go build, go vet, go test to verify everything compiles and passes
    status: completed
isProject: false
---

# Pre-PR Cleanup for mcp-server-planton

## Scope

Files to **leave untouched**: `_changelog/`, `_kustomize/`, `_projects/`, `.cursor/`, `SECURITY.md`, `.gitignore`

---

## 1. Fix GitHub org name mismatch

The module path says `plantoncloud` but the repo lives at `plantonhq`. Three different org names exist across the project (`plantonhq`, `plantoncloud`, `plantoncloud-inc`). All must become `plantonhq`.

**Hand-written files** (mechanical `plantoncloud` -> `plantonhq` replacement):

- [go.mod](go.mod) -- module path on line 1
- [tools/codegen/generator/main.go](tools/codegen/generator/main.go) -- `--module` flag default on line 110
- [internal/server/server.go](internal/server/server.go) -- 4 import references
- [internal/domains/cloudresource/tools.go](internal/domains/cloudresource/tools.go) -- 2 imports
- [internal/domains/cloudresource/delete.go](internal/domains/cloudresource/delete.go) -- 1 import
- [internal/domains/cloudresource/get.go](internal/domains/cloudresource/get.go) -- 1 import
- [internal/domains/cloudresource/identifier.go](internal/domains/cloudresource/identifier.go) -- 1 import
- [internal/domains/cloudresource/resources.go](internal/domains/cloudresource/resources.go) -- 1 import
- [internal/domains/cloudresource/schema.go](internal/domains/cloudresource/schema.go) -- 1 import
- [docs/configuration.md](docs/configuration.md) -- 1 reference

**Generated files** (367 files in `gen/cloudresource/`):

- Re-run `make codegen-types` after updating the generator default. This regenerates all files with the correct `plantonhq` import path.
- No need for manual edits in `gen/`.

**Verification**: `go mod tidy && go build ./... && go vet ./... && go test ./...`

---

## 2. .github cleanup

**Delete** (premature for a project at this stage; CODEOWNERS references stale paths like `/internal/mcp/`):

- `.github/ISSUE_TEMPLATE/` (3 files: `bug_report.yml`, `feature_request.yml`, `config.yml`)
- `.github/pull_request_template.md`
- `.github/CODEOWNERS`

**Keep and update**:

- [.github/workflows/ci.yml](.github/workflows/ci.yml) -- update `go-version: '1.24'` to `'1.25'` (matches go.mod)
- [.github/workflows/release.yml](.github/workflows/release.yml) -- update `go-version: '1.24'` to `'1.25'`; strip 30+ trailing blank lines

---

## 3. Dockerfile improvements

Align with [stigmer/mcp-server/Dockerfile](../stigmer/stigmer/mcp-server/Dockerfile):

- **Runtime base**: `debian:bookworm-slim` -> `alpine:3.19` (smaller image, matches stigmer)
- **Non-root user**: Add `addgroup`/`adduser`, `USER planton` (security best practice; stigmer does this)
- **Go version**: `golang:1.24-alpine` -> `golang:1.25-alpine` (match go.mod)
- **Build ldflags**: Embed version via `-ldflags="-s -w"` (strip debug symbols like stigmer)
- **Remove**: Verbose comments, `COPY README.md LICENSE` (unnecessary at runtime), `GOTOOLCHAIN=auto` env var
- **Keep**: Multi-stage build, static binary, health check, port 8080 expose

---

## 4. .goreleaser.yaml cleanup

Align with [stigmer/mcp-server/.goreleaser.yaml](../stigmer/stigmer/mcp-server/.goreleaser.yaml):

- **Remove** `before.hooks` block (go mod tidy in goreleaser is unnecessary)
- **Fix** `release.github.owner`: `plantoncloud` -> `plantonhq`
- **Fix** archive format: `format: tar.gz` -> `formats: [tar.gz]` (goreleaser v2 syntax)
- **Simplify** `name_template`: use `{{ .ProjectName }}_{{ .Os }}_{{ .Arch }}` (matches stigmer; current template has odd `_Version_` in the middle)
- **Remove** changelog groups (follow stigmer's simpler filter-only approach)
- **Remove** `name_template` from release section
- **Add** `ldflags: -s -w` (strip debug symbols)
- **Strip** 30+ trailing blank lines

---

## 5. Makefile cleanup

Align with [stigmer/mcp-server/Makefile](../stigmer/stigmer/mcp-server/Makefile):

- **Fix** `GHCR_IMAGE`: `ghcr.io/plantoncloud-inc/mcp-server-planton` -> `ghcr.io/plantonhq/mcp-server-planton`
- **Add** `VERSION` variable from git describe (like stigmer)
- **Add** `-ldflags` to build target (embed version)
- **Add** `vet` target that excludes `gen/` (like stigmer -- gen/ has false struct tag positives)
- **Add** `-race -timeout 30s` to test target (like stigmer)
- **Add** `tidy` target
- **Simplify**: Remove excessive `@echo` chatter, follow stigmer's minimal style
- **Keep**: `codegen-schemas`, `codegen-types`, `codegen`, `release`, `docker-build`, `docker-run`, `clean`, `help`

---

## 6. README.md rewrite

Following the **Lead Technical Document Writer** role from `_roles/002_document_writer`:

**Audience**: MCP client users who want to connect their AI IDE to Planton Cloud.

**Key improvements** (using stigmer README as the structural reference):

- **Opening paragraph**: One clear sentence explaining what the server is and what it connects. Add the ASCII architecture diagram like stigmer: `AI IDE <-> MCP protocol <-> mcp-server-planton <-> gRPC <-> Planton Cloud`.
- **Key Concepts**: Define `cloud resource`, `kind`, `org`, `env`, `slug`, `apply`, `api_version` for users who are new to Planton. Stigmer does this and it eliminates assumed knowledge.
- **Installation**: Reorder to: Prerequisites (Planton Cloud account + API key) -> Go install (one-liner) -> Pre-built binary -> Docker. Currently Quick Start jumps straight into MCP client config before the binary is even installed.
- **MCP Client Configuration**: Consolidate the 4 duplicate JSON blocks (Cursor STDIO, Cursor HTTP, Claude Desktop, LangGraph) into one canonical config with a per-client location table (like stigmer). The Cursor STDIO and Claude Desktop blocks are identical -- redundancy hurts scanability.
- **Configuration Reference**: Keep the env var table.
- **Tools**: Document each tool properly with parameter tables (like stigmer does for its tools). Current README only has a 3-row summary table.
- **MCP Resources**: Document the two resources with URI patterns and what they return.
- **HTTP Mode**: Dedicated section with Docker networking note (currently missing -- stigmer has this).
- **Security**: Keep the existing bullet points.
- **Development**: Keep brief, link to `docs/development.md`.
- **Remove**: Architecture file tree (belongs in `docs/development.md`, not user README), CodeQL badge (no such workflow exists), Go Report Card badge (not verified)
- **Fix**: All URLs from `plantoncloud` to `plantonhq`

---

## 7. docs/configuration.md and docs/development.md

- Fix any `plantoncloud` references to `plantonhq`
- No structural changes

---

## Verification

After all changes:

```
go mod tidy
go build ./...
go vet ./...
go test ./...
```

