package preset

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	infrahubsearch "github.com/plantonhq/planton/apis/stubs/go/ai/planton/search/v1/infrahub"
	"google.golang.org/grpc"
)

// SearchInput holds the validated parameters for searching cloud object presets.
type SearchInput struct {
	Org        string
	Kind       string
	SearchText string
}

// Search queries cloud object presets via the InfraHubSearchQueryController.
//
// When Org is provided, the org-context RPC is called with both official and
// organization presets included. When Org is empty, only official (public)
// presets are searched.
//
// The optional Kind filter is resolved to a CloudResourceKind enum value
// before being sent to the backend.
func Search(ctx context.Context, serverAddress string, input SearchInput) (string, error) {
	var kind cloudresourcekind.CloudResourceKind
	if input.Kind != "" {
		k, err := domains.ResolveKind(input.Kind)
		if err != nil {
			return "", err
		}
		kind = k
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrahubsearch.NewInfraHubSearchQueryControllerClient(conn)

			if input.Org != "" {
				resp, err := client.SearchCloudObjectPresetsByOrgContext(ctx, &infrahubsearch.SearchCloudObjectPresetsByOrgContextInput{
					Org:                          input.Org,
					SearchText:                   input.SearchText,
					CloudResourceKind:            kind,
					IsIncludeOfficial:            true,
					IsIncludeOrganizationPresets: true,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("cloud object presets in org %q", input.Org))
				}
				return domains.MarshalJSON(resp)
			}

			resp, err := client.SearchOfficialCloudObjectPresets(ctx, &infrahubsearch.SearchOfficialCloudObjectPresetsInput{
				SearchText:        input.SearchText,
				CloudResourceKind: kind,
			})
			if err != nil {
				return "", domains.RPCError(err, "official cloud object presets")
			}
			return domains.MarshalJSON(resp)
		})
}
