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
