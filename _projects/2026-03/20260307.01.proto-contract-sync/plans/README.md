# Implementation Plans

Plans created for this project. Each plan documents a specific implementation effort.

## Plan Registry

| Plan | Status | Created | Completed | Description |
|------|--------|---------|-----------|-------------|
| `T01_0_plan.md` | ✅ Completed | 2026-03-07 | 2026-03-07 | Full migration plan: credential→connection rename, proto import sync, redaction removal |
| `T02_phase2_connect_tools.plan.md` | 🔵 In Progress | 2026-03-08 | - | Phase 2: Enrich connect tools — T02.1–T02.4 completed, T02.5 pending decision |
| `T03_phase3_new_resources.plan.md` | ✅ Completed | 2026-03-08 | 2026-03-08 | Phase 3: New resources — secretbackend, variablegroup, serviceaccount, iacprovisionermapping (23 tools) |

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

### T03_phase3_new_resources.plan.md
- **Objective**: Implement MCP tool packages for 4 new proto resources totaling 23 tools, with clear architectural decisions on apply style, security boundaries, and API surface
- **Tasks**: T03.4, T03.1, T03.3, T03.2 (all completed)
- **Current**: All phases complete
- **Design Decisions**: Envelope Apply for complex specs, explicit params for entry operations, sensitive field redaction for SecretBackend, read-modify-write for ServiceAccount update, scope enum resolver for VariableGroup

---

*Last updated: 2026-03-08*
