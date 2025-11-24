"""gRPC client for Environment queries."""

import logging
from typing import List

import grpc
from blintora_apis_protocolbuffers_python.ai.planton.resourcemanager.environment.v1.query_pb2_grpc import (
    EnvironmentQueryControllerStub,
)
from blintora_apis_protocolbuffers_python.ai.planton.resourcemanager.environment.v1.api_pb2 import Environment
from blintora_apis_protocolbuffers_python.ai.planton.resourcemanager.organization.v1.io_pb2 import OrganizationId

from mcp_server_planton.auth.user_token_interceptor import UserTokenAuthClientInterceptor

logger = logging.getLogger(__name__)


class EnvironmentClient:
    """
    gRPC client for querying Planton Cloud Environment resources.
    
    This client uses the user's JWT token (not machine account) to make
    authenticated gRPC calls to Planton Cloud APIs. The APIs validate the
    JWT and enforce Fine-Grained Authorization (FGA) checks based on the
    user's actual permissions.
    """
    
    def __init__(self, grpc_endpoint: str, user_token: str):
        """
        Initialize Environment gRPC client.
        
        Args:
            grpc_endpoint: Planton Cloud APIs endpoint (e.g., "localhost:8080")
            user_token: User's JWT token from environment variable
        """
        self.grpc_endpoint = grpc_endpoint
        self.user_token = user_token
        
        # Create gRPC channel
        self.channel = grpc.aio.insecure_channel(grpc_endpoint)
        
        # Create interceptor that attaches user JWT to all requests
        interceptor = UserTokenAuthClientInterceptor(user_token)
        
        # Intercept the channel
        self.intercepted_channel = grpc.aio.intercept_channel(
            self.channel, interceptor
        )
        
        # Create stub with intercepted channel
        self.stub = EnvironmentQueryControllerStub(self.intercepted_channel)
        
        logger.info(f"EnvironmentClient initialized for endpoint: {grpc_endpoint}")
    
    async def find_by_org(self, org_id: str) -> List[Environment]:
        """
        Query all environments for an organization.
        
        This method makes an authenticated gRPC call to Planton Cloud APIs
        using the user's JWT token. The API validates the JWT and checks
        FGA permissions to ensure the user has access to view environments
        in the specified organization.
        
        Args:
            org_id: Organization ID to query environments for
            
        Returns:
            List of Environment protobuf objects
            
        Raises:
            grpc.RpcError: If the gRPC call fails (authentication, authorization, etc.)
        """
        logger.info(f"Querying environments for org: {org_id}")
        
        try:
            # Create request protobuf
            request = OrganizationId(value=org_id)
            
            # Make gRPC call (interceptor attaches JWT automatically)
            response = await self.stub.findByOrg(request)
            
            # Extract environments from response
            environments = list(response.entries)
            
            logger.info(f"Found {len(environments)} environments for org: {org_id}")
            return environments
            
        except grpc.RpcError as e:
            logger.error(
                f"gRPC error querying environments for org {org_id}: "
                f"code={e.code()}, details={e.details()}"
            )
            raise
    
    async def close(self):
        """Close the gRPC channel."""
        if self.channel:
            await self.channel.close()
            logger.info("EnvironmentClient channel closed")

