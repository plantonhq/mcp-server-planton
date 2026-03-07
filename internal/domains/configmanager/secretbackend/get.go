package secretbackend

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	secretbackendv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/configmanager/secretbackend/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
)

// Get retrieves a secret backend by ID or by org+slug.
func Get(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			sb, err := resolveSecretBackend(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}
			RedactSecretBackend(sb)
			return domains.MarshalJSON(sb)
		})
}

// resolveSecretBackend fetches the full SecretBackend proto by ID or by org+slug.
func resolveSecretBackend(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (*secretbackendv1.SecretBackend, error) {
	client := secretbackendv1.NewSecretBackendQueryControllerClient(conn)

	if id != "" {
		resp, err := client.Get(ctx, &secretbackendv1.SecretBackendId{Value: id})
		if err != nil {
			return nil, domains.RPCError(err, fmt.Sprintf("secret backend %q", id))
		}
		return resp, nil
	}

	resp, err := client.GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{
		Org:  org,
		Slug: slug,
	})
	if err != nil {
		return nil, domains.RPCError(err, fmt.Sprintf("secret backend %q in org %q", slug, org))
	}
	return resp, nil
}

// resolveSecretBackendID resolves identification inputs to a system-assigned ID.
func resolveSecretBackendID(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (string, error) {
	if id != "" {
		return id, nil
	}

	sb, err := resolveSecretBackend(ctx, conn, id, org, slug)
	if err != nil {
		return "", err
	}

	resourceID := sb.GetMetadata().GetId()
	if resourceID == "" {
		return "", fmt.Errorf("resolved secret backend %q in org %q but it has no ID — this indicates a backend issue", slug, org)
	}
	return resourceID, nil
}

// describeSecretBackend returns a human-readable description for error messages.
func describeSecretBackend(id, org, slug string) string {
	if id != "" {
		return fmt.Sprintf("secret backend %q", id)
	}
	return fmt.Sprintf("secret backend %q in org %q", slug, org)
}
