# Tools & Resources Reference

Complete reference for all MCP tools and resources exposed by the Planton Cloud
MCP Server. For installation and configuration, see the
[README](../README.md).

---

## Overview

### Tools

**Cloud Resource Lifecycle**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`apply_cloud_resource`](#apply_cloud_resource) | Write | Create or update a cloud resource (idempotent) |
| [`get_cloud_resource`](#get_cloud_resource) | Read | Retrieve a cloud resource by ID or by coordinates |
| [`delete_cloud_resource`](#delete_cloud_resource) | Write | Delete a cloud resource record (does not tear down infrastructure) |
| [`list_cloud_resources`](#list_cloud_resources) | Read | List resources in an org with optional filters |
| [`destroy_cloud_resource`](#destroy_cloud_resource) | Write | Tear down cloud infrastructure while keeping the record |
| [`check_slug_availability`](#check_slug_availability) | Read | Verify a slug is available before creating a resource |
| [`rename_cloud_resource`](#rename_cloud_resource) | Write | Change a resource's display name |
| [`list_cloud_resource_locks`](#list_cloud_resource_locks) | Read | Show lock status, holder, and wait queue |
| [`remove_cloud_resource_locks`](#remove_cloud_resource_locks) | Write | Force-clear stuck locks on a resource |
| [`get_env_var_map`](#get_env_var_map) | Read | Extract env vars and secrets from a manifest |
| [`resolve_value_references`](#resolve_value_references) | Read | Resolve all valueFrom references in a resource's spec |

**Stack Job Observability**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`get_stack_job`](#get_stack_job) | Read | Retrieve a stack job by ID |
| [`get_latest_stack_job`](#get_latest_stack_job) | Read | Get the most recent stack job for a resource |
| [`list_stack_jobs`](#list_stack_jobs) | Read | List stack jobs with filters |

**Stack Job Commands**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`rerun_stack_job`](#rerun_stack_job) | Write | Retry a failed or completed stack job |
| [`cancel_stack_job`](#cancel_stack_job) | Write | Cancel a running stack job |
| [`resume_stack_job`](#resume_stack_job) | Write | Approve a stack job awaiting approval |
| [`check_stack_job_essentials`](#check_stack_job_essentials) | Read | Pre-validate deployment prerequisites |

**InfraChart Templates**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`list_infra_charts`](#list_infra_charts) | Read | List reusable infrastructure chart templates |
| [`get_infra_chart`](#get_infra_chart) | Read | Get full chart content (templates, params, values) |
| [`build_infra_chart`](#build_infra_chart) | Read | Preview rendered output with parameter overrides |

**InfraProject Lifecycle**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`search_infra_projects`](#search_infra_projects) | Read | Search infrastructure projects in an org |
| [`get_infra_project`](#get_infra_project) | Read | Retrieve a project by ID or org+slug |
| [`apply_infra_project`](#apply_infra_project) | Write | Create or update an infra project (idempotent) |
| [`delete_infra_project`](#delete_infra_project) | Write | Remove the project record |
| [`check_infra_project_slug`](#check_infra_project_slug) | Read | Check slug availability within an org |
| [`undeploy_infra_project`](#undeploy_infra_project) | Write | Tear down deployed resources, keep the record |

**InfraPipeline Monitoring & Control**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`list_infra_pipelines`](#list_infra_pipelines) | Read | List deployment pipelines in an org |
| [`get_infra_pipeline`](#get_infra_pipeline) | Read | Get full pipeline details and status |
| [`get_latest_infra_pipeline`](#get_latest_infra_pipeline) | Read | Most recent pipeline for a project |
| [`run_infra_pipeline`](#run_infra_pipeline) | Write | Trigger a new pipeline run |
| [`cancel_infra_pipeline`](#cancel_infra_pipeline) | Write | Cancel a running pipeline |
| [`resolve_infra_pipeline_env_gate`](#resolve_infra_pipeline_env_gate) | Write | Approve or reject a manual gate for an environment |
| [`resolve_infra_pipeline_node_gate`](#resolve_infra_pipeline_node_gate) | Write | Approve or reject a manual gate for a DAG node |

**Dependency Graph**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`get_organization_graph`](#get_organization_graph) | Read | Full resource topology for an org |
| [`get_environment_graph`](#get_environment_graph) | Read | Resource graph scoped to an environment |
| [`get_service_graph`](#get_service_graph) | Read | Service-centric graph showing all related resources |
| [`get_cloud_resource_graph`](#get_cloud_resource_graph) | Read | Graph centered on a specific cloud resource |
| [`get_dependencies`](#get_dependencies) | Read | What does a resource depend on? (upstream) |
| [`get_dependents`](#get_dependents) | Read | What depends on a resource? (downstream) |
| [`get_impact_analysis`](#get_impact_analysis) | Read | Blast radius for a delete or update operation |

**Config Manager — Variables**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`list_variables`](#list_variables) | Read | List config variables in an org |
| [`get_variable`](#get_variable) | Read | Get a variable by ID or org+scope+slug |
| [`apply_variable`](#apply_variable) | Write | Create or update a config variable (idempotent) |
| [`delete_variable`](#delete_variable) | Write | Permanently delete a variable |
| [`resolve_variable`](#resolve_variable) | Read | Look up a variable's current value by slug |

**Config Manager — Secrets**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`list_secrets`](#list_secrets) | Read | List secrets in an org (metadata only, no values) |
| [`get_secret`](#get_secret) | Read | Get secret metadata by ID or org+scope+slug |
| [`apply_secret`](#apply_secret) | Write | Create or update secret metadata (idempotent) |
| [`delete_secret`](#delete_secret) | Write | Permanently delete a secret and all its versions |
| [`create_secret_version`](#create_secret_version) | Write | Store a new encrypted key-value version |
| [`list_secret_versions`](#list_secret_versions) | Read | List version history for a secret (no values) |

**Audit & Version History**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`list_resource_versions`](#list_resource_versions) | Read | Paginated change history for any platform resource |
| [`get_resource_version`](#get_resource_version) | Read | Full version details with unified diff |
| [`get_resource_version_count`](#get_resource_version_count) | Read | Count of versions for a resource |

**Context Discovery**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`list_organizations`](#list_organizations) | Read | List organizations the caller belongs to |
| [`list_environments`](#list_environments) | Read | List environments within an organization |

**Presets**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`search_cloud_object_presets`](#search_cloud_object_presets) | Read | Search for preset templates |
| [`get_cloud_object_preset`](#get_cloud_object_preset) | Read | Get full preset content by ID |

**Catalog**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`search_deployment_components`](#search_deployment_components) | Read | Browse the cloud resource type catalog |
| [`get_deployment_component`](#get_deployment_component) | Read | Get full component details by ID or kind |
| [`search_iac_modules`](#search_iac_modules) | Read | Find IaC modules by kind, provisioner, or provider |
| [`get_iac_module`](#get_iac_module) | Read | Get full IaC module details by ID |

**Service Lifecycle**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`search_services`](#search_services) | Read | Search services within an organization |
| [`get_service`](#get_service) | Read | Retrieve a service by ID or org+slug |
| [`apply_service`](#apply_service) | Write | Create or update a service (idempotent) |
| [`delete_service`](#delete_service) | Write | Delete a service record |
| [`disconnect_service_git_repo`](#disconnect_service_git_repo) | Write | Remove the webhook from the connected Git repository |
| [`configure_service_webhook`](#configure_service_webhook) | Write | Create or refresh the webhook on the connected Git repository |
| [`list_service_branches`](#list_service_branches) | Read | List Git branches from the connected repository |

**Service CI/CD Pipelines**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`list_pipelines`](#list_pipelines) | Read | List CI/CD pipelines with optional service and environment filters |
| [`get_pipeline`](#get_pipeline) | Read | Retrieve full pipeline details by ID |
| [`get_last_pipeline`](#get_last_pipeline) | Read | Most recent pipeline for a service |
| [`run_pipeline`](#run_pipeline) | Write | Trigger a new pipeline run for a branch |
| [`rerun_pipeline`](#rerun_pipeline) | Write | Re-run a previously executed pipeline |
| [`cancel_pipeline`](#cancel_pipeline) | Write | Cancel a running pipeline |
| [`resolve_pipeline_gate`](#resolve_pipeline_gate) | Write | Approve or reject a manual deployment gate |
| [`list_pipeline_files`](#list_pipeline_files) | Read | List Tekton pipeline YAML files in the service repository |
| [`update_pipeline_file`](#update_pipeline_file) | Write | Create or update a pipeline YAML file in the service repository |

**Service Variables Groups**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`search_variables`](#search_variables) | Read | Search variable entries across all groups in an org |
| [`get_variables_group`](#get_variables_group) | Read | Retrieve a variables group by ID or org+slug |
| [`apply_variables_group`](#apply_variables_group) | Write | Create or update a variables group (idempotent) |
| [`delete_variables_group`](#delete_variables_group) | Write | Delete a variables group |
| [`upsert_variable`](#upsert_variable) | Write | Add or update a single variable entry in a group |
| [`delete_variable`](#servicehub-delete_variable) | Write | Remove a single variable entry from a group |
| [`get_variable_value`](#get_variable_value) | Read | Resolve the value of a specific variable |
| [`transform_variables`](#transform_variables) | Read | Batch-resolve `$variables-group/` references |

**Service Secrets Groups**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`search_secrets`](#search_secrets) | Read | Search secret entries across all groups in an org |
| [`get_secrets_group`](#get_secrets_group) | Read | Retrieve a secrets group by ID or org+slug |
| [`apply_secrets_group`](#apply_secrets_group) | Write | Create or update a secrets group (idempotent) |
| [`delete_secrets_group`](#delete_secrets_group) | Write | Delete a secrets group |
| [`upsert_secret`](#upsert_secret) | Write | Add or update a single secret entry in a group |
| [`delete_secret`](#servicehub-delete_secret) | Write | Remove a single secret entry from a group |
| [`get_secret_value`](#get_secret_value) | Read | Resolve the plaintext value of a specific secret |
| [`transform_secrets`](#transform_secrets) | Read | Batch-resolve `$secrets-group/` references |

**DNS Domains**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`get_dns_domain`](#get_dns_domain) | Read | Retrieve a DNS domain by ID or org+slug |
| [`apply_dns_domain`](#apply_dns_domain) | Write | Create or update a DNS domain (idempotent) |
| [`delete_dns_domain`](#delete_dns_domain) | Write | Delete a DNS domain record |

**Tekton Pipelines**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`get_tekton_pipeline`](#get_tekton_pipeline) | Read | Retrieve a Tekton pipeline template by ID or org+slug |
| [`apply_tekton_pipeline`](#apply_tekton_pipeline) | Write | Create or update a Tekton pipeline template (idempotent) |
| [`delete_tekton_pipeline`](#delete_tekton_pipeline) | Write | Delete a Tekton pipeline template |

**Tekton Tasks**

| Tool | Operation | Description |
|------|-----------|-------------|
| [`get_tekton_task`](#get_tekton_task) | Read | Retrieve a Tekton task template by ID or org+slug |
| [`apply_tekton_task`](#apply_tekton_task) | Write | Create or update a Tekton task template (idempotent) |
| [`delete_tekton_task`](#delete_tekton_task) | Write | Delete a Tekton task template |

### MCP Resources

| Resource URI | Description |
|--------------|-------------|
| [`cloud-resource-kinds://catalog`](#cloud-resource-kindscatalog) | Catalog of all 362 supported kinds grouped by 17 cloud providers |
| [`cloud-resource-schema://{kind}`](#cloud-resource-schemakind) | Full JSON schema for a specific kind |

---

## Resource Identification Pattern

Many cloud resource tools accept two identification strategies. Use whichever
you have available:

- **By ID**: provide `id` alone.
- **By coordinates**: provide all four of `kind`, `org`, `env`, and `slug`
  together.

Mixing a partial set of coordinates without `id` returns a validation error.
This pattern applies to: `get_cloud_resource`, `delete_cloud_resource`,
`destroy_cloud_resource`, `list_cloud_resource_locks`,
`remove_cloud_resource_locks`, and `rename_cloud_resource`.

The parameters for the coordinate path are:

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | string | System-assigned resource ID. Provide this alone, or provide all of `kind`, `org`, `env`, `slug`. |
| `kind` | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | string | Organization identifier. |
| `env` | string | Environment identifier. |
| `slug` | string | Immutable unique resource name within the `(org, env, kind)` scope. |

Individual tool sections below note "Accepts the standard resource
identification pattern" rather than repeating this table.

---

## Cloud Resource Lifecycle Tools

### apply_cloud_resource

Create or update a cloud resource on Planton Cloud. The operation is
idempotent: if the resource already exists it is updated in-place; if it does
not exist it is created. This matches `kubectl apply` semantics.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `cloud_object` | **yes** | object | Full OpenMCF cloud resource object. Must contain `api_version`, `kind`, `metadata` (with `name`, `org`, `env`), and `spec`. |

#### cloud_object structure

```json
{
  "api_version": "ai.planton.provider.aws.v1",
  "kind": "AwsEksCluster",
  "metadata": {
    "name": "my-cluster",
    "org": "my-org",
    "env": "production"
  },
  "spec": {
    "..."
  }
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `api_version` | yes | Versioned API namespace for the provider (e.g. `ai.planton.provider.aws.v1`). Obtain from the kind catalog. |
| `kind` | yes | PascalCase resource type (e.g. `AwsEksCluster`). Must match a kind in the catalog. |
| `metadata.name` | yes | Human-readable display name for the resource. |
| `metadata.org` | yes | Organization identifier — the tenant that will own the resource. |
| `metadata.env` | yes | Environment identifier (e.g. `production`, `staging`). |
| `spec` | yes | Provider-specific configuration. Structure is defined by the kind's JSON schema. |

#### Agent workflow

Before calling `apply_cloud_resource`, an agent must resolve the schema for the
target kind. The correct sequence is:

1. Read `cloud-resource-kinds://catalog` to discover all supported kinds and
   their `api_version` values.
2. Read `cloud-resource-schema://{kind}` to get the full spec definition,
   including field types, validation constraints, and default values.
3. Assemble the `cloud_object` using the resolved `api_version`, `kind`, and a
   `spec` that satisfies the schema.
4. Call `apply_cloud_resource` with the assembled `cloud_object`.

Skipping steps 1–2 and guessing the spec structure will produce validation
errors from the backend.

---

### get_cloud_resource

Retrieve a cloud resource from Planton Cloud. Returns the full resource
including metadata, spec, and status. Accepts the standard
[resource identification pattern](#resource-identification-pattern).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | System-assigned resource ID. Provide this alone, or provide all of `kind`, `org`, `env`, `slug`. |
| `kind` | conditional | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | conditional | string | Organization identifier. |
| `env` | conditional | string | Environment identifier. |
| `slug` | conditional | string | Immutable unique resource name within the `(org, env, kind)` scope. |

---

### delete_cloud_resource

Delete a cloud resource record from Planton Cloud. This removes the record
only — it does **not** tear down the actual cloud infrastructure. Use
[`destroy_cloud_resource`](#destroy_cloud_resource) first to tear down
infrastructure, then this tool to remove the record. Accepts the standard
[resource identification pattern](#resource-identification-pattern).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | System-assigned resource ID. Provide this alone, or provide all of `kind`, `org`, `env`, `slug`. |
| `kind` | conditional | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | conditional | string | Organization identifier. |
| `env` | conditional | string | Environment identifier. |
| `slug` | conditional | string | Immutable unique resource name within the `(org, env, kind)` scope. |

---

### list_cloud_resources

List cloud resources in an organization on Planton Cloud. Returns resources
grouped by environment and kind. Use
[`list_organizations`](#list_organizations) to discover available organization
identifiers.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `envs` | no | string[] | Environment slugs to filter by. |
| `kinds` | no | string[] | PascalCase cloud resource kinds to filter by. Read `cloud-resource-kinds://catalog` for valid values. |
| `search_text` | no | string | Free-text search query. |

---

### destroy_cloud_resource

Destroy the cloud infrastructure (Terraform/Pulumi destroy) for a resource
while keeping the resource record on Planton Cloud. This tears down the actual
cloud resources (VPCs, clusters, databases, etc.). Use
[`delete_cloud_resource`](#delete_cloud_resource) to remove the record itself.
Use [`get_latest_stack_job`](#get_latest_stack_job) to monitor the destroy
operation's progress. Accepts the standard
[resource identification pattern](#resource-identification-pattern).

> **WARNING:** This is a destructive operation that will destroy real cloud
> infrastructure.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | System-assigned resource ID. Provide this alone, or provide all of `kind`, `org`, `env`, `slug`. |
| `kind` | conditional | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | conditional | string | Organization identifier. |
| `env` | conditional | string | Environment identifier. |
| `slug` | conditional | string | Immutable unique resource name within the `(org, env, kind)` scope. |

---

### check_slug_availability

Check whether a cloud resource slug is available within the scoped composite
key `(org, env, kind)`. Slugs must be unique within this scope. Use this before
[`apply_cloud_resource`](#apply_cloud_resource) to verify that the desired slug
is not already taken.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `env` | **yes** | string | Environment identifier. |
| `kind` | **yes** | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `slug` | **yes** | string | The slug to check for availability. |

---

### rename_cloud_resource

Rename a cloud resource on Planton Cloud. Changes the human-readable display
name; the immutable slug is unaffected. Returns the updated resource. Accepts
the standard
[resource identification pattern](#resource-identification-pattern).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | System-assigned resource ID. Provide this alone, or provide all of `kind`, `org`, `env`, `slug`. |
| `kind` | conditional | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | conditional | string | Organization identifier. |
| `env` | conditional | string | Environment identifier. |
| `slug` | conditional | string | Immutable unique resource name within the `(org, env, kind)` scope. |
| `new_name` | **yes** | string | The new display name for the cloud resource. |

---

### list_cloud_resource_locks

List lock information for a cloud resource on Planton Cloud. Returns whether
the resource is locked, current lock holder details (workflow ID, acquired
timestamp, TTL remaining), and any workflows waiting in the lock queue. Use
[`remove_cloud_resource_locks`](#remove_cloud_resource_locks) to force-clear
stuck locks. Accepts the standard
[resource identification pattern](#resource-identification-pattern).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | System-assigned resource ID. Provide this alone, or provide all of `kind`, `org`, `env`, `slug`. |
| `kind` | conditional | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | conditional | string | Organization identifier. |
| `env` | conditional | string | Environment identifier. |
| `slug` | conditional | string | Immutable unique resource name within the `(org, env, kind)` scope. |

---

### remove_cloud_resource_locks

Remove all locks (active lock and wait queue) for a cloud resource on Planton
Cloud. Returns details about what was removed (active lock removed, queue
entries cleared). Use
[`list_cloud_resource_locks`](#list_cloud_resource_locks) to inspect the
current lock state and
[`get_latest_stack_job`](#get_latest_stack_job) to verify no jobs are running
before removing locks. Accepts the standard
[resource identification pattern](#resource-identification-pattern).

> **WARNING:** Removing locks on a resource with an active stack job may cause
> IaC state corruption.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | System-assigned resource ID. Provide this alone, or provide all of `kind`, `org`, `env`, `slug`. |
| `kind` | conditional | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | conditional | string | Organization identifier. |
| `env` | conditional | string | Environment identifier. |
| `slug` | conditional | string | Immutable unique resource name within the `(org, env, kind)` scope. |

---

### get_env_var_map

Extract the environment variable map from a cloud resource manifest. Provide
the raw YAML content of a cloud resource manifest (OpenMCF format). The server
parses the YAML, identifies the resource kind, extracts environment variables
and secrets, and resolves valueFrom references to plain string values. Returns
separate maps for variables and secrets.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `yaml_content` | **yes** | string | Raw YAML content of the cloud resource manifest (OpenMCF format with `api_version`, `kind`, `metadata`, and `spec`). |

---

### resolve_value_references

Resolve all valueFrom references in a cloud resource's spec. The server loads
the resource, finds all valueFrom references in its specification, resolves
them to concrete values, and returns the fully transformed cloud resource as
YAML. The response includes resolution status, any errors, and diagnostics.

The `kind` field is always required (used for both authorization and resource
transformation). The resource is identified by `id` alone, or by all of `org`,
`env`, and `slug` — note that `kind` is separate from the coordinate path here
unlike other tools where it is part of the four-field set.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `kind` | **yes** | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). Always required. |
| `id` | conditional | string | System-assigned resource ID. Provide this alone, or provide all of `org`, `env`, `slug`. |
| `org` | conditional | string | Organization identifier. Required with `env`, `slug` when `id` is not provided. |
| `env` | conditional | string | Environment identifier. Required with `org`, `slug` when `id` is not provided. |
| `slug` | conditional | string | Immutable unique resource slug. Required with `org`, `env` when `id` is not provided. |

---

## Stack Job Observability Tools

Stack jobs track the outcome of infrastructure operations (apply, destroy,
refresh, etc.). After calling `apply_cloud_resource` or
`destroy_cloud_resource`, use these tools to monitor progress and verify
success.

### get_stack_job

Retrieve a specific stack job by its ID. Returns the full job including
operation type, progress, result, timestamps, errors, and IaC resource counts.
Use when you have a stack job ID from a previous response or from the user.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The stack job ID. |

---

### get_latest_stack_job

Retrieve the most recent stack job for a cloud resource. This is the primary
tool to check whether an `apply_cloud_resource` or `destroy_cloud_resource`
operation completed successfully. Returns the full stack job including
operation type, progress, result, timestamps, errors, and IaC resource counts.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `cloud_resource_id` | **yes** | string | The cloud resource ID to look up the most recent stack job for. |

---

### list_stack_jobs

List stack jobs matching the given filters. Requires an organization ID.
Supports filtering by environment, cloud resource kind, resource ID, operation
type, execution status, and result. Returns a paginated list. Use to find
failed deployments, audit provisioning history, or discover jobs across
resources.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `env` | no | string | Environment name to filter by. |
| `cloud_resource_kind` | no | string | PascalCase cloud resource kind to filter by. |
| `cloud_resource_id` | no | string | Cloud resource ID to filter by. |
| `operation_type` | no | string | One of: `init`, `refresh`, `update_preview`, `update`, `destroy_preview`, `destroy`. |
| `status` | no | string | One of: `queued`, `running`, `completed`, `awaiting_approval`. |
| `result` | no | string | One of: `tbd`, `succeeded`, `failed`, `cancelled`, `skipped`. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

## Context Discovery Tools

Use these tools to discover the operating context (organizations and
environments) before working with cloud resources.

### list_organizations

List all organizations the caller is a member of. Returns the full organization
objects including id, name, and slug. This is often the first tool an agent
calls to establish the operating context.

#### Parameters

This tool takes no input parameters. The server returns organizations scoped to
the authenticated caller's membership.

---

### list_environments

List environments the caller can access within an organization. Returns only
environments where the caller has at least view permission. Use
[`list_organizations`](#list_organizations) first to discover available
organization identifiers.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |

---

## Preset Tools

Presets are pre-configured cloud resource manifests that serve as starting
points for [`apply_cloud_resource`](#apply_cloud_resource). The two-step
workflow is: search for presets, then get the full content of the one you want.

### search_cloud_object_presets

Search for cloud object preset templates. When `org` is provided, results
include both official platform presets and organization-specific presets. When
`org` is omitted, only official presets are returned. Use
[`get_cloud_object_preset`](#get_cloud_object_preset) with the preset ID from
the results to retrieve the full YAML content.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | no | string | Organization identifier. When provided, includes org-specific presets. |
| `kind` | no | string | PascalCase cloud resource kind to filter by. |
| `search_text` | no | string | Free-text search query to filter presets by name or description. |

---

### get_cloud_object_preset

Get the full content of a cloud object preset by ID. Returns the complete
preset including YAML manifest content, markdown documentation, cloud resource
kind, rank, and provider metadata. Use the YAML content as a template for
[`apply_cloud_resource`](#apply_cloud_resource), replacing placeholder values
as needed.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The preset ID obtained from `search_cloud_object_presets` results. |

---

## Stack Job Command Tools

These tools complete the stack job lifecycle. Use the observability tools
above to monitor, and these command tools to take action.

### rerun_stack_job

Retry a previously executed stack job without re-triggering an apply. Useful
after a transient failure (network timeout, provider outage) when the cloud
resource spec has not changed. The new run uses the same parameters as the
original. Returns the updated stack job. Use
[`get_stack_job`](#get_stack_job) to monitor progress.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The stack job ID to re-run. |

---

### cancel_stack_job

Gracefully cancel a running stack job. The currently executing IaC operation
completes fully before cancellation takes effect — remaining operations are
skipped and marked as cancelled. There is no automatic rollback of completed
operations. The resource lock is released after cancellation, allowing queued
jobs to proceed.

> **Note:** The job must be in `running` status.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The stack job ID to cancel. |

---

### resume_stack_job

Approve and resume a stack job in `awaiting_approval` status. Stack jobs enter
this state when a flow control policy requires manual approval before IaC
execution proceeds. This tool unblocks the job. To reject instead, use
[`cancel_stack_job`](#cancel_stack_job). Returns the updated stack job.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The stack job ID to approve and resume. |

---

### check_stack_job_essentials

Pre-validate that all prerequisites for running a stack job are in place for a
given cloud resource kind and organization. Returns four preflight checks:

| Check | What it verifies |
|-------|-----------------|
| `iac_module` | An IaC module is resolved for this resource kind |
| `backend_credential` | A state backend is configured for the org |
| `flow_control` | An approval policy is resolved |
| `provider_credential` | Cloud provider credentials are available |

Use before `apply_cloud_resource` to catch missing configuration early rather
than discovering failures mid-deployment.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `cloud_resource_kind` | **yes** | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). Read `cloud-resource-kinds://catalog` for valid kinds. |
| `org` | **yes** | string | Organization identifier. |
| `env` | no | string | Environment name. Provide when the resource will be deployed to a specific environment. |

---

## InfraChart Tools

Infra charts are reusable infrastructure-as-code templates that define cloud
resource compositions. The typical workflow is: list charts → get the one you
want → build (preview) with your parameter values → create an infra project.

### list_infra_charts

List available infra chart templates. All parameters are optional — an empty
call returns the first page of all charts.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | no | string | Organization identifier. Scopes results to charts available in this org. |
| `env` | no | string | Environment identifier. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_infra_chart

Retrieve the full details of an infra chart. Returns template YAML files,
`values.yaml`, parameter definitions, description, and web links. Use
[`build_infra_chart`](#build_infra_chart) to preview rendered output before
creating a project.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The infra chart ID from `list_infra_charts` results. |

---

### build_infra_chart

Preview the rendered output of a chart without persisting anything. Fetches
the chart by ID, merges your parameter overrides with the chart defaults, and
returns the rendered YAML and cloud resource DAG. Use this to validate
parameter choices before creating an infra project.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `chart_id` | **yes** | string | The infra chart ID to build. |
| `params` | no | object | Parameter overrides as a name-to-value map. Keys must match parameter names visible in `get_infra_chart` output. Omitted params keep chart defaults. |

---

## InfraProject Tools

Infra projects are deployable infrastructure compositions sourced from infra
charts or Git repositories. They represent a concrete instantiation of a chart
with specific parameter values applied to a target org.

### InfraProject identification

`get_infra_project`, `delete_infra_project`, and `undeploy_infra_project` each
accept two identification paths:

- **By ID**: provide `id` alone.
- **By org + slug**: provide both `org` and `slug` together.

### search_infra_projects

Search infra projects within an organization. Returns lightweight records with
project IDs and metadata. Use
[`get_infra_project`](#get_infra_project) to retrieve full details.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `env` | no | string | Environment slug. When provided, filters to chart-sourced projects deployed to this environment. |
| `search_text` | no | string | Free-text search query. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_infra_project

Retrieve the full details of an infra project by its ID or by org+slug.
Returns metadata, spec (source type, chart or Git config, parameters), and
status (rendered YAML, cloud resource DAG, pipeline ID). The output can be
modified and passed directly to
[`apply_infra_project`](#apply_infra_project).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Infra project ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Project slug. Required with `org` when `id` is not provided. |

---

### apply_infra_project

Create or update an infra project. The operation is idempotent. For new
projects, provide `metadata` (name, org) and `spec` (source type with chart or
Git config). For updates, retrieve the project with
[`get_infra_project`](#get_infra_project), modify the desired fields, and pass
the result here. Returns the applied project.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `infra_project` | **yes** | object | Full InfraProject resource as a JSON object. Must include `metadata` (with `name` and `org`) and `spec`. The output of `get_infra_project` can be passed directly. |

---

### delete_infra_project

Remove an infra project record. This removes the database record only — it
does **not** tear down deployed cloud resources. Use
[`undeploy_infra_project`](#undeploy_infra_project) first to tear down
infrastructure.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Infra project ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Project slug. Required with `org` when `id` is not provided. |

---

### check_infra_project_slug

Check whether an infra project slug is available within an organization.
Returns `true` if no project with the given slug exists. Use before
[`apply_infra_project`](#apply_infra_project) to avoid slug conflicts.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `slug` | **yes** | string | The slug to check for availability. |

---

### undeploy_infra_project

Tear down all cloud resources deployed by an infra project while keeping the
project record. Triggers an undeploy pipeline that destroys the
infrastructure. The project record is preserved and can be redeployed later
via [`apply_infra_project`](#apply_infra_project).

> **WARNING:** This destroys real cloud infrastructure.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Infra project ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Project slug. Required with `org` when `id` is not provided. |

---

## InfraPipeline Tools

Infra pipelines represent deployment runs triggered by apply or run operations
on infra projects. Use these tools to monitor and control the pipeline
lifecycle.

### list_infra_pipelines

List infra pipelines within an organization. Optionally filter by infra
project ID.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `infra_project_id` | no | string | Filter to pipelines for a specific infra project. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_infra_pipeline

Retrieve the full details of an infra pipeline. Returns status, environment
stages, DAG nodes, timestamps, and any errors.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The infra pipeline ID. |

---

### get_latest_infra_pipeline

Retrieve the most recent infra pipeline for a project. This is the primary
tool to check whether an
[`apply_infra_project`](#apply_infra_project) or
[`run_infra_pipeline`](#run_infra_pipeline) call completed successfully.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `infra_project_id` | **yes** | string | The infra project ID to look up the most recent pipeline for. |

---

### run_infra_pipeline

Trigger a new deployment pipeline run for an infra project. For chart-sourced
projects, omit `commit_sha`. For git-repo sourced projects, provide
`commit_sha` to deploy a specific commit. Returns the new pipeline ID. Use
[`get_infra_pipeline`](#get_infra_pipeline) to monitor progress.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `infra_project_id` | **yes** | string | The infra project ID. |
| `commit_sha` | no | string | Git commit SHA to deploy. Required for git-repo sourced projects. Omit for chart-sourced projects. |

---

### cancel_infra_pipeline

Cancel a running infra pipeline. Returns the updated pipeline with its final
status.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The infra pipeline ID to cancel. |

---

### resolve_infra_pipeline_env_gate

Approve or reject a manual gate for an entire deployment environment within a
pipeline. Manual gates pause pipeline execution until explicitly resolved. Use
[`get_infra_pipeline`](#get_infra_pipeline) to inspect which environments have
pending gates.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `infra_pipeline_id` | **yes** | string | The infra pipeline ID. |
| `env` | **yes** | string | Environment name where the gate is pending (e.g. `staging`, `production`). |
| `decision` | **yes** | string | `approve` or `reject`. |

---

### resolve_infra_pipeline_node_gate

Approve or reject a manual gate for a specific DAG node within a pipeline.
DAG nodes represent individual cloud resources in the deployment graph. Use
[`get_infra_pipeline`](#get_infra_pipeline) to inspect which nodes have
pending gates.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `infra_pipeline_id` | **yes** | string | The infra pipeline ID. |
| `env` | **yes** | string | Environment name where the node exists. |
| `node_id` | **yes** | string | Node identifier in the format `CloudResourceKind/slug` (e.g. `KubernetesOpenFga/fga-gcp-dev`). Visible in `get_infra_pipeline` output. |
| `decision` | **yes** | string | `approve` or `reject`. |

---

## Dependency Graph Tools

The graph tools expose the resource topology of your organization. Use them to
understand relationships between resources, plan deployment order, and assess
the impact of changes before making them.

### get_organization_graph

Retrieve the complete resource topology for an organization. Returns all nodes
(organizations, environments, services, cloud resources, credentials, infra
projects) and their relationships. Use this as the starting point for
understanding an infrastructure landscape. Use
[`get_cloud_resource_graph`](#get_cloud_resource_graph) or
[`get_service_graph`](#get_service_graph) to drill into specific resources.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `envs` | no | string[] | Environment slugs to restrict the graph to. When omitted, all environments are included. |
| `node_types` | no | string[] | Node types to include. Valid values: `organization`, `environment`, `service`, `cloud_resource`, `credential`, `infra_project`. When omitted, all types are included. |
| `include_topological_order` | no | boolean | When `true`, includes a topological ordering of node IDs (roots first). Useful for determining deployment order. |
| `max_depth` | no | integer | Maximum traversal depth. `0` or omitted means unlimited. |

---

### get_environment_graph

Retrieve the resource graph scoped to a specific environment. Returns the
environment node, its parent organization, and all resources deployed in the
environment with their relationships.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `env_id` | **yes** | string | Environment identifier. Use `list_environments` to discover available IDs. |
| `node_types` | no | string[] | Node types to include. Same valid values as `get_organization_graph`. |
| `include_topological_order` | no | boolean | When `true`, includes topological ordering of node IDs. |

---

### get_service_graph

Retrieve a service-centric subgraph showing the service and all related
resources. Returns the service node, its cloud resource deployments per
environment, and optionally upstream dependencies and downstream dependents.
Service IDs are discoverable from
[`get_organization_graph`](#get_organization_graph) results.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `service_id` | **yes** | string | Service identifier. |
| `envs` | no | string[] | Environment slugs to restrict results to. |
| `include_upstream` | no | boolean | Include upstream dependencies (what the service depends on). |
| `include_downstream` | no | boolean | Include downstream dependents (what depends on the service). |
| `max_depth` | no | integer | Maximum traversal depth. `0` or omitted means unlimited. |

---

### get_cloud_resource_graph

Retrieve a cloud-resource-centric subgraph. Returns the resource node,
services deployed as it, credentials it uses, and connected nodes. Enable
`include_upstream` and `include_downstream` to traverse beyond immediate
neighbors.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `cloud_resource_id` | **yes** | string | Cloud resource ID. Use `get_cloud_resource` to look up the ID if needed. |
| `include_upstream` | no | boolean | Include upstream dependencies (what this resource depends on). |
| `include_downstream` | no | boolean | Include downstream dependents (what depends on this resource). |
| `max_depth` | no | integer | Maximum traversal depth. `0` or omitted means unlimited. |

---

### get_dependencies

Find all resources that a given resource depends on (upstream traversal). For
example, an EKS cluster might depend on a VPC and an IAM credential. Useful
for understanding deployment prerequisites. Use
[`get_dependents`](#get_dependents) for the reverse direction.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `resource_id` | **yes** | string | Resource ID to find dependencies for. Accepts any resource type. |
| `max_depth` | no | integer | Maximum traversal depth. `0` or omitted means unlimited. |
| `relationship_types` | no | string[] | Relationship types to include. Valid values: `belongs_to_org`, `belongs_to_env`, `deployed_as`, `uses_credential`, `depends_on`, `runs_on`, `managed_by`, `uses`, `service_depends_on`, `owned_by`. When omitted, all types are included. |

---

### get_dependents

Find all resources that depend on a given resource (downstream traversal). For
example, a VPC might have EKS clusters and RDS instances depending on it. Use
before deleting or modifying a resource to understand what might be affected.
For a full blast-radius report with counts and type breakdown, use
[`get_impact_analysis`](#get_impact_analysis) instead.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `resource_id` | **yes** | string | Resource ID to find dependents for. Accepts any resource type. |
| `max_depth` | no | integer | Maximum traversal depth. `0` or omitted means unlimited. |
| `relationship_types` | no | string[] | Relationship types to include. Same valid values as `get_dependencies`. |

---

### get_impact_analysis

Analyze the impact of modifying or deleting a resource. Returns directly
affected resources, transitively affected resources, total affected count, and
a breakdown by type. Use before destructive operations to understand the blast
radius.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `resource_id` | **yes** | string | Resource ID to analyze. Accepts any resource type. |
| `change_type` | no | string | `delete` or `update`. When omitted, general impact analysis is returned. |

---

## Config Manager Tools

Config Manager stores plaintext configuration variables and encrypted secrets.
Both are scoped to either an organization (shared across all environments) or
a specific environment.

### Variables

#### list_variables

List configuration variables in an organization. Optionally filter by
environment.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `env` | no | string | Environment slug to filter by. When omitted, all variables in the org are returned. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

#### get_variable

Retrieve the full details of a configuration variable by its ID or by
org+scope+slug. Variables are uniquely identified within `(org, scope, slug)`.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Variable ID. Provide this alone, or provide all of `org`, `scope`, and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `scope` and `slug` when `id` is not provided. |
| `scope` | conditional | string | `organization` or `environment`. Required with `org` and `slug` when `id` is not provided. |
| `slug` | conditional | string | Variable slug. Required with `org` and `scope` when `id` is not provided. |

---

#### apply_variable

Create or update a configuration variable. The operation is idempotent — if a
variable with the same `(org, scope, slug)` exists it is updated, otherwise it
is created. When `scope` is `environment`, `env` is required. Returns the
applied variable.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `name` | **yes** | string | Display name. Also used to derive the slug on first create. |
| `org` | **yes** | string | Organization identifier. |
| `scope` | **yes** | string | `organization` (shared across all environments) or `environment` (scoped to one env). |
| `env` | conditional | string | Environment slug. Required when `scope` is `environment`. |
| `value` | **yes** | string | The plaintext variable value. |
| `description` | no | string | Human-readable description. |

---

#### delete_variable

Permanently delete a configuration variable.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Variable ID. Provide this alone, or provide all of `org`, `scope`, and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `scope` and `slug` when `id` is not provided. |
| `scope` | conditional | string | `organization` or `environment`. Required with `org` and `slug` when `id` is not provided. |
| `slug` | conditional | string | Variable slug. Required with `org` and `scope` when `id` is not provided. |

---

#### resolve_variable

Look up a variable's current value by slug. Returns only the plain string
value — no metadata or wrapper. Faster than
[`get_variable`](#get_variable) when you only need the value.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `scope` | **yes** | string | `organization` or `environment`. |
| `slug` | **yes** | string | Variable slug within the org and scope. |

---

### Secrets

Secrets are metadata containers for encrypted key-value pairs. The tools here
manage secret metadata only — actual values are managed through secret
versions.

> **Security boundary:** Agents can write secret values via
> `create_secret_version` but cannot read them back. This is intentional.

#### list_secrets

List secrets in an organization. Only metadata is returned — no values are
exposed.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `env` | no | string | Environment slug to filter by. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

#### get_secret

Retrieve secret metadata by its ID or by org+scope+slug. Returns spec (scope,
description, backend) and audit status. No secret values are returned.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Secret ID. Provide this alone, or provide all of `org`, `scope`, and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `scope` and `slug` when `id` is not provided. |
| `scope` | conditional | string | `organization` or `environment`. Required with `org` and `slug` when `id` is not provided. |
| `slug` | conditional | string | Secret slug. Required with `org` and `scope` when `id` is not provided. |

---

#### apply_secret

Create or update a secret's metadata. This manages the secret record only —
use [`create_secret_version`](#create_secret_version) to store actual values.
When `scope` is `environment`, `env` is required. The `backend` cannot be
changed after creation.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `name` | **yes** | string | Display name. Also used to derive the slug on first create. |
| `org` | **yes** | string | Organization identifier. |
| `scope` | **yes** | string | `organization` (shared) or `environment` (scoped to one env). |
| `env` | conditional | string | Environment slug. Required when `scope` is `environment`. |
| `description` | no | string | Human-readable description. |
| `backend` | no | string | Slug of the SecretBackend resource to use for encryption. When omitted, the org's default backend is used. Cannot be changed after creation. |

---

#### delete_secret

Delete a secret and **all its versions**. This permanently destroys the
encrypted data stored in the backend and cannot be undone.

> **WARNING:** Use `list_secret_versions` first to understand what will be
> destroyed.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Secret ID. Provide this alone, or provide all of `org`, `scope`, and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `scope` and `slug` when `id` is not provided. |
| `scope` | conditional | string | `organization` or `environment`. Required with `org` and `slug` when `id` is not provided. |
| `slug` | conditional | string | Secret slug. Required with `org` and `scope` when `id` is not provided. |

---

#### create_secret_version

Store a new version of encrypted key-value data for a secret. Each call
creates an immutable version — previous versions are preserved. The data is
encrypted via envelope encryption and stored in the secret's configured
backend. Use [`apply_secret`](#apply_secret) first to create the parent secret
if it does not exist.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `secret_id` | **yes** | string | The parent secret's ID. Use `get_secret` or `list_secrets` to find it. |
| `data` | **yes** | object | Key-value pairs to store. Values are encrypted. Example: `{"DB_PASSWORD": "s3cret", "API_KEY": "abc123"}`. |

---

#### list_secret_versions

List all versions of a secret. Returns version metadata only (timestamps,
backend version ID) — encrypted data is never returned. Use to understand
version history or to verify that
[`create_secret_version`](#create_secret_version) succeeded.

##### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `secret_id` | **yes** | string | The parent secret's ID. |

---

## Audit Tools

The audit tools provide version history for any platform resource — cloud
resources, infra projects, variables, secrets, and more.

### list_resource_versions

List the version history for a platform resource. Returns a paginated list of
version entries with metadata, event type, and timestamps. Use
[`get_resource_version`](#get_resource_version) with a version ID to retrieve
full details and diffs.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `resource_id` | **yes** | string | The resource ID to retrieve history for. |
| `kind` | **yes** | string | Platform resource kind. Common values: `cloud_resource`, `infra_project`, `infra_chart`, `infra_pipeline`, `variable`, `secret`, `environment`, `organization`, `service`, `stack_job`. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_resource_version

Retrieve a specific resource version with full change details. Returns the
original and new state as YAML, a unified diff, the event type (`create`,
`update`, `delete`), linked stack job ID, and cloud object version details
when applicable. The `context_size` parameter controls how many surrounding
lines appear in the diff, analogous to `git diff -U<n>`.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `version_id` | **yes** | string | The resource version ID from `list_resource_versions` results. |
| `context_size` | no | integer | Surrounding diff lines. Defaults to 3. |

---

### get_resource_version_count

Get the count of versions for a resource without transferring any version
data. Use to quickly check whether a resource has change history, or to
estimate pagination before calling
[`list_resource_versions`](#list_resource_versions).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `resource_id` | **yes** | string | The resource ID to count versions for. |
| `kind` | **yes** | string | Platform resource kind. Same valid values as `list_resource_versions`. |

---

## Catalog Tools

The catalog tools let you browse the types of cloud resources available on the
platform and the IaC modules that provision them.

### search_deployment_components

Browse the deployment component catalog. Deployment components represent the
types of cloud resources that can be provisioned (e.g. `AwsEksCluster`,
`GcpCloudRunService`, `ConfluentKafkaCluster`). Use
[`get_deployment_component`](#get_deployment_component) for full details.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `search_text` | no | string | Free-text search query. |
| `provider` | no | string | Cloud provider to filter by (e.g. `aws`, `gcp`, `azure`, `confluent`, `snowflake`). |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_deployment_component

Retrieve the full details of a deployment component by its ID or by cloud
resource kind. A deployment component defines a type of cloud resource
including its supported IaC modules, provider, and configuration schema.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Component ID. Provide this alone, or provide `kind` alone. |
| `kind` | conditional | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). Provide this alone, or provide `id` alone. |

---

### search_iac_modules

Search for IaC (Infrastructure as Code) modules. IaC modules are the
provisioning implementations that deploy cloud resources — each targets a
specific kind and IaC provisioner (Terraform, Pulumi, or OpenTofu). When `org`
is provided, results include both official platform modules and
organization-specific modules. Use the `kind` filter to find modules that can
provision a specific deployment component (e.g. `kind=AwsEksCluster` returns
all modules capable of deploying EKS clusters).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | no | string | Organization identifier. When provided, includes org-specific modules. |
| `search_text` | no | string | Free-text search query. |
| `kind` | no | string | PascalCase cloud resource kind to filter by. Returns only modules that can provision this type. |
| `provisioner` | no | string | IaC provisioner: `terraform`, `pulumi`, or `tofu`. |
| `provider` | no | string | Cloud provider (e.g. `aws`, `gcp`, `azure`). |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_iac_module

Retrieve the full details of an IaC module. Returns metadata, provisioner
type, cloud resource kind, Git repository URL, version, and parameter schema.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The IaC module ID from `search_iac_modules` results. |

---

## MCP Resources

MCP resources are read-only and accessed via `resources/read` — they are not
tools and do not modify state. Use them to discover kinds and schemas before
issuing write tool calls.

### cloud-resource-kinds://catalog

Returns a JSON object that maps each of the 17 supported cloud providers to its
`api_version` string and a sorted list of PascalCase kind identifiers.

**MIME type:** `application/json`

**Example response shape:**

```json
{
  "aws": {
    "api_version": "ai.planton.provider.aws.v1",
    "kinds": ["AwsEksCluster", "AwsRdsInstance", "AwsVpc", "..."]
  },
  "gcp": {
    "api_version": "ai.planton.provider.gcp.v1",
    "kinds": ["GcpCloudSqlInstance", "GcpGkeCluster", "..."]
  },
  "..."
}
```

Use the `api_version` from this catalog as the `api_version` field in every
`cloud_object` you pass to `apply_cloud_resource`.

### cloud-resource-schema://{kind}

Returns the JSON schema for a specific kind, including all `spec` fields with
their types, validation constraints, and default values.

**URI parameter:** Replace `{kind}` with the exact PascalCase kind string from
the catalog (e.g. `cloud-resource-schema://AwsEksCluster`).

**MIME type:** `application/json`

Use the schema to understand what fields are required in `spec` before
assembling a `cloud_object`.

---

## Service Lifecycle Tools

Services connect Git repositories to Planton Cloud's CI/CD pipeline system. Each service defines where code lives, how to build it, and which environments to deploy to.

### Service identification pattern

`get_service`, `delete_service`, `disconnect_service_git_repo`, `configure_service_webhook`, and `list_service_branches` each accept two identification paths:

- **By ID**: provide `id` alone.
- **By org + slug**: provide both `org` and `slug` together.

### search_services

Search services within an organization. Returns lightweight search records with service IDs, names, and metadata. Use [`get_service`](#get_service) with a service ID from the results to retrieve full details.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. Use `list_organizations` to discover available organizations. |
| `search_text` | no | string | Free-text search query to filter services by name or description. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_service

Retrieve the full details of a service by its ID or by org+slug. Returns the complete service including metadata, spec (Git repo, pipeline config, ingress, deployment targets), and status (per-environment deployment tracking). The output JSON can be modified and passed to [`apply_service`](#apply_service) for updates.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Service ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Service slug. Required with `org` when `id` is not provided. |

---

### apply_service

Create or update a service (idempotent). Accepts the full Service resource as a JSON object. For new services, provide `api_version`, `kind`, `metadata` (name, org), and `spec` (git_repo, pipeline_configuration). For updates, retrieve the service with [`get_service`](#get_service), modify the desired fields, and pass the result here. Returns the applied service with server-assigned ID and audit information.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `service` | **yes** | object | Full Service resource as a JSON object. Must include `api_version` (`service-hub.planton.ai/v1`), `kind` (`Service`), `metadata` (with `name` and `org`), and `spec` (with `git_repo` and `pipeline_configuration`). The output of `get_service` can be modified and passed directly here. |

---

### delete_service

Delete a service record from the platform. Removes the service definition and disconnects webhooks from the Git provider. Does **not** delete deployed cloud resources — those must be removed separately.

> **WARNING:** Removes the service definition and all webhook registrations on the connected Git provider.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Service ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Service slug. Required with `org` when `id` is not provided. |

---

### disconnect_service_git_repo

Remove the webhook from the connected GitHub or GitLab repository. After disconnection, new commits no longer trigger pipelines. The service definition remains in Planton Cloud and can be reconnected later via [`configure_service_webhook`](#configure_service_webhook).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Service ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Service slug. Required with `org` when `id` is not provided. |

---

### configure_service_webhook

Create or refresh the webhook on GitHub or GitLab for a service. Registers (or updates) the webhook so that push and pull_request events trigger pipelines. Useful for recovering from accidentally deleted webhooks or troubleshooting webhook delivery issues.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Service ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Service slug. Required with `org` when `id` is not provided. |

---

### list_service_branches

List all Git branches in the repository connected to a service. Uses the GitHub or GitLab API to fetch branch information. Useful for selecting a branch before triggering a pipeline run or validating that a branch exists.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Service ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Service slug. Required with `org` when `id` is not provided. |

---

## Service CI/CD Pipeline Tools

Service pipelines represent build-and-deploy runs triggered by Git pushes or manual actions. Use these tools to monitor runs, trigger new runs, and manage the pipeline YAML files stored in the service repository.

### list_pipelines

List CI/CD pipelines within an organization. Optionally filter by service ID and environment names. Returns a paginated list. Use [`get_pipeline`](#get_pipeline) with a pipeline ID from the results to retrieve full details.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `service_id` | no | string | Service ID to filter pipelines by. When omitted, all pipelines in the org are returned. |
| `envs` | no | string[] | Environment names to filter by. Returns pipelines where any deployment task targets one of these environments. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_pipeline

Retrieve the full details of a CI/CD pipeline by its ID. Returns the complete pipeline including status, build stage, deployment tasks, timestamps, and any errors.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The pipeline ID. |

---

### get_last_pipeline

Retrieve the most recent CI/CD pipeline for a service. This is the primary tool to check whether a [`run_pipeline`](#run_pipeline) operation completed successfully. Returns the full pipeline including status, build stage, deployment tasks, timestamps, and any errors.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `service_id` | **yes** | string | The service ID to look up the most recent pipeline for. |

---

### run_pipeline

Trigger a new CI/CD pipeline run for a service on a specific Git branch. Optionally provide a commit SHA to deploy that exact commit; when omitted, the branch HEAD is used. Use [`get_last_pipeline`](#get_last_pipeline) with the service ID to monitor the triggered pipeline.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `service_id` | **yes** | string | The service ID to run the pipeline for. |
| `branch` | **yes** | string | Git branch name to run the pipeline on. Use `list_service_branches` to discover available branches. |
| `commit_sha` | no | string | Git commit SHA to deploy. When omitted, the branch HEAD is used. |

---

### rerun_pipeline

Re-run a previously executed CI/CD pipeline using the same service, branch, and commit configuration as the original. Useful for retrying pipelines that failed due to transient issues. Returns the newly created pipeline.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The pipeline ID to re-run. |

---

### cancel_pipeline

Cancel a running CI/CD pipeline. During the build stage, Tekton PipelineRun resources are deleted and running build pods are terminated. During the deploy stage, the current deployment task receives a cancellation signal and remaining tasks are skipped. Returns the updated pipeline with its final status.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | **yes** | string | The pipeline ID to cancel. |

---

### resolve_pipeline_gate

Approve or reject a manual gate for a deployment task within a CI/CD pipeline. Manual gates pause pipeline execution at a specific deployment task until explicitly resolved. Use [`get_pipeline`](#get_pipeline) to inspect which deployment tasks have pending gates.

> **WARNING:** Approving a gate for a production deployment task will deploy to production.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `pipeline_id` | **yes** | string | The pipeline ID containing the manual gate. |
| `deployment_task_name` | **yes** | string | Name of the deployment task with the pending gate. Visible in `get_pipeline` output. |
| `decision` | **yes** | string | `approve` or `reject`. |

---

### list_pipeline_files

List Tekton pipeline YAML files in the Git repository connected to a service. Discovers pipeline YAMLs under Planton conventions (`.planton/`, `.tekton/`, `tekton/`) and returns their paths, content, and blob SHAs. Use this before making changes with [`update_pipeline_file`](#update_pipeline_file).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `service_id` | **yes** | string | The service ID whose repository will be scanned. |
| `branch` | no | string | Git branch to scan. When omitted, the service's default branch is used. |

---

### update_pipeline_file

Create or update a pipeline YAML file in the Git repository connected to a service. Commits the change directly to the specified branch (or the service's default branch). Provide `expected_base_sha` from [`list_pipeline_files`](#list_pipeline_files) for optimistic locking — the write is rejected if the file has been modified since. Returns the new blob SHA, commit SHA, and branch name.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `service_id` | **yes** | string | The service ID whose repository will be updated. |
| `path` | **yes** | string | Path to write relative to the repository root (e.g. `.planton/pipeline.yaml`). |
| `content` | **yes** | string | New file content (plain text, typically YAML). |
| `expected_base_sha` | no | string | Current blob SHA from `list_pipeline_files`. When set, the write is rejected if the file has changed since. Use for safe concurrent editing. |
| `commit_message` | no | string | Custom commit message. When omitted, a default message is generated. |
| `branch` | no | string | Target branch. When omitted, the service's default branch is used. |

---

## Service Variables Group Tools

Variables groups are named collections of plaintext environment variables. Services reference entries in a group using `$variables-group/<group-name>/<entry-name>` references. Use these tools to manage groups, individual entries, and resolve references.

### Variables group identification pattern

`get_variables_group`, `delete_variables_group` accept two identification paths:

- **By ID**: provide `id` alone.
- **By org + slug**: provide both `org` and `slug` together.

For entry-level tools (`upsert_variable`, `delete_variable`), the group can be identified by `group_id` alone, or by both `org` and `group_slug`.

### search_variables

Search variable entries across all variables groups within an organization. Returns individual variable entries (not whole groups) with their group context — useful for finding where a specific variable like `DATABASE_HOST` is defined. Each result includes the group name, group ID, variable name, value, and description.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `search_text` | no | string | Free-text search query to filter variable entries by name, value, or description. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_variables_group

Retrieve the full details of a variables group by its ID or by org+slug. Returns the complete group including metadata, all variable entries (names, values, descriptions), and audit status. The output JSON can be modified and passed to [`apply_variables_group`](#apply_variables_group) for updates.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Variables group ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Group slug. Required with `org` when `id` is not provided. |

---

### apply_variables_group

Create or update a variables group (idempotent). Accepts the full VariablesGroup resource as a JSON object. For new groups, provide `api_version`, `kind`, `metadata` (name, org), and `spec` (description, entries). For updates, retrieve the group with [`get_variables_group`](#get_variables_group), modify the desired fields, and pass here.

> **WARNING:** This replaces **all** entries in the group. To modify a single variable without affecting others, use [`upsert_variable`](#upsert_variable) instead.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `variables_group` | **yes** | object | Full VariablesGroup resource as a JSON object. Must include `api_version` (`service-hub.planton.ai/v1`), `kind` (`VariablesGroup`), `metadata` (with `name` and `org`), and `spec` (with optional `description` and `entries`). |

---

### delete_variables_group

Delete a variables group from the platform.

> **WARNING:** Ensure no services reference this group before deleting — services using `$variables-group/` references to this group will fail during deployment.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Variables group ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Group slug. Required with `org` when `id` is not provided. |

---

### upsert_variable

Add or update a single variable in a variables group. If a variable with the same name already exists it is updated; otherwise it is added. This is safer than [`apply_variables_group`](#apply_variables_group) when modifying a single variable because it does not affect other entries. Returns the updated variables group.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `group_id` | conditional | string | Variables group ID. Provide this alone, or provide both `org` and `group_slug`. |
| `org` | conditional | string | Organization identifier. Required with `group_slug` when `group_id` is not provided. |
| `group_slug` | conditional | string | Group slug. Required with `org` when `group_id` is not provided. |
| `entry` | **yes** | object | Variable entry with `name` (required), optional `description`, and either `value` (a literal string) or `value_from` (a reference object). |

---

### <a name="servicehub-delete_variable"></a>delete_variable (Variables Group)

Remove a single variable entry from a variables group. Other variables in the group remain unchanged. Identify the target group by `group_id` or by `org`+`group_slug`.

> **WARNING:** Services referencing this variable via `$variables-group/` will fail during deployment.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `group_id` | conditional | string | Variables group ID. Provide this alone, or provide both `org` and `group_slug`. |
| `org` | conditional | string | Organization identifier. Required with `group_slug` when `group_id` is not provided. |
| `group_slug` | conditional | string | Group slug. Required with `org` when `group_id` is not provided. |
| `entry_name` | **yes** | string | Name of the variable to remove. |

---

### get_variable_value

Retrieve the resolved value of a specific variable from a variables group. If the variable uses a `value_from` reference, the reference is resolved to its current value. Returns the plain string value without metadata.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `group_name` | **yes** | string | Name (slug) of the variables group. |
| `entry_name` | **yes** | string | Name of the variable whose value to retrieve. |

---

### transform_variables

Batch-resolve `$variables-group/` references in a set of environment variables. Accepts a map of key-value pairs where values may contain `$variables-group/<group-name>/<entry-name>` references. References are resolved to their actual values; literal values pass through unchanged. Returns two maps: successfully transformed entries and any entries that failed resolution with error messages. Useful for debugging configuration issues or previewing resolved values before deployment.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. Determines which variables groups are available for reference resolution. |
| `entries` | **yes** | object | Map of environment variable names to values. Values starting with `$variables-group/` are resolved; all others pass through unchanged. |

---

## Service Secrets Group Tools

Secrets groups are named collections of encrypted secrets. Services reference entries using `$secrets-group/<group-name>/<entry-name>` references. Use these tools to manage groups, individual entries, and resolve references.

### Secrets group identification pattern

`get_secrets_group`, `delete_secrets_group` accept two identification paths:

- **By ID**: provide `id` alone.
- **By org + slug**: provide both `org` and `slug` together.

For entry-level tools (`upsert_secret`, `delete_secret`), the group can be identified by `group_id` alone, or by both `org` and `group_slug`.

### search_secrets

Search secret entries across all secrets groups within an organization. Returns individual secret entries (not whole groups) with their group context — useful for finding where a specific secret like `DB_PASSWORD` is defined. Each result includes the group name, group ID, secret name, and description. Secret values are never returned.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `search_text` | no | string | Free-text search query to filter secret entries by name or description. |
| `page_num` | no | integer | Page number (1-based). Defaults to 1. |
| `page_size` | no | integer | Results per page. Defaults to 20. |

---

### get_secrets_group

Retrieve the full details of a secrets group by its ID or by org+slug. Returns the complete group including metadata, all secret entry names and descriptions (no values), and audit status. The output JSON can be modified and passed to [`apply_secrets_group`](#apply_secrets_group) for updates.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Secrets group ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Group slug. Required with `org` when `id` is not provided. |

---

### apply_secrets_group

Create or update a secrets group (idempotent). Accepts the full SecretsGroup resource as a JSON object. For new groups, provide `api_version`, `kind`, `metadata` (name, org), and `spec` (description, entries). For updates, retrieve the group with [`get_secrets_group`](#get_secrets_group), modify the desired fields, and pass here.

> **WARNING:** This replaces **all** entries in the group. To modify a single secret without affecting others, use [`upsert_secret`](#upsert_secret) instead.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `secrets_group` | **yes** | object | Full SecretsGroup resource as a JSON object. Must include `api_version` (`service-hub.planton.ai/v1`), `kind` (`SecretsGroup`), `metadata` (with `name` and `org`), and `spec` (with optional `description` and `entries`). |

---

### delete_secrets_group

Delete a secrets group from the platform.

> **WARNING:** Ensure no services reference this group before deleting — services using `$secrets-group/` references to this group will fail during deployment.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Secrets group ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Group slug. Required with `org` when `id` is not provided. |

---

### upsert_secret

Add or update a single secret in a secrets group. If a secret with the same name already exists it is updated; otherwise it is added. This is safer than [`apply_secrets_group`](#apply_secrets_group) when modifying a single secret because it does not affect other entries. Returns the updated secrets group.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `group_id` | conditional | string | Secrets group ID. Provide this alone, or provide both `org` and `group_slug`. |
| `org` | conditional | string | Organization identifier. Required with `group_slug` when `group_id` is not provided. |
| `group_slug` | conditional | string | Group slug. Required with `org` when `group_id` is not provided. |
| `entry` | **yes** | object | Secret entry with `name` (required), optional `description`, and either `value` (a literal string) or `value_from` (a reference object). |

---

### <a name="servicehub-delete_secret"></a>delete_secret (Secrets Group)

Remove a single secret entry from a secrets group. Other secrets in the group remain unchanged.

> **WARNING:** Services referencing this secret via `$secrets-group/` will fail during deployment.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `group_id` | conditional | string | Secrets group ID. Provide this alone, or provide both `org` and `group_slug`. |
| `org` | conditional | string | Organization identifier. Required with `group_slug` when `group_id` is not provided. |
| `group_slug` | conditional | string | Group slug. Required with `org` when `group_id` is not provided. |
| `entry_name` | **yes** | string | Name of the secret to remove. |

---

### get_secret_value

Retrieve the resolved value of a specific secret from a secrets group. If the secret uses a `value_from` reference, it is resolved to its current value. Returns the plain string value.

> **WARNING:** This returns the secret value in **PLAINTEXT**. Only use when the user explicitly requests to see a secret value. Never log or display secret values unless specifically asked.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `group_name` | **yes** | string | Name (slug) of the secrets group. |
| `entry_name` | **yes** | string | Name of the secret whose value to retrieve. |

---

### transform_secrets

Batch-resolve `$secrets-group/` references in a set of environment variables. Accepts a map of key-value pairs where values may contain `$secrets-group/<group-name>/<entry-name>` references. References are resolved to their actual values; literal values pass through unchanged. Returns two maps: successfully transformed entries and any entries that failed resolution.

> **WARNING:** Resolved values are returned in **PLAINTEXT**. Use only when debugging or when the user explicitly requests resolved secret values.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `org` | **yes** | string | Organization identifier. |
| `entries` | **yes** | object | Map of environment variable names to values. Values starting with `$secrets-group/` are resolved; all others pass through unchanged. |

---

## DNS Domain Tools

DNS domains register custom domain names within an organization, making them available for services to use in their ingress configuration.

### get_dns_domain

Retrieve the full details of a DNS domain by its ID or by org+slug. Returns the complete domain including metadata, spec (domain_name, description), and audit status. The output JSON can be modified and passed to [`apply_dns_domain`](#apply_dns_domain) for updates.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | DNS domain ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Domain slug. Required with `org` when `id` is not provided. |

---

### apply_dns_domain

Create or update a DNS domain (idempotent). For new domains, provide `api_version`, `kind`, `metadata` (name, org), and `spec` (domain_name). For updates, retrieve the domain with [`get_dns_domain`](#get_dns_domain), modify the desired fields, and pass the result here. Returns the applied domain with server-assigned ID and audit information.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `dns_domain` | **yes** | object | Full DnsDomain resource as a JSON object. Must include `api_version` (`service-hub.planton.ai/v1`), `kind` (`DnsDomain`), `metadata` (with `name` and `org`), and `spec` (with `domain_name`). |

---

### delete_dns_domain

Delete a DNS domain record from the platform.

> **WARNING:** Services using this domain for ingress hostnames will lose their custom domain routing. Ensure no services reference this domain before deleting.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | DNS domain ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Domain slug. Required with `org` when `id` is not provided. |

---

## Tekton Pipeline Tools

Tekton pipelines are reusable CI/CD pipeline definitions that services reference to orchestrate build and deployment stages.

### get_tekton_pipeline

Retrieve the full details of a Tekton pipeline template by its ID or by org+slug. Returns the complete pipeline including metadata, spec (description, git_repo, yaml_content, overview_markdown), and audit status. The output JSON can be modified and passed to [`apply_tekton_pipeline`](#apply_tekton_pipeline) for updates.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Tekton pipeline ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Pipeline slug (name). Required with `org` when `id` is not provided. |

---

### apply_tekton_pipeline

Create or update a Tekton pipeline template (idempotent). For new pipelines, provide `api_version`, `kind`, `metadata` (name, org), and `spec` (selector, and optionally description, git_repo, yaml_content). For updates, retrieve the pipeline with [`get_tekton_pipeline`](#get_tekton_pipeline), modify the desired fields, and pass here.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `tekton_pipeline` | **yes** | object | Full TektonPipeline resource as a JSON object. Must include `api_version` (`service-hub.planton.ai/v1`), `kind` (`TektonPipeline`), `metadata` (with `name` and `org`), and `spec` (with `selector` and optionally `description`, `git_repo`, `yaml_content`, `overview_markdown`). |

---

### delete_tekton_pipeline

Delete a Tekton pipeline template from the platform.

> **WARNING:** Services referencing this pipeline will lose their CI/CD pipeline definition. Ensure no services reference this pipeline before deleting.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Tekton pipeline ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Pipeline slug (name). Required with `org` when `id` is not provided. |

---

## Tekton Task Tools

Tekton tasks are reusable CI/CD step definitions (e.g., git-clone, docker-build, deploy) that Tekton pipelines reference as individual build steps.

### get_tekton_task

Retrieve the full details of a Tekton task template by its ID or by org+slug. Returns the complete task including metadata, spec (description, git_repo, yaml_content, overview_markdown), and audit status. The output JSON can be modified and passed to [`apply_tekton_task`](#apply_tekton_task) for updates.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Tekton task ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Task slug (name). Required with `org` when `id` is not provided. |

---

### apply_tekton_task

Create or update a Tekton task template (idempotent). For new tasks, provide `api_version`, `kind`, `metadata` (name, org), and `spec` (selector, and optionally description, git_repo, yaml_content). For updates, retrieve the task with [`get_tekton_task`](#get_tekton_task), modify the desired fields, and pass here.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `tekton_task` | **yes** | object | Full TektonTask resource as a JSON object. Must include `api_version` (`service-hub.planton.ai/v1`), `kind` (`TektonTask`), `metadata` (with `name` and `org`), and `spec` (with `selector` and optionally `description`, `git_repo`, `yaml_content`, `overview_markdown`). |

---

### delete_tekton_task

Delete a Tekton task template from the platform.

> **WARNING:** Tekton pipelines referencing this task will fail during execution. Ensure no pipelines reference this task before deleting.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | Tekton task ID. Provide this alone, or provide both `org` and `slug`. |
| `org` | conditional | string | Organization identifier. Required with `slug` when `id` is not provided. |
| `slug` | conditional | string | Task slug (name). Required with `org` when `id` is not provided. |

---

## Error Handling

All tools translate gRPC errors from the Planton backend into user-facing
messages. The raw gRPC status code is never surfaced directly.

| gRPC Status | What It Means | What to Do |
|-------------|---------------|------------|
| `NotFound` | The resource does not exist | Verify the `id` or `(kind, org, env, slug)` coordinates |
| `PermissionDenied` | The API key lacks permission for this operation | Check API key permissions in the Planton Console |
| `Unauthenticated` | The API key is invalid or missing | Verify `PLANTON_API_KEY` is set correctly |
| `Unavailable` | The Planton backend is unreachable | Check connectivity and that the backend is running |
| `InvalidArgument` | The request failed schema validation | The server's validation message is returned directly — fix the field it identifies |

---

## Agent Cheat Sheet

A quick decision guide for agents working with this server.

### Getting started — discover your context

```
Don't know the org?          →  list_organizations
Don't know the environment?  →  list_environments (pass org)
Don't know available kinds?  →  read cloud-resource-kinds://catalog
```

### Creating a resource

```
1. list_organizations             →  get your org
2. list_environments (org)        →  get your env
3. read cloud-resource-kinds://catalog  →  find the kind + api_version
4. read cloud-resource-schema://{kind}  →  get the spec definition
5. (optional) search_cloud_object_presets (kind)  →  find a preset template
6. (optional) get_cloud_object_preset (id)        →  use preset as starting point
7. (optional) check_slug_availability (org, env, kind, slug)
8. apply_cloud_resource (cloud_object)
9. get_latest_stack_job (cloud_resource_id)  →  verify success
```

### Reading or modifying a resource

```
Have an ID?                    →  pass id to get/delete/destroy/rename/locks tools
Have org + env + kind + slug?  →  pass all four as coordinates
Have a partial set?            →  use list_cloud_resources to find the resource first
```

### Destroy vs Delete

```
Want to tear down infrastructure?  →  destroy_cloud_resource (infra gone, record stays)
Want to remove the record?         →  delete_cloud_resource  (record gone)
Full cleanup?                      →  destroy first, then delete
```

### Monitoring operations

```
Just called apply or destroy?
  ├─ Have the resource ID?    →  get_latest_stack_job (cloud_resource_id)
  └─ Have a stack job ID?     →  get_stack_job (id)

Need to audit history?        →  list_stack_jobs (org, filters)
```

### Troubleshooting locks

```
Resource stuck / can't apply?
  1. list_cloud_resource_locks   →  check who holds the lock
  2. get_latest_stack_job        →  verify no jobs are running
  3. remove_cloud_resource_locks →  clear stuck locks (only if safe)
```

### Which identification strategy?

| Situation | Use |
|-----------|-----|
| You just called `apply_cloud_resource` and have the returned `id` | `id` alone |
| You know the resource coordinates but not the system ID | `kind` + `org` + `env` + `slug` together |
| You have a partial set of coordinates | Use `list_cloud_resources` to find the resource |

### Minimal apply_cloud_resource call shape

```json
{
  "cloud_object": {
    "api_version": "<from catalog>",
    "kind": "<PascalCase kind>",
    "metadata": {
      "name": "<display name>",
      "org": "<org id>",
      "env": "<env id>"
    },
    "spec": {
      "<fields from schema>"
    }
  }
}
```

### Pre-flight check before deploying

```
Want to verify prerequisites before apply_cloud_resource?
  →  check_stack_job_essentials (cloud_resource_kind, org)
     Checks: iac_module, backend_credential, flow_control, provider_credential
```

### Stack job stuck or failed

```
Job failed?        →  rerun_stack_job (id)
Job running?       →  cancel_stack_job (id)
Awaiting approval? →  resume_stack_job (id)  or  cancel_stack_job (id) to reject
```

### Working with infra charts and projects

```
Find a chart template:
  1. list_infra_charts (org)              →  browse available templates
  2. get_infra_chart (id)                 →  see params and template YAML
  3. build_infra_chart (chart_id, params) →  preview rendered output (no-op)
  4. apply_infra_project (infra_project)  →  create the project from the chart
  5. get_latest_infra_pipeline (infra_project_id)  →  monitor deployment

Manage an existing project:
  get_infra_project (id or org+slug)      →  retrieve current state
  apply_infra_project (modified object)   →  update
  undeploy_infra_project (id or org+slug) →  tear down infra, keep record
  delete_infra_project (id or org+slug)   →  remove record (undeploy first)
```

### Understanding resource relationships

```
What does this resource depend on?     →  get_dependencies (resource_id)
What depends on this resource?         →  get_dependents (resource_id)
What breaks if I delete/change this?   →  get_impact_analysis (resource_id, change_type)
Full org topology:                     →  get_organization_graph (org)
Just one environment:                  →  get_environment_graph (env_id)
```

### Config variables and secrets

```
Variables (plaintext):
  list_variables (org)                   →  see all variables
  resolve_variable (org, scope, slug)    →  quick value lookup
  apply_variable (name, org, scope, value)  →  create or update

Secrets (encrypted):
  list_secrets (org)                     →  see metadata only, never values
  apply_secret (name, org, scope)        →  create the secret container
  create_secret_version (secret_id, data) →  store encrypted key-value pairs
  list_secret_versions (secret_id)       →  see version history
```

### Auditing changes

```
How many times was this changed?  →  get_resource_version_count (resource_id, kind)
Full change history:              →  list_resource_versions (resource_id, kind)
What exactly changed in v3?       →  get_resource_version (version_id)
```

### Working with services

```
Find a service:
  search_services (org)              →  find by name/description
  get_service (id or org+slug)       →  retrieve full spec and status

Manage Git integration:
  list_service_branches (id or org+slug)           →  discover available branches
  configure_service_webhook (id or org+slug)        →  register or refresh webhook
  disconnect_service_git_repo (id or org+slug)      →  pause CI/CD automation

Create or update a service:
  get_service (id or org+slug)       →  retrieve current spec
  apply_service (service)            →  create or update (idempotent)
  delete_service (id or org+slug)    →  remove record (does NOT remove deployed resources)
```

### Triggering and monitoring service pipelines

```
Trigger a run:
  list_service_branches (service)      →  confirm target branch exists
  run_pipeline (service_id, branch)    →  trigger; returns immediately
  get_last_pipeline (service_id)       →  poll for status and result

Inspect a specific run:
  get_pipeline (id)                    →  full stage-by-stage details

Retry or stop:
  rerun_pipeline (id)                  →  retry same branch/commit
  cancel_pipeline (id)                 →  stop in-progress run

Resolve a manual gate:
  get_pipeline (id)                         →  find deployment tasks with pending gates
  resolve_pipeline_gate (pipeline_id, deployment_task_name, decision)
```

### Managing pipeline YAML files in a service repository

```
list_pipeline_files (service_id)                           →  see current files + blob SHAs
update_pipeline_file (service_id, path, content)           →  create or update a file
update_pipeline_file (..., expected_base_sha)              →  with optimistic locking
```

### Managing variables and secrets groups

```
Variables (plaintext, use $variables-group/ references):
  search_variables (org)                          →  find entries across all groups
  get_variables_group (id or org+slug)            →  see full group with all entries
  apply_variables_group (variables_group)         →  replace entire group (idempotent)
  upsert_variable (group_id or org+group_slug, entry)   →  add/update single entry safely
  delete_variable (group_id or org+group_slug, entry_name)  →  remove single entry
  get_variable_value (org, group_name, entry_name)         →  quick value lookup
  transform_variables (org, entries)              →  batch-resolve references

Secrets (encrypted, use $secrets-group/ references):
  search_secrets (org)                            →  find entries across all groups
  get_secrets_group (id or org+slug)              →  see full group (no values returned)
  apply_secrets_group (secrets_group)             →  replace entire group (idempotent)
  upsert_secret (group_id or org+group_slug, entry)      →  add/update single entry safely
  delete_secret (group_id or org+group_slug, entry_name)  →  remove single entry
  get_secret_value (org, group_name, entry_name)          →  plaintext value (use sparingly)
  transform_secrets (org, entries)                →  batch-resolve references (plaintext output)
```

### Managing DNS domains

```
get_dns_domain (id or org+slug)      →  view domain spec
apply_dns_domain (dns_domain)        →  create or update
delete_dns_domain (id or org+slug)   →  remove (verify no services use it first)
```

### Managing Tekton pipeline and task templates

```
Pipelines (orchestrate build + deploy stages):
  get_tekton_pipeline (id or org+slug)           →  view pipeline spec
  apply_tekton_pipeline (tekton_pipeline)        →  create or update
  delete_tekton_pipeline (id or org+slug)        →  remove (verify no services reference it)

Tasks (individual build steps within a pipeline):
  get_tekton_task (id or org+slug)               →  view task spec
  apply_tekton_task (tekton_task)                →  create or update
  delete_tekton_task (id or org+slug)            →  remove (verify no pipelines reference it)
```
