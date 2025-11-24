"""gRPC client interceptor for attaching user JWT to requests."""

import grpc
from typing import Any, Callable, Awaitable, Tuple, Optional, Sequence


class _ClientCallDetails:
    """Custom ClientCallDetails to allow metadata modification."""
    
    def __init__(
        self,
        method: str,
        timeout: Optional[float],
        metadata: Optional[Sequence[Tuple[str, str]]],
        credentials: Optional[grpc.CallCredentials],
        wait_for_ready: Optional[bool],
    ):
        self.method = method
        self.timeout = timeout
        self.metadata = metadata
        self.credentials = credentials
        self.wait_for_ready = wait_for_ready


class UserTokenAuthClientInterceptor(
    grpc.aio.UnaryUnaryClientInterceptor,
    grpc.aio.UnaryStreamClientInterceptor,
    grpc.aio.StreamUnaryClientInterceptor,
    grpc.aio.StreamStreamClientInterceptor,
):
    """
    gRPC client interceptor that attaches user JWT to all requests.
    
    This interceptor passes through the user's JWT token (from environment)
    to Planton Cloud APIs, enabling Fine-Grained Authorization (FGA) checks
    using the user's actual permissions.
    
    Key Difference from agent-fleet-worker's AuthClientInterceptor:
    - agent-fleet-worker: Fetches machine account token from Auth0
    - MCP server: Uses user JWT directly (no token fetching)
    """
    
    def __init__(self, user_token: str):
        """
        Initialize interceptor with user JWT token.
        
        Args:
            user_token: User's JWT token (from environment variable)
        """
        self.user_token = user_token
    
    def _augment_call_details(
        self, 
        client_call_details: grpc.aio.ClientCallDetails
    ) -> _ClientCallDetails:
        """
        Add user JWT to call metadata as Authorization header.
        
        Args:
            client_call_details: Original call details
            
        Returns:
            Modified call details with Authorization header
        """
        # Get current metadata or create empty list
        metadata = list(client_call_details.metadata or [])
        
        # Add authorization header with user JWT
        metadata.append(("authorization", f"Bearer {self.user_token}"))
        
        # Create new call details with updated metadata
        return _ClientCallDetails(
            method=client_call_details.method,
            timeout=client_call_details.timeout,
            metadata=tuple(metadata),
            credentials=client_call_details.credentials,
            wait_for_ready=client_call_details.wait_for_ready,
        )
    
    async def intercept_unary_unary(
        self,
        continuation: Callable[[grpc.aio.ClientCallDetails, Any], Awaitable[Any]],
        client_call_details: grpc.aio.ClientCallDetails,
        request: Any,
    ) -> Any:
        """Intercept unary-unary calls."""
        new_details = self._augment_call_details(client_call_details)
        return await continuation(new_details, request)
    
    async def intercept_unary_stream(
        self,
        continuation: Callable[[grpc.aio.ClientCallDetails, Any], Awaitable[Any]],
        client_call_details: grpc.aio.ClientCallDetails,
        request: Any,
    ) -> Any:
        """Intercept unary-stream calls."""
        new_details = self._augment_call_details(client_call_details)
        return await continuation(new_details, request)
    
    async def intercept_stream_unary(
        self,
        continuation: Callable[[grpc.aio.ClientCallDetails, Any], Awaitable[Any]],
        client_call_details: grpc.aio.ClientCallDetails,
        request_iterator: Any,
    ) -> Any:
        """Intercept stream-unary calls."""
        new_details = self._augment_call_details(client_call_details)
        return await continuation(new_details, request_iterator)
    
    async def intercept_stream_stream(
        self,
        continuation: Callable[[grpc.aio.ClientCallDetails, Any], Awaitable[Any]],
        client_call_details: grpc.aio.ClientCallDetails,
        request_iterator: Any,
    ) -> Any:
        """Intercept stream-stream calls."""
        new_details = self._augment_call_details(client_call_details)
        return await continuation(new_details, request_iterator)

