package environment

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	environmentv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/resourcemanager/environment/v1"
	"google.golang.org/grpc"
)

// Create provisions a new environment within an organization via the
// EnvironmentCommandController.Create RPC.
//
// The caller supplies user-facing fields; the function assembles the full
// Environment proto with the required api_version and kind constants.
func Create(ctx context.Context, serverAddress, org, slug, name, description string) (string, error) {
	env := &environmentv1.Environment{
		ApiVersion: "resource-manager.planton.ai/v1",
		Kind:       "Environment",
		Metadata:   &apiresource.ApiResourceMetadata{Org: org, Slug: slug, Name: name},
		Spec:       &environmentv1.EnvironmentSpec{Description: description},
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := environmentv1.NewEnvironmentCommandControllerClient(conn)
			resp, err := client.Create(ctx, env)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("environment %q in org %q", slug, org))
			}
			return domains.MarshalJSON(resp)
		})
}
