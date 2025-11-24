"""Configuration for Planton Cloud MCP Server."""

import os
from pydantic_settings import BaseSettings


class MCPServerConfig(BaseSettings):
    """
    Configuration loaded from environment variables.
    
    Unlike agent-fleet-worker (which uses machine account), this server
    expects USER_JWT_TOKEN to be passed via environment by LangGraph.
    """
    
    # User authentication (passed by LangGraph via environment)
    user_jwt_token: str
    
    # Planton Cloud APIs endpoint
    planton_apis_grpc_endpoint: str = "localhost:8080"
    
    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
    
    @classmethod
    def load_from_env(cls) -> "MCPServerConfig":
        """
        Load configuration from environment variables.
        
        Raises:
            ValueError: If USER_JWT_TOKEN is missing
        """
        user_jwt = os.environ.get("USER_JWT_TOKEN")
        
        if not user_jwt:
            raise ValueError(
                "USER_JWT_TOKEN environment variable required. "
                "This should be set by LangGraph when spawning MCP server."
            )
        
        return cls(
            user_jwt_token=user_jwt,
            planton_apis_grpc_endpoint=os.environ.get(
                "PLANTON_APIS_GRPC_ENDPOINT",
                "localhost:8080"
            )
        )

