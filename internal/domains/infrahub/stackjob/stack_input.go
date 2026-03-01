package stackjob

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
	"google.golang.org/grpc"
)

// GetStackInput retrieves the safe (credential-free) stack input for a stack
// job via the StackJobQueryController.GetCloudObjectStackInput RPC.
//
// The response contains the exact data that was fed to the Pulumi or Terraform
// module — target spec, provider config, and docker config — but excludes
// platform-level backend credentials.
func GetStackInput(ctx context.Context, serverAddress, stackJobID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobQueryControllerClient(conn)
			resp, err := client.GetCloudObjectStackInput(ctx, &stackjobv1.StackJobId{Value: stackJobID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("stack input for stack job %q", stackJobID))
			}
			return domains.MarshalJSON(resp)
		})
}
