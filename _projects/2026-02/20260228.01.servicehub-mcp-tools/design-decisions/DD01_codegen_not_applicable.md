# DD01: Codegen Pipeline Is Not Applicable to ServiceHub Tools

**Date**: 2026-02-28
**Status**: Accepted
**Context**: Should we extend or reuse the existing `proto2schema` + `generator` codegen pipeline for ServiceHub MCP tools?

## Decision

Hand-write all 35 ServiceHub MCP tools. The existing codegen pipeline is not applicable.

## Rationale

### What the codegen does

The two-stage pipeline (`tools/codegen/proto2schema/` → `tools/codegen/generator/`) solves a specific problem:

- **362+ cloud resource kinds** from the OpenMCF provider proto definitions
- Each kind has a `Spec` message that varies by provider (AWS, GCP, Azure, etc.)
- The MCP server receives these specs as opaque `map[string]any` (`cloud_object`)
- Codegen generates typed Go input structs with `validate()`, `applyDefaults()`, `toMap()`, and `Parse{Kind}()` for each kind
- The generated `GetParser(kind)` registry dispatches to the correct parser at runtime

This is a classic "schema-driven code generation at scale" solution.

### Why it doesn't fit ServiceHub

| Aspect | Cloud Resources | ServiceHub |
|--------|----------------|------------|
| Entity count | 362+ | 7 |
| API pattern | Single generic `apply(CloudResource)` with embedded `cloud_object` | Entity-specific RPCs with typed protobuf messages |
| Schema shape | Varies wildly per kind (each provider has unique spec) | Fixed, well-known specs |
| Operations | Uniform CRUD on all kinds | Unique domain ops per entity (`disconnectGitRepo`, `resolveManualGate`, `upsertEntry`) |
| Validation | Schema-driven (field types, enums, required) | Business-rule-driven (mutual exclusion, conditional) |

The codegen's value proposition is *scale* — generating boilerplate for 362+ kinds that humans shouldn't maintain by hand. ServiceHub has 7 entities, each with unique behavior that requires human judgment in tool descriptions, input validation, and error messages.

### The one reuse point

`apply_service` handles inline deployment targets where `deployment_targets[].cloud_object` is a standard OpenMCF cloud resource spec. The existing `cloudresource.GetParser(kind)` registry can validate these cloud objects without any codegen changes. This is reuse of the codegen *output*, not extension of the codegen *pipeline*.

## Consequences

- All 35 tools are hand-written in `internal/domains/servicehub/`
- Each entity package follows the established pattern: `register.go`, `tools.go`, individual operation files
- No new codegen pipeline or generator modifications required
- The `apply_service` handler imports `gen/infrahub/cloudresource` to validate inline deployment target cloud objects
