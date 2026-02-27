---
name: H4 get.go refactor close-out
overview: Refactor `get.go` to eliminate duplicated dual-path resolution logic by delegating to `resolveResource`, then close out the project.
todos:
  - id: h4-refactor
    content: Refactor get.go to delegate to resolveResource, remove unused imports
    status: completed
  - id: verify-build
    content: Run go build ./... and go test ./... to verify correctness
    status: completed
  - id: close-out
    content: Update next-task.md status to Complete, write final checkpoint
    status: completed
isProject: false
---

# H4: Refactor `get.go` + Project Close-Out

## Why This Matters

`get.go` contains 25 lines of inline dual-path resolution (ID vs slug) that is a near-exact duplicate of `[resolveResource](internal/domains/cloudresource/identifier.go)` in `identifier.go`. Every other domain function that needs the full resource proto (`Destroy`) already delegates to `resolveResource`. If this resolution logic ever changes — new identification path, different error handling — `get.go` would silently drift.

This is the only remaining structural inconsistency in the `cloudresource` package.

## The Duplication

**Current `get.go`** (lines 19-49) does this inline:

```go
if id.ID != "" {
    cr, err := client.Get(ctx, &cloudresourcev1.CloudResourceId{Value: id.ID})
    // error handling + marshal
}
kind, err := domains.ResolveKind(id.Kind)
cr, err := client.GetByOrgByEnvByKindBySlug(ctx, ...)
// error handling + marshal
```

`**resolveResource**` in `identifier.go` (lines 120-147) does the exact same thing, and `destroy.go` already consumes it cleanly:

```go
cr, err := resolveResource(ctx, conn, id)
if err != nil {
    return "", err
}
return domains.MarshalJSON(cr)
```

## The Change

Replace the body of `Get()` in `[internal/domains/cloudresource/get.go](internal/domains/cloudresource/get.go)` with:

```go
func Get(ctx context.Context, serverAddress string, id ResourceIdentifier) (string, error) {
    return domains.WithConnection(ctx, serverAddress,
        func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
            cr, err := resolveResource(ctx, conn, id)
            if err != nil {
                return "", err
            }
            return domains.MarshalJSON(cr)
        })
}
```

- The import of `domains.ResolveKind` and `cloudresourcev1` proto types becomes unnecessary (they are used inside `resolveResource`).
- The doc comment stays accurate — it still describes the two identification paths, just notes delegation to `resolveResource`.

## Verification

- `go build ./...` must pass
- `go test ./...` must pass
- Zero linter errors

## Project Close-Out

After the refactor, update `[_projects/2026-02/20260227.01.expand-cloud-resource-tools/next-task.md](_projects/2026-02/20260227.01.expand-cloud-resource-tools/next-task.md)`:

- Mark H4 as complete
- Set project status to "Complete"
- Write a final checkpoint

