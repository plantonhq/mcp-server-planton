# Tools & Resources Reference

Complete reference for all MCP tools and resources exposed by the Planton Cloud
MCP Server. For installation and configuration, see the
[README](../README.md).

---

## Overview

| Tool | Operation | Description |
|------|-----------|-------------|
| [`apply_cloud_resource`](#apply_cloud_resource) | Write | Create or update a cloud resource (idempotent) |
| [`get_cloud_resource`](#get_cloud_resource) | Read | Retrieve a cloud resource by ID or by coordinates |
| [`delete_cloud_resource`](#delete_cloud_resource) | Write | Delete a cloud resource by ID or by coordinates |

| Resource URI | Description |
|--------------|-------------|
| [`cloud-resource-kinds://catalog`](#cloud-resource-kindscatalog) | Catalog of all 362 supported kinds grouped by 17 cloud providers |
| [`cloud-resource-schema://{kind}`](#cloud-resource-schemakind) | Full JSON schema for a specific kind |

---

## Tools

### apply_cloud_resource

Create or update a cloud resource on the Planton platform. The operation is
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

Retrieve a cloud resource from the Planton platform. Two identification
strategies are supported — use whichever you have available:

- **By ID**: provide `id` alone.
- **By coordinates**: provide all four of `kind`, `org`, `env`, and `slug`
  together.

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | System-assigned resource ID. Provide this alone, or provide all of `kind`, `org`, `env`, `slug`. |
| `kind` | conditional | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | conditional | string | Organization identifier. |
| `env` | conditional | string | Environment identifier. |
| `slug` | conditional | string | Immutable unique resource name within the `(org, env, kind)` scope. |

Either `id` alone, or all four of `kind` + `org` + `env` + `slug`, must be
provided. Mixing a partial set of coordinates without `id` will return a
validation error.

---

### delete_cloud_resource

Delete a cloud resource from the Planton platform. Accepts the same
identification strategies as [`get_cloud_resource`](#get_cloud_resource).

#### Parameters

| Parameter | Required | Type | Description |
|-----------|----------|------|-------------|
| `id` | conditional | string | System-assigned resource ID. Provide this alone, or provide all of `kind`, `org`, `env`, `slug`. |
| `kind` | conditional | string | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | conditional | string | Organization identifier. |
| `env` | conditional | string | Environment identifier. |
| `slug` | conditional | string | Immutable unique resource name within the `(org, env, kind)` scope. |

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

### Which resource to read first?

```
Need to create or update a resource?
  ├─ Don't know the kind catalog?  →  read cloud-resource-kinds://catalog
  ├─ Know the kind, need the spec? →  read cloud-resource-schema://{kind}
  └─ Have both?                    →  call apply_cloud_resource

Need to read or delete a resource?
  ├─ Have an ID?                   →  pass id to get_cloud_resource or delete_cloud_resource
  └─ Have org + env + slug + kind? →  pass all four as coordinates
```

### Which identification strategy for get/delete?

| Situation | Use |
|-----------|-----|
| You just called `apply_cloud_resource` and have the returned `id` | `id` alone |
| You know the resource coordinates but not the system ID | `kind` + `org` + `env` + `slug` together |
| You have a partial set of coordinates | Fetch the catalog first to confirm the full coordinates |

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
