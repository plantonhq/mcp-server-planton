# Implementation Plans

Plans created for this project. Each plan documents a specific implementation effort.

## Plan Registry

| Plan | Status | Created | Completed | Description |
|------|--------|---------|-----------|-------------|
| `T01_0_plan.md` | ✅ Completed | 2026-03-07 | 2026-03-07 | Full migration plan: credential→connection rename, proto import sync, redaction removal |
| `T02_phase2_connect_tools.plan.md` | 🔵 In Progress | 2026-03-08 | - | Phase 2: Enrich connect tools — T02.1–T02.4 completed, T02.5 pending decision |

### Status Legend

- 🟡 **Pending**: Plan created, not yet started
- 🔵 **In Progress**: Currently being executed
- ✅ **Completed**: All phases/tasks finished
- ❌ **Abandoned**: Plan was cancelled or superseded

## Plan Details

### T01_0_plan.md
- **Objective**: Migrate MCP server tools to match restructured protobuf contracts — fix the broken build, rename credential→connection, sync all import paths
- **Phases**: 5 total (Phase 1 completed, Phases 2–5 pending future sessions)
- **Current Phase**: Phase 1 complete; Phase 2 is next
- **Design Decisions**: Redaction removal (secret slugs are not sensitive), tool name rename to match protos

### T02_phase2_connect_tools.plan.md
- **Objective**: Wire unwired gRPC methods as MCP tools across 4 connect sub-packages, fix Resolve bug, design provider-specific controllers
- **Tasks**: T02.1–T02.5 (T02.1–T02.4 completed, T02.5 pending user decision)
- **Current**: T02.5 deferred — needs decision on OAuth callback scope
- **Design Decisions**: Separate tools for org/env-level operations, enhanced delete with semantic key support, Find methods skipped (operator-only)

---

*Last updated: 2026-03-08*
