# Configuration

All settings are read from environment variables with a `PLANTON_` prefix. Reasonable defaults are provided for development use.

## Environment Variables

### Authentication

#### `PLANTON_API_KEY`

API key for authenticating with the Planton Cloud backend.

- **STDIO mode**: Required. Loaded once at startup and used for all API calls.
- **HTTP mode**: Not required in the environment. Each HTTP request carries its own key via the `Authorization: Bearer` header.
- **Both mode**: Required (for the STDIO transport).

Obtain your key from the [Planton Cloud Console](https://console.planton.cloud) under Profile > API Keys.

### Backend Endpoint

#### `PLANTON_CLOUD_ENVIRONMENT`

Selects a preset gRPC endpoint:

| Value | Endpoint |
|-------|----------|
| `live` (default) | `api.live.planton.ai:443` |
| `test` | `api.test.planton.cloud:443` |
| `local` | `localhost:8080` |

#### `PLANTON_APIS_GRPC_ENDPOINT`

Explicit gRPC dial target. Takes precedence over `PLANTON_CLOUD_ENVIRONMENT` when set.

```bash
export PLANTON_APIS_GRPC_ENDPOINT="custom-api.example.com:443"
```

### Transport

#### `PLANTON_MCP_TRANSPORT`

Communication mode between MCP clients and the server.

| Value | Description |
|-------|-------------|
| `stdio` (default) | Communicates over stdin/stdout. The MCP client spawns the server as a child process. |
| `http` | Streamable HTTP transport for remote/shared deployments. Each request carries its own Bearer token. |
| `both` | Runs STDIO and HTTP simultaneously. |

#### `PLANTON_MCP_HTTP_PORT`

TCP port for the HTTP transport. Default: `8080`.

#### `PLANTON_MCP_HTTP_AUTH_ENABLED`

Whether HTTP requests require a valid `Authorization: Bearer` token. Default: `true`.

Set to `false` when running behind a trusted reverse proxy that already validates tokens.

### Logging

#### `PLANTON_MCP_LOG_FORMAT`

Structured log output encoding. Default: `text`.

| Value | Description |
|-------|-------------|
| `text` | Human-readable key=value format |
| `json` | JSON lines (for log aggregation pipelines) |

#### `PLANTON_MCP_LOG_LEVEL`

Minimum severity for emitted log records. Default: `info`.

| Value | Description |
|-------|-------------|
| `debug` | Verbose output including internal state |
| `info` | Normal operational messages |
| `warn` | Warnings and recoverable errors |
| `error` | Errors only |

All log output goes to stderr so that stdout remains available for the STDIO MCP transport.

## Transport Security

The gRPC client determines transport security by convention:

- Endpoints on port **443** use TLS with the system root CA pool
- All other ports use plaintext (suitable for localhost and internal networks)

## Examples

### Local Development (STDIO)

```bash
export PLANTON_API_KEY="your-key"
export PLANTON_CLOUD_ENVIRONMENT="local"
mcp-server-planton stdio
```

### Remote Deployment (HTTP)

```bash
export PLANTON_MCP_TRANSPORT="http"
export PLANTON_MCP_HTTP_PORT="8080"
export PLANTON_MCP_HTTP_AUTH_ENABLED="true"
export PLANTON_MCP_LOG_FORMAT="json"
export PLANTON_MCP_LOG_LEVEL="info"
mcp-server-planton http
```

### Docker

```bash
docker run -p 8080:8080 \
  -e PLANTON_MCP_TRANSPORT="http" \
  -e PLANTON_MCP_HTTP_AUTH_ENABLED="true" \
  -e PLANTON_MCP_LOG_FORMAT="json" \
  ghcr.io/plantoncloud/mcp-server-planton:latest
```
