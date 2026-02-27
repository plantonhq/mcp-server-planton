package iacmodule

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	iacmodulev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/iacmodule/v1"
	"google.golang.org/grpc"
)

// Get retrieves an IaC module by ID via the
// IacModuleQueryController.Get RPC.
//
// The returned module includes the full spec with provisioner type,
// cloud resource kind, Git repository URL, version, and parameter schema.
func Get(ctx context.Context, serverAddress, moduleID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := iacmodulev1.NewIacModuleQueryControllerClient(conn)
			resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: moduleID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("iac module %q", moduleID))
			}
			return domains.MarshalJSON(resp)
		})
}
