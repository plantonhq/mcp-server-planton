---
name: proto2schema codegen tool
overview: Build the proto2schema codegen tool that parses OpenMCF provider .proto files from the local filesystem and outputs JSON schemas for each provider kind. This is Phase 2, Stage 1 of the mcp-server-planton refactoring.
todos:
  - id: schema-types
    content: Define schema type definitions in schema.go (ProviderSchema, SpecSchema, FieldSchema, TypeSpec, Validation, NestedType)
    status: completed
  - id: options-extractor
    content: Build OpenMCF custom option extraction in options.go (default_kind, default_kind_field_path, default, recommended_default via protowire)
    status: completed
  - id: parser
    content: Build proto parser in parser.go (parse api.proto for kind/apiVersion constants, parse spec.proto for fields, handle StringValueOrRef simplification, collect nested types)
    status: completed
  - id: registry
    content: Build registry generation in registry.go (scan all generated schemas, produce registry.json with kind-to-path mapping)
    status: completed
  - id: main-cli
    content: Build CLI entry point in main.go (flag parsing, SCM_ROOT resolution, provider directory scanning, buf cache detection, orchestrate parsing)
    status: completed
  - id: shared-metadata
    content: Generate shared/metadata.json from CloudResourceMetadata proto
    status: completed
  - id: makefile
    content: Add codegen-schemas target to Makefile
    status: completed
  - id: dependencies
    content: Add jhump/protoreflect and buf.validate deps to go.mod
    status: completed
  - id: verify
    content: Run full codegen, spot-check output schemas against source protos for correctness
    status: completed
isProject: false
---

# Phase 2, Stage 1: proto2schema for OpenMCF Providers

## Context

Adapting Stigmer's two-stage codegen pipeline for Planton. Stage 1 (`proto2schema`) parses OpenMCF provider `.proto` files and emits JSON schemas that will be consumed by Stage 2 (generator) and by MCP resource templates.

**Reference implementation**: [tools/codegen/proto2schema/main.go](stigmer/tools/codegen/proto2schema/main.go) (~1070 lines, single file)

**Key difference from Stigmer**: Stigmer parses protos from the same monorepo. Planton parses from a separate repo (`openmcf`) resolved via local filesystem paths following the SCM root convention (`$SCM_ROOT/github.com/{org}/{repo}/`).

---

## Design Decision: StringValueOrRef Handling

**Decision**: Option C -- Simplify `StringValueOrRef` to `string` in schemas, but capture `referenceKind` and `referenceFieldPath` as metadata.

- Agents provide literal string values at the MCP boundary
- The schema preserves cross-resource relationship information for documentation and future expansion
- Stage 2 `ToProto()` methods wrap strings into `StringValueOrRef{Value: s}`
- Respects bounded context boundaries (specification layer vs provisioning layer)

---

## Proto File Resolution

OpenMCF protos live in a separate repo. Resolution strategy:

- **Environment variable**: `SCM_ROOT` (defaults to `$HOME/scm`)
- **CLI flag**: `--openmcf-apis-dir` (defaults to `$SCM_ROOT/github.com/plantonhq/openmcf/apis`)
- **Buf cache**: Auto-detected from `~/.cache/buf/v3/modules/` for `buf.validate` dependencies
- The openmcf `apis/` directory serves as the proto include path (matches `buf.yaml` module root)

---

## Output Schema Format

One JSON file per provider kind, plus a shared metadata schema and a provider registry.

### Directory structure

```
tools/codegen/schemas/
  shared/
    metadata.json                 # CloudResourceMetadata fields (extracted once)
  providers/
    registry.json                 # Kind-to-path index for all providers
    aws/
      awsalb.json
      awscertmanagercert.json
      ...
    gcp/
      gcpgkecluster.json
      ...
    (one directory per cloud provider)
```

### Provider schema shape (e.g., `awsalb.json`)

```json
{
  "name": "AwsAlb",
  "kind": "AwsAlb",
  "cloudProvider": "aws",
  "apiVersion": "aws.openmcf.org/v1",
  "description": "...",
  "protoPackage": "org.openmcf.provider.aws.awsalb.v1",
  "protoFiles": {
    "api": "org/openmcf/provider/aws/awsalb/v1/api.proto",
    "spec": "org/openmcf/provider/aws/awsalb/v1/spec.proto"
  },
  "spec": {
    "name": "AwsAlbSpec",
    "fields": [
      {
        "name": "Subnets",
        "jsonName": "subnets",
        "protoField": "subnets",
        "type": { "kind": "array", "elementType": { "kind": "string" } },
        "description": "list of subnet IDs...",
        "required": true,
        "validation": { "required": true, "minItems": 2 },
        "referenceKind": "AwsVpc",
        "referenceFieldPath": ""
      }
    ]
  },
  "nestedTypes": [
    {
      "name": "AwsAlbDns",
      "description": "...",
      "fields": [...]
    }
  ]
}
```

Key design choices in the schema format:

- `StringValueOrRef` fields are typed as `string` with `referenceKind`/`referenceFieldPath` metadata
- `spec` contains only the `{Kind}Spec` message fields (not api_version, kind, metadata -- those are constant envelope fields)
- `nestedTypes` captures provider-specific sub-messages (e.g., `AwsAlbDns`, `AwsAlbSsl`)
- `status` is excluded (backend-managed, per resolved decision #1)

### Registry shape (`registry.json`)

```json
{
  "providers": {
    "AwsAlb": { "cloudProvider": "aws", "apiVersion": "aws.openmcf.org/v1", "schemaFile": "aws/awsalb.json" },
    "GcpGkeCluster": { "cloudProvider": "gcp", "apiVersion": "gcp.openmcf.org/v1", "schemaFile": "gcp/gcpgkecluster.json" }
  }
}
```

Used by Stage 2 generator and MCP resource template handlers to discover available kinds.

---

## Source File Structure

Unlike Stigmer's single-file approach, we split into focused files for maintainability. Each file has a single responsibility.

```
tools/codegen/proto2schema/
  main.go          # CLI entry point, flag parsing, provider directory scanning
  schema.go        # Schema type definitions (the data model for JSON output)
  parser.go        # Proto file parsing: proto descriptors -> schema structs
  options.go       # OpenMCF custom proto option extraction (default_kind, default, recommended_default)
  registry.go      # Registry generation (kind -> schema file mapping)
```

### Why split (vs Stigmer's single file)?

Stigmer's proto2schema is ~1070 lines and growing. Planton's will be larger due to:

- Cross-repo proto resolution logic
- OpenMCF-specific custom option extraction (4 custom extensions vs Stigmer's 3)
- Provider directory scanning (150 providers across 12 clouds vs Stigmer's ~15 resources across 3 namespaces)
- Registry generation

Single responsibility per file makes each piece independently testable and reviewable.

---

## Dependencies

Add to `go.mod`:

- `github.com/jhump/protoreflect` -- proto file parsing with source info (same as Stigmer)
- `buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate` -- for extracting buf.validate rules via `proto.GetExtension` (same as Stigmer)

---

## Custom Proto Options to Extract

OpenMCF defines these custom extensions that we extract from proto unknown fields using `protowire`:


| Extension                 | Field Number | Source Proto                      | Type          | Purpose                                  |
| ------------------------- | ------------ | --------------------------------- | ------------- | ---------------------------------------- |
| `default_kind`            | 200001       | `foreignkey/v1/foreign_key.proto` | varint (enum) | Which kind a StringValueOrRef references |
| `default_kind_field_path` | 200002       | `foreignkey/v1/foreign_key.proto` | string        | Output field path for references         |
| `default`                 | 60001        | `options/options.proto`           | string        | Default value for a field                |
| `recommended_default`     | 60002        | `options/options.proto`           | string        | Recommended default value                |


These are stored as metadata in the schema, not as validation rules.

---

## Makefile Target

```makefile
codegen-schemas:
	go run tools/codegen/proto2schema/main.go --all
```

The `--all` flag scans all provider directories. For development/testing, `--provider aws/awsalb` processes a single provider.

---

## Verification

After implementation, verify with:

1. `go build ./tools/codegen/proto2schema/` -- compiles cleanly
2. `go vet ./tools/codegen/proto2schema/` -- no issues
3. `make codegen-schemas` -- generates schemas for all 150 providers
4. Spot-check 3-4 generated schemas against their source protos (AwsAlb, GcpGkeCluster, KubernetesDeployment, ConfluentKafka) for correctness
5. Validate registry.json has entries for all providers

