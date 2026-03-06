---
name: T12 MCP Resources
overview: Add a single, curated MCP resource -- `api-resource-kinds://catalog` -- that serves as the agent's navigational index of all platform resource types, grouped by bounded context. Follows the exact static-embedded pattern established by `cloud-resource-kinds://catalog` and `credential-types://catalog`.
todos:
  - id: catalog-json
    content: Hand-craft schemas/apiresourcekinds/catalog.json with curated kinds grouped by domain
    status: completed
  - id: embed
    content: Add ApiResourceKindFS embed directive to schemas/embed.go
    status: completed
  - id: discovery-pkg
    content: Create internal/domains/discovery/ package (doc.go, resources.go, register.go)
    status: completed
  - id: server-register
    content: Register discovery.RegisterResources in internal/server/server.go
    status: completed
  - id: verify
    content: Build and verify the server compiles cleanly; check lints
    status: completed
isProject: false
---

# T12: Expand MCP Resources -- `api-resource-kinds://catalog`

## Architectural Surprise (Documented)

The original T12 plan proposed 5 new resources. Investigation revealed:

- `**credential-types://catalog**` -- already delivered in T05.
- `**cloud-object-presets://{kind}**`, `**deployment-components://catalog**`, `**iac-modules://catalog**` -- all three already have MCP **tools** (`search_cloud_object_presets`, `search_deployment_components`, `search_iac_modules`) registered in `server.go`. These are dynamic database records, not static type-system metadata. Adding static MCP resources would be either stale or redundant, and would break the established pattern where `RegisterResources(srv)` takes no `serverAddress` (all existing resources are static/embedded).

**Remaining deliverable:** `api-resource-kinds://catalog` -- a curated, static catalog of all platform API resource types. This is the "table of contents" for the platform, grouped by domain. No existing tool or resource covers this.

## What Gets Built

One new MCP resource:

```
URI:  api-resource-kinds://catalog
Type: Static resource (not a template)
MIME: application/json
```

Minimal metadata per kind: `kind` (snake_case enum name) and `display_name`. Kinds grouped by domain (bounded context).

## Data Source

The `ApiResourceKind` proto enum at `planton/apis/.../apiresourcekind/api_resource_kind.proto` defines ~~60+ values. We hand-curate a subset (~~30-35 user-relevant kinds), excluding internal/system values (`unspecified`, `test_api_resource`, `platform`, `session`, `execution`, `chat`, `chat_message`, `billing_account`, etc.). Only kinds that have corresponding MCP tools or are meaningful to platform users are included.

Domain groupings (derived from server.go package structure):

- **resource_manager**: organization, environment, promotion_policy
- **infra_hub**: cloud_resource, cloud_object_preset, deployment_component, iac_module, stack_job, infra_pipeline, infra_chart, infra_project, flow_control_policy
- **service_hub**: service, variables_group, secrets_group, dns_domain, tekton_pipeline, tekton_task
- **config_manager**: secret, secret_version
- **connect**: runner_registration, default_provider_connection, default_runner_binding, provider_connection_authorization (plus 19 credential kinds referenced via `credential-types://catalog`)
- **iam**: identity_account, team, iam_role, api_key

## Catalog JSON Structure

```json
{
  "description": "All Planton API resource types accessible through MCP tools, grouped by domain.",
  "total_kinds": 30,
  "credential_catalog_uri": "credential-types://catalog",
  "cloud_resource_catalog_uri": "cloud-resource-kinds://catalog",
  "domains": {
    "resource_manager": {
      "display_name": "Resource Manager",
      "kinds": [
        { "kind": "organization", "display_name": "Organization" },
        { "kind": "environment", "display_name": "Environment" },
        { "kind": "promotion_policy", "display_name": "Promotion Policy" }
      ]
    }
  }
}
```

Cross-references to existing catalogs (`credential-types://catalog`, `cloud-resource-kinds://catalog`) avoid duplicating their content.

## File Layout

```
schemas/
  apiresourcekinds/
    catalog.json              # NEW -- hand-crafted catalog data

internal/domains/discovery/
  doc.go                      # NEW -- package documentation
  resources.go                # NEW -- CatalogResource() + CatalogHandler()
  register.go                 # NEW -- RegisterResources(srv)
```

## Files Modified

- [schemas/embed.go](schemas/embed.go) -- add `ApiResourceKindFS` embed directive for `apiresourcekinds/`
- [internal/server/server.go](internal/server/server.go) -- import `discovery` package, add `discovery.RegisterResources(srv)` call

## Implementation Pattern

Follows the exact pattern from [internal/domains/infrahub/cloudresource/resources.go](internal/domains/infrahub/cloudresource/resources.go) and [internal/domains/connect/credential/resources.go](internal/domains/connect/credential/resources.go):

- `CatalogResource()` returns `*mcp.Resource` with URI, name, title, description, MIME type
- `CatalogHandler()` returns `mcp.ResourceHandler` that reads from embedded FS via `schemas.ApiResourceKindFS`
- `RegisterResources(srv *mcp.Server)` calls `srv.AddResource()`
- Handler uses `domains.ResourceResult()` helper
- No `serverAddress` needed -- purely static, no gRPC calls

Key difference from cloudresource/credential: no `sync.Once` catalog builder needed. The catalog is a single JSON file read directly from the embedded FS, not assembled from a registry + per-kind schemas. Simpler.

## Curation Decisions to Make During Implementation

- Exact set of kinds to include (cross-reference with registered tool packages in `server.go`)
- Whether `pipeline` (enum value 25) maps to infra_pipeline, service_hub pipeline, or both
- Whether agent_fleet kinds (agent, skill, agent_test_suite, mcp_server) should be included if they don't have MCP tools yet

