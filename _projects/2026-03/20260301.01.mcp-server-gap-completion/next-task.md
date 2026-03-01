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
| T08 | IAM Domain: Access Control (Team, Policy, Role, ApiKey, Identity) | 12-15 | NOT STARTED |
| T09 | InfraPipeline: Missing Trigger Variants | 2 | NOT STARTED |
| T10 | PromotionPolicy: Cross-Environment Deployment Governance | 4 | NOT STARTED |
| T11 | FlowControlPolicy: Change Approval Workflows | 3 | NOT STARTED |

### TIER 3 -- MCP Resources

| Task | Description | Est. Resources | Status |
|------|-------------|----------------|--------|
| T12 | Expand MCP Resources (api-resource-kinds, credential-types, presets, catalogs) | 5+ | NOT STARTED |

### TIER 4 -- Explore / Deferred

| Task | Description | Status |
|------|-------------|--------|
| T13 | Investigation: Runner Domain Accessibility | NOT STARTED |
| T14 | Investigation: Non-Streaming Log Retrieval | NOT STARTED |
| T15 | MCP Prompts (Exploratory) | NOT STARTED |

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

## Current Status

**Created**: 2026-03-01
**Current Task**: T05 (COMPLETED)
**Next Task**: T08 (IAM Domain: Access Control)
**Status**: Ready for next task — all Tier 1 tasks complete

**Current step:**
- ✅ COMPLETED T02: Architecture Decision DD-01 (2026-03-01)
- ✅ COMPLETED T07: CloudResource Lifecycle Completion — purge_cloud_resource (2026-03-01)
- ✅ COMPLETED T03: Organization Full CRUD — get, create, update, delete (2026-03-01)
- ✅ COMPLETED T04: Environment Full CRUD — get (dual-resolution), create, update, delete (2026-03-01)
- ✅ COMPLETED T06: StackJob AI-Native Tools — 5 tools (IaC resources, stack input, service env status, error recommendation) (2026-03-01)
- ✅ **COMPLETED T05: Connect Domain — 22 tools + 2 MCP resources across 5 sub-packages (2026-03-01)**
- 🔵 Next: **T08** (IAM Domain: Access Control, 12-15 tools) or choose from Tier 2 tasks

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

## Objectives for Next Conversation

**Recommended: T08 — IAM Domain: Access Control (12-15 tools)**
- Team, Policy, Role, ApiKey, Identity management
- Includes the deferred ProviderConnectionAuthorization from T05

**Alternatives:**
- T09 — InfraPipeline: Missing Trigger Variants (2 tools, quick win)
- T10 — PromotionPolicy: Cross-Environment Deployment Governance (4 tools)
- T11 — FlowControlPolicy: Change Approval Workflows (3 tools)

---

*This file provides direct paths to all project resources for quick context loading.*
