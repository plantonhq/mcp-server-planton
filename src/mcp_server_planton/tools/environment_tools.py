"""MCP tools for Environment resource queries."""

import json
import logging
from typing import Dict, Any, List

import grpc
from mcp.types import Tool, TextContent

from mcp_server_planton.grpc_clients.environment_client import EnvironmentClient
from mcp_server_planton.config import MCPServerConfig

logger = logging.getLogger(__name__)


def create_environment_tool() -> Tool:
    """
    Create MCP tool definition for listing environments.
    
    Returns:
        MCP Tool definition with name, description, and input schema
    """
    return Tool(
        name="list_environments_for_org",
        description=(
            "List all environments available in an organization. "
            "Returns environment details including id, slug, name, and description. "
            "Only returns environments the user has permission to view."
        ),
        inputSchema={
            "type": "object",
            "properties": {
                "org_id": {
                    "type": "string",
                    "description": "Organization ID to query environments for"
                }
            },
            "required": ["org_id"]
        }
    )


async def handle_list_environments_for_org(
    arguments: Dict[str, Any],
    config: MCPServerConfig
) -> List[TextContent]:
    """
    Handle MCP tool invocation for listing environments.
    
    This function:
    1. Validates the org_id argument
    2. Creates EnvironmentClient with user JWT
    3. Queries Planton Cloud APIs for environments
    4. Converts protobuf responses to JSON-serializable dicts
    5. Returns formatted response or error message
    
    Args:
        arguments: Tool arguments containing org_id
        config: MCP server configuration with user JWT and gRPC endpoint
        
    Returns:
        List containing single TextContent with JSON response
    """
    org_id = arguments.get("org_id")
    
    if not org_id:
        error_response = {
            "error": "INVALID_ARGUMENT",
            "message": "org_id is required"
        }
        return [TextContent(
            type="text",
            text=json.dumps(error_response, indent=2)
        )]
    
    logger.info(f"Tool invoked: list_environments_for_org, org_id={org_id}")
    
    # Create gRPC client with user JWT
    client = EnvironmentClient(
        grpc_endpoint=config.planton_apis_grpc_endpoint,
        user_token=config.user_jwt_token
    )
    
    try:
        # Query environments
        environments = await client.find_by_org(org_id)
        
        # Convert protobuf objects to JSON-serializable dicts
        environment_list = [
            {
                "id": env.metadata.id,
                "slug": env.metadata.slug,
                "name": env.metadata.name,
                "description": env.spec.description if env.spec.description else "",
            }
            for env in environments
        ]
        
        logger.info(
            f"Tool completed: list_environments_for_org, "
            f"returned {len(environment_list)} environments"
        )
        
        # Return formatted JSON response
        return [TextContent(
            type="text",
            text=json.dumps(environment_list, indent=2)
        )]
        
    except grpc.RpcError as e:
        # Handle gRPC errors with user-friendly messages
        error_code = e.code()
        error_details = e.details()
        
        logger.error(
            f"Tool error: list_environments_for_org, "
            f"org_id={org_id}, code={error_code}, details={error_details}"
        )
        
        # Map gRPC error codes to user-friendly messages
        if error_code == grpc.StatusCode.UNAUTHENTICATED:
            error_response = {
                "error": "UNAUTHENTICATED",
                "message": (
                    "Authentication failed. Your session may have expired. "
                    "Please refresh and try again."
                ),
                "org_id": org_id
            }
        elif error_code == grpc.StatusCode.PERMISSION_DENIED:
            error_response = {
                "error": "PERMISSION_DENIED",
                "message": (
                    f"You don't have permission to view environments "
                    f"for organization '{org_id}'. Please contact your "
                    f"organization administrator."
                ),
                "org_id": org_id
            }
        elif error_code == grpc.StatusCode.UNAVAILABLE:
            error_response = {
                "error": "UNAVAILABLE",
                "message": (
                    "Planton Cloud APIs are currently unavailable. "
                    "Please try again in a moment."
                ),
                "org_id": org_id
            }
        elif error_code == grpc.StatusCode.NOT_FOUND:
            error_response = {
                "error": "NOT_FOUND",
                "message": f"Organization '{org_id}' not found.",
                "org_id": org_id
            }
        else:
            error_response = {
                "error": str(error_code.name),
                "message": error_details or "An unexpected error occurred.",
                "org_id": org_id
            }
        
        return [TextContent(
            type="text",
            text=json.dumps(error_response, indent=2)
        )]
        
    except Exception as e:
        # Handle unexpected errors
        logger.error(
            f"Unexpected error in list_environments_for_org: {e}",
            exc_info=True
        )
        
        error_response = {
            "error": "INTERNAL_ERROR",
            "message": f"An unexpected error occurred: {str(e)}",
            "org_id": org_id
        }
        
        return [TextContent(
            type="text",
            text=json.dumps(error_response, indent=2)
        )]
        
    finally:
        # Clean up gRPC client
        await client.close()

