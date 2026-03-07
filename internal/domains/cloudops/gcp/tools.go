// Package gcp provides MCP tools for GCP cloud operations via the Planton control plane.
package gcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	gcpcompute "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/gcp/v1/compute"
	gcpstorage "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/gcp/v1/storage"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudopsctx "github.com/plantonhq/mcp-server-planton/internal/domains/cloudops"
	"google.golang.org/grpc"
)

type listComputeInstancesInput struct {
	Org               string `json:"org"                        jsonschema:"required,Organization slug."`
	Env               string `json:"env,omitempty"              jsonschema:"Environment slug. Required for cloud_resource access mode."`
	CloudResourceKind string `json:"cloud_resource_kind,omitempty" jsonschema:"Cloud resource kind (PascalCase). Use with cloud_resource_slug for cloud resource access mode."`
	CloudResourceSlug string `json:"cloud_resource_slug,omitempty" jsonschema:"Cloud resource slug. Use with cloud_resource_kind for cloud resource access mode."`
	Connection        string `json:"connection,omitempty"       jsonschema:"Provider connection slug for direct access. Mutually exclusive with cloud resource fields."`
	Project           string `json:"project,omitempty"          jsonschema:"GCP project ID to query. Optional."`
	Zone              string `json:"zone,omitempty"             jsonschema:"GCP zone to query. Use '-' for all zones. Optional."`
	Filter            string `json:"filter,omitempty"           jsonschema:"GCP filter expression. Optional."`
}

func ListComputeInstancesTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_gcp_compute_instances",
		Description: "List GCP Compute Engine instances in a project and optionally a zone. Returns instance details including name, status, machine type, and network configuration. Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListComputeInstancesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *listComputeInstancesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *listComputeInstancesInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := gcpcompute.NewComputeInstanceQueryControllerClient(conn).List(ctx, &gcpcompute.ListComputeInstancesRequest{
					Context: opsCtx,
					Project: input.Project,
					Zone:    input.Zone,
					Filter:  input.Filter,
				})
				if err != nil {
					return "", domains.RPCError(err, "list GCP compute instances")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type listStorageBucketsInput struct {
	Org               string `json:"org"                        jsonschema:"required,Organization slug."`
	Env               string `json:"env,omitempty"              jsonschema:"Environment slug. Required for cloud_resource access mode."`
	CloudResourceKind string `json:"cloud_resource_kind,omitempty" jsonschema:"Cloud resource kind (PascalCase). Use with cloud_resource_slug for cloud resource access mode."`
	CloudResourceSlug string `json:"cloud_resource_slug,omitempty" jsonschema:"Cloud resource slug. Use with cloud_resource_kind for cloud resource access mode."`
	Connection        string `json:"connection,omitempty"       jsonschema:"Provider connection slug for direct access. Mutually exclusive with cloud resource fields."`
	Project           string `json:"project,omitempty"          jsonschema:"GCP project ID to query. Optional."`
}

func ListStorageBucketsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "list_gcp_storage_buckets",
		Description: "List GCP Cloud Storage buckets in a project. Returns bucket metadata including name, location, and storage class. Use cloud resource or connection access mode to specify credentials.",
	}
}

func ListStorageBucketsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *listStorageBucketsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *listStorageBucketsInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				resp, err := gcpstorage.NewStorageBucketQueryControllerClient(conn).List(ctx, &gcpstorage.ListStorageBucketsRequest{
					Context: opsCtx,
					Project: input.Project,
				})
				if err != nil {
					return "", domains.RPCError(err, "list GCP storage buckets")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
