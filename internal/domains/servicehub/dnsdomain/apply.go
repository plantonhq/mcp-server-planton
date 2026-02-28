package dnsdomain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	dnsdomainv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/dnsdomain/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// Apply creates or updates a DNS domain via the
// DnsDomainCommandController.Apply RPC.
//
// The input is a raw JSON map matching the DnsDomain proto shape. It is
// serialized to JSON bytes and deserialized into the typed proto using
// protojson, which handles proto field-name conventions, enums, and
// well-known types correctly.
func Apply(ctx context.Context, serverAddress string, raw map[string]any) (string, error) {
	jsonBytes, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("failed to serialize DNS domain input: %w", err)
	}

	domain := &dnsdomainv1.DnsDomain{}
	if err := protojson.Unmarshal(jsonBytes, domain); err != nil {
		return "", fmt.Errorf("invalid DNS domain structure: %w", err)
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := dnsdomainv1.NewDnsDomainCommandControllerClient(conn)
			result, err := client.Apply(ctx, domain)
			if err != nil {
				desc := "DNS domain"
				if md := domain.GetMetadata(); md != nil {
					name := md.GetName()
					if name == "" {
						name = md.GetSlug()
					}
					if name != "" {
						desc = fmt.Sprintf("DNS domain %q in org %q", name, md.GetOrg())
					}
				}
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}
