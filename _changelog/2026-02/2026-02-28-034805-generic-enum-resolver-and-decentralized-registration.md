# Generic Enum Resolver and Decentralized Tool Registration

**Date**: February 28, 2026

## Summary

Eliminated ~240 lines of repetitive boilerplate across the MCP server's domain layer by introducing a generic `EnumResolver[T]` type (Go generics) and decentralizing tool registration from a monolithic 140-line function in `server.go` to per-domain `Register()` functions. All 63 tools continue to work identically — this is a pure structural refactoring with zero behavior change.

## Problem Statement

An architectural review of the `internal/domains/` layer revealed two maintainability gaps that would compound as the tool surface grows beyond 63 tools.

### Pain Points

- **12 hand-written enum resolver functions** across 6 `enum.go` files, all structurally identical — differing only in the proto enum type, label, and zero-value key
- Two packages (`configmanager/variable`, `configmanager/secret`) independently reinvented `joinScopeValues()` instead of using the shared `domains.JoinEnumValues`, creating copy-paste drift
- `server.go` contained a **63-line `registerTools()` function** plus a **63-element string slice** for logging — every new tool required changes in 3 locations (domain package + 2 places in server.go)
- The hardcoded `"count", 63` in the log statement would silently lie whenever a tool was added without updating it

## Solution

Two independent, zero-risk refactorings:

1. **Generic `EnumResolver[T ~int32]`** — a single type with `Resolve(string)` and `ResolveSlice([]string)` methods that replaces all 12 hand-written functions
2. **Per-domain `register.go` files** — each of the 15 domain packages now owns its own `Register(srv, serverAddress)` function, reducing `server.go` to a 15-line delegation list

## Implementation Details

### Generic Enum Resolver

Added `EnumResolver[T]` to `internal/domains/enum.go`:

```go
type EnumResolver[T ~int32] struct { ... }
func NewEnumResolver[T ~int32](values map[string]int32, typeName, excludeKey string) EnumResolver[T]
func (r EnumResolver[T]) Resolve(s string) (T, error)
func (r EnumResolver[T]) ResolveSlice(ss []string) ([]T, error)
```

Each domain now declares resolver vars instead of functions:

```go
var operationTypeResolver = domains.NewEnumResolver[stackjobv1.StackJobOperationType](
    stackjobv1.StackJobOperationType_value, "stack job operation type", "stack_job_operation_type_unspecified")
```

Files rewritten: `graph/enum.go`, `stackjob/enum.go`, `audit/enum.go`, `variable/enum.go`, `secret/enum.go`, `provider.go`

### Decentralized Registration

Created 15 `register.go` files — one per domain package. Each contains a `Register(srv *mcp.Server, serverAddress string)` function that calls `mcp.AddTool` for its own tools.

`server.go` now delegates:

```go
func registerTools(srv *mcp.Server, serverAddress string) {
    cloudresource.Register(srv, serverAddress)
    stackjob.Register(srv, serverAddress)
    // ... one line per domain package
}
```

### Special cases preserved

- `domains.ResolveKind()` was left unchanged — it has a custom error message pointing to a catalog resource instead of listing 362+ enum values
- `domains.ResolveProvider()` and `domains.ResolveProvisioner()` use resolvers internally but keep their exported function signatures for backward compatibility

## Benefits

- **~240 net lines removed** (132 added, 368 deleted across 14 modified files + 15 new files)
- **Adding a new tool now touches exactly one package** — no more shotgun surgery across `server.go`
- **Adding a new enum resolver** is a single `var` declaration instead of a 7-line function
- **Duplicate `joinScopeValues()`** in variable and secret packages eliminated
- **Hardcoded tool count** removed — no more stale log data
- **Fully type-safe** — Go generics enforce correct proto enum types at compile time

## Impact

- All 63 MCP tools: no behavior change
- All existing tests pass unchanged (enum_test.go updated to call resolver methods)
- `go build`, `go test`, `go vet` all clean
- Future domain package additions follow a clear, consistent pattern

## Related Work

- Builds on the existing two-stage codegen pipeline (`proto2schema` → `generator`) which already handles cloud resource input types
- The `EnumResolver` pattern complements codegen — it handles the cases where full code generation would be overkill

---

**Status**: ✅ Production Ready
**Timeline**: Single session
