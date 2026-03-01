package credential

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	connectsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/connect"
	"google.golang.org/grpc"
)

// CheckSlugAvailability verifies whether a slug is available for a given
// credential kind within an organization.
func CheckSlugAvailability(ctx context.Context, serverAddress, org, kind, slug string) (string, error) {
	apiKind, ok := credentialKindToAPIResourceKind[kind]
	if !ok {
		return "", fmt.Errorf("unknown credential kind %q — valid values: %s", kind, validCredentialKindNames())
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := connectsearch.NewConnectSearchQueryControllerClient(conn)
			resp, err := client.CheckConnectionSlugAvailability(ctx, &connectsearch.ConnectionSlugAvailabilityCheckRequest{
				Org:  org,
				Kind: apiKind,
				Slug: slug,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("slug %q for %s in org %q", slug, kind, org))
			}
			return domains.MarshalJSON(resp)
		})
}
