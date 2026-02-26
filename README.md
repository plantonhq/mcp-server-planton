# Planton Cloud MCP Server

A stateless [MCP](https://modelcontextprotocol.io) server that connects
AI-powered IDEs to [Planton Cloud](https://planton.cloud). It translates
MCP tool calls and resource reads into gRPC requests against the Planton
backend, letting Cursor, Claude Desktop, VS Code, Windsurf, and any
MCP-compliant client create, read, and delete cloud resources across 17 cloud
providers without leaving the editor.

```
AI IDE (Cursor / Claude Desktop / VS Code / Windsurf)
     |  MCP protocol (stdio or Streamable HTTP)
mcp-server-planton
     |  gRPC (TLS on :443, plaintext otherwise)
Planton Cloud Backend
```

This server does not store state. It is a protocol bridge: every tool call
opens a short-lived gRPC connection, performs the RPC, and returns the result.
It can serve both STDIO and HTTP transports concurrently from a single process.

---

## Key Concepts

| Term | Definition |
|------|------------|
| **cloud resource** | Any infrastructure component managed by Planton Cloud (e.g. an EKS cluster, a GCP VPC, an Azure database). |
| **kind** | PascalCase type identifier for a cloud resource (e.g. `AwsEksCluster`, `GcpCloudSqlInstance`). |
| **org** | Organization identifier -- the tenant-level namespace that owns a resource. |
| **env** | Environment identifier (e.g. `production`, `staging`). |
| **slug** | URL-safe unique name for a resource within an (org, env, kind) scope. |
| **api_version** | Versioned API namespace for a cloud provider (e.g. `ai.planton.provider.aws.v1`). |
| **apply** | Idempotent create-or-update. Same semantics as `kubectl apply` -- if the resource exists it is updated, otherwise it is created. |

---

## Installation

### Prerequisites

1. A [Planton Cloud](https://planton.cloud) account
2. An API key (Console > Profile > **API Keys** > **Create Key**)
3. A compatible MCP host (Cursor, Claude Desktop, VS Code, Windsurf, or any
   MCP-compliant client)

### Go Install

```bash
go install github.com/plantonhq/mcp-server-planton/cmd/mcp-server-planton@latest
```

### Pre-built Binary

```bash
# macOS (ARM64)
curl -L https://github.com/plantonhq/mcp-server-planton/releases/latest/download/mcp-server-planton_Darwin_arm64.tar.gz | tar xz
sudo mv mcp-server-planton /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/plantonhq/mcp-server-planton/releases/latest/download/mcp-server-planton_Darwin_x86_64.tar.gz | tar xz
sudo mv mcp-server-planton /usr/local/bin/

# Linux (AMD64)
curl -L https://github.com/plantonhq/mcp-server-planton/releases/latest/download/mcp-server-planton_Linux_x86_64.tar.gz | tar xz
sudo mv mcp-server-planton /usr/local/bin/
```

### Docker

```bash
docker run -i --rm \
  -e PLANTON_API_KEY=your-api-key \
  ghcr.io/plantonhq/mcp-server-planton
```

> **Docker networking:** `localhost` inside a container refers to the
> container's own loopback, not the host machine. To reach a custom endpoint
> running on the host, use `host.docker.internal` on Docker Desktop
> (macOS / Windows) or add `--network host` on Linux.

---

## MCP Client Configuration

All MCP clients use the same JSON structure. The differences are the config file
location and the top-level key.

### Using the standalone binary or Go install

```json
{
  "mcpServers": {
    "planton-cloud": {
      "command": "mcp-server-planton",
      "env": {
        "PLANTON_API_KEY": "YOUR_PLANTON_API_KEY"
      }
    }
  }
}
```

### Using Docker

```json
{
  "mcpServers": {
    "planton-cloud": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "PLANTON_API_KEY",
        "ghcr.io/plantonhq/mcp-server-planton"
      ],
      "env": {
        "PLANTON_API_KEY": "YOUR_PLANTON_API_KEY"
      }
    }
  }
}
```

### Where to put the config

| Client | Config file | Top-level key |
|--------|-------------|---------------|
| Cursor | `.cursor/mcp.json` (workspace) or global settings | `mcpServers` |
| Claude Desktop / Claude Code | `claude_desktop_config.json` | `mcpServers` |
| VS Code / GitHub Copilot | `.vscode/mcp.json` (workspace) or user settings | `servers` |
| Windsurf | Windsurf MCP settings | `mcpServers` |
| LangGraph | `langgraph.json` | `mcp_servers` |

> VS Code uses `"servers"` instead of `"mcpServers"` as the top-level key.

---

## Configuration Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `PLANTON_API_KEY` | *(required for stdio/both)* | API key for Planton Cloud. |
| `PLANTON_CLOUD_ENVIRONMENT` | `live` | Target environment: `live`, `test`, or `local`. |
| `PLANTON_APIS_GRPC_ENDPOINT` | *(from environment)* | Explicit gRPC endpoint override. |
| `PLANTON_MCP_TRANSPORT` | `stdio` | Transport mode: `stdio`, `http`, or `both`. |
| `PLANTON_MCP_HTTP_PORT` | `8080` | HTTP listen port. |
| `PLANTON_MCP_HTTP_AUTH_ENABLED` | `true` | Require Bearer token for HTTP requests. |
| `PLANTON_MCP_LOG_FORMAT` | `text` | Log encoding: `text` or `json`. |
| `PLANTON_MCP_LOG_LEVEL` | `info` | Minimum log level: `debug`, `info`, `warn`, `error`. |

See [docs/configuration.md](docs/configuration.md) for details.

**TLS:** Connections to endpoints on port `443` automatically use TLS with the
system root CA pool. All other ports use plaintext. There is no separate TLS
configuration flag.

---

## Tools

### apply_cloud_resource

Create or update a cloud resource on the Planton platform (idempotent).

| Parameter | Required | Description |
|-----------|----------|-------------|
| `cloud_object` | **yes** | Full OpenMCF cloud resource object. Must contain `api_version`, `kind`, `metadata` (with `name`, `org`, `env`), and `spec`. |

**Agent workflow:**

1. Read `cloud-resource-kinds://catalog` to discover supported kinds and their `api_version` values
2. Read `cloud-resource-schema://{kind}` to get the full spec definition for the desired kind
3. Call `apply_cloud_resource` with the assembled `cloud_object`

### get_cloud_resource

Retrieve a cloud resource from the Planton platform. Identify the resource by
`id` alone, or by all of `kind`, `org`, `env`, and `slug` together.

| Parameter | Required | Description |
|-----------|----------|-------------|
| `id` | conditional | System-assigned resource ID. Provide this alone OR provide all of kind, org, env, slug. |
| `kind` | conditional | PascalCase cloud resource kind (e.g. `AwsEksCluster`). |
| `org` | conditional | Organization identifier. |
| `env` | conditional | Environment identifier. |
| `slug` | conditional | Immutable unique resource slug within (org, env, kind). |

### delete_cloud_resource

Delete a cloud resource from the Planton platform. Same identification options
as `get_cloud_resource`.

| Parameter | Required | Description |
|-----------|----------|-------------|
| `id` | conditional | System-assigned resource ID. Provide this alone OR provide all of kind, org, env, slug. |
| `kind` | conditional | PascalCase cloud resource kind. |
| `org` | conditional | Organization identifier. |
| `env` | conditional | Environment identifier. |
| `slug` | conditional | Immutable unique resource slug within (org, env, kind). |

### Error handling

All tools translate gRPC errors into user-friendly messages:

| gRPC Status | Tool Error Message |
|-------------|-------------------|
| `NotFound` | Resource not found. Verify the identifier is correct. |
| `PermissionDenied` | Permission denied. Check your API key permissions. |
| `Unauthenticated` | Authentication failed. Check your API key. |
| `Unavailable` | Planton backend is unavailable. Ensure it is running and reachable. |
| `InvalidArgument` | The server's validation message is returned directly. |

---

## Resources

MCP clients can read Planton resources directly by URI via `resources/read`.

| URI | Description | MIME Type |
|-----|-------------|-----------|
| `cloud-resource-kinds://catalog` | Catalog of all 362 supported kinds grouped by 17 cloud providers | `application/json` |
| `cloud-resource-schema://{kind}` | JSON schema for a specific kind with field types, validation rules, and defaults | `application/json` |

The kind catalog returns a JSON object with each cloud provider entry containing
an `api_version` and a sorted list of PascalCase kind strings. Use these kind
values with the schema template to fetch the full spec definition before calling
`apply_cloud_resource`.

---

## HTTP Mode

For shared or remote deployments, set the transport to `http`. This runs the
MCP Streamable HTTP transport -- not a REST API.

```bash
PLANTON_MCP_TRANSPORT=http \
PLANTON_API_KEY=your-api-key \
  mcp-server-planton
```

Or with Docker:

```bash
docker run --rm \
  -e PLANTON_MCP_TRANSPORT=http \
  -e PLANTON_MCP_HTTP_AUTH_ENABLED=true \
  -p 8080:8080 \
  ghcr.io/plantonhq/mcp-server-planton
```

Connect your MCP client to `http://host:8080` with an
`Authorization: Bearer <token>` header, where `<token>` is a valid Planton
API key. Each HTTP request carries its own API key, so multiple users can share
a single server instance.

**Auth:** HTTP auth is enabled by default. Set
`PLANTON_MCP_HTTP_AUTH_ENABLED=false` only for trusted internal networks where
all callers are already authenticated at the network level.

**Dual transport:** Set `PLANTON_MCP_TRANSPORT=both` to serve STDIO and HTTP
simultaneously from a single process. This is useful in development when you
want local IDE access (STDIO) and remote access (HTTP) at the same time.

**TLS:** The HTTP transport does not terminate TLS natively. For production
deployments, place a TLS-terminating reverse proxy (e.g. nginx, Envoy, or a
cloud load balancer) in front of the MCP server.

---

## Security

- **STDIO mode**: API key loaded once from `PLANTON_API_KEY` at startup
- **HTTP mode**: Each request carries its own key via `Authorization: Bearer` header -- true multi-tenant support
- Keys are held in memory only during request execution and never persisted
- All API calls are validated and logged with the caller's identity
- Fine-grained authorization is enforced by the Planton backend

---

## Supported Cloud Providers

AWS, GCP, Azure, Kubernetes, AliCloud, DigitalOcean, Civo, Cloudflare,
Confluent, Auth0, OpenFGA, Snowflake, MongoDB Atlas, Hetzner Cloud, OCI,
OpenStack, Scaleway -- 362 resource kinds total.

---

## Development

### Build and test

```bash
make build          # Build binary to bin/mcp-server-planton
make test           # Run tests with race detection
make lint           # Run golangci-lint (falls back to go vet)
make vet            # Run go vet (excludes gen/)
make fmt            # Format all Go source files
make tidy           # Run go mod tidy
```

### Code generation

MCP input types are auto-generated from OpenMCF protobuf definitions via a
two-stage pipeline:

```bash
make codegen-schemas   # Stage 1: Proto -> JSON schemas
make codegen-types     # Stage 2: JSON schemas -> Go input types in gen/
make codegen           # Both stages
```

The `gen/` directory is entirely machine-generated. **Never edit files in `gen/`
by hand** -- they will be overwritten on the next `make codegen` run.

See [docs/development.md](docs/development.md) for codegen details and project
structure.

---

## License

Apache License 2.0. See [LICENSE](LICENSE).
