# Installation Guide

This guide covers different ways to install and run the Planton Cloud MCP Server.

## Prerequisites

- Python 3.11 or higher
- pip or Poetry package manager
- Access to Planton Cloud (account and JWT token)

## Installation Methods

### Method 1: Install from PyPI (Recommended)

The simplest way to install:

```bash
pip install mcp-server-planton
```

Verify installation:

```bash
mcp-server-planton --help
```

### Method 2: Install with Poetry

If you prefer Poetry:

```bash
poetry add mcp-server-planton
```

### Method 3: Install from Source

For development or latest features:

```bash
# Clone repository
git clone https://github.com/plantoncloud-inc/mcp-server-planton.git
cd mcp-server-planton

# Install with Poetry
poetry install

# Or install with pip
pip install -e .
```

## Configuration

### 1. Obtain JWT Token

Get your Planton Cloud JWT token:

**Option A: From Web Console**
1. Log in to Planton Cloud web console
2. Open browser developer tools (F12)
3. Navigate to Application → Local Storage
4. Find and copy the authentication token

**Option B: From CLI**
```bash
planton auth login
planton auth token
```

### 2. Set Environment Variables

Create a `.env` file or export variables:

```bash
export USER_JWT_TOKEN="your-jwt-token-here"
export PLANTON_APIS_GRPC_ENDPOINT="apis.planton.cloud:443"
```

For local development against a local Planton Cloud instance:

```bash
export USER_JWT_TOKEN="your-jwt-token-here"
export PLANTON_APIS_GRPC_ENDPOINT="localhost:8080"
```

### 3. Verify Installation

Test the server:

```bash
# With environment variables set
mcp-server-planton
```

You should see log output indicating the server started successfully.

## Integration Setup

### LangGraph Integration

Add to your `langgraph.json`:

```json
{
  "mcp_servers": {
    "planton-cloud": {
      "command": "mcp-server-planton",
      "env": {
        "USER_JWT_TOKEN": "${USER_JWT_TOKEN}",
        "PLANTON_APIS_GRPC_ENDPOINT": "apis.planton.cloud:443"
      }
    }
  }
}
```

### Claude Desktop Integration

Add to Claude Desktop MCP configuration:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

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

### Cursor Integration

Add to your Cursor MCP settings:

```json
{
  "mcp": {
    "servers": {
      "planton-cloud": {
        "command": "mcp-server-planton",
        "env": {
          "USER_JWT_TOKEN": "your-jwt-token",
          "PLANTON_APIS_GRPC_ENDPOINT": "apis.planton.cloud:443"
        }
      }
    }
  }
}
```

## Troubleshooting

### Import Errors

If you encounter import errors related to `blintora_apis_protocolbuffers_python`:

```bash
# Clear pip cache and reinstall
pip cache purge
pip uninstall mcp-server-planton
pip install --no-cache-dir mcp-server-planton
```

### Connection Issues

If the server can't connect to Planton Cloud APIs:

1. Verify the endpoint is correct:
   - Production: `apis.planton.cloud:443`
   - Local: `localhost:8080`

2. Check network connectivity:
```bash
ping apis.planton.cloud
```

3. Verify JWT token is valid:
```bash
# Check token expiration
echo $USER_JWT_TOKEN | cut -d'.' -f2 | base64 -d
```

### Permission Errors

If you get permission denied errors:

1. Verify your JWT token is current (not expired)
2. Check that your user has permissions in the organization
3. Contact your Planton Cloud administrator

### Version Conflicts

If you have dependency conflicts:

```bash
# Use a virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate
pip install mcp-server-planton
```

## Updating

### Update from PyPI

```bash
pip install --upgrade mcp-server-planton
```

### Update from Source

```bash
cd mcp-server-planton
git pull
poetry install
```

## Uninstallation

```bash
pip uninstall mcp-server-planton
```

## Next Steps

- [Configuration Guide](configuration.md) - Detailed configuration options
- [Development Guide](development.md) - Contributing and development setup
- [README](../README.md) - Back to main documentation

