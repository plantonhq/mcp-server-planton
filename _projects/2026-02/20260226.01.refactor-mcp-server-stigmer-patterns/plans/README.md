# Implementation Plans

Plans created for the mcp-server-planton Stigmer patterns refactoring project.

## Plan Registry

| Plan | Status | Created | Completed | Description |
|------|--------|---------|-----------|-------------|
| `phase-1-foundation.plan.md` | ✅ Completed | 2026-02-26 | 2026-02-26 | Clean slate + shared utilities — SDK migration, auth, config, gRPC, domains, server, entry point |
| `proto2schema-codegen-tool.plan.md` | ✅ Completed | 2026-02-26 | 2026-02-26 | Phase 2 Stage 1 — proto2schema tool parsing 362 OpenMCF provider protos to JSON schemas |
| `schema2go-generator.plan.md` | ✅ Completed | 2026-02-26 | 2026-02-26 | Phase 2 Stage 2 — schema2go generator producing typed Go input structs for 362 providers |
| `phase-3-apply-cloud-resource.plan.md` | ✅ Completed | 2026-02-26 | 2026-02-26 | Phase 3 — apply_cloud_resource tool, MCP resource templates, schema promotion |
| `cloud-resource-kinds-catalog.plan.md` | ✅ Completed | 2026-02-26 | 2026-02-26 | Static MCP resource for kind discovery — agents discover all 362 kinds before using schema template |
| `phase-5-hardening.plan.md` | ✅ Completed | 2026-02-27 | 2026-02-27 | Phase 5 — Dead code cleanup, 50 unit tests for pure domain logic, README + docs rewrite |

### Status Legend

- 🟡 **Pending**: Plan created, not yet started
- 🔵 **In Progress**: Currently being executed
- ✅ **Completed**: All phases/tasks finished
- ❌ **Abandoned**: Plan was cancelled or superseded

## Plan Details

### phase-1-foundation.plan.md
- **Objective**: Rebuild mcp-server-planton foundation following Stigmer MCP server patterns
- **Phases**: 1 (single-phase plan)
- **Current Phase**: Complete
- **Key outcome**: 12 new Go files forming the zero-tool skeleton; MCP SDK migrated to official go-sdk

### proto2schema-codegen-tool.plan.md
- **Objective**: Build proto2schema codegen tool for OpenMCF provider protos → JSON schemas
- **Phases**: 1 (single-phase plan, 9 tasks)
- **Current Phase**: Complete
- **Key outcome**: 362 provider schemas generated across 17 cloud platforms; shared metadata schema; provider registry; `make codegen-schemas` target

### schema2go-generator.plan.md
- **Objective**: Build schema2go codegen generator transforming JSON schemas into typed Go input structs
- **Phases**: 1 (single-phase plan, 10 tasks)
- **Current Phase**: Complete
- **Key outcome**: 367 generated Go files (362 providers, 5 shared types, 1 registry) across 17 cloud packages; `make codegen-types` and `make codegen` targets

### phase-3-apply-cloud-resource.plan.md
- **Objective**: Implement apply_cloud_resource MCP tool with typed validation and MCP resource templates for schema discovery
- **Phases**: 1 (single-phase plan, 6 tasks: prereq + 5 stages)
- **Current Phase**: Complete
- **Key outcome**: First working MCP tool registered; `cloud-resource-schema://{kind}` resource template serving 362 provider schemas; schemas promoted to top-level `schemas/` package with `go:embed`

### cloud-resource-kinds-catalog.plan.md
- **Objective**: Add static MCP resource for cloud resource kind discovery, enabling agents to enumerate all valid kinds
- **Phases**: 1 (single-phase plan, 5 tasks)
- **Current Phase**: Complete
- **Key outcome**: `cloud-resource-kinds://catalog` static resource serving 362 kinds grouped by 17 cloud providers; tool descriptions updated with 3-step discovery workflow

### phase-5-hardening.plan.md
- **Objective**: Clean up dead code, write comprehensive unit tests for all pure domain logic, rewrite README and docs
- **Phases**: 3 stages (dead code cleanup, unit tests, documentation)
- **Current Phase**: Complete
- **Key outcome**: 8 test files with 50 test cases; dead code removed; README rewritten; stale docs replaced

---

*Last updated: 2026-02-27*
