# Configuration Guide

Comprehensive configuration options for the Planton Cloud MCP Server.

## Environment Variables

### Required Variables

#### USER_JWT_TOKEN

User's JWT token for authentication with Planton Cloud APIs.

```bash
export USER_JWT_TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**How to obtain:**
- Web Console: Developer Tools → Application → Local Storage
- CLI: `planton auth token`

**Important:** This token represents the user's identity and permissions. Keep it secure and never commit it to version control.

### Optional Variables

#### PLANTON_APIS_GRPC_ENDPOINT

gRPC endpoint for Planton Cloud APIs.

```bash
export PLANTON_APIS_GRPC_ENDPOINT="apis.planton.cloud:443"
```

**Default:** `localhost:8080`

**Common values:**
- Production: `apis.planton.cloud:443`
- Staging: `staging.apis.planton.cloud:443`
- Local development: `localhost:8080`

## Configuration Files

### .env File

Create a `.env` file in your project root:

```env
# Required
USER_JWT_TOKEN=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...

# Optional (defaults to localhost:8080)
PLANTON_APIS_GRPC_ENDPOINT=apis.planton.cloud:443
```

The server automatically loads `.env` files using `pydantic-settings`.

**Security note:** Add `.env` to your `.gitignore`:

```gitignore
.env
.env.local
.env.*.local
```

## Integration-Specific Configuration

### LangGraph Configuration

In `langgraph.json`:

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

**Note:** LangGraph supports environment variable expansion with `${VAR}` syntax.

### Claude Desktop Configuration

In `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "planton-cloud": {
      "command": "mcp-server-planton",
      "env": {
        "USER_JWT_TOKEN": "actual-token-here",
        "PLANTON_APIS_GRPC_ENDPOINT": "apis.planton.cloud:443"
      }
    }
  }
}
```

**Note:** Claude Desktop requires actual values, not environment variable references.

## Advanced Configuration

### Logging

The server uses Python's standard logging module. Configure logging level:

```python
import logging

logging.basicConfig(
    level=logging.DEBUG,  # Or INFO, WARNING, ERROR
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
```

Or via environment variable:

```bash
export LOG_LEVEL=DEBUG
```

### TLS/SSL Configuration

For production deployments with TLS:

```bash
# Use secure endpoint (port 443)
export PLANTON_APIS_GRPC_ENDPOINT="apis.planton.cloud:443"
```

The gRPC client automatically detects and uses TLS for port 443.

For custom TLS certificates:

```python
from mcp_server_planton.grpc_clients.environment_client import EnvironmentClient
import grpc

# Load custom certificates
with open('ca.pem', 'rb') as f:
    ca_cert = f.read()

credentials = grpc.ssl_channel_credentials(ca_cert)
# Configure client with credentials
```

### Timeout Configuration

Default gRPC timeout is 30 seconds. To customize:

```python
from mcp_server_planton.grpc_clients.environment_client import EnvironmentClient

client = EnvironmentClient(
    grpc_endpoint="apis.planton.cloud:443",
    user_token="your-token",
    timeout=60  # 60 seconds
)
```

## Security Best Practices

### JWT Token Management

1. **Never commit tokens** - Use environment variables or secret managers
2. **Rotate regularly** - JWT tokens should have expiration times
3. **Scope appropriately** - Use tokens with minimum required permissions
4. **Monitor usage** - Track token usage for audit purposes

### Environment-Specific Tokens

Use different tokens for different environments:

```bash
# Development
export USER_JWT_TOKEN_DEV="dev-token..."

# Staging
export USER_JWT_TOKEN_STAGING="staging-token..."

# Production
export USER_JWT_TOKEN_PROD="prod-token..."
```

### Secret Management

For production deployments, use secret managers:

**AWS Secrets Manager:**
```bash
export USER_JWT_TOKEN=$(aws secretsmanager get-secret-value \
  --secret-id planton/jwt-token \
  --query SecretString \
  --output text)
```

**HashiCorp Vault:**
```bash
export USER_JWT_TOKEN=$(vault kv get -field=token secret/planton/jwt)
```

**Kubernetes Secrets:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: planton-mcp-config
type: Opaque
stringData:
  USER_JWT_TOKEN: "your-jwt-token"
  PLANTON_APIS_GRPC_ENDPOINT: "apis.planton.cloud:443"
```

## Validation

Verify your configuration:

```python
from mcp_server_planton.config import MCPServerConfig

try:
    config = MCPServerConfig.load_from_env()
    print(f"Configuration valid!")
    print(f"Endpoint: {config.planton_apis_grpc_endpoint}")
except ValueError as e:
    print(f"Configuration error: {e}")
```

## Troubleshooting

### Missing Token Error

```
ValueError: USER_JWT_TOKEN environment variable required
```

**Solution:** Set the `USER_JWT_TOKEN` environment variable.

### Invalid Token Error

```
grpc.RpcError: code=UNAUTHENTICATED, details=Invalid authentication credentials
```

**Solutions:**
1. Verify token is not expired
2. Ensure token is complete (no truncation)
3. Re-authenticate and get a new token

### Connection Refused

```
grpc.RpcError: code=UNAVAILABLE, details=failed to connect to all addresses
```

**Solutions:**
1. Verify endpoint is correct
2. Check network connectivity
3. Ensure firewall allows gRPC traffic
4. Try pinging the endpoint

## Next Steps

- [Installation Guide](installation.md) - Installation instructions
- [Development Guide](development.md) - Development setup
- [README](../README.md) - Back to main documentation

