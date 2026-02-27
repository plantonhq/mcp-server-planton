---
name: Cloud Resource Kinds Catalog
overview: Add a static MCP resource that serves a catalog of all 362 supported cloud resource kinds, grouped by cloud provider, enabling agents to discover valid kind values before using the schema template or calling tools.
todos:
  - id: catalog-data
    content: Add buildKindCatalog() to schema.go -- transforms loadRegistry() output into grouped JSON, cached with sync.Once
    status: completed
  - id: catalog-resource
    content: Add KindCatalogResource() and KindCatalogHandler() to resources.go -- static MCP resource returning cached catalog JSON
    status: completed
  - id: register-catalog
    content: Register the static catalog resource in server.go registerResources() alongside existing schema template
    status: completed
  - id: update-tool-descriptions
    content: Update apply_cloud_resource tool/handler descriptions to reference catalog resource as the discovery starting point
    status: completed
  - id: verify-build
    content: Run go build, go vet, verify catalog JSON output is well-formed with correct grouping and counts
    status: completed
isProject: false
---

# Cloud Resource Kinds Catalog Resource

## Problem

Agents cannot discover valid kind values for `cloud-resource-schema://{kind}` or `apply_cloud_resource`. The resource template answers "what does this kind look like?" but nothing answers "what kinds exist?"

## Kind Case Analysis (Resolved)

The entire system consistently uses PascalCase for kind values across all layers (proto enum, Go stubs, codegen registry, JSON schemas, StringValueOrRef options). No case conversion is needed at any boundary. The 362 production kinds in the codegen registry are a perfect subset of the proto enum (only 3 test fixtures are excluded). No action needed here.

## Design: Static MCP Resource

**URI**: `cloud-resource-kinds://catalog`

Rationale: Follows the same `{concept}://{identifier}` pattern as `cloud-resource-schema://{kind}`. Static (not parameterized) because there's only one catalog.

**Response format** (JSON):

```json
{
  "schema_uri_template": "cloud-resource-schema://{kind}",
  "total_kinds": 362,
  "providers": {
    "aws": {
      "api_version": "aws.openmcf.org/v1",
      "kinds": ["AwsAlb", "AwsEksCluster", "AwsVpc", "..."]
    },
    "gcp": {
      "api_version": "gcp.openmcf.org/v1",
      "kinds": ["GcpGkeCluster", "GcpVpc", "..."]
    }
  }
}
```

Design choices:

- **Grouped by cloud provider**: Agents can narrow by provider first (e.g., user says "deploy to AWS")
- `**api_version` per provider group**: Agent needs this for `cloud_object.api_version` -- getting it here avoids a second lookup
- `**schema_uri_template` at top level**: Tells the agent exactly how to fetch schema for a chosen kind, documented once rather than per-entry
- **Sorted kind lists**: Deterministic output, easier for agents to scan
- **No per-kind schema URIs**: Redundant -- the agent has the template and the kind string, it can construct the URI trivially

## Implementation

### Files to modify

1. **[internal/domains/cloudresource/resources.go](internal/domains/cloudresource/resources.go)** -- Add `KindCatalogResource()` returning `*mcp.Resource` and `KindCatalogHandler()` returning `mcp.ResourceHandler`
2. **[internal/domains/cloudresource/schema.go](internal/domains/cloudresource/schema.go)** -- Add a `buildKindCatalog()` function that transforms `loadRegistry()` data into the grouped JSON structure, cached with `sync.Once`
3. **[internal/server/server.go](internal/server/server.go)** -- Register the static resource in `registerResources()` alongside the existing template

### Data flow

The registry data is already loaded and cached via `loadRegistry()` in `schema.go`. The new `buildKindCatalog()` function:

- Reads from `loadRegistry()` (already cached with `sync.Once`)
- Groups entries by `cloudProvider`
- Extracts `apiVersion` per provider group (all entries within a provider share the same `apiVersion`)
- Sorts kind lists alphabetically
- Marshals to JSON once via its own `sync.Once`
- Returns cached bytes on every subsequent call

### Registration pattern

```go
// In server.go registerResources():
srv.AddResource(cloudresource.KindCatalogResource(), cloudresource.KindCatalogHandler())
srv.AddResourceTemplate(cloudresource.SchemaTemplate(), cloudresource.SchemaHandler())
```

Uses `srv.AddResource()` (static) vs `srv.AddResourceTemplate()` (parameterized) -- both use the same `ResourceHandler` signature in the MCP SDK.

## What this does NOT change

- No changes to the codegen pipeline
- No changes to tools (`apply_cloud_resource` behavior unchanged)
- No changes to the schema template (`cloud-resource-schema://{kind}` unchanged)
- No changes to kind case conventions (PascalCase remains canonical everywhere)
- No new dependencies

## Agent workflow after this change

```
Agent                          MCP Server
  |                                |
  |-- resources/list ------------->|  Discovers catalog + schema template
  |<-- cloud-resource-kinds://catalog
  |    cloud-resource-schema://{kind}
  |                                |
  |-- read cloud-resource-kinds://catalog -->|  Gets all 362 kinds by provider
  |<-- { providers: { aws: { kinds: [...] } } }
  |                                |
  |-- read cloud-resource-schema://AwsEksCluster -->|  Gets spec schema
  |<-- { properties: { ... } }     |
  |                                |
  |-- call apply_cloud_resource -->|  Creates resource with validated input
  |<-- { status: "applied" }       |
```

