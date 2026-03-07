package gcp

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ListComputeInstancesTool(), ListComputeInstancesHandler(serverAddress))
	mcp.AddTool(srv, ListStorageBucketsTool(), ListStorageBucketsHandler(serverAddress))
}
