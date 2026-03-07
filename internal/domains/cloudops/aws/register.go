package aws

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, ListEc2InstancesTool(), ListEc2InstancesHandler(serverAddress))
	mcp.AddTool(srv, ListVpcsTool(), ListVpcsHandler(serverAddress))
	mcp.AddTool(srv, ListSubnetsTool(), ListSubnetsHandler(serverAddress))
	mcp.AddTool(srv, ListSecurityGroupsTool(), ListSecurityGroupsHandler(serverAddress))
	mcp.AddTool(srv, ListAvailabilityZonesTool(), ListAvailabilityZonesHandler(serverAddress))
	mcp.AddTool(srv, ListS3BucketsTool(), ListS3BucketsHandler(serverAddress))
}
