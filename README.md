# Planton Cloud MCP Server

MCP (Model Context Protocol) server for Planton Cloud that enables AI agents to query cloud resources using user permissions.

## Overview

The Planton Cloud MCP Server provides tools for LangGraph agents, Claude Desktop, and other MCP clients to interact with Planton Cloud resources. Unlike typical MCP servers that use API keys or machine accounts, this server uses **user JWT tokens**, ensuring that all resource queries respect Fine-Grained Authorization (FGA) based on the user's actual permissions.

### Key Features

- **User-scoped permissions** - Queries respect the user's actual permissions via JWT
- **Environment queries** - List and filter environments by organization
- **Extensible** - More resource types coming soon (organizations, projects, cloud resources)
- **MCP standard** - Works with any MCP client (LangGraph, Claude Desktop, Cursor, etc.)

## Installation

### From PyPI (recommended)

```bash
pip install mcp-server-planton
```

### From Source

```bash
git clone https://github.com/plantoncloud-inc/mcp-server-planton.git
cd mcp-server-planton
poetry install
```

## Quick Start

### Standalone Usage

Set required environment variables:

```bash
export USER_JWT_TOKEN="your-jwt-token"
export PLANTON_APIS_GRPC_ENDPOINT="apis.planton.cloud:443"
```

Run the server:

```bash
mcp-server-planton
```

Or with Python module:

```bash
python -m mcp_server_planton.server
```

### Integration with LangGraph

Add to your `langgraph.json`:

```json
{
  "mcp_servers": {
    "planton-cloud": {
      "command": "mcp-server-planton",
      "env": {
        "USER_JWT_TOKEN": "${USER_JWT_TOKEN}",
        "PLANTON_APIS_GRPC_ENDPOINT": "${PLANTON_APIS_GRPC_ENDPOINT}"
      }
    }
  }
}
```

### Integration with Claude Desktop

Add to your Claude Desktop MCP settings:

```json
{
  "mcpServers": {
    "planton-cloud": {
      "command": "mcp-server-planton",
      "env": {
        "USER_JWT_TOKEN": "your-jwt-token",
        "PLANTON_APIS_GRPC_ENDPOINT": "apis.planton.cloud:443"
      }
    }
  }
}
```

## Configuration

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `USER_JWT_TOKEN` | Yes | - | User's JWT token for authentication |
| `PLANTON_APIS_GRPC_ENDPOINT` | No | `localhost:8080` | Planton Cloud APIs gRPC endpoint |

### Getting a JWT Token

To obtain a JWT token:

1. Log in to Planton Cloud web console
2. Open browser developer tools (F12)
3. Go to Application/Storage → Local Storage
4. Find the authentication token
5. Copy the JWT value

For programmatic access, use the Planton Cloud CLI:

```bash
planton auth login
planton auth token
```

## Available Tools

### list_environments_for_org

List all environments available in an organization that the user has permission to view.

**Input:**
```json
{
  "org_id": "your-organization-id"
}
```

**Output:**
```json
[
  {
    "id": "env-123",
    "slug": "production",
    "name": "Production Environment",
    "description": "Production deployment environment"
  },
  {
    "id": "env-456",
    "slug": "staging",
    "name": "Staging Environment",
    "description": "Pre-production staging"
  }
]
```

**Error Handling:**

The tool returns user-friendly error messages for common issues:

- `UNAUTHENTICATED` - JWT token invalid or expired
- `PERMISSION_DENIED` - User lacks permission to view environments in the organization
- `NOT_FOUND` - Organization doesn't exist
- `UNAVAILABLE` - Planton Cloud APIs are temporarily unavailable

## Security Architecture

### User JWT Propagation

This MCP server follows a unique security pattern:

```
User → LangGraph/MCP Client → MCP Server (with JWT) → Planton Cloud APIs
```

**Key security properties:**

1. **No JWT persistence** - JWT is only held in memory during execution
2. **User permissions enforced** - APIs validate JWT and check FGA on every call
3. **Short-lived process** - MCP server exits when agent execution completes
4. **Audit trail** - All API calls are logged with user identity

This is different from machine account patterns where a service account has broad permissions. With user JWT propagation, the MCP server can only access what the user can access.

## Development

### Prerequisites

- Python 3.11+
- Poetry
- Access to Planton Cloud APIs (local or remote)

### Setup

```bash
# Clone repository
git clone https://github.com/plantoncloud-inc/mcp-server-planton.git
cd mcp-server-planton

# Install dependencies
poetry install

# Activate virtual environment
poetry shell
```

### Running Locally

```bash
export USER_JWT_TOKEN="your-jwt-token"
export PLANTON_APIS_GRPC_ENDPOINT="localhost:8080"
python src/mcp_server_planton/server.py
```

### Code Quality

```bash
# Run linter
poetry run ruff check src/

# Run type checker
poetry run mypy src/

# Auto-fix linting issues
poetry run ruff check --fix src/
```

### Project Structure

```
mcp-server-planton/
├── src/
│   └── mcp_server_planton/
│       ├── server.py              # MCP server entry point
│       ├── config.py              # Configuration management
│       ├── auth/
│       │   └── user_token_interceptor.py  # gRPC auth interceptor
│       ├── grpc_clients/
│       │   └── environment_client.py      # Environment API client
│       └── tools/
│           └── environment_tools.py       # Environment query tools
├── docs/                          # Documentation
├── pyproject.toml                 # Poetry dependencies
└── README.md                      # This file
```

## Roadmap

- [x] Environment query tools
- [ ] Organization query tools
- [ ] Project query tools
- [ ] Cloud resource query tools
- [ ] Resource mutation tools (create, update, delete)
- [ ] Caching and performance optimization
- [ ] Comprehensive test suite

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests and linting
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to your fork (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Support

- **Documentation**: [docs/](docs/)
- **Issues**: [GitHub Issues](https://github.com/plantoncloud-inc/mcp-server-planton/issues)
- **Discussions**: [GitHub Discussions](https://github.com/plantoncloud-inc/mcp-server-planton/discussions)

## License

Apache-2.0 - see [LICENSE](LICENSE) for details.

## Related Projects

- [Planton Cloud](https://planton.cloud) - Cloud infrastructure management platform
- [Project Planton](https://github.com/project-planton/project-planton) - Open-source deployment components
- [LangGraph](https://langchain-ai.github.io/langgraph/) - Framework for building stateful AI agents
- [Model Context Protocol](https://modelcontextprotocol.io) - Protocol for connecting AI models to data sources

---

Built with ❤️ by [Planton Cloud](https://planton.cloud)

