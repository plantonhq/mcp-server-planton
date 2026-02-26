# Planton Cloud MCP Server

[![CI](https://github.com/plantoncloud/mcp-server-planton/actions/workflows/ci.yml/badge.svg)](https://github.com/plantoncloud/mcp-server-planton/actions/workflows/ci.yml)
[![CodeQL](https://github.com/plantoncloud/mcp-server-planton/actions/workflows/codeql.yml/badge.svg)](https://github.com/plantoncloud/mcp-server-planton/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/plantoncloud/mcp-server-planton)](https://goreportcard.com/report/github.com/plantoncloud/mcp-server-planton)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-ghcr.io-blue?logo=docker)](https://github.com/plantoncloud/mcp-server-planton/pkgs/container/mcp-server-planton)

MCP (Model Context Protocol) server for [Planton Cloud](https://planton.cloud) that enables AI agents to create, read, and delete cloud resources across 17 cloud providers using the [OpenMCF](https://github.com/plantonhq/openmcf) specification.

## Overview

The server exposes three MCP tools and two MCP resources that give any MCP-capable AI client (Cursor, Claude Desktop, Windsurf, LangGraph, etc.) full CRUD access to Planton-managed cloud infrastructure. All operations go through the Planton backend and respect per-user API key permissions.

### Tools

| Tool | Description |
|------|-------------|
| `apply_cloud_resource` | Create or update a cloud resource (idempotent) |
| `get_cloud_resource` | Retrieve a cloud resource by ID or by kind+org+env+slug |
| `delete_cloud_resource` | Delete a cloud resource by ID or by kind+org+env+slug |

### MCP Resources

| Resource | Description |
|----------|-------------|
| `cloud-resource-kinds://catalog` | Static catalog of all 362 supported kinds grouped by 17 cloud providers |
| `cloud-resource-schema://{kind}` | Per-kind JSON schema with field types, descriptions, validation rules, and defaults |

### Agent Workflow

Agents follow a 3-step discovery pattern:

1. **Discover** -- Read `cloud-resource-kinds://catalog` to find available kinds and their `api_version` values
2. **Learn** -- Read `cloud-resource-schema://{kind}` to get the full spec definition for a specific kind
3. **Act** -- Call `apply_cloud_resource` with the assembled `cloud_object`

## Quick Start

### Get Your API Key

1. Log in to [Planton Cloud Console](https://console.planton.cloud)
2. Click your profile icon > **API Keys**
3. Click **Create Key** and copy the generated key

### Cursor (STDIO)

Add to `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "planton-cloud": {
      "command": "mcp-server-planton",
      "env": {
        "PLANTON_API_KEY": "YOUR_PLANTON_API_KEY",
        "PLANTON_CLOUD_ENVIRONMENT": "live"
      }
    }
  }
}
```

### Cursor (Remote HTTP)

```json
{
  "mcpServers": {
    "planton-cloud": {
      "type": "http",
      "url": "https://mcp.planton.ai/",
      "headers": {
        "Authorization": "Bearer YOUR_PLANTON_API_KEY"
      }
    }
  }
}
```

### Claude Desktop

Add to your Claude Desktop MCP settings:

```json
{
  "mcpServers": {
    "planton-cloud": {
      "command": "mcp-server-planton",
      "env": {
        "PLANTON_API_KEY": "YOUR_PLANTON_API_KEY",
        "PLANTON_CLOUD_ENVIRONMENT": "live"
      }
    }
  }
}
```

### LangGraph

Add to your `langgraph.json`:

```json
{
  "mcp_servers": {
    "planton-cloud": {
      "command": "mcp-server-planton",
      "env": {
        "PLANTON_API_KEY": "${PLANTON_API_KEY}",
        "PLANTON_CLOUD_ENVIRONMENT": "live"
      }
    }
  }
}
```

### Docker (HTTP mode)

```bash
docker run -p 8080:8080 \
  -e PLANTON_MCP_TRANSPORT="http" \
  -e PLANTON_MCP_HTTP_AUTH_ENABLED="true" \
  ghcr.io/plantoncloud/mcp-server-planton:latest
```

Each user provides their own API key in the `Authorization: Bearer` header. The server does not store keys.

## Installation

### Pre-built Binary

```bash
# macOS (ARM64)
curl -L https://github.com/plantoncloud/mcp-server-planton/releases/latest/download/mcp-server-planton_Darwin_arm64.tar.gz | tar xz
sudo mv mcp-server-planton /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/plantoncloud/mcp-server-planton/releases/latest/download/mcp-server-planton_Darwin_x86_64.tar.gz | tar xz
sudo mv mcp-server-planton /usr/local/bin/

# Linux (AMD64)
curl -L https://github.com/plantoncloud/mcp-server-planton/releases/latest/download/mcp-server-planton_Linux_x86_64.tar.gz | tar xz
sudo mv mcp-server-planton /usr/local/bin/
```

### From Source

```bash
go install github.com/plantoncloud/mcp-server-planton/cmd/mcp-server-planton@latest
```

## Configuration

All settings are read from environment variables with a `PLANTON_` prefix.

| Variable | Default | Description |
|----------|---------|-------------|
| `PLANTON_API_KEY` | *(required for stdio/both)* | API key for Planton Cloud |
| `PLANTON_CLOUD_ENVIRONMENT` | `live` | Target environment: `live`, `test`, or `local` |
| `PLANTON_APIS_GRPC_ENDPOINT` | *(from environment)* | Explicit gRPC endpoint override |
| `PLANTON_MCP_TRANSPORT` | `stdio` | Transport mode: `stdio`, `http`, or `both` |
| `PLANTON_MCP_HTTP_PORT` | `8080` | HTTP listen port |
| `PLANTON_MCP_HTTP_AUTH_ENABLED` | `true` | Require Bearer token for HTTP requests |
| `PLANTON_MCP_LOG_FORMAT` | `text` | Log encoding: `text` or `json` |
| `PLANTON_MCP_LOG_LEVEL` | `info` | Minimum log level: `debug`, `info`, `warn`, `error` |

See [docs/configuration.md](docs/configuration.md) for details.

## Security

- **STDIO mode**: API key loaded once from `PLANTON_API_KEY` at startup
- **HTTP mode**: Each request carries its own key via `Authorization: Bearer` header -- true multi-tenant support
- Keys are held in memory only during request execution and never persisted
- All API calls are validated and logged with the caller's identity
- Fine-grained authorization is enforced by the Planton backend

## Architecture

```
cmd/mcp-server-planton/        CLI entry point (stdio | http | both)
pkg/mcpserver/                  Public embedding API (Config, Run)
internal/
  auth/                         Context-based API key propagation
  config/                       Env-var configuration with validation
  grpc/                         Authenticated gRPC client factory
  server/                       MCP server init, STDIO + HTTP transports
  domains/
    cloudresource/              Tool handlers, resource templates, schema lookup
  parse/                        Shared helpers for generated parsers
gen/cloudresource/              Generated typed input structs (362 providers, 17 clouds)
schemas/                        Embedded JSON schemas (go:embed)
tools/codegen/
  proto2schema/                 Stage 1: OpenMCF .proto -> JSON schema
  generator/                    Stage 2: JSON schema -> Go input types
```

## Development

### Prerequisites

- Go 1.25+
- Access to Planton Cloud APIs (local or remote)
- OpenMCF and Planton API repos (for codegen only)

### Build and Test

```bash
make build          # Build binary to bin/
make test           # Run all tests
make lint           # Run golangci-lint
make fmt            # Format Go code
```

### Codegen Pipeline

The two-stage codegen pipeline generates typed Go input structs from OpenMCF provider proto definitions:

```bash
make codegen-schemas    # Stage 1: proto -> JSON schemas (requires openmcf repo)
make codegen-types      # Stage 2: JSON schemas -> Go input types
make codegen            # Full pipeline (Stage 1 + Stage 2)
```

See [docs/development.md](docs/development.md) for codegen details.

## Supported Cloud Providers

AWS, GCP, Azure, Kubernetes, AliCloud, DigitalOcean, Civo, Cloudflare, Confluent, Auth0, OpenFGA, Snowflake, MongoDB Atlas, Hetzner Cloud, OCI, OpenStack, Scaleway -- 362 resource kinds total.

## License

Apache-2.0 -- see [LICENSE](LICENSE) for details.
