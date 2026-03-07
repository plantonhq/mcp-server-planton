package stackjob

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	stackjobv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/infrahub/stackjob/v1"
	"google.golang.org/grpc"
)

// FindIacResourcesByStackJob retrieves all IaC resources (Pulumi/Terraform state
// entries) associated with a specific stack job via the
// StackJobQueryController.FindIacResourcesByStackJobId RPC.
func FindIacResourcesByStackJob(ctx context.Context, serverAddress, stackJobID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobQueryControllerClient(conn)
			resp, err := client.FindIacResourcesByStackJobId(ctx, &stackjobv1.StackJobId{Value: stackJobID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("IaC resources for stack job %q", stackJobID))
			}
			return domains.MarshalJSON(resp)
		})
}

// FindIacResourcesByApiResource retrieves all IaC resources from the most recent
// stack job for a given API resource via the
// StackJobQueryController.FindIacResourcesByApiResourceId RPC.
func FindIacResourcesByApiResource(ctx context.Context, serverAddress, apiResourceID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobQueryControllerClient(conn)
			resp, err := client.FindIacResourcesByApiResourceId(ctx, &apiresource.ApiResourceId{Value: apiResourceID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("IaC resources for API resource %q", apiResourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
