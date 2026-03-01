---
name: T09 T10 T11 implementation
overview: "Add 8 MCP tools across 3 tasks: T09 (delete_infra_pipeline, 1 tool), T10 (PromotionPolicy CRUD + query, 4 tools in new domain), T11 (FlowControlPolicy CRUD, 3 tools in new domain). Uses `apply` pattern for policy singletons and typed input structs."
todos:
  - id: t09-delete
    content: "T09: Add delete_infra_pipeline tool to existing infrapipeline package (delete.go, update tools.go + register.go)"
    status: completed
  - id: t10-promotionpolicy
    content: "T10: Create promotionpolicy domain with 4 tools (apply, get, which, delete) -- 7 new files under internal/domains/resourcemanager/promotionpolicy/"
    status: completed
  - id: t11-flowcontrolpolicy
    content: "T11: Create flowcontrolpolicy domain with 3 tools (apply, get, delete) -- 6 new files under internal/domains/infrahub/flowcontrolpolicy/"
    status: completed
  - id: server-registration
    content: Wire T10 + T11 into server.go (2 imports + 2 Register calls)
    status: completed
  - id: build-verify
    content: go build verification after each phase
    status: completed
isProject: false
---

# T09 + T10 + T11: Remaining Tier 2 Tools

## Decisions (confirmed)

- **T09 reduced to 1 tool** -- gap analysis said 2 trigger variants, but all triggers already exist. Only `delete_infra_pipeline` adds real value.
- `**apply` pattern** for both policies -- selector-scoped singletons, not standalone entities. Consistent with Connect credentials.
- `**whichFlowControlPolicy` excluded** -- cross-domain meta-response, overlaps with `check_stack_job_essentials`.

## T09: delete_infra_pipeline (1 tool)

Extends existing `internal/domains/infrahub/infrapipeline/` package.

**Files modified:**

- [infrapipeline/tools.go](internal/domains/infrahub/infrapipeline/tools.go) -- add `DeleteInfraPipelineInput`, `DeleteTool()`, `DeleteHandler()`, update package doc (7 -> 8 tools)
- [infrapipeline/register.go](internal/domains/infrahub/infrapipeline/register.go) -- add `mcp.AddTool(srv, DeleteTool(), DeleteHandler(serverAddress))`

**File created:**

- `internal/domains/infrahub/infrapipeline/delete.go` -- `Delete(ctx, serverAddress, pipelineID)` calling `InfraPipelineCommandController.Delete`

**RPC:** `InfraPipelineCommandController.Delete(ApiResourceDeleteInput) -> InfraPipeline`

**Pattern reference:** follows `cancel.go` exactly -- single ID input, command controller call, marshal response.

---

## T10: PromotionPolicy (4 tools, new domain)

New package: `internal/domains/resourcemanager/promotionpolicy/`

**Bounded context:** ResourceManager (same as Organization, Environment). The proto lives at `ai.planton.resourcemanager.promotionpolicy.v1`.

### Tools


| Tool                      | RPC                     | Controller | Notes                                                                                                                |
| ------------------------- | ----------------------- | ---------- | -------------------------------------------------------------------------------------------------------------------- |
| `apply_promotion_policy`  | `Apply`                 | Command    | Create-or-update. Constructs full `PromotionPolicy` proto from typed input. If `policy_id` provided, it's an update. |
| `get_promotion_policy`    | `Get` / `GetBySelector` | Query      | Dual-resolution: by `policy_id` or by `selector_kind + selector_id`.                                                 |
| `which_promotion_policy`  | `WhichPolicy`           | Query      | Resolves effective policy with inheritance (org-specific -> platform default). Takes `selector_kind + selector_id`.  |
| `delete_promotion_policy` | `Delete`                | Command    | Delete by policy ID.                                                                                                 |


**Excluded RPCs:**

- `create/update` -- covered by `apply`
- `find` -- requires `platform/operator` authorization

### Files (7 new)

- `doc.go` -- package documentation (4 tools, backed by PromotionPolicyQueryController + CommandController)
- `register.go` -- `Register(srv, serverAddress)` with 4 `mcp.AddTool` calls
- `tools.go` -- 4 input structs, 4 tool definitions, 4 handlers
- `apply.go` -- constructs `PromotionPolicy` proto (api_version, kind, metadata, spec with selector + graph), calls `CommandController.Apply`
- `get.go` -- dual-resolution: `Query.Get(ApiResourceId)` or `Query.GetBySelector(ApiResourceSelector)`
- `which.go` -- calls `Query.WhichPolicy(ApiResourceSelector)`, returns full policy
- `delete.go` -- calls `CommandController.Delete(ApiResourceId)`

### Key input struct: ApplyPromotionPolicyInput

```go
type ApplyPromotionPolicyInput struct {
    PolicyID     string                 `json:"policy_id,omitempty"`
    Name         string                 `json:"name,omitempty"`
    SelectorKind string                 `json:"selector_kind"`
    SelectorID   string                 `json:"selector_id"`
    Strict       bool                   `json:"strict,omitempty"`
    Nodes        []EnvironmentNodeInput `json:"nodes"`
    Edges        []PromotionEdgeInput   `json:"edges"`
}

type EnvironmentNodeInput struct {
    Name string `json:"name"`
}

type PromotionEdgeInput struct {
    From           string `json:"from"`
    To             string `json:"to"`
    ManualApproval bool   `json:"manual_approval,omitempty"`
}
```

### Selector kind resolution

Uses `domains.NewEnumResolver[apiresourcekind.ApiResourceKind]` (existing pattern from [iam/policy/list_principals.go](internal/domains/iam/policy/list_principals.go)):

```go
var selectorKindResolver = domains.NewEnumResolver[apiresourcekind.ApiResourceKind](
    apiresourcekind.ApiResourceKind_value,
    "selector kind",
    "api_resource_kind_unspecified",
)
```

Valid values for PromotionPolicy: `organization`, `platform`.

---

## T11: FlowControlPolicy (3 tools, new domain)

New package: `internal/domains/infrahub/flowcontrolpolicy/`

**Bounded context:** InfraHub. The proto lives at `ai.planton.infrahub.flowcontrolpolicy.v1`.

### Tools


| Tool                         | RPC                     | Controller | Notes                                                                |
| ---------------------------- | ----------------------- | ---------- | -------------------------------------------------------------------- |
| `apply_flow_control_policy`  | `Apply`                 | Command    | Create-or-update. Flat boolean flags for flow control settings.      |
| `get_flow_control_policy`    | `Get` / `GetBySelector` | Query      | Dual-resolution: by `policy_id` or by `selector_kind + selector_id`. |
| `delete_flow_control_policy` | `Delete`                | Command    | Delete by policy ID.                                                 |


**Excluded RPCs:**

- `create/update` -- covered by `apply`
- `whichFlowControlPolicy` -- lives in StackJobEssentials, returns meta-response, overlaps with `check_stack_job_essentials`

### Files (6 new)

- `doc.go` -- package documentation (3 tools)
- `register.go` -- `Register(srv, serverAddress)` with 3 `mcp.AddTool` calls
- `tools.go` -- 3 input structs, 3 tool definitions, 3 handlers
- `apply.go` -- constructs `FlowControlPolicy` proto with `StackJobFlowControl` spec
- `get.go` -- dual-resolution: by ID or by selector
- `delete.go` -- calls `CommandController.Delete(ApiResourceId)`

### Key input struct: ApplyFlowControlPolicyInput

```go
type ApplyFlowControlPolicyInput struct {
    PolicyID                              string `json:"policy_id,omitempty"`
    Name                                  string `json:"name,omitempty"`
    SelectorKind                          string `json:"selector_kind"`
    SelectorID                            string `json:"selector_id"`
    IsManual                              bool   `json:"is_manual,omitempty"`
    DisableOnLifecycleEvents              bool   `json:"disable_on_lifecycle_events,omitempty"`
    SkipRefresh                           bool   `json:"skip_refresh,omitempty"`
    PreviewBeforeUpdateOrDestroy          bool   `json:"preview_before_update_or_destroy,omitempty"`
    PauseBetweenPreviewAndUpdateOrDestroy bool   `json:"pause_between_preview_and_update_or_destroy,omitempty"`
}
```

### Selector kind resolution

Same resolver pattern as T10. Valid values for FlowControlPolicy: `organization`, `environment`, `platform`, or any cloud resource kind (wider scope than PromotionPolicy).

---

## Server Registration

Modify [internal/server/server.go](internal/server/server.go):

- 2 new imports: `promotionpolicy` and `flowcontrolpolicy`
- 2 new `Register()` calls in `registerTools()`

T09 requires no server.go changes (infrapipeline already registered).

---

## Execution Phases

1. **T09** -- smallest change, extends existing package. Build + verify.
2. **T10** -- new PromotionPolicy domain (4 tools, 7 files). Build + verify.
3. **T11** -- new FlowControlPolicy domain (3 tools, 6 files). Same structure as T10. Build + verify.
4. **Server registration** -- wire T10 + T11. Final `go build`.

---

## Totals

- **8 new tools** (1 + 4 + 3)
- **15 new Go files** (1 + 7 + 6 + 1 test build)
- **3 modified files** (infrapipeline tools.go, infrapipeline register.go, server.go)

