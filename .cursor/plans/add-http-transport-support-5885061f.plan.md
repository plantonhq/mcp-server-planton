<!-- 5885061f-dea9-448a-82da-1f36463ccf5f 75ec4590-7be8-4584-8cb7-bbab383acf8a -->
# Add HTTP Transport Support to MCP Server

## Industry Context & Recommendations

Based on research of GitHub's MCP server and other production implementations, here's how the ecosystem works:

### Transport Strategy: **Support Both STDIO + HTTP**

**Recommendation:** Keep STDIO and add HTTP as a configurable option.

**Why both?**

- **STDIO** → Local AI clients (Claude Desktop, Cursor, LangGraph local agents)
- **HTTP** → Production deployments (cloud runners, webhooks, remote access, serverless)

**Industry pattern:** Most mature MCP servers (AWS, GitHub) support both transports with environment variable selection.

### Authentication Architecture: **Two-Layer Security**

**Recommendation:** Implement bearer token auth for HTTP endpoints + keep PLANTON_API_KEY for backend.

**Security model:**

```
Client → Bearer Token → MCP HTTP Server → PLANTON_API_KEY → Planton APIs
         (Who can access?)              (What can they access?)
```

**Why two layers?**

1. **Bearer token** - Controls WHO can use your MCP server (prevents unauthorized access to your endpoint)
2. **PLANTON_API_KEY** - Controls WHAT resources the user can access (FGA enforcement from Planton Cloud)

**Industry standard:** OAuth 2.1 with bearer tokens is the MCP spec standard for HTTP transport.

### Session Management: **Stateless Mode**

**Recommendation:** Use stateless HTTP transport for cloud deployments.

**Why stateless?**

- Scales horizontally (multiple instances)
- Works perfectly with serverless (AWS Lambda, Cloud Run)
- No session storage overhead
- Simpler operations

**Trade-off:** Basic MCP features only (tools, resources) - but that's all you need for your use case.

### Deployment Strategy: **Multi-Environment Support**

**Three deployment modes:**

1. **Local Development/Desktop AI** (Current)

   - Transport: STDIO
   - Use case: Claude Desktop, Cursor, local testing
   - Auth: PLANTON_API_KEY only

2. **Long-Running Service** (New)

   - Transport: HTTP server on port (e.g., 8080)
   - Use case: Team shared service, webhook endpoints, API integrations
   - Auth: Bearer token + PLANTON_API_KEY
   - Deployment: Docker container on cloud runner (GitHub Actions self-hosted, GitLab Runner)

3. **Serverless/Webhooks** (Future-ready)

   - Transport: HTTP (stateless)
   - Use case: Event-driven, webhook callbacks, auto-scaling
   - Auth: Bearer token + PLANTON_API_KEY
   - Deployment: AWS Lambda, Cloud Functions, Cloud Run

## Implementation Plan

### 1. Update Configuration (`internal/config/config.go`)

Add new environment variables:

```
PLANTON_MCP_TRANSPORT=stdio|http|both (default: stdio)
PLANTON_MCP_HTTP_PORT=8080
PLANTON_MCP_HTTP_AUTH_ENABLED=true
PLANTON_MCP_HTTP_BEARER_TOKEN=<secret-token>
```

### 2. Create HTTP Server Handler (`internal/mcp/http_server.go`)

Implement `StreamableHTTPServer` using `mark3labs/mcp-go`:

- Bearer token authentication middleware
- Health check endpoint (`/health`)
- Metrics endpoint (`/metrics`) 
- Stateless session management
- Graceful shutdown

### 3. Update Main Server (`internal/mcp/server.go`)

Add `ServeHTTP()` method alongside existing `Serve()` (STDIO):

- Check `PLANTON_MCP_TRANSPORT` config
- Support running both transports simultaneously when `transport=both`

### 4. Update Main Entry Point (`cmd/mcp-server-planton/main.go`)

Add transport selection logic:

- If `http` or `both`: start HTTP server (optionally in goroutine)
- If `stdio` or `both`: start STDIO server
- Handle graceful shutdown for both

### 5. Update Documentation

**README.md:**

- Add HTTP transport usage section
- Document bearer token authentication
- Add deployment examples (Docker, serverless)
- Update integration examples for remote access

**docs/configuration.md:**

- Document all new environment variables
- Explain authentication model
- Provide security best practices

**docs/deployment.md:**

- Add cloud deployment guide
- AWS Lambda example
- Google Cloud Run example
- Webhook integration patterns

### 6. Update Docker Configuration

**Dockerfile:**

- Expose port 8080
- Add HEALTHCHECK directive

**.github/workflows/release.yml:**

- Update Docker image with proper port exposure

### 7. Add Examples

Create `examples/` directory:

- `http-server/` - Standalone HTTP server setup
- `webhook-lambda/` - AWS Lambda webhook handler
- `cloud-run/` - Google Cloud Run deployment

## Security Considerations

1. **Bearer token generation:** Recommend using strong random tokens (32+ chars)
2. **HTTPS requirement:** Document that HTTP mode should only be used behind TLS termination
3. **Token rotation:** Document bearer token rotation procedures
4. **Rate limiting:** Consider adding rate limiting middleware (future enhancement)

## Testing Strategy

1. Manual testing: HTTP server with curl/Postman
2. Integration tests: Both STDIO and HTTP transports
3. Security testing: Verify bearer token validation
4. Load testing: Stateless HTTP performance

## Files to Modify

- `internal/config/config.go` - Add HTTP config
- `internal/mcp/server.go` - Update server initialization
- `internal/mcp/http_server.go` - New HTTP handler
- `cmd/mcp-server-planton/main.go` - Transport selection
- `README.md` - Usage documentation
- `docs/configuration.md` - Config reference
- `docs/deployment.md` - Deployment guide
- `Dockerfile` - Port exposure
- `.github/workflows/release.yml` - Docker update

### To-dos

- [ ] Add HTTP transport configuration to internal/config/config.go
- [ ] Create internal/mcp/http_server.go with StreamableHTTPServer and bearer token auth
- [ ] Add ServeHTTP method to internal/mcp/server.go
- [ ] Update cmd/mcp-server-planton/main.go with transport selection logic
- [ ] Expose port 8080 and add health check to Dockerfile
- [ ] Add HTTP transport section to README.md
- [ ] Document new environment variables in docs/configuration.md
- [ ] Add cloud deployment examples to docs/deployment.md