# AD-02: Restructure Generated Code Under Domain Directories

**Date**: 2026-02-27
**Status**: Accepted
**Decision**: Move generated code from flat `gen/cloudresource/` to domain-scoped `gen/infrahub/cloudresource/`

## Context

Currently, generated code lives at:
```
gen/
└── cloudresource/
    ├── registry_gen.go
    ├── aws/
    ├── gcp/
    ├── azure/
    └── ... (17 provider directories)
```

Hand-written code follows the domain hierarchy:
```
internal/domains/
├── infrahub/
│   ├── cloudresource/
│   ├── stackjob/
│   └── preset/
└── resourcemanager/
    ├── organization/
    └── environment/
```

The generated code is flat (`gen/cloudresource/`) while the domain code is nested (`internal/domains/infrahub/cloudresource/`). As we add more domains (infrachart, infraproject, graph, configmanager, audit), the gen/ directory will become a confusing mix of unrelated domains at the same level.

## Decision

Restructure gen/ to mirror the domain hierarchy:
```
gen/
└── infrahub/
    └── cloudresource/
        ├── registry_gen.go
        ├── aws/
        ├── gcp/
        └── ...
```

## Impact

1. All import paths change from `gen/cloudresource` → `gen/infrahub/cloudresource`
2. All import paths for providers change from `gen/cloudresource/aws` → `gen/infrahub/cloudresource/aws`
3. The code generator configuration must be updated to output to the new path
4. This is a Phase 0 prerequisite — must be done before adding new domain tools

## Migration Steps

1. Create `gen/infrahub/cloudresource/` directory
2. Move all files and subdirectories
3. Update Go package declarations
4. Update all import paths across the codebase
5. Update code generator config
6. Verify build and tests pass
7. Remove old `gen/cloudresource/` directory
