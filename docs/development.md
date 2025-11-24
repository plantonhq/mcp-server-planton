# Development Guide

Guide for contributing to and developing the Planton Cloud MCP Server.

## Development Setup

### Prerequisites

- Python 3.11+
- Poetry 1.5+
- Git
- Access to Planton Cloud APIs (local or remote)

### Initial Setup

1. Fork and clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/mcp-server-planton.git
cd mcp-server-planton
```

2. Install dependencies:

```bash
poetry install
```

3. Activate virtual environment:

```bash
poetry shell
```

4. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

## Development Workflow

### Running the Server

```bash
# Set environment variables
export USER_JWT_TOKEN="your-jwt-token"
export PLANTON_APIS_GRPC_ENDPOINT="localhost:8080"

# Run server
python src/mcp_server_planton/server.py
```

### Code Quality Tools

#### Linting with Ruff

```bash
# Check for issues
poetry run ruff check src/

# Auto-fix issues
poetry run ruff check --fix src/

# Check specific file
poetry run ruff check src/mcp_server_planton/server.py
```

#### Type Checking with MyPy

```bash
# Run type checker
poetry run mypy src/

# Check specific file
poetry run mypy src/mcp_server_planton/server.py
```

#### Running All Checks

```bash
# Run linting and type checking
poetry run ruff check src/ && poetry run mypy src/
```

### Testing

Currently, the project doesn't have automated tests. This is a great area for contributions!

**Planned test structure:**

```
tests/
├── unit/
│   ├── test_config.py
│   ├── test_auth/
│   │   └── test_interceptor.py
│   ├── test_grpc_clients/
│   │   └── test_environment_client.py
│   └── test_tools/
│       └── test_environment_tools.py
└── integration/
    └── test_server.py
```

## Project Structure

```
mcp-server-planton/
├── src/
│   └── mcp_server_planton/         # Main package
│       ├── __init__.py             # Package initialization
│       ├── server.py               # MCP server entry point
│       ├── config.py               # Configuration management
│       ├── auth/                   # Authentication modules
│       │   ├── __init__.py
│       │   └── user_token_interceptor.py  # gRPC auth interceptor
│       ├── grpc_clients/           # gRPC client implementations
│       │   ├── __init__.py
│       │   └── environment_client.py      # Environment API client
│       └── tools/                  # MCP tool implementations
│           ├── __init__.py
│           └── environment_tools.py       # Environment query tools
├── docs/                           # Documentation
│   ├── installation.md
│   ├── configuration.md
│   └── development.md
├── .github/
│   └── workflows/                  # CI/CD pipelines
│       ├── ci.yml
│       └── publish.yml
├── pyproject.toml                  # Poetry configuration
├── poetry.toml                     # Poetry settings
├── README.md                       # Main documentation
├── LICENSE                         # Apache-2.0 license
├── CONTRIBUTING.md                 # Contribution guidelines
└── .gitignore                      # Git ignore rules
```

## Adding New Features

### Adding a New MCP Tool

1. **Create or update gRPC client** (if needed):

```python
# src/mcp_server_planton/grpc_clients/organization_client.py
class OrganizationClient:
    def __init__(self, grpc_endpoint: str, user_token: str):
        # Initialize gRPC client
        pass
    
    async def list_organizations(self):
        # Implement API call
        pass
```

2. **Implement the tool**:

```python
# src/mcp_server_planton/tools/organization_tools.py
from mcp.types import Tool, TextContent

def create_organization_tool() -> Tool:
    return Tool(
        name="list_organizations",
        description="List all organizations the user has access to",
        inputSchema={"type": "object", "properties": {}}
    )

async def handle_list_organizations(arguments, config):
    # Implement tool handler
    pass
```

3. **Register in server**:

```python
# src/mcp_server_planton/server.py
from mcp_server_planton.tools.organization_tools import (
    create_organization_tool,
    handle_list_organizations,
)

@mcp_server.list_tools()
async def list_tools():
    return [
        create_environment_tool(),
        create_organization_tool(),  # Add new tool
    ]

@mcp_server.call_tool()
async def call_tool(name, arguments):
    if name == "list_organizations":
        return await handle_list_organizations(arguments, server_config)
    # ... other tools
```

4. **Update documentation**:
   - Add tool description to README.md
   - Document input/output schema
   - Provide usage examples

### Code Style Guidelines

- **Type hints**: Use type hints for all function signatures
- **Docstrings**: Document all public functions and classes
- **Error handling**: Handle gRPC errors gracefully
- **Logging**: Use appropriate log levels (INFO, ERROR, DEBUG)
- **Naming**: Use descriptive variable and function names

**Example:**

```python
async def fetch_resource_by_id(
    resource_id: str,
    user_token: str
) -> Optional[Resource]:
    """
    Fetch a resource by its ID.
    
    Args:
        resource_id: Unique identifier of the resource
        user_token: User's JWT token
        
    Returns:
        Resource object if found, None otherwise
        
    Raises:
        grpc.RpcError: If the API call fails
    """
    logger.info(f"Fetching resource: {resource_id}")
    try:
        # Implementation
        pass
    except grpc.RpcError as e:
        logger.error(f"Failed to fetch resource: {e}")
        raise
```

## Debugging

### Enable Debug Logging

```python
import logging
logging.basicConfig(level=logging.DEBUG)
```

### Debugging gRPC Calls

Enable gRPC debug logging:

```bash
export GRPC_VERBOSITY=DEBUG
export GRPC_TRACE=all
```

### Using Python Debugger

```python
import pdb; pdb.set_trace()  # Breakpoint
```

Or use IDE debuggers in VS Code, PyCharm, etc.

## Releasing

### Version Bumping

Update version in `pyproject.toml`:

```toml
[tool.poetry]
version = "0.2.0"
```

And in `src/mcp_server_planton/__init__.py`:

```python
__version__ = "0.2.0"
```

### Creating a Release

1. Create a git tag:

```bash
git tag v0.2.0
git push origin v0.2.0
```

2. GitHub Actions will automatically:
   - Build the package
   - Run tests (when available)
   - Publish to PyPI
   - Create GitHub release

### Manual Publishing (if needed)

```bash
# Build package
poetry build

# Publish to PyPI
poetry publish
```

## Continuous Integration

The project uses GitHub Actions for CI/CD:

- **ci.yml**: Runs on every push/PR
  - Linting with ruff
  - Type checking with mypy
  - Tests (when available)

- **publish.yml**: Runs on tag push
  - Builds package
  - Publishes to PyPI
  - Creates GitHub release

## Common Issues

### Import Errors

If you get import errors for `blintora_apis_protocolbuffers_python`:

```bash
poetry lock --no-update
poetry install
```

### gRPC Connection Issues

Test gRPC connection:

```python
import grpc

channel = grpc.insecure_channel('localhost:8080')
grpc.channel_ready_future(channel).result(timeout=10)
```

## Resources

- [MCP Protocol Documentation](https://modelcontextprotocol.io)
- [gRPC Python Guide](https://grpc.io/docs/languages/python/)
- [Poetry Documentation](https://python-poetry.org/docs/)
- [Planton Cloud Documentation](https://docs.planton.cloud)

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/plantoncloud-inc/mcp-server-planton/issues)
- **Discussions**: [GitHub Discussions](https://github.com/plantoncloud-inc/mcp-server-planton/discussions)
- **Contributing**: See [CONTRIBUTING.md](../CONTRIBUTING.md)

