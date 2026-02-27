---
name: Phase 0 gen restructure
overview: Move generated code from flat `gen/cloudresource/` to domain-scoped `gen/infrahub/cloudresource/` by editing 3 hand-written files, re-running the code generator, and verifying build + tests.
todos:
  - id: edit-generator-default
    content: Update default --output-dir in tools/codegen/generator/main.go from gen/cloudresource to gen/infrahub/cloudresource (+ doc comment)
    status: completed
  - id: edit-makefile
    content: Update codegen-types target in Makefile to use gen/infrahub/cloudresource
    status: completed
  - id: edit-consumer-import
    content: Update import path in internal/domains/infrahub/cloudresource/tools.go
    status: completed
  - id: update-doc-comments
    content: Update path references in codegen.go and registry.go comments
    status: completed
  - id: regenerate
    content: Delete old gen/cloudresource/, run make codegen-types to regenerate at new path
    status: completed
  - id: verify-build-tests
    content: Run go build ./... and go test ./... to verify nothing broke
    status: completed
isProject: false
---

# Phase 0: Restructure Generated Code Under Domain Directory

## Scope Analysis

This is a surgical, low-risk restructuring. The blast radius is narrow:

- **367 generated files** move from `gen/cloudresource/` to `gen/infrahub/cloudresource/`
- **Only 1 hand-written consumer** imports from `gen/cloudresource` ([tools.go](internal/domains/infrahub/cloudresource/tools.go) line 27)
- The other importer is `registry_gen.go` itself (generated, will be regenerated)
- The code generator already accepts `--output-dir` as a flag, so the change is config-level

## Strategy

Delete the old directory + regenerate into the new path (rather than `git mv`). Rationale: every generated file has a timestamp header that changes on regeneration, so rename-tracking offers no meaningful git history value for machine-generated code.

## Files to Edit (3 hand-written + 2 doc comments)

### 1. Code generator default path

[tools/codegen/generator/main.go](tools/codegen/generator/main.go) line 109:

```go
// Before:
outputDir := flag.String("output-dir", "gen/cloudresource", ...)
// After:
outputDir := flag.String("output-dir", "gen/infrahub/cloudresource", ...)
```

Also update the package doc comment (line 7) and usage example.

### 2. Makefile codegen target

[Makefile](Makefile) lines 92-94:

```makefile
# Before:
codegen-types:
	rm -rf gen/cloudresource/
	go run ./tools/codegen/generator/ --schemas-dir=schemas --output-dir=gen/cloudresource
# After:
codegen-types:
	rm -rf gen/infrahub/cloudresource/
	go run ./tools/codegen/generator/ --schemas-dir=schemas --output-dir=gen/infrahub/cloudresource
```

### 3. Consumer import path

[internal/domains/infrahub/cloudresource/tools.go](internal/domains/infrahub/cloudresource/tools.go) line 27:

```go
// Before:
"github.com/plantonhq/mcp-server-planton/gen/cloudresource"
// After:
"github.com/plantonhq/mcp-server-planton/gen/infrahub/cloudresource"
```

The alias `cloudresource` will still resolve correctly since Go uses the last path segment.

### 4. Documentation comments (non-functional, but accurate)

- [tools/codegen/generator/codegen.go](tools/codegen/generator/codegen.go) line 20: comment says `gen/cloudresource/aws/`
- [tools/codegen/generator/registry.go](tools/codegen/generator/registry.go) line 13: comment says `gen/cloudresource/registry_gen.go`

## Execution Sequence

1. Edit the 3 source files + 2 comments listed above
2. Delete old generated directory: `rm -rf gen/cloudresource/`
3. Create parent directory: `mkdir -p gen/infrahub/`
4. Regenerate: `make codegen-types`
5. Verify build: `go build ./...`
6. Verify tests: `go test ./...`
7. Verify vet still works: `make vet` (the `grep -v '/gen/'` filter still matches `gen/infrahub/...`)

## What Does NOT Change

- **Package names**: All generated packages keep their current names (`cloudresource`, `aws`, `gcp`, etc.) since Go derives package names from the last path segment
- `**vet` target**: The `grep -v '/gen/'` filter in the Makefile still excludes `gen/infrahub/...`
- **Consumer code logic**: `tools.go` still calls `cloudresource.GetParser(kindStr)` with the same alias
- **No hand-written domain code moves** -- only the generated `gen/` tree relocates

## Risk Assessment

- **Risk**: Very low. Pure path/import restructuring with no logic changes.
- **Rollback**: Revert the 3 source edits + re-run `make codegen-types` to regenerate at old path.
- **Key verification**: `go build ./...` is the primary gate. If imports resolve, it works.

