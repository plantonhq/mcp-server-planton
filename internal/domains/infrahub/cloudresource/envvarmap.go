package cloudresource

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"google.golang.org/grpc"
)

// GetEnvVarMap retrieves the environment variable map from a cloud resource
// manifest via the CloudResourceQueryController.GetEnvVarMap RPC.
//
// The server parses the provided YAML to identify the resource kind, extract
// environment variables and secrets, and resolve valueFrom references.
// Authorization is handled server-side by extracting the cloud resource
// identity from the YAML content.
func GetEnvVarMap(ctx context.Context, serverAddress, yamlContent string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := cloudresourcev1.NewCloudResourceQueryControllerClient(conn)
			resp, err := client.GetEnvVarMap(ctx, &cloudresourcev1.GetEnvVarMapRequest{
				YamlContent: yamlContent,
			})
			if err != nil {
				return "", domains.RPCError(err, "environment variable map")
			}
			return domains.MarshalJSON(resp)
		})
}
