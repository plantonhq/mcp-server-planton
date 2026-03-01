# Task T01: Gap Completion Master Plan

**Created**: 2026-03-01
**Status**: PENDING REVIEW
**Type**: Feature Development

> This plan requires your review before execution.

## Background

A comprehensive gap analysis compared the MCP server (100+ tools, 22 domain packages, 2 MCP resources) against the full Planton Cloud API surface (~564 proto files, 30+ domains). The analysis identified ~60-70 missing tools across 8+ domains and significant underutilization of MCP resources.

**Reference**: [MCP Server Gap Analysis Plan](../../../.cursor/plans/mcp_server_gap_analysis_1c322248.plan.md)

---

## Task Breakdown

Each task below is designed as a single conversation. Tasks within a tier can be done in any order, but tiers should be completed in sequence (Tier 1 before Tier 2, etc.).

### TIER 1 -- Critical Gaps

#### T02: Architecture Decision -- Generic vs Per-Type Credential Tools

**Type**: Design Decision (no code)
**Scope**: Decide the tool architecture for the Connect domain before building it.

The Connect domain has 20+ credential types (AWS, GCP, Azure, Kubernetes, GitHub, GitLab, Docker, Terraform Backend, Pulumi Backend, Auth0, MongoDB Atlas, Snowflake, Confluent, Cloudflare, DigitalOcean, Civo, OpenFGA, Maven, NPM), each with identical CRUD operations but type-specific spec fields.

**Option A**: Per-type tools (`apply_aws_credential`, `apply_gcp_credential`, ...) -- ~60+ tools. Follows existing pattern but creates tool proliferation.

**Option B**: Generic credential tools (`apply_credential(type, ...)`, `get_credential(type, ...)`) -- ~5-6 tools with a `credential_type` parameter. More elegant, requires `credential-types://catalog` MCP resource for discovery.

**Deliverable**: Design decision document in `design-decisions/`.

---

#### T03: ResourceManager -- Organization Full CRUD

**Type**: Implementation
**Scope**: Expand Organization from read-only (list) to full lifecycle management.
**Estimated tools**: 4

New tools:
- `get_organization` -- Get by ID
- `apply_organization` -- Create or update
- `delete_organization` -- Delete
- `check_organization_slug_availability` -- Slug pre-validation

**API surface** (from `apis/ai/planton/resourcemanager/organization/v1/`):
- `OrganizationCommandController`: `apply`, `delete`
- `OrganizationQueryController`: `get`, `checkSlugAvailability`

**Pattern**: Follow existing ResourceManager domain structure in `internal/domains/resourcemanager/organization/`.

---

#### T04: ResourceManager -- Environment Full CRUD

**Type**: Implementation
**Scope**: Expand Environment from read-only (list) to full lifecycle management.
**Estimated tools**: 4

New tools:
- `get_environment` -- Get by ID or org+slug
- `apply_environment` -- Create or update
- `delete_environment` -- Delete
- `check_environment_slug_availability` -- Slug pre-validation

**API surface** (from `apis/ai/planton/resourcemanager/environment/v1/`):
- `EnvironmentCommandController`: `apply`, `delete`
- `EnvironmentQueryController`: `get`, `getByOrgBySlug`, `checkSlugAvailability`

**Pattern**: Follow existing ResourceManager domain structure in `internal/domains/resourcemanager/environment/`.

---

#### T05: Connect Domain -- Credential Management (depends on T02)

**Type**: Implementation (largest single task)
**Scope**: Implement the entire Connect bounded context.
**Estimated tools**: 25-30 (depending on T02 decision)

**Credential types to cover** (all follow same CRUD pattern):
- Cloud providers: AWS, GCP, Azure, Kubernetes, DigitalOcean, Civo, Cloudflare
- SaaS: Auth0, Confluent, MongoDB Atlas, OpenFGA, Snowflake
- SCM: GitHub, GitLab
- Package registries: Docker, Maven, NPM
- IaC backends: Terraform Backend, Pulumi Backend

**Per-type operations** (from `apis/ai/planton/connect/*credential/v1/`):
- `apply` -- Create or update credential
- `get` -- Get by ID
- `getByOrgBySlug` -- Get by org and slug
- `delete` -- Delete credential

**Platform connection resources** (additional tools):
- `DefaultProviderConnection`: `get`, `apply`, `resolve`
- `DefaultRunnerBinding`: `get`, `apply`, `resolve`
- `RunnerRegistration`: `get`, `list`, `apply`, `delete`

**Security boundary**: Credential secret values must be write-only (same pattern as ConfigManager secrets).

**New MCP resource**: `credential-types://catalog` -- Static resource listing all credential types.

---

#### T06: StackJob -- AI-Native Tools

**Type**: Implementation
**Scope**: Add AI-specific and diagnostic tools to StackJob domain.
**Estimated tools**: 5

New tools:
- `get_error_resolution_recommendation` -- AI-powered analysis of failed stack jobs with actionable fix suggestions. Highest-ROI single tool for AI agents.
- `find_iac_resources_by_stack_job` -- View actual IaC resources (e.g., AWS resources) created/modified by a stack job.
- `find_iac_resources_by_api_resource` -- Same query but by API resource (cloud resource) ID.
- `get_last_stack_job_by_cloud_resource` -- Shortcut to get the most recent execution for a cloud resource.
- `find_service_stack_jobs_by_env` -- Find all service-related stack jobs in an environment.

**API surface** (from `apis/ai/planton/infrahub/stackjob/v1/`):
- `StackJobQueryController`: `getErrorResolutionRecommendation`, `findIacResourcesByStackJobId`, `findIacResourcesByApiResourceId`, `getLastStackJobByCloudResourceId`, `findServiceStackJobsByEnv`

**Pattern**: Extend existing `internal/domains/infrahub/stackjob/` package.

---

#### T07: CloudResource -- Lifecycle Completion

**Type**: Implementation
**Scope**: Add missing lifecycle operations to CloudResource.
**Estimated tools**: 2

New tools:
- `purge_cloud_resource` -- Atomic destroy + delete (the "clean removal" users want most often).
- `cleanup_cloud_resource` -- Fix resources stuck in bad state.

**API surface** (from `apis/ai/planton/infrahub/cloudresource/v1/`):
- `CloudResourceCommandController`: `purge`, `cleanup`

**Pattern**: Extend existing `internal/domains/infrahub/cloudresource/` package.

---

### TIER 2 -- Important Gaps

#### T08: IAM Domain -- Access Control

**Type**: Implementation (new bounded context)
**Scope**: Add IAM domain for access control management.
**Estimated tools**: 12-15

New tools (recommended subset):
- **Team**: `list_teams`, `get_team`, `apply_team`, `delete_team`
- **IamPolicy v2**: `list_resource_policies`, `grant_access`, `revoke_access`, `check_authorization`
- **IamRole**: `list_iam_roles` (read-only reference)
- **ApiKey**: `list_api_keys`, `create_api_key`, `delete_api_key`
- **IdentityAccount**: `whoami`, `get_identity_account` (read-only)

**API surface**:
- `apis/ai/planton/iam/team/v1/`
- `apis/ai/planton/iam/iampolicy/v2/`
- `apis/ai/planton/iam/iamrole/v1/`
- `apis/ai/planton/iam/apikey/v1/`
- `apis/ai/planton/iam/identityaccount/v1/`

**Pattern**: New `internal/domains/iam/` package with sub-packages per entity.

---

#### T09: InfraPipeline -- Missing Trigger Variants

**Type**: Implementation
**Scope**: Add missing trigger operations to InfraPipeline.
**Estimated tools**: 2

New tools:
- `run_infrapipeline_git_commit` -- Run pipeline for a specific Git commit (GitOps workflow).
- `run_infrapipeline_chart_source` -- Run pipeline from infra chart source.

**API surface** (from `apis/ai/planton/infrahub/infrapipeline/v1/`):
- `InfraPipelineCommandController`: `runGitCommit`, `runInfraProjectChartSourcePipeline`

**Pattern**: Extend existing `internal/domains/infrahub/infrapipeline/` package.

---

#### T10: PromotionPolicy

**Type**: Implementation (new domain)
**Scope**: Add promotion policy management for cross-environment deployment governance.
**Estimated tools**: 4

New tools:
- `get_promotion_policy` -- Get policy by ID or selector
- `apply_promotion_policy` -- Create or update policy
- `delete_promotion_policy` -- Delete policy
- `which_promotion_policy` -- Resolve which policy applies to a given context

**API surface** (from `apis/ai/planton/resourcemanager/promotionpolicy/v1/`):
- `PromotionPolicyCommandController`: `apply`, `delete`
- `PromotionPolicyQueryController`: `get`, `getBySelector`, `whichPolicy`

**Pattern**: New `internal/domains/resourcemanager/promotionpolicy/` package.

---

#### T11: FlowControlPolicy

**Type**: Implementation (new domain)
**Scope**: Add flow control policy management for infrastructure change approval workflows.
**Estimated tools**: 3

New tools:
- `get_flow_control_policy` -- Get policy details
- `apply_flow_control_policy` -- Create or update policy
- `delete_flow_control_policy` -- Delete policy

**API surface**: From `apis/ai/planton/infrahub/flowcontrolpolicy/v1/` (verify exact path).

**Pattern**: New `internal/domains/infrahub/flowcontrolpolicy/` package.

---

### TIER 3 -- MCP Resources

#### T12: Expand MCP Resources

**Type**: Implementation
**Scope**: Add MCP resources for improved agent discovery.
**Estimated resources**: 5+

Current (2 resources):
- `cloud-resource-kinds://catalog` (static)
- `cloud-resource-schema://{kind}` (template)

New resources:
- `api-resource-kinds://catalog` -- All ~50 platform resource types with metadata.
- `credential-types://catalog` -- All credential types with required fields summary (built with T05).
- `cloud-object-presets://{kind}` -- Presets for a specific cloud resource kind (template).
- `deployment-components://catalog` -- Browsable catalog of deployment components.
- `iac-modules://catalog` -- Browsable catalog of IaC modules.

**Pattern**: Extend `internal/domains/infrahub/cloudresource/resources.go` or create a dedicated `internal/domains/discovery/` package.

---

### TIER 4 -- Explore / Deferred

#### T13: Investigation -- Runner Domain Accessibility

**Type**: Research (no code)
**Scope**: Determine if Runner domain cloud API wrappers (VPC lookup, subnet discovery, pod listing) are accessible via the same gRPC server used by the MCP server.

**If accessible**: Plan tools for cloud API queries (e.g., "list VPCs in my AWS account").
**If not**: Document the limitation and defer.

---

#### T14: Investigation -- Non-Streaming Log Retrieval

**Type**: Research (no code)
**Scope**: Determine if non-streaming log/progress retrieval endpoints exist or could be added server-side for StackJob and Pipeline logs.

**If endpoints exist**: Plan snapshot-based log retrieval tools.
**If not**: Document the need for backend work and defer.

---

#### T15: MCP Prompts (Exploratory)

**Type**: Research (no code)
**Scope**: Evaluate whether MCP Prompts (the third MCP primitive) would add value for guided agent workflows.

Potential prompts:
- `deploy-cloud-resource` -- Guided deployment workflow
- `debug-failed-stack-job` -- Diagnostic workflow
- `setup-new-environment` -- Onboarding workflow

---

## Execution Order

```
T02 (Credential Architecture Decision)
 |
 v
T03, T04, T06, T07 (can be parallel -- independent domains)
 |
 v
T05 (Connect Domain -- depends on T02)
 |
 v
T08, T09, T10, T11 (Tier 2 -- can be parallel)
 |
 v
T12 (MCP Resources -- benefits from all prior work)
 |
 v
T13, T14, T15 (Tier 4 -- research tasks, any order)
```

## Estimated Total

| Tier | Tasks | Tools | Resources |
|------|-------|-------|-----------|
| Tier 1 | T02-T07 | ~40 | 1 |
| Tier 2 | T08-T11 | ~21 | 0 |
| Tier 3 | T12 | 0 | 5 |
| Tier 4 | T13-T15 | TBD | TBD |
| **Total** | **14 tasks** | **~61 tools** | **6 resources** |

## Open Decisions (Require Your Input)

1. **T02**: Generic vs per-type credential tools?
2. **T13**: Is the Runner domain accessible via the same gRPC server?
3. **T14**: Do non-streaming log endpoints exist?
4. **T15**: Are MCP Prompts worth implementing now?

---

**Please review this plan and provide your feedback. I will not proceed to execution until you explicitly approve.**
