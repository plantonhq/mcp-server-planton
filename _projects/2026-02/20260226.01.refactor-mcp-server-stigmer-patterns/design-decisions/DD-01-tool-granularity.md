# DD-01: Generic Cloud Resource Tools vs. Per-Resource-Kind Tools

**Date**: 2026-02-26
**Status**: PROPOSED

## Context

Planton has 50+ cloud resource kinds (AwsAlb, GcpGkeCluster, KubernetesDeployment, etc.) defined in OpenMCF. Each has its own proto spec. We need to decide how to expose these via MCP tools.

## Options

### Option A: Per-Kind Tools (apply_aws_alb, delete_aws_alb, etc.)
- Pros: Type-safe inputs, better discoverability, AI gets schema for each kind
- Cons: 100+ tools needed, overwhelming for MCP clients, massive codegen scope

### Option B: Generic Tools (apply_cloud_resource, delete_cloud_resource)
- Pros: 2-3 tools total, matches backend API, simple to implement
- Cons: No per-kind schema validation at the tool level, AI needs to know the provider spec format

## Decision

**Option B: Generic cloud resource tools.**

The backend API already works this way — `CloudResourceCommandController.Apply` accepts any `CloudResource` with the provider spec in `cloud_object`. The AI can reference OpenMCF documentation for provider-specific schemas.

## Consequences

- Simple tool set that scales with zero additional code per new resource kind
- AI assistants will need context about OpenMCF provider specs to construct valid inputs
- Validation happens server-side, not at tool input level
