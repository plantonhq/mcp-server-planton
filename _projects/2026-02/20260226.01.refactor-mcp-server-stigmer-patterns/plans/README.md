# Implementation Plans

Plans created for the mcp-server-planton Stigmer patterns refactoring project.

## Plan Registry

| Plan | Status | Created | Completed | Description |
|------|--------|---------|-----------|-------------|
| `phase-1-foundation.plan.md` | ✅ Completed | 2026-02-26 | 2026-02-26 | Clean slate + shared utilities — SDK migration, auth, config, gRPC, domains, server, entry point |
| `proto2schema-codegen-tool.plan.md` | ✅ Completed | 2026-02-26 | 2026-02-26 | Phase 2 Stage 1 — proto2schema tool parsing 362 OpenMCF provider protos to JSON schemas |
| `schema2go-generator.plan.md` | ✅ Completed | 2026-02-26 | 2026-02-26 | Phase 2 Stage 2 — schema2go generator producing typed Go input structs for 362 providers |
| `phase-3-apply-cloud-resource.plan.md` | ✅ Completed | 2026-02-26 | 2026-02-26 | Phase 3 — apply_cloud_resource tool, MCP resource templates, schema promotion |

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

---

*Last updated: 2026-02-26*
