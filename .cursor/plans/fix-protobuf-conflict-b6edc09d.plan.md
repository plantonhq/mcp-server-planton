<!-- b6edc09d-36f1-4737-a69a-9d41ac3f9118 89886e4e-792f-448e-98f5-bd5bc27c9233 -->
# Fix Protobuf Namespace Conflict

## Problem

The MCP server crashes with a proto file registration conflict because the codebase imports the same proto files from two different sources:

- `buf.build/gen/go/project-planton/apis/protocolbuffers/go` (correct)
- `github.com/project-planton/project-planton` (incorrect)

## Root Cause

`internal/domains/infrahub/cloudresource/kinds.go` line 11 imports:

```go
cloudresourcekind "github.com/project-planton/project-planton/apis/org/project_planton/shared/cloudresourcekind"
```

While other files (`unwrap.go`, `cloudresource_client.go`) correctly use:

```go
cloudresourcekind "buf.build/gen/go/project-planton/apis/protocolbuffers/go/org/project_planton/shared/cloudresourcekind"
```

## Solution

### 1. Update Import in kinds.go

Replace the github.com import with the buf.build import to match the rest of the codebase:

**File:** `internal/domains/infrahub/cloudresource/kinds.go`

Change line 11 from:

```go
cloudresourcekind "github.com/project-planton/project-planton/apis/org/project_planton/shared/cloudresourcekind"
```

To:

```go
cloudresourcekind "buf.build/gen/go/project-planton/apis/protocolbuffers/go/org/project_planton/shared/cloudresourcekind"
```

### 2. Remove Conflicting Dependency

Remove the direct github.com dependency from `go.mod`:

**File:** `go.mod`

Remove line 10:

```go
github.com/project-planton/project-planton v0.2.245
```

### 3. Clean Up Dependencies

Run `go mod tidy` to ensure all dependencies are correct and remove any unused transitive dependencies.

## Verification

After the changes, the MCP server should start without the proto registration panic. The buf.build modules provide the same functionality without conflicts.

### To-dos

- [ ] Update import in kinds.go to use buf.build path
- [ ] Remove github.com/project-planton dependency from go.mod
- [ ] Run go mod tidy to clean up dependencies