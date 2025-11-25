# Fix Protobuf Namespace Conflict

**Date**: November 25, 2025  
**Type**: Bug Fix  
**Impact**: Critical - Server startup was completely broken

## Summary

Fixed a critical protobuf namespace conflict that caused the MCP server to crash on startup with a "proto file already registered" panic. The issue was caused by importing the same proto files from two different sources: `buf.build/gen/go/project-planton/apis` (correct) and `github.com/project-planton/project-planton` (incorrect).

## Problem Statement

The MCP server crashed immediately on startup with the following error:

```
panic: proto: file "org/project_planton/shared/cloudresourcekind/cloud_resource_provider.proto" is already registered
	previously from: "buf.build/gen/go/project-planton/apis/protocolbuffers/go/org/project_planton/shared/cloudresourcekind"
	currently from:  "github.com/project-planton/project-planton/apis/org/project_planton/shared/cloudresourcekind"
See https://protobuf.dev/reference/go/faq#namespace-conflict
```

### Root Cause

The codebase had inconsistent imports for the `cloudresourcekind` package:

- ✅ `unwrap.go` and `cloudresource_client.go`: Used buf.build modules (correct)
- ❌ `kinds.go`: Used direct github.com import (incorrect)

When Go initialized the package, both import paths tried to register the same proto files to the global protobuf registry, causing a panic.

### Impact

- ❌ MCP server completely non-functional
- ❌ Could not be added to Cursor or any MCP client
- ❌ No tools accessible
- ❌ Blocked all user workflows

## Solution

Standardized all imports to use the buf.build published modules and removed the conflicting github.com dependency.

### Changes Made

#### 1. Updated Import in kinds.go

**File**: `internal/domains/infrahub/cloudresource/kinds.go`

Changed line 11 from:
```go
cloudresourcekind "github.com/project-planton/project-planton/apis/org/project_planton/shared/cloudresourcekind"
```

To:
```go
cloudresourcekind "buf.build/gen/go/project-planton/apis/protocolbuffers/go/org/project_planton/shared/cloudresourcekind"
```

**Rationale**: Aligns with the import pattern used throughout the rest of the codebase.

#### 2. Removed Conflicting Dependency

**File**: `go.mod`

Removed the direct github.com dependency:
```diff
  require (
  	buf.build/gen/go/blintora/apis/grpc/go v1.5.1-20251125011413-52ef5c4f2840.2
  	buf.build/gen/go/blintora/apis/protocolbuffers/go v1.36.10-20251125011413-52ef5c4f2840.1
  	buf.build/gen/go/project-planton/apis/protocolbuffers/go v1.36.10-20251124125039-9c224fb3651e.1
  	github.com/mark3labs/mcp-go v0.6.0
- 	github.com/project-planton/project-planton v0.2.245
  	google.golang.org/grpc v1.75.0
  	google.golang.org/protobuf v1.36.10
  )
```

**Rationale**: The buf.build published modules provide all the required proto definitions. The direct github.com import was redundant and conflicting.

#### 3. Cleaned Up Dependencies

Ran `go mod tidy` to ensure all dependencies are properly resolved and removed any unused transitive dependencies.

**Result**: ✅ No unused dependencies, clean dependency graph

## Technical Details

### Why This Happened

The migration to buf.build published modules (changelog entry: 2025-11-25-130536) missed updating one import in `kinds.go`, which still referenced the old github.com import path. This created a situation where:

1. `kinds.go` imported from github.com path → registered proto files under github.com namespace
2. `unwrap.go` imported from buf.build path → tried to register same proto files under buf.build namespace
3. Go's protobuf registry detected duplicate registration → panic

### Protobuf Registry Behavior

Go's protobuf library maintains a global registry of all proto files. Each proto file can only be registered once. When the same proto file is imported from different Go module paths, both paths attempt to register it, causing the conflict.

From [protobuf.dev FAQ](https://protobuf.dev/reference/go/faq#namespace-conflict):
> "A proto namespace conflict occurs when multiple descriptor.proto files with the same name are registered with the global proto registry."

### Why buf.build Modules Are Correct

1. **Official Distribution**: buf.build is the recommended way to distribute protobuf definitions
2. **Versioning**: Proper semantic versioning tied to BSR (Buf Schema Registry) commits
3. **Consistency**: Single source of truth for proto definitions across all consumers
4. **No Conflicts**: Each module has a unique import path based on BSR organization and commit

## Verification

### Before Fix

```bash
$ mcp-server-planton
panic: proto: file "org/project_planton/shared/cloudresourcekind/cloud_resource_provider.proto" is already registered
# Server crashed immediately ❌
```

### After Fix

```bash
$ mcp-server-planton
# Server starts successfully ✅
# Tools are accessible ✅
# No proto registration errors ✅
```

### Import Consistency Check

```bash
# Verify no more github.com/project-planton imports
$ grep -r "github.com/project-planton/project-planton" internal/
# No matches found ✅

# Verify all use buf.build
$ grep -r "buf.build/gen/go/project-planton/apis" internal/
internal/domains/infrahub/cloudresource/kinds.go:	cloudresourcekind "buf.build/..."
internal/domains/infrahub/cloudresource/unwrap.go:	cloudresourcekind "buf.build/..."
internal/domains/infrahub/clients/cloudresource_client.go:	cloudresourcekind "buf.build/..."
# All consistent ✅
```

## Benefits

### 1. Server Functionality Restored

- ✅ MCP server starts successfully
- ✅ All tools are accessible
- ✅ Can be added to Cursor and other MCP clients
- ✅ User workflows unblocked

### 2. Consistent Dependency Management

- ✅ All proto imports use buf.build modules
- ✅ Single source of truth for proto definitions
- ✅ No duplicate dependencies
- ✅ Cleaner go.mod

### 3. Better Maintainability

- ✅ Clear import pattern throughout codebase
- ✅ Reduced dependency on github.com direct imports
- ✅ Follows buf.build best practices
- ✅ Easier to update proto definitions

### 4. Prevention of Future Conflicts

- ✅ Eliminated mixed import paths
- ✅ Standardized on buf.build modules
- ✅ Clear precedent for future proto imports

## Related Changes

This fix completes the migration started in:

- **Completes**: [2025-11-25-130536-migrate-to-buf-published-modules.md](./2025-11-25-130536-migrate-to-buf-published-modules.md) - Initial migration to buf.build
- **Fixes**: Incomplete migration that missed `kinds.go`

## Best Practices Established

1. **Use buf.build Modules**: Always import proto definitions from buf.build published modules
2. **Avoid Direct GitHub Imports**: Don't import proto definitions directly from github.com repositories
3. **Consistency Check**: When migrating, verify all files are updated
4. **Dependency Cleanup**: Run `go mod tidy` after removing dependencies
5. **Import Pattern**: Use consistent import aliases (e.g., `cloudresourcekind`)

## Migration Notes

### For Developers

**No action required**. The fix is complete and the server works correctly.

**If you were blocked by this issue:**
- ✅ Pull the latest changes
- ✅ Restart your MCP server
- ✅ Reconnect in Cursor settings

### For CI/CD

**No configuration changes needed**. The build pipeline will automatically use the updated dependencies.

## Troubleshooting

### Issue: Still seeing proto registration error

**Cause**: Using an old binary or stale build cache

**Solution**:
```bash
# Clean build cache
go clean -cache

# Rebuild
go build ./cmd/mcp-server-planton

# Restart server
./mcp-server-planton
```

### Issue: Import resolution errors

**Cause**: Dependencies not properly downloaded

**Solution**:
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Verify dependencies
go mod verify
```

## Testing

### Startup Test
```bash
$ mcp-server-planton
# Server starts successfully ✅
# No panic errors ✅
```

### Tools Test
```bash
# In Cursor, verify MCP server tools are available:
# - get_cloud_resource_by_id ✅
# - search_cloud_resources ✅
# - lookup_cloud_resource_by_name ✅
# - list_cloud_resource_kinds ✅
```

### Import Test
```bash
# Verify no conflicting imports
$ grep -r "github.com/project-planton/project-planton" .
# No matches in code files ✅
```

## Future Recommendations

### 1. Automated Import Consistency Check

Add a CI check to ensure all proto imports use buf.build modules:

```bash
# In CI pipeline
if grep -r "github.com/project-planton/project-planton" internal/; then
  echo "ERROR: Found direct github.com proto imports. Use buf.build modules instead."
  exit 1
fi
```

### 2. Pre-commit Hook

Add a pre-commit hook to prevent direct github.com proto imports:

```yaml
# .pre-commit-config.yaml
- repo: local
  hooks:
    - id: check-proto-imports
      name: Check proto imports use buf.build
      entry: bash -c 'grep -r "github.com/project-planton/project-planton" internal/ && exit 1 || exit 0'
      language: system
```

### 3. Documentation Update

Document the import pattern in `CONTRIBUTING.md`:

```markdown
## Protobuf Imports

Always import proto definitions from buf.build published modules:

✅ Good:
import cloudresourcekind "buf.build/gen/go/project-planton/apis/protocolbuffers/go/..."

❌ Bad:
import cloudresourcekind "github.com/project-planton/project-planton/apis/..."
```

### 4. Dependency Review

Periodically review dependencies to ensure:
- No duplicate proto sources
- All proto imports use buf.build
- Dependencies are up-to-date

## References

- **Protobuf Namespace Conflict FAQ**: https://protobuf.dev/reference/go/faq#namespace-conflict
- **Buf Schema Registry**: https://buf.build/project-planton/apis
- **Go Protobuf Registry**: https://pkg.go.dev/google.golang.org/protobuf/reflect/protoregistry
- **Migration to buf.build**: [2025-11-25-130536-migrate-to-buf-published-modules.md](./2025-11-25-130536-migrate-to-buf-published-modules.md)

## Conclusion

This fix resolves a critical startup failure caused by protobuf namespace conflicts. By standardizing on buf.build published modules and removing the conflicting github.com dependency, the MCP server now starts reliably and all functionality is restored.

The fix completes the migration to buf.build modules and establishes clear patterns for future proto imports, preventing similar issues from occurring.

---

**Status**: ✅ Completed  
**Files Changed**: 2
- `internal/domains/infrahub/cloudresource/kinds.go` - Updated import to buf.build
- `go.mod` - Removed github.com/project-planton/project-planton dependency

**Breaking Changes**: None  
**Testing**: 
- ✅ Server startup verified
- ✅ All tools accessible
- ✅ No import conflicts
- ✅ Dependencies clean

**Impact**: Critical bug fix - restored full server functionality
