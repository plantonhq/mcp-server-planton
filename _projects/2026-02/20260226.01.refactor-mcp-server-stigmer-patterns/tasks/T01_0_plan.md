# Task T01: Refactor mcp-server-planton to Follow Stigmer MCP Server Patterns

**Created**: 2026-02-26
**Status**: APPROVED
**Type**: Refactoring

## Executive Summary

Completely gut and rebuild `mcp-server-planton` following the proven architecture of `stigmer/mcp-server`. The first milestone targets **three MCP tools**: `apply_cloud_resource`, `delete_cloud_resource`, and `get_cloud_resource` for the infrahub domain. This includes building a codegen pipeline adapted from Stigmer's two-stage approach from day one — no hand-written types.

---

## Current State Analysis

### What Stigmer Does Well (Reference Architecture)

1. **Domain-driven tool structure**: `internal/domains/{domain}/{resource}/` with consistent files:
   - `tools.go` — tool definitions + handler wiring
   - `apply.go` — apply (create-or-update) logic
   - `delete.go` — delete logic  
   - `fetch.go` — get/read logic
   - `resources.go` — MCP resource templates

2. **Two-stage codegen pipeline**:
   - **Stage 1** (`proto2schema`): Proto definitions → JSON schemas (field types, validation, docs)
   - **Stage 2** (`generator --target=mcp`): JSON schemas → Go input structs with `ToProto()` methods
   - Generated code lives in `gen/{domain}/{resource}/` with `_gen.go` suffix

3. **Shared domain utilities**: `internal/domains/` has `conn.go`, `marshal.go`, `rpcerr.go`, `toolresult.go`, `resourcehandler.go`

4. **Typed tool handlers**: Uses `mcp.AddTool(srv, tool, handler)` with typed input structs (not `map[string]interface{}`)

5. **Central registration**: `internal/server/server.go` has `registerTools()` and `registerResources()`

### What mcp-server-planton Currently Has

1. **55 Go files** with manual tool implementations (no codegen)
2. **Domain structure exists** but is more fragmented: `internal/domains/{domain}/{resource}/`
3. **Dynamic schema extraction** using protobuf reflection — clever but not maintainable
4. **Complex wrap/unwrap** logic for `cloud_object` (`google.protobuf.Struct`)
5. **8 cloud resource tools + 1 resource** — more than needed initially
6. **Inconsistent handler signatures** — uses `map[string]interface{}` for arguments

### The CloudResource Wrapping Challenge

Unlike Stigmer's flat resources (Agent, McpServer, Skill, Workflow), Planton's CloudResource has a unique architecture:

```
CloudResource (Planton envelope)
├── api_version: "infra-hub.planton.ai/v1"
├── kind: "CloudResource"  
├── metadata: ApiResourceMetadata
├── spec: CloudResourceSpec
│   ├── kind: CloudResourceKind (e.g., aws_alb, gcp_gke_cluster)
│   ├── cloud_object: google.protobuf.Struct  ← OpenMCF spec lives here
│   ├── provisioner_info: ...
│   ├── provider_info: ...
│   └── ...
└── status: CloudResourceStatus

OpenMCF Provider Spec (e.g., AwsAlb) — what the USER sees:
├── api_version: "aws.openmcf.org/v1"
├── kind: "AwsAlb"
├── metadata: CloudResourceMetadata
├── spec: AwsAlbSpec
└── status: AwsAlbStatus
```

**Key insight**: The user thinks in terms of OpenMCF provider specs (AwsAlb, GcpGkeCluster, etc.), but the backend API accepts `CloudResource` with the provider spec serialized into `cloud_object` as a `google.protobuf.Struct`.

---

## Design Decisions

### DD-01: Tool Granularity — Generic vs. Per-Resource-Kind

**Decision**: Use **generic cloud resource tools** (`apply_cloud_resource`, `delete_cloud_resource`) that accept the OpenMCF YAML/JSON directly, not per-kind tools.

**Rationale**:
- The backend `CloudResourceCommandController.apply` already accepts any `CloudResource` — the `spec.kind` + `cloud_object` determine the provider type
- Creating per-kind tools (apply_aws_alb, apply_gcp_gke_cluster, etc.) would require codegen for every OpenMCF provider — massive scope
- The user's MCP client (Claude, Cursor, etc.) would be overwhelmed with hundreds of tools
- Generic tools match how the backend API actually works

**Apply tool input**: The tool accepts the OpenMCF-formatted resource (with `api_version`, `kind`, `metadata`, `spec` from the provider-specific schema). The tool handler wraps this into the `CloudResource` envelope before calling the backend.

**Delete tool input**: `org`, `env`, `cloud_resource_kind`, `slug` — uses `getByOrgByEnvByKindBySlug` to resolve the resource, then calls `delete`.

### DD-02: Codegen Strategy — Adapt Stigmer's Pipeline for Planton

**Decision**: Build a codegen pipeline in `tools/codegen/` following Stigmer's two-stage approach, but adapted for Planton's needs.

**Differences from Stigmer**:
- Stigmer codegen generates **input structs for each resource** (AgentInput, McpServerInput)
- For Planton cloud resources, we generate a **single CloudResourceApplyInput** that handles the envelope wrapping, plus uses `google.protobuf.Struct` for the provider-specific `cloud_object`
- The codegen for Planton is simpler initially — we don't need per-provider-kind input structs since the cloud_object is passed as JSON

**What gets generated** (in `gen/infrahub/cloudresource/`):
- `cloud_resource_gen.go`: `CloudResourceApplyInput` struct with `ToProto()` method
- The `ToProto()` handles: constructing `CloudResource`, setting `api_version`/`kind` constants, wrapping the cloud_object

### DD-03: Directory Structure — Mirror Stigmer

**Decision**: Adopt the same directory structure as Stigmer MCP server.

```
mcp-server-planton/
├── cmd/mcp-server-planton/main.go          # Entry point (simplified)
├── gen/                                     # Generated code (DO NOT EDIT)
│   └── infrahub/
│       └── cloudresource/
│           └── cloud_resource_gen.go
├── internal/
│   ├── auth/                                # API key management
│   ├── config/                              # Configuration (env vars)
│   ├── domains/                             # Domain-specific tools
│   │   ├── conn.go                          # gRPC connection helper
│   │   ├── marshal.go                       # JSON marshaling
│   │   ├── rpcerr.go                        # Error translation
│   │   ├── toolresult.go                    # Tool result helpers
│   │   └── infrahub/
│   │       └── cloudresource/
│   │           ├── tools.go                 # Tool definitions + registration
│   │           ├── apply.go                 # Apply handler
│   │           ├── delete.go                # Delete handler
│   │           └── fetch.go                 # Get handler
│   ├── grpc/                                # gRPC client factory
│   └── server/                              # MCP server init + registration
│       └── server.go
├── pkg/mcpserver/                           # Public API for embedding
│   ├── config.go
│   └── run.go
├── tools/                                   # Codegen pipeline
│   └── codegen/
│       ├── proto2schema/main.go             # Stage 1: Proto → Schema
│       ├── generator/                       # Stage 2: Schema → Go
│       │   ├── main.go
│       │   └── mcp.go
│       └── schemas/                         # Intermediate JSON schemas
│           └── infrahub/
│               └── cloudresource/
│                   └── cloudresource.json
├── Makefile
├── go.mod
└── README.md
```

### DD-04: Remove Everything, Start Clean

**Decision**: Delete all existing `internal/domains/` code and rebuild from scratch.

**Rationale**:
- Current code uses `map[string]interface{}` handlers (old mcp-go pattern)
- Current code has complex dynamic protobuf reflection that won't be needed
- Current code has inconsistent patterns across domains
- It's faster to write clean code following Stigmer patterns than to refactor the existing mess
- We're only building 2 tools initially, so the scope is small

### DD-05: Authentication — Match Stigmer Pattern

**Decision**: Use Stigmer's context-based auth pattern.

- STDIO mode: API key from env var, injected into base context
- HTTP mode: API key from `Authorization: Bearer` header per request
- Propagated via `auth.WithAPIKey()` / `auth.APIKey()` context helpers
- Passed as gRPC metadata to Planton backend

---

## Phase Breakdown

### Phase 1: Clean Slate + Shared Utilities

**Goal**: Remove existing tools, set up the Stigmer-style foundation.

1. Delete all files under `internal/domains/`
2. Delete `internal/common/` (will be replaced by domain utilities)
3. Rewrite `internal/server/server.go` following Stigmer's registration pattern
4. Create shared domain utilities:
   - `internal/domains/conn.go` — gRPC connection with auth
   - `internal/domains/marshal.go` — protojson marshaling
   - `internal/domains/rpcerr.go` — gRPC error translation
   - `internal/domains/toolresult.go` — MCP tool result helpers
5. Rewrite `internal/auth/` following Stigmer's credentials pattern
6. Simplify `internal/config/` to match Stigmer's env-based config
7. Rewrite `cmd/mcp-server-planton/main.go` following Stigmer's entry point

### Phase 2: Codegen Pipeline (Adapted from Stigmer)

**Goal**: Build the two-stage codegen that produces Go input types.

1. Create `tools/codegen/proto2schema/main.go` — adapted from Stigmer
   - Parse `CloudResource`, `CloudResourceSpec` proto definitions
   - Extract fields, types, validation rules, documentation
   - Output JSON schemas to `tools/codegen/schemas/infrahub/cloudresource/`

2. Create `tools/codegen/generator/main.go` + `mcp.go` — adapted from Stigmer
   - Read JSON schemas
   - Generate `CloudResourceApplyInput` struct
   - Generate `ToProto()` method for CloudResource envelope wrapping
   - Output to `gen/infrahub/cloudresource/cloud_resource_gen.go`

3. Add Makefile targets:
   - `codegen-schemas`: Proto → JSON schemas
   - `codegen-mcp`: JSON schemas → Go code
   - `codegen`: Full pipeline

### Phase 3: Implement apply_cloud_resource Tool

**Goal**: First working MCP tool.

1. Create `internal/domains/infrahub/cloudresource/tools.go`:
   - `ApplyTool()` — MCP tool definition
   - `ApplyHandler()` — wired to generated input type

2. Create `internal/domains/infrahub/cloudresource/apply.go`:
   - Accepts generated `CloudResourceApplyInput`
   - Calls `ToProto()` to build `CloudResource` protobuf
   - Calls `CloudResourceCommandController.Apply` via gRPC
   - Returns JSON response

3. Register in `internal/server/server.go`

### Phase 4: Implement delete_cloud_resource + get_cloud_resource Tools

**Goal**: Complete the tool set with delete and get.

1. Add to `internal/domains/infrahub/cloudresource/tools.go`:
   - `DeleteTool()` + `DeleteHandler()` — delete tool
   - `GetTool()` + `GetHandler()` — get tool

2. Create `internal/domains/infrahub/cloudresource/delete.go`:
   - Input: `org`, `env`, `cloud_resource_kind`, `slug`
   - Step 1: Call `getByOrgByEnvByKindBySlug` to resolve resource ID
   - Step 2: Call `CloudResourceCommandController.Delete` with `ApiResourceDeleteInput`
   - Return JSON response

3. Create `internal/domains/infrahub/cloudresource/fetch.go`:
   - Input: `org`, `env`, `cloud_resource_kind`, `slug`
   - Call `getByOrgByEnvByKindBySlug` to fetch the resource
   - Return JSON-marshaled response

4. Register all tools in server

### Phase 5: Testing + Documentation

1. Test both tools end-to-end against Planton APIs
2. Update README.md with new architecture documentation
3. Update docs/ with new tool reference
4. Ensure `make build` works

---

## Tool Specifications

### apply_cloud_resource

```
Name: apply_cloud_resource
Description: Create or update a cloud resource. Accepts the OpenMCF-formatted 
             resource specification. The resource kind determines the cloud provider 
             and infrastructure type (e.g., AwsAlb, GcpGkeCluster, KubernetesDeployment).

Input Schema:
  - org (string, required): Organization slug
  - env (string, required): Environment slug  
  - cloud_resource_kind (string, required): OpenMCF cloud resource kind (e.g., "aws_alb")
  - cloud_object (object, required): Provider-specific resource spec as JSON
    (matches the OpenMCF provider API message, e.g., AwsAlb)
  - container_registry (string, optional): Docker credential slug
  - connection (string, optional): Provider credential slug

Backend RPC: CloudResourceCommandController.Apply(CloudResource)
```

### delete_cloud_resource

```
Name: delete_cloud_resource
Description: Delete a cloud resource by its organization, environment, kind, and slug.

Input Schema:
  - org (string, required): Organization slug
  - env (string, required): Environment slug
  - cloud_resource_kind (string, required): Cloud resource kind (e.g., "aws_alb") 
  - slug (string, required): Resource slug

Backend RPC: 
  1. CloudResourceQueryController.GetByOrgByEnvByKindBySlug → resolve ID
  2. CloudResourceCommandController.Delete(ApiResourceDeleteInput)
```

### get_cloud_resource

```
Name: get_cloud_resource
Description: Get a cloud resource by its organization, environment, kind, and slug.

Input Schema:
  - org (string, required): Organization slug
  - env (string, required): Environment slug
  - cloud_resource_kind (string, required): Cloud resource kind (e.g., "aws_alb")
  - slug (string, required): Resource slug

Backend RPC: CloudResourceQueryController.GetByOrgByEnvByKindBySlug
```

---

## Resolved Decisions (from review)

1. **Cloud object format**: Accept the full OpenMCF message (`api_version`, `kind`, `metadata` with name/org) but NOT `status` (backend-managed). Metadata fields (name, org) are required for creating resources. The tool handler wraps the OpenMCF message into the CloudResource envelope before calling the backend.

2. **Tool naming**: Keep `apply_cloud_resource` / `delete_cloud_resource` / `get_cloud_resource`.

3. **Codegen scope**: Build the full codegen pipeline from the start — no hand-written types. This matches Stigmer's approach and simplifies adding new domains going forward. Many more domains will be added after infrahub/cloudresource.

4. **get_cloud_resource**: Included in scope — implemented alongside apply and delete.

---

## Notes

- Preserve working functionality at all times
- **IMPORTANT**: Only document in knowledge folders after ASKING for permission
- Task logs can be updated freely
- The Stigmer MCP server at `/Users/suresh/scm/github.com/stigmer/stigmer/mcp-server/` is the reference implementation

## Review Process

**What happens next**:
1. **You review this plan** — especially the Design Decisions and Open Questions
2. **Provide feedback** — any concerns, changes to scope, answers to open questions
3. **I'll revise the plan** — incorporate your feedback
4. **You approve** — give explicit approval to proceed
5. **Execution begins** — tracked in T01_3_execution.md
