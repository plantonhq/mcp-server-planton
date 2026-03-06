# Rename "Planton Cloud" to "Planton" Across Codebase

**Date**: March 7, 2026

## Summary

Renamed all "Planton Cloud" references to "Planton" across the entire MCP server codebase, aligning with the updated product branding. Also renamed the `PLANTON_CLOUD_ENVIRONMENT` environment variable to `PLANTON_ENVIRONMENT`. This affects Go source code, documentation, schemas, Kustomize manifests, changelogs, and project plans — 70 files in total.

## Problem Statement

The product name changed from "Planton Cloud" to "Planton", but references to the old name persisted throughout the codebase — in tool descriptions, package documentation, README, configuration docs, Kubernetes manifests, and historical records.

### Pain Points

- Tool descriptions shown to AI assistants still said "Planton Cloud"
- The `PLANTON_CLOUD_ENVIRONMENT` environment variable carried the old branding
- README, configuration docs, and contribution guidelines referenced "Planton Cloud"
- Historical changelogs and project plans used inconsistent naming

## Solution

Performed a systematic, category-by-category rename across the repository:

1. **Go source files** — Updated tool descriptions, package doc comments, prompt text, and error messages
2. **Documentation** — Updated README, `docs/configuration.md`, `docs/development.md`, `docs/tools.md`, `CONTRIBUTING.md`
3. **Configuration** — Renamed `PLANTON_CLOUD_ENVIRONMENT` → `PLANTON_ENVIRONMENT` in `internal/config/config.go` (the runtime `os.Getenv` call), `pkg/mcpserver/config.go`, and `_kustomize/base/service.yaml`
4. **Schemas** — Updated `schemas/apiresourcekinds/catalog.json` description
5. **Historical records** — Updated all `_changelog/`, `_projects/`, and `.cursor/plans/` files

### What Was Preserved

- `org: planton-cloud` and `org_id="planton-cloud"` values (actual backend identifiers)
- `planton-cloud-resource-*` GCP label prefix
- `gen/` auto-generated protobuf code (regenerated separately)

## Implementation Details

### Environment Variable Rename

The core runtime change in `internal/config/config.go`:

```go
// Before
switch strings.ToLower(os.Getenv("PLANTON_CLOUD_ENVIRONMENT")) {

// After
switch strings.ToLower(os.Getenv("PLANTON_ENVIRONMENT")) {
```

### Go Source Files Updated (12 files)

| File | Change |
|------|--------|
| `cmd/mcp-server-planton/main.go` | Server description |
| `internal/config/config.go` | Env var name, comments, error messages |
| `internal/domains/prompts/debug_deployment.go` | Prompt text |
| `internal/domains/discovery/resources.go` | Tool description |
| `internal/domains/discovery/doc.go` | Package doc |
| `internal/domains/connect/credential/doc.go` | Package doc |
| `internal/domains/infrahub/cloudresource/tools.go` | Tool descriptions |
| `internal/domains/resourcemanager/organization/tools.go` | Tool descriptions |
| `internal/domains/resourcemanager/environment/tools.go` | Tool descriptions |
| `internal/domains/servicehub/service/tools.go` | Tool descriptions |
| `internal/domains/servicehub/service/disconnect.go` | Tool description |
| `pkg/mcpserver/config.go` | Doc comment |

### Documentation Files Updated (5 files)

`README.md`, `CONTRIBUTING.md`, `docs/configuration.md`, `docs/development.md`, `docs/tools.md`

### Historical Records Updated (53 files)

- 22 changelog files across `_changelog/2025-11/`, `2025-12/`, `2026-02/`, `2026-03/`
- 13 project plan files in `_projects/`
- 9 cursor plan files in `.cursor/plans/`
- `schemas/apiresourcekinds/catalog.json`
- `_kustomize/base/service.yaml`

## Benefits

- Consistent branding across all user-facing tool descriptions and documentation
- Simpler, cleaner environment variable name (`PLANTON_ENVIRONMENT`)
- Historical records updated to reduce confusion for future contributors
- AI assistants using this MCP server now see the correct product name

## Impact

- **MCP tool consumers**: Tool descriptions now show "Planton" instead of "Planton Cloud"
- **Operators**: Must use `PLANTON_ENVIRONMENT` instead of `PLANTON_CLOUD_ENVIRONMENT` in deployments
- **Contributors**: All documentation reflects current branding

---

**Status**: ✅ Production Ready
**Files Changed**: 70
