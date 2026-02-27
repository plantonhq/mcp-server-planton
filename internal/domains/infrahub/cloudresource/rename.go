package cloudresource

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"google.golang.org/grpc"
)

// Rename changes the display name of a cloud resource via the
// CloudResourceCommandController.Rename RPC.
//
// The slug (immutable identifier) is unaffected — only the human-readable
// display name is updated.
//
// The caller must validate the ResourceIdentifier before calling this function.
func Rename(ctx context.Context, serverAddress string, id ResourceIdentifier, newName string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveResourceID(ctx, conn, id)
			if err != nil {
				return "", err
			}

			desc := describeIdentifier(id)
			client := cloudresourcev1.NewCloudResourceCommandControllerClient(conn)
			cr, err := client.Rename(ctx, &cloudresourcev1.RenameCloudResourceRequest{
				Id:   resourceID,
				Name: newName,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(cr)
		})
}
