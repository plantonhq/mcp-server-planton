package preset

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	cloudobjectpresetv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudobjectpreset/v1"
	"google.golang.org/grpc"
)

// Get retrieves a cloud object preset by ID via the
// CloudObjectPresetQueryController.Get RPC.
//
// The returned preset includes the full spec with YAML content, markdown
// documentation, kind, rank, and provider metadata.
func Get(ctx context.Context, serverAddress, presetID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := cloudobjectpresetv1.NewCloudObjectPresetQueryControllerClient(conn)
			resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: presetID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("cloud object preset %q", presetID))
			}
			return domains.MarshalJSON(resp)
		})
}
