package audit

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresourceversionv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/audit/apiresourceversion/v1"
	"github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource/apiresourcekind"
	"google.golang.org/grpc"
)

// Count returns the number of versions that exist for a specific resource via
// the ApiResourceVersionQueryController.GetCount RPC.
//
// This is a lightweight check — no version data is transferred. Useful for
// quick "has anything changed?" queries or for deciding whether to paginate
// through the full version list.
func Count(ctx context.Context, serverAddress string, kind apiresourcekind.ApiResourceKind, resourceID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := apiresourceversionv1.NewApiResourceVersionQueryControllerClient(conn)
			resp, err := client.GetCount(ctx, &apiresourceversionv1.GetApiResourceVersionCountInput{
				Kind: kind,
				Id:   resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("version count for %s %q", kind, resourceID))
			}
			return domains.MarshalJSON(resp)
		})
}
