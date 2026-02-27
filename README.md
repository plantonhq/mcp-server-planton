# Planton Cloud MCP Server

A stateless [MCP](https://modelcontextprotocol.io) server that connects
AI-powered IDEs to [Planton Cloud](https://planton.cloud). It translates
MCP tool calls and resource reads into gRPC requests against the Planton
backend, letting Cursor, Claude Desktop, VS Code, Windsurf, and any
MCP-compliant client manage cloud resources across 17 cloud providers without
leaving the editor — from discovering organizations and environments through
creating, listing, and destroying resources to observing provisioning outcomes
via stack jobs.

```mermaid
flowchart TD
    IDE["AI IDE\n(Cursor / Claude Desktop / VS Code / Windsurf)"]
    MCP["mcp-server-planton\n(stdio or Streamable HTTP)"]
    Backend["Planton Cloud Backend\n(gRPC · TLS on :443)"]
    IDE -->|"MCP protocol"| MCP
    MCP -->|"gRPC"| Backend
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

## Tools & Resources

18 tools cover the full cloud resource lifecycle:

**Cloud Resource Lifecycle**

| Tool | What It Does |
|------|--------------|
| `apply_cloud_resource` | Create or update a resource (idempotent — same semantics as `kubectl apply`) |
| `get_cloud_resource` | Retrieve a resource by ID or by `(kind, org, env, slug)` |
| `delete_cloud_resource` | Delete a resource record (does not tear down infrastructure) |
| `list_cloud_resources` | List resources in an org, with optional environment/kind/text filters |
| `destroy_cloud_resource` | Tear down cloud infrastructure (Terraform/Pulumi destroy) while keeping the record |
| `check_slug_availability` | Verify a slug is available within `(org, env, kind)` before creating |
| `rename_cloud_resource` | Change a resource's display name (slug is immutable) |
| `list_cloud_resource_locks` | Show lock status, holder, and wait queue for a resource |
| `remove_cloud_resource_locks` | Force-clear stuck locks on a resource |
| `get_env_var_map` | Extract environment variables and secrets from a resource manifest |
| `resolve_value_references` | Resolve all valueFrom references in a resource's spec |

**Stack Job Observability**

| Tool | What It Does |
|------|--------------|
| `get_stack_job` | Retrieve a stack job by ID |
| `get_latest_stack_job` | Get the most recent stack job for a resource (primary polling tool after apply/destroy) |
| `list_stack_jobs` | List stack jobs with filters (org, env, kind, status, result) |

**Context Discovery**

| Tool | What It Does |
|------|--------------|
| `list_organizations` | List organizations the caller belongs to |
| `list_environments` | List environments the caller can access within an organization |

**Presets**

| Tool | What It Does |
|------|--------------|
| `search_cloud_object_presets` | Search for preset templates (official and org-scoped) |
| `get_cloud_object_preset` | Get full preset content by ID, for use as a template with `apply_cloud_resource` |

Two read-only MCP resources drive schema discovery before any tool call:

| URI | What It Returns |
|-----|-----------------|
| `cloud-resource-kinds://catalog` | All 362 supported kinds grouped by 17 cloud providers |
| `cloud-resource-schema://{kind}` | Full JSON schema for a specific kind — field types, validation rules, and defaults |

For full parameter reference, agent workflow guidance, and error handling, see
[docs/tools.md](docs/tools.md).

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
