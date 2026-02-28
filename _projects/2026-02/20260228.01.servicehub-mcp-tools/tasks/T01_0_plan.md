# Task T01: ServiceHub MCP Tools — Architecture & Implementation Plan

**Created**: 2026-02-28
**Status**: Planning

## Objective

Implement 35 MCP tools across 7 ServiceHub bounded contexts, following the established infrahub tool patterns. These tools expose the ServiceHub gRPC APIs (Service, Pipeline, VariablesGroup, SecretsGroup, DnsDomain, TektonPipeline, TektonTask) to AI agents via the Model Context Protocol.

---

## Codegen Applicability Assessment

### Existing Codegen Pipeline

The repo has a two-stage codegen pipeline:

1. **Stage 1 — `proto2schema`**: Parses OpenMCF provider protobuf definitions (`openmcf/apis/org/openmcf/provider/{cloud}/{resource}/v1/api.proto`) and generates JSON schemas.
2. **Stage 2 — `generator`**: Reads those JSON schemas and generates Go input structs with `validate()`, `applyDefaults()`, `toMap()`, and `Parse{Kind}()` functions.

The output is consumed by `apply_cloud_resource` in `internal/domains/infrahub/cloudresource/tools.go`, which calls `cloudresource.GetParser(kind)` to validate the opaque `cloud_object` map before sending it to the backend.

### Why the Codegen Does NOT Apply to ServiceHub

| Factor | Cloud Resources (InfraHub) | ServiceHub Entities |
|--------|---------------------------|---------------------|
| **Count** | 362+ cloud resource kinds | 7 API resources |
| **Schema source** | `openmcf/provider/{cloud}/{kind}/v1/spec.proto` | `planton/servicehub/{entity}/v1/spec.proto` |
| **Transport format** | `google.protobuf.Struct` (opaque `cloud_object`) | Typed protobuf messages (Service, Pipeline, etc.) |
| **API pattern** | Generic `apply(CloudResource)` with `cloud_object` field | Entity-specific RPCs (`apply(Service)`, `rerun(PipelineId)`, etc.) |
| **Domain operations** | All uniform CRUD | Each entity has unique ops (`disconnectGitRepo`, `resolveManualGate`, `upsertEntry`, etc.) |
| **Tool descriptions** | Formulaic (all follow same pattern) | Each tool needs domain-specific context and warnings |
| **Validation** | Schema-driven (field types, enums, required flags) | Business-rule-driven (mutual exclusion, conditional requirements) |

**Verdict: Codegen is not applicable.** The existing codegen exists to solve a scale problem (362+ kinds with uniform schemas). ServiceHub has 7 entities with unique operations and business rules — hand-writing is the correct approach.

### One Reuse Point

The `apply_service` tool CAN reuse the existing `cloudresource.GetParser(kind)` when the Service uses `deployment_config_source = inline`. In that case, `deployment_targets[].cloud_object` IS a cloud resource spec, and the generated parsers can validate it. This is reuse of the codegen *output*, not extension of the codegen *pipeline*.

---

## Tool Catalogue (35 tools)

### Tier 1 — Service (7 tools)

| # | Tool Name | RPC | Notes |
|---|-----------|-----|-------|
| 1 | `search_services` | Search index | Free-text search within org |
| 2 | `get_service` | `get` / `getByOrgBySlug` | By ID or org+slug |
| 3 | `apply_service` | `apply` | Idempotent create-or-update |
| 4 | `delete_service` | `delete` | Removes definition, NOT deployed resources |
| 5 | `disconnect_service_git_repo` | `disconnectGitRepo` | Removes webhook |
| 6 | `configure_service_webhook` | `configureWebhook` | Re-registers webhook |
| 7 | `list_service_branches` | `listBranches` | Lists Git branches from connected repo |

### Tier 2 — Pipeline (9 tools)

| # | Tool Name | RPC | Notes |
|---|-----------|-----|-------|
| 8 | `list_pipelines` | `listByFilters` | By org, optionally filtered by service/env |
| 9 | `get_pipeline` | `get` | Full pipeline details |
| 10 | `get_last_pipeline` | `getLastPipelineByServiceId` | Most recent pipeline for a service |
| 11 | `run_pipeline` | `runGitCommit` | Trigger pipeline for specific branch/commit |
| 12 | `rerun_pipeline` | `rerun` | Re-execute a failed pipeline |
| 13 | `cancel_pipeline` | `cancel` | Gracefully cancel running pipeline |
| 14 | `resolve_pipeline_gate` | `resolveManualGate` | Approve/reject deployment gate |
| 15 | `list_pipeline_files` | `listServiceRepoPipelineFiles` | Discover Tekton pipeline YAMLs in repo |
| 16 | `update_pipeline_file` | `updateServiceRepoPipelineFile` | Modify pipeline file in repo |

### Tier 3a — VariablesGroup (6 tools)

| # | Tool Name | RPC | Notes |
|---|-----------|-----|-------|
| 17 | `get_variables_group` | `get` / `getByOrgBySlug` | By ID or org+slug |
| 18 | `apply_variables_group` | `apply` | Idempotent create-or-update |
| 19 | `delete_variables_group` | `delete` | Remove group |
| 20 | `upsert_variable` | `upsertEntry` | Add/update single variable |
| 21 | `delete_variable` | `deleteEntry` | Remove single variable |
| 22 | `get_variable_value` | `getValue` | Resolve specific variable value |

### Tier 3b — SecretsGroup (6 tools)

| # | Tool Name | RPC | Notes |
|---|-----------|-----|-------|
| 23 | `get_secrets_group` | `get` / `getByOrgBySlug` | By ID or org+slug |
| 24 | `apply_secrets_group` | `apply` | Idempotent create-or-update |
| 25 | `delete_secrets_group` | `delete` | Remove group |
| 26 | `upsert_secret` | `upsertEntry` | Add/update single secret |
| 27 | `delete_secret` | `deleteEntry` | Remove single secret |
| 28 | `get_secret_value` | `getValue` | Resolve specific secret value (sensitive!) |

### Tier 4 — DnsDomain (3 tools)

| # | Tool Name | RPC | Notes |
|---|-----------|-----|-------|
| 29 | `get_dns_domain` | `get` / `getByOrgBySlug` | By ID or org+slug |
| 30 | `apply_dns_domain` | `apply` | Idempotent create-or-update |
| 31 | `delete_dns_domain` | `delete` | Remove domain |

### Tier 5a — TektonPipeline (2 tools)

| # | Tool Name | RPC | Notes |
|---|-----------|-----|-------|
| 32 | `get_tekton_pipeline` | `get` / `getByOrgAndName` | Retrieve pipeline template |
| 33 | `apply_tekton_pipeline` | `apply` | Create/update pipeline template |

### Tier 5b — TektonTask (2 tools)

| # | Tool Name | RPC | Notes |
|---|-----------|-----|-------|
| 34 | `get_tekton_task` | `get` / `getByOrgAndName` | Retrieve task template |
| 35 | `apply_tekton_task` | `apply` | Create/update task template |

---

## Package Structure

```
internal/domains/servicehub/
├── service/
│   ├── register.go
│   ├── tools.go
│   ├── search.go
│   ├── get.go
│   ├── apply.go
│   ├── delete.go
│   ├── disconnect.go
│   ├── webhook.go
│   └── branches.go
├── pipeline/
│   ├── register.go
│   ├── tools.go
│   ├── list.go
│   ├── get.go
│   ├── run.go
│   ├── rerun.go
│   ├── cancel.go
│   ├── gate.go
│   ├── files.go
│   └── update_file.go
├── variablesgroup/
│   ├── register.go
│   ├── tools.go
│   ├── get.go
│   ├── apply.go
│   ├── delete.go
│   ├── upsert_entry.go
│   ├── delete_entry.go
│   └── get_value.go
├── secretsgroup/
│   ├── register.go
│   ├── tools.go
│   ├── get.go
│   ├── apply.go
│   ├── delete.go
│   ├── upsert_entry.go
│   ├── delete_entry.go
│   └── get_value.go
├── dnsdomain/
│   ├── register.go
│   ├── tools.go
│   ├── get.go
│   ├── apply.go
│   └── delete.go
├── tektonpipeline/
│   ├── register.go
│   ├── tools.go
│   ├── get.go
│   └── apply.go
└── tektontask/
    ├── register.go
    ├── tools.go
    ├── get.go
    └── apply.go
```

Server registration in `internal/server/server.go`:

```go
servicehubservice.Register(srv, serverAddress)
pipeline.Register(srv, serverAddress)
variablesgroup.Register(srv, serverAddress)
secretsgroup.Register(srv, serverAddress)
dnsdomain.Register(srv, serverAddress)
tektonpipeline.Register(srv, serverAddress)
tektontask.Register(srv, serverAddress)
```

---

## Excluded RPCs (and why)

| RPC | Reason |
|-----|--------|
| `getStatusStream` / `getLogStream` | MCP doesn't support server-side streaming |
| `streamPipelinesByOrg` | Streaming + platform-operator-only |
| `find` (on all entities) | Platform-operator-only admin queries |
| `findByWebhookId` | Internal system use (temporal-worker) |
| `Pipeline.create/update/delete` | Pipelines are system-created; user ops are `run/rerun/cancel` |
| `TektonPipeline.delete` / `TektonTask.delete` | Low-priority admin ops; not needed for AI agents |

---

## Implementation Order

1. **Service** — Foundation; all other entities reference Services
2. **Pipeline** — Highest operational value ("what's happening with my deployment?")
3. **VariablesGroup + SecretsGroup** — Configuration management
4. **DnsDomain** — Simple CRUD, quick win
5. **TektonPipeline + TektonTask** — Reference data, lowest urgency

Each tier is independently shippable.

---

## Key Architectural Decisions

1. **`apply` over `create/update` split**: Consistent with InfraHub pattern. AI agents work declaratively.
2. **Entry-level ops for groups**: `upsert_variable`/`delete_variable` protect against accidentally clobbering entries.
3. **No streaming in MCP**: Use polling (`get_pipeline` with backoff) instead.
4. **`resolve_pipeline_gate` requires confirmation**: Tool description must warn that approving deploys to production.
5. **`get_secret_value` sensitivity**: Tool description must warn about plaintext secret exposure.

---

## Next Steps

1. [ ] Review this plan and provide feedback
2. [ ] Approve to proceed with Tier 1 (Service tools) implementation
3. [ ] Each tier will produce a checkpoint document upon completion
