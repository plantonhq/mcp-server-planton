---
name: Domain bounded-context refactor
overview: Restructure `internal/domains/` from a flat layout to one level of bounded-context grouping that mirrors the Planton API's top-level domain taxonomy (`infrahub`, `resourcemanager`), aligning ubiquitous language and preparing the codebase for scalable expansion.
todos:
  - id: move-infrahub
    content: git mv cloudresource, stackjob, preset into internal/domains/infrahub/
    status: completed
  - id: move-resourcemanager
    content: git mv environment, organization into internal/domains/resourcemanager/
    status: completed
  - id: create-doc-files
    content: Create doc.go files for infrahub and resourcemanager packages
    status: completed
  - id: update-server-imports
    content: Update 5 import paths in internal/server/server.go
    status: completed
  - id: update-docs
    content: Update paths in docs/development.md (test table + project structure)
    status: completed
  - id: verify-build-test
    content: Run go build, go vet, go test to verify correctness
    status: completed
isProject: false
---

# Domain Bounded-Context Refactor

## Context

The MCP server's 5 domain packages currently live flat under `internal/domains/`. The canonical Planton API organizes these same concepts into bounded contexts: `infrahub` (cloudresource, stackjob, preset) and `resourcemanager` (environment, organization). This refactor introduces one intermediate directory level to match.

## Blast Radius Analysis

**Files that MOVE (35 files across 5 directories):**

- `internal/domains/cloudresource/` (22 files) --> `internal/domains/infrahub/cloudresource/`
- `internal/domains/stackjob/` (6 files) --> `internal/domains/infrahub/stackjob/`
- `internal/domains/preset/` (3 files) --> `internal/domains/infrahub/preset/`
- `internal/domains/environment/` (2 files) --> `internal/domains/resourcemanager/environment/`
- `internal/domains/organization/` (2 files) --> `internal/domains/resourcemanager/organization/`

**Files that get EDITED (import path updates):**

- [internal/server/server.go](internal/server/server.go) -- 5 import paths change
- [docs/development.md](docs/development.md) -- test table paths and project structure diagram update

**Files that get CREATED:**

- `internal/domains/infrahub/doc.go` -- package doc for the infrahub bounded context
- `internal/domains/resourcemanager/doc.go` -- package doc for the resourcemanager bounded context

**Files that DO NOT change:**

- Root shared utilities (`conn.go`, `marshal.go`, `rpcerr.go`, `toolresult.go`, `kind.go`, `kind_test.go`, `rpcerr_test.go`) stay at `internal/domains/` -- their import path is unchanged
- All 26 domain implementation files import `internal/domains` (the root package), NOT sibling domains -- their internal import statements need zero edits
- Generated code (`gen/cloudresource/`) has no dependency on `internal/domains/`
- `go.mod`, `Makefile`, CI workflows, Dockerfile -- all use `./...` patterns, unaffected
- `pkg/mcpserver/` imports `internal/server`, not domain packages directly

## Execution Strategy

Use `git mv` for all directory moves to preserve rename tracking in git history. Perform the refactor in this strict order:

1. **Move directories** -- `git mv` the 5 domain packages into their bounded-context parents
2. **Create doc.go files** -- thin package documentation for `infrahub` and `resourcemanager`
3. **Update server.go** -- rewrite the 5 import paths
4. **Update documentation** -- fix paths in `docs/development.md`
5. **Verify** -- `go build ./...`, `go vet ./...`, `go test ./...`

## Import Path Changes in server.go

```
-- BEFORE --                                          -- AFTER --
internal/domains/cloudresource          -->  internal/domains/infrahub/cloudresource
internal/domains/stackjob               -->  internal/domains/infrahub/stackjob
internal/domains/preset                 -->  internal/domains/infrahub/preset
internal/domains/environment            -->  internal/domains/resourcemanager/environment
internal/domains/organization           -->  internal/domains/resourcemanager/organization
```

## Risk Assessment

- **Low risk:** No domain package imports another domain package (verified: zero cross-domain imports). The only consumer of domain package symbols is `server.go`.
- **No functional change:** Tool names, handler signatures, gRPC calls, MCP resources -- all unchanged. This is purely a code organization refactor.
- **Git history:** `git mv` ensures `git log --follow` works for all moved files.

## Pause Points

I will pause and consult you if:

- Any unexpected import dependency surfaces during the move
- `go build` or `go test` reveals a dependency I missed
- Anything in the CI or build tooling references domain paths beyond `./...`

