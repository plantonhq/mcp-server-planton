package infrachart

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	infrachartv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/infrachart/v1"
	"google.golang.org/grpc"
)

// Get retrieves an infra chart by ID via the
// InfraChartQueryController.Get RPC.
//
// The returned chart includes the full spec with template YAML files,
// values.yaml, parameter definitions, description, and web links.
func Get(ctx context.Context, serverAddress, chartID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := infrachartv1.NewInfraChartQueryControllerClient(conn)
			resp, err := client.Get(ctx, &apiresource.ApiResourceId{Value: chartID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("infra chart %q", chartID))
			}
			return domains.MarshalJSON(resp)
		})
}
