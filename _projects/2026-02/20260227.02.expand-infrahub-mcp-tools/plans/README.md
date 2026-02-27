# Implementation Plans

Plans created for the InfraHub MCP Tools Expansion project. Each plan documents a specific implementation effort.

## Plan Registry

| Plan | Status | Created | Completed | Description |
|------|--------|---------|-----------|-------------|
| `phase-0-gen-restructure.md` | ✅ Completed | 2026-02-27 | 2026-02-27 | Move gen/cloudresource/ to gen/infrahub/cloudresource/ |
| `phase-1a-infrachart-tools.md` | ✅ Completed | 2026-02-27 | 2026-02-27 | Add list, get, build tools for InfraChart domain |

### Status Legend

- 🟡 **Pending**: Plan created, not yet started
- 🔵 **In Progress**: Currently being executed
- ✅ **Completed**: All phases/tasks finished
- ❌ **Abandoned**: Plan was cancelled or superseded

## Plan Details

### phase-0-gen-restructure.md
- **Objective**: Restructure generated code directory to mirror domain hierarchy
- **Phases**: Single phase (prerequisite)
- **Current Phase**: Complete
- **ADR**: AD-02 (Restructure Generated Code Under Domain Directories)

### phase-1a-infrachart-tools.md
- **Objective**: Add 3 MCP tools for InfraChart discovery and rendering
- **Phases**: Single phase
- **Current Phase**: Complete
- **Key Decisions**: `list_` naming over `search_`, simplified build input with chart_id + params map

---

*Last updated: 2026-02-27*
