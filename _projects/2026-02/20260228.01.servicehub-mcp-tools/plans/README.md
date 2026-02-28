# Implementation Plans

Plans created for the ServiceHub MCP Tools project. Each plan documents a specific implementation effort.

## Plan Registry

| Plan | Status | Created | Completed | Description |
|------|--------|---------|-----------|-------------|
| `tier-1-service-mcp-tools.md` | Completed | 2026-02-28 | 2026-02-28 | 7 MCP tools for the Service entity (search, get, apply, delete, disconnect, webhook, branches) |
| `tier-2-pipeline-mcp-tools.md` | Completed | 2026-02-28 | 2026-02-28 | 9 MCP tools for the Pipeline entity (list, get, get_last, run, rerun, cancel, gate, list_files, update_file) |

### Status Legend

- Pending: Plan created, not yet started
- In Progress: Currently being executed
- Completed: All phases/tasks finished
- Abandoned: Plan was cancelled or superseded

## Plan Details

### tier-1-service-mcp-tools.md
- **Objective**: Implement 7 MCP tools for the ServiceHub Service entity
- **Tools**: search_services, get_service, apply_service, delete_service, disconnect_service_git_repo, configure_service_webhook, list_service_branches
- **Key Decision**: Used generic `ApiResourceSearchQueryController.searchByKind` instead of domain-specific search RPC (none exists for Service)

### tier-2-pipeline-mcp-tools.md
- **Objective**: Implement 9 MCP tools for the ServiceHub Pipeline entity
- **Tools**: list_pipelines, get_pipeline, get_last_pipeline, run_pipeline, rerun_pipeline, cancel_pipeline, resolve_pipeline_gate, list_pipeline_files, update_pipeline_file
- **Key Decisions**: DD-T2-1 (branch required for run), DD-T2-2 (bytes-to-string decode for pipeline files), DD-T2-3 (single gate tool)

---

*Last updated: 2026-02-28*
