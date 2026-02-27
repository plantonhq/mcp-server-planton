# Implementation Plans

Plans created for the InfraHub MCP Tools Expansion project. Each plan documents a specific implementation effort.

## Plan Registry

| Plan | Status | Created | Completed | Description |
|------|--------|---------|-----------|-------------|
| `phase-0-gen-restructure.md` | ✅ Completed | 2026-02-27 | 2026-02-27 | Move gen/cloudresource/ to gen/infrahub/cloudresource/ |
| `phase-1a-infrachart-tools.md` | ✅ Completed | 2026-02-27 | 2026-02-27 | Add list, get, build tools for InfraChart domain |
| `phase-1b-infraproject-tools.md` | ✅ Completed | 2026-02-27 | 2026-02-27 | Add search, get, apply, delete, slug, undeploy tools for InfraProject domain |
| `phase-1c-infrapipeline-tools.md` | ✅ Completed | 2026-02-28 | 2026-02-28 | Add list, get, latest, run, cancel, gate resolution tools for InfraPipeline domain |
| `phase-2a-graph-tools.md` | ✅ Completed | 2026-02-28 | 2026-02-28 | Add 7 graph/dependency intelligence tools (org graph, env graph, service graph, resource graph, dependencies, dependents, impact analysis) |
| `phase-2b-configmanager-tools.md` | ✅ Completed | 2026-02-28 | 2026-02-28 | Add 11 configmanager tools (5 variable, 4 secret, 2 secret version) — second non-infrahub bounded context |

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

### phase-1b-infraproject-tools.md
- **Objective**: Add 6 MCP tools for InfraProject lifecycle management (search, get, apply, delete, slug check, undeploy)
- **Phases**: Single phase
- **Current Phase**: Complete
- **Key Decisions**: Simpler identification (org+slug vs 4-field), full JSON passthrough for apply, Purge excluded, Search over Find

### phase-1c-infrapipeline-tools.md
- **Objective**: Add 7 MCP tools for InfraPipeline observability, execution control, and manual gate resolution
- **Phases**: Single phase
- **Current Phase**: Complete
- **Key Decisions**: Unified run tool (chart + git), manual gate tools included, streaming RPCs excluded, user-friendly approve/reject decisions

### phase-2a-graph-tools.md
- **Objective**: Add 7 MCP tools for dependency intelligence and impact analysis (first non-infrahub bounded context)
- **Phases**: Single phase
- **Current Phase**: Complete
- **Key Decisions**: Expanded from planned 4 to 7 tools (added environment graph, service graph, dependents), new `internal/domains/graph/` bounded context, shared DependencyInput struct

### phase-2b-configmanager-tools.md
- **Objective**: Add 11 MCP tools for configuration lifecycle management (variables, secrets, secret versions)
- **Phases**: Single phase
- **Current Phase**: Complete
- **Key Decisions**: Expanded from planned 5 to 11 tools, write-only secret values (DD-2), explicit params for apply (DD-3), delete_secret with destructive warning (DD-5)

---

*Last updated: 2026-02-28*
