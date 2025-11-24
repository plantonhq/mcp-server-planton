# Next Steps for Publishing

The Planton Cloud MCP Server has been successfully extracted from the monorepo into its own public repository. All code, documentation, and CI/CD workflows are in place.

## What's Been Completed ✅

1. **Repository Setup**
   - ✅ Created public GitHub repository: [plantoncloud-inc/mcp-server-planton](https://github.com/plantoncloud-inc/mcp-server-planton)
   - ✅ Set up proper Python package structure with `src/mcp_server_planton/`
   - ✅ Migrated all source code with updated imports

2. **Dependencies**
   - ✅ Updated to use `blintora-apis-protocolbuffers-python` from buf.build registry
   - ✅ All imports updated from monorepo stubs to public stubs
   - ✅ Poetry configuration complete with PyPI metadata

3. **Documentation**
   - ✅ External-facing README (removed monorepo references)
   - ✅ Installation guide (`docs/installation.md`)
   - ✅ Configuration guide (`docs/configuration.md`)
   - ✅ Development guide (`docs/development.md`)
   - ✅ Contributing guidelines (`CONTRIBUTING.md`)

4. **Quality & CI/CD**
   - ✅ GitHub Actions workflows for CI and publishing
   - ✅ Linting (ruff) and type checking (mypy) passing
   - ✅ LICENSE (Apache-2.0)
   - ✅ .gitignore configured

5. **Code Pushed**
   - ✅ Initial commit made and pushed to GitHub
   - ✅ All 21 files committed successfully

6. **Monorepo Cleanup**
   - ✅ Removed `backend/services/planton-cloud-mcp-server/`
   - ✅ Added migration note to changelog

## Remaining: Publishing to PyPI 🚀

To publish the initial `v0.1.0` release to PyPI, follow these steps:

### 1. Configure PyPI Token

First, generate a PyPI API token:

1. Go to [PyPI Account Settings](https://pypi.org/manage/account/token/)
2. Create a new API token
3. Set scope to "Entire account" or specific to `mcp-server-planton` (after first upload)
4. Copy the token (starts with `pypi-`)

Then, add it to GitHub repository secrets:

1. Go to repository Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Name: `PYPI_TOKEN`
4. Value: Your PyPI token
5. Save

### 2. Create and Push Release Tag

Once the token is configured, create the release:

```bash
cd /Users/suresh/scm/github.com/plantoncloud-inc/mcp-server-planton

# Create annotated tag
git tag -a v0.1.0 -m "Release v0.1.0 - Initial public release

- User JWT authentication for Planton Cloud APIs
- Environment query tool (list_environments_for_org)
- Comprehensive documentation and examples
- CI/CD with GitHub Actions
- Apache-2.0 license"

# Push tag to trigger publish workflow
git push origin v0.1.0
```

### 3. Monitor Release

The GitHub Actions workflow will:

1. Run linting and type checking
2. Build the package with Poetry
3. Publish to PyPI
4. Create a GitHub release

Monitor the workflow at:
- Actions tab in GitHub repository
- Check for any errors in the publish workflow

### 4. Verify Publication

Once published, verify:

```bash
# Install from PyPI
pip install mcp-server-planton

# Verify version
python -c "from mcp_server_planton import __version__; print(__version__)"

# Should print: 0.1.0
```

### 5. Update Graph-Fleet (When Ready to Integrate)

When you're ready to integrate the MCP server with graph-fleet, update `graph-fleet/pyproject.toml`:

```toml
[tool.poetry.dependencies]
mcp-server-planton = "^0.1.0"
```

Or for pre-release testing, use Git dependency:

```toml
[tool.poetry.dependencies]
mcp-server-planton = { git = "https://github.com/plantoncloud-inc/mcp-server-planton.git", branch = "master" }
```

## Repository Links

- **GitHub**: https://github.com/plantoncloud-inc/mcp-server-planton
- **PyPI** (after publishing): https://pypi.org/project/mcp-server-planton/
- **Documentation**: See README.md and docs/ folder

## Future Enhancements

Consider these improvements for future releases:

1. **Testing**
   - Add pytest framework
   - Unit tests for all modules
   - Integration tests with mock gRPC server

2. **Additional Tools**
   - Organization query tools
   - Project query tools
   - Cloud resource query tools

3. **Performance**
   - Caching layer for frequently accessed data
   - Connection pooling for gRPC clients

4. **Observability**
   - Metrics for tool invocation latency
   - Structured logging with correlation IDs
   - OpenTelemetry integration

5. **Security**
   - Token refresh mechanism
   - Rate limiting
   - Audit logging

## Questions or Issues?

- **GitHub Issues**: https://github.com/plantoncloud-inc/mcp-server-planton/issues
- **GitHub Discussions**: https://github.com/plantoncloud-inc/mcp-server-planton/discussions

---

**Status**: Ready for publishing to PyPI ✨  
**Next Action**: Configure PYPI_TOKEN secret and push v0.1.0 tag

