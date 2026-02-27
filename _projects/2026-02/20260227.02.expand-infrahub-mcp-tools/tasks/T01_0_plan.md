# Task T01: Expand InfraHub MCP Tools — Master Plan

**Created**: 2026-02-27
**Status**: Planning — Pending Review

## Context

The MCP server currently has **18 tools** across 5 subdomains:
- `cloudresource` (11 tools) — full CRUD lifecycle, locks, references, env vars
- `stackjob` (3 tools) — read-only observability
- `preset` (2 tools) — discovery and retrieval
- `organization` (1 tool) — list only
- `environment` (1 tool) — list only

This plan expands the server to **~55+ tools** across the full Planton Cloud product surface.

## Decision: Credentials Excluded

Credential management (Connect domain) is intentionally excluded:
- **Security concern**: Exposing credential CRUD through MCP means AI agents could read/write secrets
- **Not needed**: Cloud resources already reference credentials by slug; the platform resolves them at deploy time
- A read-only `list_credentials` (slugs only, no secrets) may be added later if needed for discoverability

---

## Phase 0: Restructure Generated Code (prerequisite)

**Goal**: Move `gen/cloudresource/` → `gen/infrahub/cloudresource/` so generated code mirrors the domain structure of `internal/domains/infrahub/cloudresource/`.

### Tasks
1. Create `gen/infrahub/cloudresource/` directory structure
2. Move all provider subdirectories (`aws/`, `gcp/`, `azure/`, etc.) under `gen/infrahub/cloudresource/`
3. Move `registry_gen.go` to `gen/infrahub/cloudresource/`
4. Update package name from `package cloudresource` to match new path
5. Update all import paths in `internal/domains/infrahub/cloudresource/` that reference `gen/cloudresource`
6. Update the code generator (`schema2go` or whatever produces these files) config to output to the new path
7. Verify `go build ./...` passes
8. Verify existing tests pass

**Why first**: Every subsequent phase may add generated code. Getting the directory structure right now avoids a painful rename later.

---

## Phase 1: InfraChart + InfraProject + InfraPipeline (Highest Impact)

This phase turns the MCP server from a "single resource manager" into a "full infrastructure orchestrator."

### Phase 1A: InfraChart (3 tools)

| Tool | Backend RPC | Purpose |
|---|---|---|
| `search_infra_charts` | `InfraChartQueryController.find` | Browse available infrastructure chart templates |
| `get_infra_chart` | `InfraChartQueryController.get` | Retrieve full chart (templates, values.yaml, params) |
| `build_infra_chart` | `InfraChartQueryController.build` | Preview rendered output before committing |

**Files to create**:
- `internal/domains/infrahub/infrachart/tools.go` — tool definitions and handlers
- `internal/domains/infrahub/infrachart/search.go` — search/find RPC call
- `internal/domains/infrahub/infrachart/get.go` — get RPC call
- `internal/domains/infrahub/infrachart/build.go` — build (preview) RPC call

**Backend protos**:
- `ai.planton.infrahub.infrachart.v1.InfraChartQueryController`

### Phase 1B: InfraProject (6 tools)

| Tool | Backend RPC | Purpose |
|---|---|---|
| `search_infra_projects` | `InfraHubSearchQueryController.searchInfraProjects` | Search projects by org, env, text |
| `get_infra_project` | `InfraProjectQueryController.get` / `getByOrgBySlug` | Retrieve full project |
| `apply_infra_project` | `InfraProjectCommandController.apply` | Create or update |
| `delete_infra_project` | `InfraProjectCommandController.delete` | Remove project |
| `check_infra_project_slug` | `InfraProjectQueryController.checkSlugAvailability` | Slug uniqueness check |
| `undeploy_infra_project` | `InfraProjectCommandController.undeploy` | Tear down all deployed resources |

**Files to create**:
- `internal/domains/infrahub/infraproject/tools.go` — tool definitions and handlers
- `internal/domains/infrahub/infraproject/search.go` — search RPC call
- `internal/domains/infrahub/infraproject/get.go` — get and getByOrgBySlug RPC calls
- `internal/domains/infrahub/infraproject/apply.go` — apply RPC call
- `internal/domains/infrahub/infraproject/delete.go` — delete RPC call
- `internal/domains/infrahub/infraproject/slug.go` — slug availability check
- `internal/domains/infrahub/infraproject/undeploy.go` — undeploy RPC call

**Backend protos**:
- `ai.planton.infrahub.infraproject.v1.InfraProjectQueryController`
- `ai.planton.infrahub.infraproject.v1.InfraProjectCommandController`
- `ai.planton.search.v1.infrahub.InfraHubSearchQueryController.searchInfraProjects`

### Phase 1C: InfraPipeline (5 tools)

| Tool | Backend RPC | Purpose |
|---|---|---|
| `list_infra_pipelines` | `InfraPipelineQueryController.listByFilters` | List pipelines by project, status |
| `get_infra_pipeline` | `InfraPipelineQueryController.get` | Full pipeline status and details |
| `get_latest_infra_pipeline` | `InfraPipelineQueryController.getLastInfraPipelineByInfraProjectId` | Last pipeline for a project |
| `run_infra_pipeline` | `InfraPipelineCommandController.runInfraProjectChartSourcePipeline` | Trigger pipeline run |
| `cancel_infra_pipeline` | `InfraPipelineCommandController.cancel` | Cancel a running pipeline |

**Files to create**:
- `internal/domains/infrahub/infrapipeline/tools.go` — tool definitions and handlers
- `internal/domains/infrahub/infrapipeline/list.go` — list by filters RPC call
- `internal/domains/infrahub/infrapipeline/get.go` — get RPC call
- `internal/domains/infrahub/infrapipeline/latest.go` — get latest RPC call
- `internal/domains/infrahub/infrapipeline/run.go` — run pipeline RPC call
- `internal/domains/infrahub/infrapipeline/cancel.go` — cancel RPC call

**Backend protos**:
- `ai.planton.infrahub.infrapipeline.v1.InfraPipelineQueryController`
- `ai.planton.infrahub.infrapipeline.v1.InfraPipelineCommandController`

---

## Phase 2: Dependency Intelligence & Configuration

### Phase 2A: Graph — Dependency & Impact Analysis (4 tools)

| Tool | Backend RPC | Purpose |
|---|---|---|
| `get_organization_graph` | `GraphQueryController.getOrganizationGraph` | Full resource topology for an org |
| `get_cloud_resource_graph` | `GraphQueryController.getCloudResourceGraph` | Resource-centric dependency view |
| `get_dependencies` | `GraphQueryController.getDependencies` | What does this resource depend on? |
| `get_impact_analysis` | `GraphQueryController.getImpactAnalysis` | If I change/delete this, what breaks? |

**Files to create**:
- `internal/domains/graph/tools.go`
- `internal/domains/graph/organization.go`
- `internal/domains/graph/cloudresource.go`
- `internal/domains/graph/dependencies.go`
- `internal/domains/graph/impact.go`

**Backend protos**:
- `ai.planton.graph.v1.GraphQueryController`

### Phase 2B: ConfigManager — Variables & Secrets (5 tools)

| Tool | Backend RPC | Purpose |
|---|---|---|
| `list_variables` | `VariableQueryController.find` or list | List config variables in org/env |
| `apply_variable` | `VariableCommandController.apply` | Create or update a variable |
| `get_secret` | `SecretQueryController.get` | Retrieve secret metadata (not values) |
| `apply_secret` | `SecretCommandController.apply` | Create or update secret metadata |
| `create_secret_version` | `SecretVersionCommandController.create` | Set a new secret value |

**Files to create**:
- `internal/domains/configmanager/variable/tools.go`
- `internal/domains/configmanager/variable/list.go`
- `internal/domains/configmanager/variable/apply.go`
- `internal/domains/configmanager/secret/tools.go`
- `internal/domains/configmanager/secret/get.go`
- `internal/domains/configmanager/secret/apply.go`
- `internal/domains/configmanager/secretversion/tools.go`
- `internal/domains/configmanager/secretversion/create.go`

**Backend protos**:
- `ai.planton.configmanager.variable.v1.*`
- `ai.planton.configmanager.secret.v1.*`
- `ai.planton.configmanager.secretversion.v1.*`

---

## Phase 3: Audit, StackJob Control & Catalog

### Phase 3A: Audit — Version History (3 tools)

| Tool | Backend RPC | Purpose |
|---|---|---|
| `list_resource_versions` | `ApiResourceVersionQueryController.listByFilters` | Change history for any resource |
| `get_resource_version` | `ApiResourceVersionQueryController.getByIdWithContextSize` | Full version with diff |
| `get_resource_version_count` | `ApiResourceVersionQueryController.getCount` | How many changes made? |

**Files to create**:
- `internal/domains/audit/tools.go`
- `internal/domains/audit/list.go`
- `internal/domains/audit/get.go`
- `internal/domains/audit/count.go`

**Backend protos**:
- `ai.planton.audit.apiresourceversion.v1.ApiResourceVersionQueryController`

### Phase 3B: StackJob Commands — Lifecycle Control (3 tools)

| Tool | Backend RPC | Purpose |
|---|---|---|
| `rerun_stack_job` | `StackJobCommandController.rerun` | Retry a failed deployment |
| `cancel_stack_job` | `StackJobCommandController.cancel` | Cancel a stuck/running job |
| `check_stack_job_essentials` | `StackJobEssentialsQueryController.check` | Pre-validate all prerequisites |

**Files to add/modify**:
- `internal/domains/infrahub/stackjob/rerun.go`
- `internal/domains/infrahub/stackjob/cancel.go`
- `internal/domains/infrahub/stackjob/essentials.go`
- Update `internal/domains/infrahub/stackjob/tools.go` with new tool registrations

**Backend protos**:
- `ai.planton.infrahub.stackjob.v1.StackJobCommandController`
- `ai.planton.infrahub.stackjob.v1.StackJobEssentialsQueryController`

### Phase 3C: Deployment Component & IaC Module Catalog (3 tools)

| Tool | Backend RPC | Purpose |
|---|---|---|
| `search_deployment_components` | `InfraHubSearchQueryController.searchDeploymentComponentsByFilter` | Browse cloud resource type catalog |
| `search_iac_modules` | `InfraHubSearchQueryController.searchIacModulesByOrgContext` | Find IaC modules by kind/provisioner |
| `get_iac_module` | `IacModuleQueryController.get` | Full module details |

**Files to create**:
- `internal/domains/infrahub/deploymentcomponent/tools.go`
- `internal/domains/infrahub/deploymentcomponent/search.go`
- `internal/domains/infrahub/iacmodule/tools.go`
- `internal/domains/infrahub/iacmodule/search.go`
- `internal/domains/infrahub/iacmodule/get.go`

**Backend protos**:
- `ai.planton.search.v1.infrahub.InfraHubSearchQueryController`
- `ai.planton.infrahub.iacmodule.v1.IacModuleQueryController`

---

## Phase Summary

| Phase | Component | New Tools | Running Total |
|---|---|---|---|
| **0** | Gen code restructuring | 0 | 18 (existing) |
| **1A** | InfraChart | 3 | 21 |
| **1B** | InfraProject | 6 | 27 |
| **1C** | InfraPipeline | 5 | 32 |
| **2A** | Graph | 4 | 36 |
| **2B** | ConfigManager | 5 | 41 |
| **3A** | Audit | 3 | 44 |
| **3B** | StackJob Commands | 3 | 47 |
| **3C** | Catalog (DeploymentComponent + IacModule) | 3 | 50 |
| | **Total new tools** | **32** | **50** |

## Execution Order

Each phase is independently deployable. Within each phase:
1. Create the domain package with `doc.go`
2. Implement RPC call functions (following existing patterns in `cloudresource/`)
3. Define tool structs and handlers in `tools.go`
4. Register tools in `cmd/server` (or wherever tool registration happens)
5. Write tests
6. Verify `go build ./...` and `go test ./...`

## Patterns to Follow

All new code should follow the existing patterns established in:
- `internal/domains/infrahub/cloudresource/` — for tool definitions, handlers, RPC calls
- `internal/domains/infrahub/stackjob/` — for simpler query-only tools
- `internal/domains/infrahub/preset/` — for search + get tool pairs
- `internal/domains/resourcemanager/` — for non-InfraHub domains

## Notes

- Each phase can be a separate PR for manageable review
- Phase 0 (gen restructuring) should be its own PR to isolate the import path changes
- Phase 1 is the highest-value delivery and should be prioritized
- Phase 2A (Graph) is the "wow factor" differentiator for a world-class product
