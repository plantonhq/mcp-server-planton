"""Planton Cloud MCP Server entry point."""

import asyncio
import logging
from typing import Optional

from mcp.server import Server
from mcp.server.stdio import stdio_server

from mcp_server_planton.config import MCPServerConfig
from mcp_server_planton.tools.environment_tools import (
    create_environment_tool,
    handle_list_environments_for_org,
)

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)


# Create MCP server instance
mcp_server = Server("planton-cloud")

# Global config instance (initialized in main())
server_config: Optional[MCPServerConfig] = None


@mcp_server.list_tools()
async def list_tools():
    """
    List available MCP tools.
    
    Returns list of available tools for querying Planton Cloud resources.
    """
    logger.info("list_tools() called")
    return [
        create_environment_tool(),
    ]


@mcp_server.call_tool()
async def call_tool(name, arguments):
    """
    Handle MCP tool invocations.
    
    Args:
        name: Tool name
        arguments: Tool arguments
        
    Returns:
        Tool execution result
        
    Raises:
        ValueError: If tool not found
    """
    logger.info(f"call_tool() called: {name} with arguments: {arguments}")
    
    if name == "list_environments_for_org":
        return await handle_list_environments_for_org(arguments, server_config)
    
    raise ValueError(f"Tool not found: {name}")


async def main():
    """
    Main entry point for MCP server.
    
    Loads configuration from environment and starts MCP server
    using stdio transport (stdin/stdout communication).
    """
    global server_config
    
    try:
        # Load configuration (validates USER_JWT_TOKEN exists)
        server_config = MCPServerConfig.load_from_env()
        logger.info(
            f"MCP server starting with endpoint: {server_config.planton_apis_grpc_endpoint}"
        )
        logger.info("User JWT token loaded from environment")
        
        # Start MCP server with stdio transport
        async with stdio_server() as (read_stream, write_stream):
            logger.info("MCP server running on stdio")
            await mcp_server.run(
                read_stream,
                write_stream,
                mcp_server.create_initialization_options()
            )
    
    except ValueError as e:
        logger.error(f"Configuration error: {e}")
        raise
    except Exception as e:
        logger.error(f"MCP server error: {e}", exc_info=True)
        raise


if __name__ == "__main__":
    asyncio.run(main())

