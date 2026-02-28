package dnsdomain

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	dnsdomainv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/dnsdomain/v1"
	"google.golang.org/grpc"
)

// Delete removes a DNS domain record via the
// DnsDomainCommandController.Delete RPC.
//
// Two identification paths are supported:
//   - ID path: calls Delete directly with the given ID.
//   - Slug path: first resolves org+slug to a domain ID via the query
//     controller, then calls Delete. Both calls share a single gRPC connection.
func Delete(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveDomainID(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}

			desc := describeDomain(id, org, slug)
			client := dnsdomainv1.NewDnsDomainCommandControllerClient(conn)
			deleted, err := client.Delete(ctx, &apiresource.ApiResourceDeleteInput{
				ResourceId: resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(deleted)
		})
}
