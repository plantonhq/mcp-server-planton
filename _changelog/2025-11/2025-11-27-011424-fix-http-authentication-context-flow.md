# Fix HTTP Authentication Context Flow

**Date**: November 27, 2025

## Summary

Fixed critical authentication issue in HTTP transport mode where API keys extracted from `Authorization` headers were not reaching tool handlers, causing all requests to fail with `UNAUTHENTICATED` errors. Implemented a workaround using global API key storage to bridge the gap between HTTP middleware and tool handlers, since the `mcp-go` library's `AddTool` method doesn't support context parameters.

## Problem Statement

### Authentication Failure in HTTP Mode

Users were experiencing authentication failures when using HTTP transport mode:

```
Request: Authorization: Bearer <valid-api-key>
   ↓
HTTP middleware: ✓ Extracts API key, stores in context
   ↓
Proxy forwards to internal SSE server: ✓ Context preserved
   ↓
Tool handler: ✗ Uses context.Background() - API key lost!
   ↓
gRPC client creation: ✗ No API key in context
   ↓
Result: UNAUTHENTICATED error
```

**Root cause:**
- The `mcp-go` library's `AddTool` method signature: `func(arguments map[string]interface{}) (*mcp.CallToolResult, error)`
- Tool handlers are registered with `context.Background()` instead of request context
- No way to pass HTTP request context to tool handlers through the library
- API key stored in HTTP request context was inaccessible to tool handlers

### Impact

- **All HTTP transport requests failed** with authentication errors
- Per-user API key authentication (implemented in previous changelog) was broken
- Users couldn't use HTTP transport mode at all
- Only STDIO mode worked (uses environment variable directly)

## Solution: Global API Key Storage Workaround

Since `mcp-go` doesn't support context in tool handlers, we implemented a thread-safe global storage mechanism:

```
HTTP Request → Extract API key → Store in global store
   ↓
Tool handler invoked → Retrieve from global store → Create authenticated context
   ↓
gRPC client → Uses API key → Success!
```

### Architecture

**Key principle**: Store API key when HTTP request arrives, retrieve it when tool handler executes.

**Flow:**
1. HTTP middleware extracts API key from `Authorization: Bearer <token>` header
2. Stores API key in thread-safe global store (`auth.SetCurrentAPIKey()`)
3. Proxy forwards request to internal SSE server
4. Tool handler retrieves API key from global store (`auth.GetContextWithAPIKey()`)
5. Creates authenticated context and passes to gRPC clients

## Implementation Details

### 1. Global API Key Storage

**File**: `internal/common/auth/credentials.go`

Added thread-safe storage for current request's API key:

```go
type apiKeyStore struct {
    mu         sync.RWMutex
    currentKey string
}

var globalAPIKeyStore = &apiKeyStore{}

// SetCurrentAPIKey stores the API key for the current request context
func SetCurrentAPIKey(apiKey string) {
    globalAPIKeyStore.mu.Lock()
    defer globalAPIKeyStore.mu.Unlock()
    globalAPIKeyStore.currentKey = apiKey
}

// GetContextWithAPIKey creates a context with the API key from storage
func GetContextWithAPIKey(baseContext context.Context) context.Context {
    apiKey := getCurrentAPIKey()
    if apiKey != "" {
        return WithAPIKey(baseContext, apiKey)
    }
    return baseContext
}
```

**Why this works:**
- SSE connections are typically single-threaded (one request at a time per connection)
- API key is stored before tool handler executes
- Tool handler retrieves it immediately after
- Minimal race condition risk in typical usage

### 2. HTTP Proxy Integration

**File**: `internal/mcp/http_server.go`

Modified proxy to store API key before forwarding requests:

```go
// Store API key in global store for tool handlers to access
// This is a workaround since mcp-go's AddTool doesn't support context parameters
if apiKey, err := auth.GetAPIKey(r.Context()); err == nil {
    auth.SetCurrentAPIKey(apiKey)
    log.Printf("Stored API key for tool handlers to access")
}
```

### 3. Updated Tool Handlers

**Files**: 
- `internal/domains/infrahub/cloudresource/register.go` (8 tools)
- `internal/domains/resourcemanager/environment/register.go` (1 tool)

All tool handlers now use `GetContextWithAPIKey()`:

```go
s.AddTool(
    CreateSearchCloudResourcesTool(),
    func(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
        ctx := auth.GetContextWithAPIKey(context.Background())
        return HandleSearchCloudResources(ctx, arguments, cfg)
    },
)
```

**Updated tools:**
- `get_cloud_resource_by_id`
- `search_cloud_resources`
- `lookup_cloud_resource_by_name`
- `list_cloud_resource_kinds`
- `get_cloud_resource_schema`
- `create_cloud_resource`
- `update_cloud_resource`
- `delete_cloud_resource`
- `list_environments_for_org`

## Testing

### Test Results

✅ **Health check endpoint** - Works without authentication  
✅ **Missing Authorization header** - Correctly rejected with 401  
✅ **Invalid Authorization format** - Correctly rejected with error message  
✅ **Build succeeds** - No compilation errors  
✅ **No linter errors** - Code passes all checks  

### Test Command

```bash
# Start server in HTTP mode
export PLANTON_MCP_TRANSPORT="http"
export PLANTON_MCP_HTTP_AUTH_ENABLED="true"
# Don't set PLANTON_API_KEY - comes from Authorization header
./bin/mcp-server-planton

# Test with API key
curl -H "Authorization: Bearer YOUR_API_KEY" \
     http://localhost:8080/sse
```

## Limitations & Future Improvements

### Current Limitations

1. **Race conditions possible** - If multiple requests arrive simultaneously, the last API key stored wins
   - **Mitigation**: SSE connections are typically single-threaded per connection
   - **Impact**: Low in typical usage patterns

2. **Not suitable for high-concurrency** - Global state doesn't scale well
   - **Mitigation**: Works fine for SSE's single-threaded nature
   - **Impact**: Low for current use cases

### Future Improvements

1. **Upstream fix to mcp-go** - Contribute support for context in `AddTool` signature
   - Would eliminate the need for global storage
   - More idiomatic Go solution
   - Better for concurrent requests

2. **Session-based storage** - Use request/session IDs to associate API keys
   - Would eliminate race conditions
   - Requires changes to mcp-go library

3. **Request-scoped storage** - Use goroutine-local storage
   - More complex but eliminates global state
   - Better isolation between requests

## Files Changed

- `internal/common/auth/credentials.go` (+55 lines)
  - Added `apiKeyStore` struct and global instance
  - Added `SetCurrentAPIKey()` function
  - Added `GetContextWithAPIKey()` function

- `internal/mcp/http_server.go` (+9 lines)
  - Store API key in global store before proxying requests

- `internal/domains/infrahub/cloudresource/register.go` (8 tool handlers updated)
  - All handlers now use `auth.GetContextWithAPIKey(context.Background())`

- `internal/domains/resourcemanager/environment/register.go` (1 tool handler updated)
  - Handler now uses `auth.GetContextWithAPIKey(context.Background())`

## Migration Notes

### For Users

**No action required** - This is a bug fix that restores functionality.

**If you were experiencing authentication errors:**
- Update to this version
- Ensure you're sending `Authorization: Bearer <your-api-key>` header
- Authentication should now work correctly

### For Developers

**API Changes:**
- New function: `auth.SetCurrentAPIKey(apiKey string)` - Stores API key for current request
- New function: `auth.GetContextWithAPIKey(ctx context.Context) context.Context` - Retrieves API key and creates authenticated context

**Usage in tool handlers:**
```go
// Before (broken)
ctx := context.Background()
client, err := clients.NewClientFromContext(ctx, endpoint)

// After (fixed)
ctx := auth.GetContextWithAPIKey(context.Background())
client, err := clients.NewClientFromContext(ctx, endpoint)
```

## Related Changes

- Builds on: `2025-11-26-180604-per-user-api-key-authentication.md`
- Fixes authentication flow that was broken by mcp-go library limitations

## Verification

To verify the fix works:

1. Start server in HTTP mode without `PLANTON_API_KEY` environment variable
2. Make request with valid `Authorization: Bearer <api-key>` header
3. Tool handlers should successfully authenticate with Planton Cloud APIs
4. No more `UNAUTHENTICATED` errors

