# Next Task: MCP Server Gap Completion

## RULES OF ENGAGEMENT - READ FIRST

**When this file is loaded in a new conversation, the AI MUST:**

1. **DO NOT AUTO-EXECUTE** - Never start implementing without explicit user approval
2. **GATHER CONTEXT SILENTLY** - Read project files without outputting
3. **PRESENT STATUS SUMMARY** - Show what's done, what's pending, agreed next steps
4. **SHOW OPTIONS** - List recommended and alternative actions
5. **WAIT FOR DIRECTION** - Do NOT proceed until user explicitly says "go" or chooses an option

---

## Quick Resume Instructions

Drop this file into your conversation to quickly resume work on this project.

## Project: 20260301.01.mcp-server-gap-completion

**Description**: Close all gaps between the MCP server (100+ tools, 22 domains, 2 MCP resources) and the full Planton Cloud API surface (~564 proto files, 30+ domains). Add missing bounded contexts: Connect (credentials), IAM, full ResourceManager CRUD, StackJob AI-native tools, CloudResource lifecycle completion, PromotionPolicy, FlowControlPolicy, and expanded MCP resources.

**Goal**: Add ~60-70 missing tools and 5+ MCP resources across 8+ new/expanded domains.

**Tech Stack**: Go/gRPC/MCP

## Task Inventory

### TIER 1 -- Critical Gaps

| Task | Description | Est. Tools | Status |
|------|-------------|------------|--------|
| T02 | Architecture Decision: Generic vs Per-Type Credential Tools | 0 (design) | COMPLETED |
| T03 | ResourceManager: Organization Full CRUD | 4 | COMPLETED |
| T04 | ResourceManager: Environment Full CRUD | 4 | COMPLETED |
| T05 | Connect Domain: Credential Management (depends on T02) | 22 + 2 resources | COMPLETED |
| T06 | StackJob: AI-Native Tools (error recommendation, IaC resources) | 5 | COMPLETED |
| T07 | CloudResource: Lifecycle Completion (purge only; cleanup excluded -- platform-operator-only) | 1 | COMPLETED |

### TIER 2 -- Important Gaps

| Task | Description | Est. Tools | Status |
|------|-------------|------------|--------|
| T08 | IAM Domain: Access Control (Team, Policy, Role, ApiKey, Identity) + ProviderConnectionAuthorization | 23 | COMPLETED |
| T09 | InfraPipeline: Pipeline Record Cleanup (reduced from 2 — all triggers already covered) | 1 | COMPLETED |
| T10 | PromotionPolicy: Cross-Environment Deployment Governance | 4 | COMPLETED |
| T11 | FlowControlPolicy: Stack Job Execution Controls | 3 | COMPLETED |

### TIER 3 -- MCP Resources

| Task | Description | Est. Resources | Status |
|------|-------------|----------------|--------|
| T12 | Expand MCP Resources (api-resource-kinds://catalog) | 1 resource | COMPLETED |

### TIER 4 -- Explore / Deferred

| Task | Description | Status |
|------|-------------|--------|
| T13 | Investigation: Runner Domain Accessibility | NOT STARTED |
| T14 | Pipeline Log Retrieval: Stream-Collect-Return Tools | COMPLETED |
| T15 | MCP Prompts (Exploratory) | COMPLETED |

## Execution Order

T02 -> T03/T04/T06/T07 (parallel) -> T05 -> T08/T09/T10/T11 (parallel) -> T12 -> T13/T14/T15

## Essential Files to Review

### 1. Gap Analysis Reference
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/.cursor/plans/mcp_server_gap_analysis_1c322248.plan.md
```

### 2. Current Task Plan
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/tasks/
```

### 3. Project Documentation
- **README**: `/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/README.md`

### 4. Existing Domain Patterns (reference implementations)
- **CloudResource**: `internal/domains/infrahub/cloudresource/`
- **StackJob**: `internal/domains/infrahub/stackjob/`
- **Service**: `internal/domains/servicehub/service/`
- **Organization (current)**: `internal/domains/resourcemanager/organization/`
- **Server registration**: `internal/server/server.go`

## Knowledge Folders to Check

### Design Decisions
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/design-decisions/
```

### Coding Guidelines
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/coding-guidelines/
```

### Wrong Assumptions
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/wrong-assumptions/
```

### Don't Dos
```
/Users/suresh/scm/github.com/plantonhq/mcp-server-planton/_projects/2026-03/20260301.01.mcp-server-gap-completion/dont-dos/
```

## Planton API References (for implementation)

Key protobuf API paths in `/Users/suresh/scm/github.com/plantonhq/planton/apis/ai/planton/`:
- **ResourceManager**: `resourcemanager/organization/v1/`, `resourcemanager/environment/v1/`
- **Connect**: `connect/*credential/v1/` (20+ credential types)
- **IAM**: `iam/team/v1/`, `iam/iampolicy/v2/`, `iam/iamrole/v1/`, `iam/apikey/v1/`, `iam/identityaccount/v1/`
- **StackJob**: `infrahub/stackjob/v1/`
- **CloudResource**: `infrahub/cloudresource/v1/`
- **InfraPipeline**: `infrahub/infrapipeline/v1/`
- **PromotionPolicy**: `resourcemanager/promotionpolicy/v1/`
- **FlowControlPolicy**: `infrahub/flowcontrolpolicy/v1/` (verify path)

## Resume Checklist

When starting a new session:

1. [ ] Read the latest checkpoint (if any) from `checkpoints/`
2. [ ] Check task inventory above for current status
3. [ ] Review any design decisions in `design-decisions/`
4. [ ] Check coding guidelines in `coding-guidelines/`
5. [ ] Review wrong assumptions and dont-dos
6. [ ] Pick the next task and begin

## Session History

### COMPLETED: T07 CloudResource Lifecycle Completion (2026-03-01)

**Added `purge_cloud_resource` MCP tool to the CloudResource domain.**

**What was delivered:**

1. **`purge_cloud_resource` tool** - Destroy infrastructure + delete record in one atomic Temporal workflow
   - `purge.go`: Domain function modeled on `delete.go` (resolveResourceID -> cmdClient.Purge)
   - Tool definition in `tools.go`: Input struct, tool description, handler
   - Registered in `register.go`

**Key Decisions Made:**
- `cleanup` RPC excluded -- it requires `platform/operator` authorization, consistent with existing exclusion of `updateOutputs`, `pipelineApply`, `pipelineDestroy`
- This establishes the pattern: MCP tools expose user-level RPCs only; platform-operator RPCs are not surfaced

**Files Changed/Created:**
- `internal/domains/infrahub/cloudresource/purge.go` - NEW: Purge domain function
- `internal/domains/infrahub/cloudresource/tools.go` - Added purge tool definition; updated package doc (11 -> 12 tools)
- `internal/domains/infrahub/cloudresource/register.go` - Registered purge tool

### COMPLETED: T02 Architecture Decision DD-01 (2026-03-01)

**Connect domain tool architecture decision (see `design-decisions/DD01-connect-domain-tool-architecture.md`).**

### COMPLETED: T03 Organization Full CRUD (2026-03-01)

**Added 4 new MCP tools to the Organization domain, completing the full CRUD lifecycle.**

**What was delivered:**

1. **`get_organization` tool** - Retrieve organization by ID via OrganizationQueryController.Get
   - `get.go`: Domain function
   - Tool definition, input struct, handler in `tools.go`

2. **`create_organization` tool** - Provision new org with slug, name, description, contact_email
   - `create.go`: Constructs full Organization proto internally (api_version, kind, metadata, spec)
   - Calls OrganizationCommandController.Create (skips auth -- any authenticated user can create)

3. **`update_organization` tool** - Read-modify-write partial update by org_id
   - `update.go`: GETs current org, merges non-empty fields, calls OrganizationCommandController.Update
   - Both RPCs share a single gRPC connection within one WithConnection callback

4. **`delete_organization` tool** - Remove organization by ID
   - `delete.go`: Calls OrganizationCommandController.Delete(OrganizationId)

**Key Decisions Made:**
- Separate `create` + `update` tools (not `apply`): different auth models, different input shapes, clearer LLM intent
- RPCs excluded: `repairFgaTuples` (platform operator), `toggleGettingStartedTasks` (UI onboarding), `find` (platform operator), `checkSlugAvailability` (deferred -- not core CRUD)
- Update uses empty-string-means-no-change semantics (omitempty); field-clearing deferred to v2 if needed

**Files Changed/Created:**
- `internal/domains/resourcemanager/organization/get.go` - NEW
- `internal/domains/resourcemanager/organization/create.go` - NEW
- `internal/domains/resourcemanager/organization/update.go` - NEW
- `internal/domains/resourcemanager/organization/delete.go` - NEW
- `internal/domains/resourcemanager/organization/tools.go` - Added 4 tool definitions; updated package doc (1 -> 5 tools)
- `internal/domains/resourcemanager/organization/register.go` - Added 4 mcp.AddTool registrations

### COMPLETED: T04 Environment Full CRUD (2026-03-01)

**Added 4 new MCP tools to the Environment domain, completing the full CRUD lifecycle.**

**What was delivered:**

1. **`get_environment` tool** — Dual-resolution: retrieve by env_id OR by org+slug
   - `get.go`: Two domain functions — `Get` (by ID via QueryController.Get) and `GetByOrgBySlug` (via QueryController.GetByOrgBySlug)
   - Handler dispatches based on which inputs are provided

2. **`create_environment` tool** — Provision a new environment within an organization
   - `create.go`: Assembles full Environment proto (api_version, kind, metadata with org/slug/name, spec with description)
   - Calls EnvironmentCommandController.Create

3. **`update_environment` tool** — Read-modify-write partial update by env_id
   - `update.go`: GETs current environment, merges non-empty fields (name, description), calls EnvironmentCommandController.Update
   - Both RPCs share a single gRPC connection within one WithConnection callback

4. **`delete_environment` tool** — Remove environment by ID with cascading cleanup warning
   - `delete.go`: Calls EnvironmentCommandController.Delete

**Key Decisions Made:**
- Dual-resolution `get_environment` (deviation from Organization's ID-only pattern): environments are child resources naturally referenced by slug within an org, and the proto API provides `getByOrgBySlug` for this use case
- Separate `create` + `update` tools (not `apply`): consistent with Organization pattern — different auth models, clearer LLM intent
- RPCs excluded: `apply` (use separate create/update), `find` (platform operator), `findByOrg` (list_environments uses superior `findAuthorized`), `checkSlugAvailability` (deferred)
- Delete tool description warns about cascading cleanup of stack-modules, microservices, secrets, and clusters

**Files Changed/Created:**
- `internal/domains/resourcemanager/environment/get.go` - NEW
- `internal/domains/resourcemanager/environment/create.go` - NEW
- `internal/domains/resourcemanager/environment/update.go` - NEW
- `internal/domains/resourcemanager/environment/delete.go` - NEW
- `internal/domains/resourcemanager/environment/tools.go` - Added 4 tool definitions; updated package doc (1 -> 5 tools)
- `internal/domains/resourcemanager/environment/register.go` - Added 4 mcp.AddTool registrations

### COMPLETED: T06 StackJob AI-Native Tools (2026-03-01)

**Added 5 new AI-native and diagnostic MCP tools to the StackJob domain, bringing the total from 7 to 12 tools.**

**What was delivered:**

1. **`find_iac_resources_by_stack_job` tool** — List Pulumi/Terraform state entries for a specific stack job
   - `iac_resources.go`: Domain function via StackJobQueryController.FindIacResourcesByStackJobId
   - Input: `stack_job_id` (required)

2. **`find_iac_resources_by_api_resource` tool** — List IaC resources for any API resource via its latest stack job
   - `iac_resources.go`: Domain function via StackJobQueryController.FindIacResourcesByApiResourceId
   - Input: `api_resource_id` (required)

3. **`get_stack_job_input` tool** — Retrieve safe (credential-free) IaC module input for debugging
   - `stack_input.go`: Domain function via StackJobQueryController.GetCloudObjectStackInput
   - Returns the exact data fed to Pulumi/Terraform with backend credentials stripped

4. **`find_service_stack_jobs_by_env` tool** — Cross-environment deployment status for a service
   - `service_stack_jobs.go`: Domain function via StackJobQueryController.FindServiceStackJobsByEnv
   - First cross-domain import (servicev1.ServiceId) in the stackjob package

5. **`get_error_resolution_recommendation` tool** — AI-generated fix recommendation for stack job errors
   - `error_recommendation.go`: Domain function via StackJobQueryController.GetErrorResolutionRecommendation
   - Returns plain-text recommendation (StringValue), not JSON — unique pattern in the domain

**Key Decisions Made:**
- Dropped `get_last_stack_job_by_cloud_resource` from plan — already exists as `get_latest_stack_job`
- Added `get_stack_job_input` (not in original plan) — credential-free debugging tool, approved during planning
- `get_error_resolution_recommendation` returns `resp.GetValue()` directly instead of `domains.MarshalJSON` — response is plain text, not a proto message to serialize
- All 5 tools are read-only Query RPCs; no Command RPCs added

**Files Changed/Created:**
- `internal/domains/infrahub/stackjob/iac_resources.go` - NEW: Two IaC resource lookup functions
- `internal/domains/infrahub/stackjob/stack_input.go` - NEW: Safe stack input retrieval
- `internal/domains/infrahub/stackjob/service_stack_jobs.go` - NEW: Service deployment status by env
- `internal/domains/infrahub/stackjob/error_recommendation.go` - NEW: AI error resolution
- `internal/domains/infrahub/stackjob/tools.go` - Added 5 input structs, 5 tool definitions, 5 handlers; updated package doc (7 -> 12 tools)
- `internal/domains/infrahub/stackjob/register.go` - Added 5 mcp.AddTool registrations

---

### COMPLETED: T15 MCP Prompts — Cross-Domain Workflow Templates (2026-03-01)

**Added 5 MCP prompts — the first implementation of the third MCP primitive alongside the existing 172 tools and 5 resources.**

**What was delivered:**

1. **`debug_failed_deployment` prompt** — Diagnose failed infrastructure deployments
   - Guides through: stack job retrieval → error analysis → AI recommendation → IaC state review → dependency check
   - Arguments: `resource_id` (optional), `stack_job_id` (optional)

2. **`assess_change_impact` prompt** — Analyze blast radius before destructive changes
   - Guides through: impact analysis → dependents → dependencies → version history → risk assessment
   - Arguments: `resource_id` (required), `change_type` (optional)

3. **`explore_infrastructure` prompt** — Top-down discovery of organization topology
   - Guides through: org selection → graph visualization → environment summary → health checks
   - Arguments: `org_id` (optional)

4. **`provision_cloud_resource` prompt** — Guided cloud infrastructure creation
   - Guides through: type selection → prerequisite checks → configuration → deployment → monitoring
   - Arguments: `kind` (optional), `org_id` (optional), `env_id` (optional)

5. **`manage_access` prompt** — IAM discovery, audit, and policy management
   - Guides through: principal discovery → access review → authorization check → role lookup → policy management
   - Arguments: `org_id` (optional), `resource_id` (optional)

**Architecture:**
- Cross-cutting `prompts` package at `internal/domains/prompts/` (follows `discovery` precedent)
- Purely static handlers — no gRPC calls, no failure modes
- Hybrid content style — goal-oriented with recommended tool sequences
- `PromptResult` and `UserMessage` helpers added to `internal/domains/toolresult.go`

**Key Decisions Made:**
- Hybrid content style chosen over prescriptive or goal-oriented — balances discoverability with LLM flexibility
- Prompts are cross-cutting (not per-domain) because each workflow spans 2-4 bounded contexts
- No gRPC calls in prompt handlers — keeps them fast, testable, and failure-free
- 5 prompts selected based on: non-obvious multi-step workflows, cross-domain tool sequences, highest user pain points

**Files Created:**
- `internal/domains/prompts/doc.go` — Package documentation
- `internal/domains/prompts/register.go` — Prompt registration
- `internal/domains/prompts/debug_deployment.go` — debug_failed_deployment prompt
- `internal/domains/prompts/assess_impact.go` — assess_change_impact prompt
- `internal/domains/prompts/explore_infrastructure.go` — explore_infrastructure prompt
- `internal/domains/prompts/provision_resource.go` — provision_cloud_resource prompt
- `internal/domains/prompts/manage_access.go` — manage_access prompt

**Files Modified:**
- `internal/domains/toolresult.go` — Added `PromptResult` and `UserMessage` helpers
- `internal/server/server.go` — Added `registerPrompts(srv)` and prompts import

---

## Current Status

**Created**: 2026-03-01
**Current Task**: T14 (COMPLETED)
**Next Task**: T13 (Tier 4 — Runner domain investigation)
**Status**: All Tier 1 + all Tier 2 + all Tier 3 + T15 + T14 complete

**Current step:**
- ✅ COMPLETED T02: Architecture Decision DD-01 (2026-03-01)
- ✅ COMPLETED T07: CloudResource Lifecycle Completion — purge_cloud_resource (2026-03-01)
- ✅ COMPLETED T03: Organization Full CRUD — get, create, update, delete (2026-03-01)
- ✅ COMPLETED T04: Environment Full CRUD — get (dual-resolution), create, update, delete (2026-03-01)
- ✅ COMPLETED T06: StackJob AI-Native Tools — 5 tools (IaC resources, stack input, service env status, error recommendation) (2026-03-01)
- ✅ COMPLETED T05: Connect Domain — 22 tools + 2 MCP resources across 5 sub-packages (2026-03-01)
- ✅ COMPLETED T08: IAM Domain — 20 tools across 5 sub-packages + 3 ProviderConnectionAuthorization tools (2026-03-01)
- ✅ COMPLETED T09/T10/T11: Remaining Tier 2 — delete_infra_pipeline + PromotionPolicy (4 tools) + FlowControlPolicy (3 tools) (2026-03-01)
- ✅ COMPLETED T12: Expand MCP Resources — api-resource-kinds://catalog (1 resource) (2026-03-01)
- ✅ COMPLETED T15: MCP Prompts — 5 cross-domain workflow templates (2026-03-01)
- ✅ **COMPLETED T14: Pipeline Log Retrieval — 2 tools + generic stream drain utility (2026-03-01)**
- 🔵 Next: **T13** (Tier 4 — Runner domain investigation)

### ✅ COMPLETED: T05 Connect Domain — Credential Management (2026-03-01)

**Implemented the entire Connect bounded context: 22 MCP tools + 2 MCP resources across 5 sub-packages, covering 19 credential types, GitHub integration extras, and 3 platform connection resource types.**

**What was delivered:**

1. **Generic Credential Tools (5 tools + 2 MCP resources)** — Dispatcher-based architecture serving 19 credential types through a single set of generic tools
   - `apply_credential` — Create/update any credential type via `kind` discriminator + protojson bridge
   - `get_credential` — Retrieve by ID or by org+slug, with automatic sensitive field redaction
   - `delete_credential` — Remove any credential by ID
   - `search_credentials` — Search by kind within an org, with enum-mapped search filters
   - `check_credential_slug` — Slug availability check for any credential type
   - `credential-types://catalog` — MCP resource listing all 19 supported credential types
   - `credential-schema://{kind}` — MCP resource template returning JSON schema per type

2. **19 Credential Type Schemas** — Hand-crafted JSON schemas with explicit `sensitive` field annotations
   - Cloud: AWS, GCP, Azure, Kubernetes, DigitalOcean, Civo
   - DevOps: GitHub, GitLab, Docker, Maven, NPM
   - SaaS: Auth0, Cloudflare, Confluent, MongoDBAtlas, Snowflake, OpenFGA
   - IaC Backends: PulumiBackend, TerraformBackend

3. **GitHub Extras (5 dedicated tools)** — GitHub-specific operations beyond CRUD
   - `configure_github_webhook` / `remove_github_webhook`
   - `get_github_installation_info` / `list_github_repositories`
   - `get_github_installation_token`

4. **Platform Connections (12 tools)** — Three resource types, 4 tools each
   - DefaultProviderConnection: apply, get, resolve, delete
   - DefaultRunnerBinding: apply, get, resolve, delete
   - RunnerRegistration: apply, get, delete, search

**Key Decisions Made:**
- ProviderConnectionAuthorization deferred to T08 (IAM) — authorization concern, not credential CRUD
- Hand-crafted JSON schemas (not codegen) — pragmatic for ~20 simple types, auditable, explicit sensitive annotations
- OAuth/CloudFormation controllers excluded — browser/infra flows, not suitable for LLM-driven MCP tools
- Secret redaction is defense-in-depth — explicit MCP-side field redaction per OWASP MCP01:2025 guidance
- 19 credential types confirmed (not 20+) — Docker, Maven, NPM validated as having proto stubs

**Files Created (25 Go files + 20 JSON files):**
- `internal/domains/connect/doc.go`
- `internal/domains/connect/credential/{doc,register,tools,apply,get,delete,search,slug,registry,redact,resources,schema}.go` (12 files)
- `internal/domains/connect/github/{doc,register,tools}.go` (3 files)
- `internal/domains/connect/defaultprovider/{doc,register,tools}.go` (3 files)
- `internal/domains/connect/defaultrunner/{doc,register,tools}.go` (3 files)
- `internal/domains/connect/runner/{doc,register,tools}.go` (3 files)
- `schemas/credentials/{registry.json + 19 credential JSON schemas}`

**Files Modified:**
- `schemas/embed.go` — Added `CredentialFS` embed directive for credential schemas
- `internal/server/server.go` — Registered all 5 Connect sub-packages (credential, github, defaultprovider, defaultrunner, runner)

---

### ✅ COMPLETED: T08 IAM Domain + ProviderConnectionAuthorization (2026-03-01)

**Implemented the entire IAM bounded context (20 MCP tools across 5 sub-packages) and ProviderConnectionAuthorization (3 MCP tools in connect/providerauth), totaling 23 new tools delivered in 7 phases.**

**What was delivered:**

1. **IAM Identity (4 tools)** — `whoami`, `get_identity_account` (dual-resolution by ID or email), `invite_member`, `list_invitations` (with enum-mapped status filter)
   - Uses `UserInvitationCommandController.Create` for invitations (not deprecated v1 IamPolicy)
   - `list_invitations` maps string status to `UserInvitationStatusType` enum via `domains.NewEnumResolver`

2. **IAM Role (2 tools, read-only)** — `get_iam_role`, `list_iam_roles_for_resource_kind`
   - Role CRUD excluded (requires `back_office_admin` — platform-operator only)

3. **IAM Team (4 tools)** — `create_team`, `get_team`, `update_team` (read-modify-write), `delete_team`
   - `list_teams` deliberately skipped — team discovery handled by `list_principals(principal_kind=team)`
   - TeamMember `member_type` mapped to `ApiResourceKind` enum via `domains.NewEnumResolver`

4. **IAM Policy v2 (7 tools)** — `create_iam_policy`, `delete_iam_policy`, `upsert_iam_policies` (declarative sync), `revoke_org_access` (nuclear revocation), `list_resource_access`, `check_authorization`, `list_principals`
   - IAM Policy v1 entirely excluded (deprecated)
   - `list_principals` replaces v1's separate findMembersByOrg/findTeamsByOrg/findMembersByEnv/findTeamsByEnv

5. **IAM API Key (3 tools)** — `create_api_key` (with one-time key visibility warning), `list_api_keys`, `delete_api_key`
   - Handles RFC3339 timestamp parsing for `expires_at`

6. **Connect ProviderConnectionAuthorization (3 tools)** — `apply_provider_connection_authorization` (protojson bridge), `get_provider_connection_authorization` (dual-resolution by ID or org+provider+connection), `delete_provider_connection_authorization`
   - Lives at `internal/domains/connect/providerauth/` per Connect bounded context decision
   - `provider` field mapped via `domains.ResolveProvider`

**Key Decisions Made:**
- IAM Policy v2 only — v1 is deprecated
- Skip `list_teams` — team discovery via `list_principals(principal_kind=team)`
- User Invitations use `UserInvitationCommandController.create` (not deprecated v1 IamPolicy)
- IAM Role is read-only — CRUD requires `back_office_admin`
- ProviderConnectionAuthorization placed in `connect/providerauth/` (proto path, Connect bounded context)
- Tool count expanded from 12-15 estimate to 23 (approved by architect)

**Files Created (38 Go files):**
- `internal/domains/iam/doc.go`
- `internal/domains/iam/identity/{doc,register,tools,whoami,get,invite,invitations}.go` (7 files)
- `internal/domains/iam/role/{doc,register,tools,get}.go` (4 files)
- `internal/domains/iam/team/{doc,register,tools,create,get,update,delete}.go` (7 files)
- `internal/domains/iam/policy/{doc,register,tools,create,delete,upsert,revoke_org,list_access,check,list_principals}.go` (10 files)
- `internal/domains/iam/apikey/{doc,register,tools,create,list,delete}.go` (6 files)
- `internal/domains/connect/providerauth/{doc,register,tools}.go` (3 files)

**Files Modified:**
- `internal/server/server.go` — 6 new imports + 6 Register calls

---

### ✅ COMPLETED: T09 + T10 + T11 — Remaining Tier 2 Tools (2026-03-01)

**Completed all remaining Tier 2 tasks in a single session: 8 new tools across 3 domains (InfraPipeline, PromotionPolicy, FlowControlPolicy).**

**What was delivered:**

1. **T09: `delete_infra_pipeline` (1 tool)** — Pipeline record cleanup
   - Gap analysis originally said "2 missing trigger variants", but upon deep analysis all trigger RPCs were already covered by `run_infra_pipeline`
   - Reduced scope to the only useful missing tool: `delete_infra_pipeline` via `InfraPipelineCommandController.Delete`
   - Extended existing infrapipeline package (delete.go + tools.go + register.go updates)

2. **T10: PromotionPolicy (4 tools, new domain)** — Cross-environment deployment governance
   - `apply_promotion_policy` — Create-or-update with typed graph input (nodes + edges with manual_approval)
   - `get_promotion_policy` — Dual-resolution: by policy_id or by selector (kind + id)
   - `which_promotion_policy` — Resolve effective policy with inheritance (org-specific → platform default)
   - `delete_promotion_policy` — Delete by policy ID
   - Package: `internal/domains/resourcemanager/promotionpolicy/` (7 files)

3. **T11: FlowControlPolicy (3 tools, new domain)** — Stack job execution controls
   - `apply_flow_control_policy` — Create-or-update with flat boolean flags (is_manual, disable_on_lifecycle, skip_refresh, preview_before_update, pause_between_preview)
   - `get_flow_control_policy` — Dual-resolution: by policy_id or by selector (kind + id)
   - `delete_flow_control_policy` — Delete by policy ID
   - Package: `internal/domains/infrahub/flowcontrolpolicy/` (6 files)

**Key Decisions Made:**
- T09 reduced from 2 to 1 tool — gap analysis "missing trigger variants" was inaccurate; all trigger RPCs already covered
- `apply` pattern chosen over separate create+update — both policies are selector-scoped singletons with identical create/update input shapes (consistent with Connect credentials, not Org/Env pattern)
- `whichFlowControlPolicy` excluded from T11 — lives in StackJobEssentialsQueryController, returns meta-response (not the policy itself), overlaps with existing `check_stack_job_essentials`
- Both policy domains use `domains.NewEnumResolver[ApiResourceKind]` for selector kind resolution (same pattern as IAM policy's `list_principals`)

**Files Created (14 Go files):**
- `internal/domains/infrahub/infrapipeline/delete.go`
- `internal/domains/resourcemanager/promotionpolicy/{doc,register,tools,apply,get,which,delete}.go` (7 files)
- `internal/domains/infrahub/flowcontrolpolicy/{doc,register,tools,apply,get,delete}.go` (6 files)

**Files Modified (3 files):**
- `internal/domains/infrahub/infrapipeline/tools.go` — Added delete tool definition, updated package doc (7 → 8 tools)
- `internal/domains/infrahub/infrapipeline/register.go` — Added delete tool registration
- `internal/server/server.go` — 2 new imports + 2 Register calls (promotionpolicy, flowcontrolpolicy)

---

### ✅ COMPLETED: T12 Expand MCP Resources — api-resource-kinds://catalog (2026-03-01)

**Added `api-resource-kinds://catalog` MCP resource — the platform's navigational index for agents.**

**Architectural Surprise:**
The original T12 plan proposed 5 new resources. Investigation revealed:
- `credential-types://catalog` — already delivered in T05
- `cloud-object-presets://{kind}`, `deployment-components://catalog`, `iac-modules://catalog` — all three already have MCP tools (`search_cloud_object_presets`, `search_deployment_components`, `search_iac_modules`). These are dynamic database records, not static type-system metadata. Adding static MCP resources would be stale or redundant.

Scope reduced to 1 genuinely missing resource: `api-resource-kinds://catalog`.

**What was delivered:**

1. **`api-resource-kinds://catalog` resource** — Curated, static catalog of 29 platform API resource types grouped by 6 bounded contexts
   - resource_manager (3): organization, environment, promotion_policy
   - infra_hub (9): cloud_resource, cloud_object_preset, deployment_component, iac_module, stack_job, infra_pipeline, infra_chart, infra_project, flow_control_policy
   - service_hub (7): service, pipeline, variables_group, secrets_group, dns_domain, tekton_pipeline, tekton_task
   - config_manager (2): secret, secret_version
   - connect (4): default_provider_connection, default_runner_binding, runner_registration, provider_connection_authorization
   - iam (4): identity_account, team, iam_role, api_key
   - Cross-references to `credential-types://catalog` and `cloud-resource-kinds://catalog`

2. **`discovery` package** — New cross-cutting domain at `internal/domains/discovery/` following the established static-embedded resource pattern

**Key Decisions Made:**
- 3 of 5 originally planned resources dropped — already covered by existing tools or delivered in T05
- `RegisterResources(srv)` signature preserved (no `serverAddress`) — static resources only
- New `discovery` package created for platform-wide resources (not domain-specific)
- 29 user-relevant kinds curated from ~60+ ApiResourceKind enum values; internal/system kinds excluded

**Files Created:**
- `schemas/apiresourcekinds/catalog.json` — Hand-crafted catalog data
- `internal/domains/discovery/doc.go` — Package documentation
- `internal/domains/discovery/resources.go` — CatalogResource() + CatalogHandler()
- `internal/domains/discovery/register.go` — RegisterResources(srv)

**Files Modified:**
- `schemas/embed.go` — Added `ApiResourceKindFS` embed directive
- `internal/server/server.go` — Imported discovery package, registered resources

---

### ✅ COMPLETED: T14 Pipeline Log Retrieval — Stream-Collect-Return Tools (2026-03-01)

**Added 2 pipeline log retrieval tools + a generic stream drain utility. These tools internally call streaming `getLogStream` RPCs, collect all entries until EOF (completed job) or a 15s timeout (running job), and return the batch as a single text response.**

**What was delivered:**

1. **`get_pipeline_logs` tool** — Retrieve raw Tekton task logs for a ServiceHub CI/CD pipeline
   - `logs.go`: Domain function using `PipelineQueryController.GetLogStream`
   - Opens streaming RPC, drains up to 1000 entries, returns formatted text

2. **`get_infra_pipeline_logs` tool** — Retrieve raw Tekton task logs for an InfraPipeline
   - `logs.go`: Domain function using `InfraPipelineQueryController.GetLogStream`
   - Same pattern as ServiceHub Pipeline

3. **Generic `DrainStream[T]` utility** — Reusable stream-to-batch collection function
   - `stream.go`: Generic function that works with any `grpc.ServerStreamingClient[T]`
   - Caller provides a `format` function, keeping proto imports out of the shared package
   - Handles EOF (completed), context deadline (running), and max-entries cap

4. **`StreamCollectTimeout` constant** — 15s timeout for stream collection, separate from the 30s `DefaultRPCTimeout`

**Investigation Findings (pre-implementation):**
- All log RPCs across the entire Planton API are server-streaming — zero unary alternatives exist
- The unary `get` RPCs return structured status (errors, diagnostics, per-task state) but NOT raw log content
- Status polling via `get_pipeline` / `get_stack_job` already works for progress monitoring
- StackJob has no `getLogStream` RPC — its structured progress events are already well-covered by existing tools
- The server replays logs from the beginning, making the stream-collect pattern viable

**Key Decisions Made:**
- Stream-collect-return pattern chosen over backend API changes — pragmatic, no server-side work needed
- 15s stream timeout — fast enough for completed jobs (EOF arrives in seconds), bounded for running jobs
- 1000 entry cap prevents massive responses and token waste
- Text output format (not JSON) — logs are inherently textual, fewer tokens for AI agents
- Generic `DrainStream[T]` — reusable for any future streaming-to-unary bridges
- Not for StackJob — StackJob has no log stream, its structured diagnostics via `get_stack_job` are sufficient

**Files Created:**
- `internal/domains/stream.go` — Generic `DrainStream[T]` utility
- `internal/domains/servicehub/pipeline/logs.go` — Pipeline log domain function
- `internal/domains/infrahub/infrapipeline/logs.go` — InfraPipeline log domain function

**Files Modified:**
- `internal/grpc/client.go` — Added `StreamCollectTimeout` constant
- `internal/domains/servicehub/pipeline/tools.go` — Added log tool definition; updated package doc (9 → 10 tools)
- `internal/domains/servicehub/pipeline/register.go` — Registered log tool
- `internal/domains/infrahub/infrapipeline/tools.go` — Added log tool definition; updated package doc (8 → 9 tools)
- `internal/domains/infrahub/infrapipeline/register.go` — Registered log tool

---

## Objectives for Next Conversation

**All Tier 1, Tier 2, Tier 3, T15, and T14 are complete.** Remaining work is T13 — exploratory/research.

**Option A: T13 — Investigation: Runner Domain Accessibility (Research)**
- Determine if Runner domain cloud API wrappers (VPC lookup, subnet discovery, pod listing) are accessible via the MCP server's gRPC connection
- If accessible: plan tools for cloud API queries

**Option B: Close out this project**
- All critical and important gaps are closed. T13 is exploratory and could be deferred to a separate project.

---

*This file provides direct paths to all project resources for quick context loading.*
