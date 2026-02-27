# Development

Guide for building, testing, and contributing to the Planton Cloud MCP Server.

## Prerequisites

- **Go 1.25+**
- **Git**
- **Docker** (optional, for container builds)
- **golangci-lint** (optional, for linting)
- Access to Planton Cloud APIs (local or remote)

For codegen only:
- Local clone of [plantonhq/openmcf](https://github.com/plantonhq/openmcf) at `~/scm/github.com/plantonhq/openmcf/`
- Local clone of [plantonhq/planton](https://github.com/plantonhq/planton) at `~/scm/github.com/plantonhq/planton/`

## Build

```bash
make build           # Build binary to bin/mcp-server-planton
make install         # Install to GOPATH/bin
make install-local   # Install to /usr/local/bin (requires sudo)
```

## Test

```bash
make test            # Run all tests (go test -v ./...)
go test ./internal/... # Run only internal package tests
```

Test files live alongside their source files following standard Go conventions:

| Test File | What It Tests |
|-----------|---------------|
| `internal/config/config_test.go` | Config validation, log level parsing |
| `internal/domains/rpcerr_test.go` | gRPC error classification |
| `internal/domains/kind_test.go` | Shared kind enum resolution |
| `internal/domains/infrahub/cloudresource/identifier_test.go` | Resource identifier validation |
| `internal/domains/infrahub/cloudresource/kind_test.go` | Kind extraction and batch resolution |
| `internal/domains/infrahub/cloudresource/list_test.go` | `resolveKinds` batch helper |
| `internal/domains/infrahub/cloudresource/metadata_test.go` | Metadata extraction from cloud_object |
| `internal/domains/infrahub/cloudresource/schema_test.go` | URI parsing, embedded FS, kind catalog |
| `internal/domains/infrahub/cloudresource/apply_test.go` | CloudResource proto assembly |
| `internal/domains/infrahub/stackjob/enum_test.go` | Stack job enum resolvers (operation, status, result, kind) |
| `internal/parse/helpers_test.go` | Shared parse utilities |

## Code Quality

```bash
make fmt             # Format all Go code
make fmt-check       # Check formatting (CI gate)
make lint            # Run golangci-lint
go vet ./...         # Run Go vet
```

## Codegen Pipeline

The two-stage codegen pipeline generates typed Go input structs from OpenMCF provider proto definitions. This produces 362 provider types across 17 cloud platforms.

### Stage 1: proto2schema

Parses OpenMCF `.proto` files and produces JSON schemas.

```bash
make codegen-schemas
```

- **Tool**: `tools/codegen/proto2schema/`
- **Input**: OpenMCF provider protos (via local `SCM_ROOT` convention or `--openmcf-apis-dir`)
- **Output**: `schemas/providers/{cloud}/{kind}.json` + `schemas/providers/registry.json`

### Stage 2: schema2go

Generates typed Go input structs from JSON schemas.

```bash
make codegen-types
```

- **Tool**: `tools/codegen/generator/`
- **Input**: `schemas/` directory
- **Output**: `gen/infrahub/cloudresource/{cloud}/{kind}_gen.go` + `gen/infrahub/cloudresource/registry_gen.go`

### Full Pipeline

```bash
make codegen         # Runs Stage 1 then Stage 2
```

### Generated Code Structure

```
gen/infrahub/cloudresource/
  registry_gen.go              Central ParseFunc registry (imports all 17 cloud packages)
  alicloud/
    alicloudapplicationloadbalancer_gen.go
    alicloudvpc_gen.go
    types_gen.go               Shared nested types for this cloud
    ...
  aws/
    awsvpc_gen.go
    awsekscluster_gen.go
    ...
  gcp/
    ...
  (17 cloud packages total)
```

Each generated file provides:
- A typed `{Kind}SpecInput` struct with JSON tags and validation
- `validate()`, `applyDefaults()`, `toMap()` methods
- A `Parse{Kind}()` function registered in the central registry

## Project Structure

```
cmd/mcp-server-planton/        CLI entry point
pkg/mcpserver/                  Public embedding API
internal/
  auth/                         Context-based API key propagation
  config/                       Environment-variable configuration
  grpc/                         gRPC client factory
  server/                       MCP server init + transports (registers all tools)
  domains/
    kind.go                     Shared kind enum resolution (used by all domains)
    infrahub/                   Infrastructure Hub bounded context
      cloudresource/            Cloud resource lifecycle tools (11 tools)
      stackjob/                 Stack job observability tools (3 tools)
      preset/                   Cloud object preset tools (2 tools)
    resourcemanager/            Resource Manager bounded context
      environment/              Environment discovery (1 tool)
      organization/             Organization discovery (1 tool)
  parse/                        Shared utilities for generated parsers
gen/infrahub/cloudresource/     Generated typed input structs
schemas/                        Embedded JSON schemas (go:embed)
tools/codegen/
  proto2schema/                 Stage 1 codegen tool
  generator/                    Stage 2 codegen tool
```

## Release

```bash
make release version=v1.0.0          # Create and push a release tag
make release version=v1.0.0 force=true  # Force-recreate an existing tag
```

GitHub Actions builds multi-platform binaries and Docker images on tag push.

## Docker

```bash
make docker-build    # Build local image
make docker-run      # Run with current env vars
```
