# Planton MCP Server

A stateless [MCP](https://modelcontextprotocol.io) server that connects
AI-powered IDEs to [Planton](https://planton.ai). It translates
MCP tool calls and resource reads into gRPC requests against the Planton
backend, letting Cursor, Claude Desktop, VS Code, Windsurf, and any
MCP-compliant client manage cloud resources across 17 cloud providers without
leaving the editor -- from discovering organizations and environments through
creating, listing, and destroying resources to observing provisioning outcomes
via stack jobs.

```mermaid
flowchart TD
    IDE["AI IDE\n(Cursor / Claude Desktop / VS Code / Windsurf)"]
    MCP["planton-mcp-server\n(stdio or Streamable HTTP)"]
    Backend["Planton Backend\n(gRPC · TLS on :443)"]
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
| **cloud resource** | Any infrastructure component managed by Planton (e.g. an EKS cluster, a GCP VPC, an Azure database). |
| **kind** | PascalCase type identifier for a cloud resource (e.g. `AwsEksCluster`, `GcpCloudSqlInstance`). |
| **org** | Organization identifier -- the tenant-level namespace that owns a resource. |
| **env** | Environment identifier (e.g. `production`, `staging`). |
| **slug** | URL-safe unique name for a resource within an (org, env, kind) scope. |
| **api_version** | Versioned API namespace for a cloud provider (e.g. `ai.planton.provider.aws.v1`). |
| **apply** | Idempotent create-or-update. Same semantics as `kubectl apply` -- if the resource exists it is updated, otherwise it is created. |

---

## Installation

### Prerequisites

1. A [Planton](https://planton.ai) account
2. An API key (Console > Profile > **API Keys** > **Create Key**)
3. A compatible MCP host (Cursor, Claude Desktop, VS Code, Windsurf, or any
   MCP-compliant client)

### Pre-built Binary

```bash
# macOS (ARM64)
curl -L https://github.com/plantonhq/planton-mcp-server/releases/latest/download/planton-mcp-server_Darwin_arm64.tar.gz | tar xz
sudo mv planton-mcp-server /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/plantonhq/planton-mcp-server/releases/latest/download/planton-mcp-server_Darwin_x86_64.tar.gz | tar xz
sudo mv planton-mcp-server /usr/local/bin/

# Linux (AMD64)
curl -L https://github.com/plantonhq/planton-mcp-server/releases/latest/download/planton-mcp-server_Linux_x86_64.tar.gz | tar xz
sudo mv planton-mcp-server /usr/local/bin/
```

### Docker

```bash
docker run -i --rm \
  -e PLANTON_API_KEY=your-api-key \
  ghcr.io/plantonhq/planton-mcp-server
```

> **Docker networking:** `localhost` inside a container refers to the
> container's own loopback, not the host machine. To reach a custom endpoint
> running on the host, use `host.docker.internal` on Docker Desktop
> (macOS / Windows) or add `--network host` on Linux.

### Hosted (Zero Install)

Planton hosts a shared MCP server at `mcp.planton.ai`. No binary or Docker
required -- connect directly from any MCP client over HTTPS:

```json
{
  "mcpServers": {
    "planton": {
      "type": "http",
      "url": "https://mcp.planton.ai/",
      "headers": {
        "Authorization": "Bearer YOUR_PLANTON_API_KEY"
      }
    }
  }
}
```

Each request carries its own API key, so multiple users share the same endpoint.

---

## MCP Client Configuration

All MCP clients use the same JSON structure. The differences are the config file
location and the top-level key.

### Using the standalone binary

```json
{
  "mcpServers": {
    "planton": {
      "command": "planton-mcp-server",
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
    "planton": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "PLANTON_API_KEY",
        "ghcr.io/plantonhq/planton-mcp-server"
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
| `PLANTON_API_KEY` | *(required for stdio/both)* | API key for Planton. |
| `PLANTON_ENVIRONMENT` | `live` | Target environment: `live`, `test`, or `local`. |
| `PLANTON_APIS_GRPC_ENDPOINT` | *(from environment)* | Explicit gRPC endpoint override. |
| `PLANTON_MCP_TRANSPORT` | `stdio` | Transport mode: `stdio`, `http`, or `both`. |
| `PLANTON_MCP_HTTP_PORT` | `8080` | HTTP listen port. |
| `PLANTON_MCP_HTTP_AUTH_ENABLED` | `true` | Require Bearer token for HTTP requests. |
| `PLANTON_MCP_LOG_FORMAT` | `text` | Log encoding: `text` or `json`. |
| `PLANTON_MCP_LOG_LEVEL` | `info` | Minimum log level: `debug`, `info`, `warn`, `error`. |

**TLS:** Connections to endpoints on port `443` automatically use TLS with the
system root CA pool. All other ports use plaintext. There is no separate TLS
configuration flag.

---

## Tools and Resources

238 tools span the full Planton product surface across 15 domain groups:

**Cloud Resource Lifecycle** -- apply, get, delete, list, destroy, purge, lock management, slug check, rename, env var extraction (12 tools)

**Stack Jobs** -- observe provisioning outcomes, retry failures, cancel or approve jobs, get essentials, IaC resources, error recommendation (12 tools)

**InfraChart Templates** -- browse and preview reusable infrastructure chart templates (3 tools)

**InfraProject Lifecycle** -- create and manage infrastructure projects sourced from charts or Git repos (6 tools)

**InfraPipeline Monitoring and Control** -- track deployment pipelines, trigger runs, resolve manual gates (9 tools)

**Dependency Graph** -- explore resource topology, trace dependencies and dependents, analyze blast radius (7 tools)

**Config Manager** -- manage plaintext variables, encrypted secrets, secret backends, and variable groups with full version history (23 tools)

**Audit and Version History** -- paginated change history and unified diffs for any platform resource (3 tools)

**Organizations and Environments** -- discover and manage orgs, environments, promotion policies, and flow control (17 tools)

**Presets and Catalog** -- find pre-configured templates, deployment components, and IaC modules (9 tools)

**Service Hub** -- manage services, CI/CD pipelines, DNS domains, Tekton pipelines and tasks, variable groups, and secret groups (46 tools)

**Connect** -- manage cloud provider connections, runners, GitHub integrations, default providers, and provider auth across 22 connection kinds (34 tools)

**IAM** -- identity management, API keys, teams, service accounts, roles, and authorization policies (28 tools)

**CloudOps** -- live queries against AWS, GCP, Azure, and Kubernetes for real-time cloud resource inspection (18 tools)

**Search and Discovery** -- text search, kind search, context hierarchy, quick actions (11 tools)

Plus **9 MCP prompts** for guided workflows (debug deployments, assess change impact,
explore infrastructure, provision resources, manage access, onboard teammates,
diagnose services, set up environments, investigate security incidents).

And **7 MCP resources** for schema discovery:

| URI | What It Returns |
|-----|-----------------|
| `cloud-resource-kinds://catalog` | All 362 supported kinds grouped by 17 cloud providers |
| `cloud-resource-schema://{kind}` | Full JSON schema for a specific kind |
| `connection-types://catalog` | All 22 connection types with descriptions |
| `api-resource-kinds://catalog` | All API resource kinds in the platform |
| `talent-pool://catalog` | Available talent profiles with expertise and rates |
| `environment://summary/{org}` | Org environments with connection and resource counts |
| `get_my_capabilities` | What the authenticated identity can do |

---

## HTTP Mode

For shared or remote deployments, set the transport to `http`. This runs the
MCP Streamable HTTP transport -- not a REST API.

```bash
PLANTON_MCP_TRANSPORT=http \
PLANTON_API_KEY=your-api-key \
  planton-mcp-server
```

Or with Docker:

```bash
docker run --rm \
  -e PLANTON_MCP_TRANSPORT=http \
  -e PLANTON_MCP_HTTP_AUTH_ENABLED=true \
  -p 8080:8080 \
  ghcr.io/plantonhq/planton-mcp-server
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
- Fine-grained authorization is enforced by the Planton backend via OpenFGA

---

## Supported Cloud Providers

AWS, GCP, Azure, Kubernetes, AliCloud, DigitalOcean, Civo, Cloudflare,
Confluent, Auth0, OpenFGA, Snowflake, MongoDB Atlas, Hetzner Cloud, OCI,
OpenStack, Scaleway -- 362 resource kinds total.

---

## License

Apache License 2.0. See [LICENSE](LICENSE).
