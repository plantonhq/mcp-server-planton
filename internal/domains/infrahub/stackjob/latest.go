package stackjob

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
	"google.golang.org/grpc"
)

// GetLatest retrieves the most recent stack job for a cloud resource via the
// StackJobQueryController.GetLastStackJobByCloudResourceId RPC.
//
// This is the primary function agents call after apply or destroy to check
// whether the operation completed successfully.
func GetLatest(ctx context.Context, serverAddress, cloudResourceID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobQueryControllerClient(conn)
			resp, err := client.GetLastStackJobByCloudResourceId(ctx, &cloudresourcev1.CloudResourceId{Value: cloudResourceID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("latest stack job for cloud resource %q", cloudResourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
