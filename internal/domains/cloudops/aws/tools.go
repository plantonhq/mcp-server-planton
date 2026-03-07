// Package aws provides MCP tools for AWS cloud operations via the Planton control plane.

package aws

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	awsec2 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/aws/v1/ec2"
	awsvpc "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/aws/v1/vpc"
	awssubnet "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/aws/v1/subnet"
	awssg "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/aws/v1/securitygroup"
	awsaz "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/aws/v1/availabilityzone"
	awss3 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/aws/v1/s3"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudopsctx "github.com/plantonhq/mcp-server-planton/internal/domains/cloudops"
	"google.golang.org/grpc"
)

type awsContextInput struct {
	Org               string `json:"org"                          jsonschema:"required,Organization slug."`
	Env               string `json:"env,omitempty"                jsonschema:"Environment slug. Required for cloud_resource access mode."`
	CloudResourceKind string `json:"cloud_resource_kind,omitempty" jsonschema:"Cloud resource kind (PascalCase). Use with cloud_resource_slug for cloud resource access mode."`
	CloudResourceSlug string `json:"cloud_resource_slug,omitempty" jsonschema:"Cloud resource slug. Use with cloud_resource_kind for cloud resource access mode."`
	Connection        string `json:"connection,omitempty"        jsonschema:"Provider connection slug for direct access. Mutually exclusive with cloud resource fields."`
}

type ListEc2InstancesInput struct {
	awsContextInput
	Region      string   `json:"region,omitempty"       jsonschema:"AWS region to query (e.g. us-west-2)."`
	InstanceIds []string `json:"instance_ids,omitempty" jsonschema:"Filter to specific instance IDs."`
}

func ListEc2InstancesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_ec2_instances",
		Description: "List EC2 instances in an AWS region via the Planton control plane. " +
			"Returns instance metadata including state, networking, and tags. " +
			"Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListEc2InstancesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListEc2InstancesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListEc2InstancesInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := awsec2.NewEc2InstanceQueryControllerClient(conn).List(ctx, &awsec2.ListEc2InstancesRequest{
					Context:    opsCtx,
					Region:     input.Region,
					InstanceIds: input.InstanceIds,
				})
				if err != nil {
					return "", domains.RPCError(err, "EC2 instances")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type ListVpcsInput struct {
	awsContextInput
	Region string `json:"region,omitempty" jsonschema:"AWS region to query (e.g. us-west-2)."`
}

func ListVpcsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_vpcs",
		Description: "List VPCs in an AWS region via the Planton control plane. " +
			"Returns VPC metadata including CIDR blocks, state, and tags. " +
			"Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListVpcsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListVpcsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListVpcsInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := awsvpc.NewVpcQueryControllerClient(conn).List(ctx, &awsvpc.ListVpcsRequest{
					Context: opsCtx,
					Region:  input.Region,
				})
				if err != nil {
					return "", domains.RPCError(err, "VPCs")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type ListSubnetsInput struct {
	awsContextInput
	Region string `json:"region,omitempty" jsonschema:"AWS region to query (e.g. us-west-2)."`
	VpcId  string `json:"vpc_id,omitempty" jsonschema:"Filter subnets to a specific VPC."`
}

func ListSubnetsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_subnets",
		Description: "List subnets in an AWS region via the Planton control plane. " +
			"Optionally filter by VPC. Returns subnet metadata including CIDR, availability zone, and tags. " +
			"Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListSubnetsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListSubnetsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListSubnetsInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := awssubnet.NewSubnetQueryControllerClient(conn).List(ctx, &awssubnet.ListSubnetsRequest{
					Context: opsCtx,
					Region:  input.Region,
					VpcId:   input.VpcId,
				})
				if err != nil {
					return "", domains.RPCError(err, "subnets")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type ListSecurityGroupsInput struct {
	awsContextInput
	Region string `json:"region,omitempty" jsonschema:"AWS region to query (e.g. us-west-2)."`
	VpcId  string `json:"vpc_id,omitempty" jsonschema:"Filter security groups to a specific VPC."`
}

func ListSecurityGroupsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_security_groups",
		Description: "List security groups in an AWS region via the Planton control plane. " +
			"Optionally filter by VPC. Returns group metadata including rules and tags. " +
			"Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListSecurityGroupsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListSecurityGroupsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListSecurityGroupsInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := awssg.NewSecurityGroupQueryControllerClient(conn).List(ctx, &awssg.ListSecurityGroupsRequest{
					Context: opsCtx,
					Region:  input.Region,
					VpcId:   input.VpcId,
				})
				if err != nil {
					return "", domains.RPCError(err, "security groups")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type ListAvailabilityZonesInput struct {
	awsContextInput
	Region string `json:"region,omitempty" jsonschema:"AWS region to query (e.g. us-west-2)."`
}

func ListAvailabilityZonesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_availability_zones",
		Description: "List availability zones in an AWS region via the Planton control plane. " +
			"Returns zone names and state. Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListAvailabilityZonesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListAvailabilityZonesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListAvailabilityZonesInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := awsaz.NewAvailabilityZoneQueryControllerClient(conn).List(ctx, &awsaz.ListAvailabilityZonesRequest{
					Context: opsCtx,
					Region:  input.Region,
				})
				if err != nil {
					return "", domains.RPCError(err, "availability zones")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type ListS3BucketsInput struct {
	awsContextInput
	Region string `json:"region,omitempty" jsonschema:"AWS region for SDK client (e.g. us-west-2). S3 ListBuckets is global but client needs a region."`
}

func ListS3BucketsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_s3_buckets",
		Description: "List S3 buckets via the Planton control plane. " +
			"Returns bucket names and optionally their regions. " +
			"Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListS3BucketsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListS3BucketsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListS3BucketsInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := awss3.NewS3BucketQueryControllerClient(conn).List(ctx, &awss3.ListS3BucketsRequest{
					Context: opsCtx,
					Region:  input.Region,
				})
				if err != nil {
					return "", domains.RPCError(err, "S3 buckets")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
