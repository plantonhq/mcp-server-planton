# Contributing to Planton Cloud MCP Server

Thank you for your interest in contributing to the Planton Cloud MCP Server! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Python 3.11 or higher
- Poetry package manager
- Git

### Setting Up Development Environment

1. Clone the repository:
```bash
git clone https://github.com/plantoncloud-inc/mcp-server-planton.git
cd mcp-server-planton
```

2. Install dependencies with Poetry:
```bash
poetry install
```

3. Activate the virtual environment:
```bash
poetry shell
```

## Development Workflow

### Running the Server Locally

Set required environment variables:
```bash
export USER_JWT_TOKEN="your-jwt-token"
export PLANTON_APIS_GRPC_ENDPOINT="localhost:8080"
```

Run the server:
```bash
python src/mcp_server_planton/server.py
```

### Code Quality

We use `ruff` for linting and `mypy` for type checking.

Run linter:
```bash
poetry run ruff check src/
```

Run type checker:
```bash
poetry run mypy src/
```

Fix auto-fixable linting issues:
```bash
poetry run ruff check --fix src/
```

### Code Style Guidelines

- Follow PEP 8 style guide
- Use type hints for all functions and methods
- Write docstrings for all public modules, classes, and functions
- Keep functions focused and single-purpose
- Use descriptive variable names

## Adding New Tools

To add a new MCP tool:

1. Create or update a client in `src/mcp_server_planton/grpc_clients/`
2. Add tool implementation in `src/mcp_server_planton/tools/`
3. Register the tool in `src/mcp_server_planton/server.py`:
   - Add tool to `list_tools()` function
   - Add handler to `call_tool()` function
4. Update documentation
5. Add tests (when test infrastructure is available)

## Submitting Changes

### Pull Request Process

1. Fork the repository
2. Create a feature branch:
```bash
git checkout -b feature/your-feature-name
```

3. Make your changes
4. Run linting and type checking
5. Commit your changes with clear, descriptive messages:
```bash
git commit -m "feat: add new tool for querying organizations"
```

6. Push to your fork:
```bash
git push origin feature/your-feature-name
```

7. Open a Pull Request against the `main` branch

### Commit Message Format

We follow conventional commit format:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

Examples:
```
feat: add organization query tool
fix: handle timeout errors in environment client
docs: update installation instructions
refactor: simplify error handling in tools
```

### Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Update documentation for any user-facing changes
- Ensure all checks pass (linting, type checking)
- Provide clear description of changes
- Reference related issues if applicable

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

- Clear description of the issue
- Steps to reproduce
- Expected behavior vs actual behavior
- Environment details (Python version, OS, etc.)
- Relevant logs or error messages

### Feature Requests

When requesting features, please include:

- Clear description of the feature
- Use case and motivation
- Example of how it would be used
- Any relevant context or alternatives considered

## Questions and Support

- **GitHub Issues**: For bug reports and feature requests
- **GitHub Discussions**: For questions and general discussion

## Code of Conduct

- Be respectful and inclusive
- Assume good intentions
- Give and accept constructive feedback gracefully
- Focus on what's best for the community

## License

By contributing, you agree that your contributions will be licensed under the Apache-2.0 License.

