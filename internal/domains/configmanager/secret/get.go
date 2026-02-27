package secret

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	secretv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/secret/v1"
	"google.golang.org/grpc"
)

// Get retrieves a secret by ID or by org+scope+slug via the
// SecretQueryController RPCs.
//
// Two identification paths are supported:
//   - ID path: calls Get(SecretId) directly.
//   - Slug path: calls GetByOrgByScopeBySlug(SecretScopeSlugRequest).
//
// Only metadata is returned — no secret values are exposed.
func Get(ctx context.Context, serverAddress, id, org string, scope secretv1.SecretSpec_Scope, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			s, err := resolveSecret(ctx, conn, id, org, scope, slug)
			if err != nil {
				return "", err
			}
			return domains.MarshalJSON(s)
		})
}

// resolveSecret fetches the full Secret proto by ID or by org+scope+slug.
// Used by operations that need the full resource before acting on it.
func resolveSecret(ctx context.Context, conn *grpc.ClientConn, id, org string, scope secretv1.SecretSpec_Scope, slug string) (*secretv1.Secret, error) {
	client := secretv1.NewSecretQueryControllerClient(conn)

	if id != "" {
		resp, err := client.Get(ctx, &secretv1.SecretId{Value: id})
		if err != nil {
			return nil, domains.RPCError(err, fmt.Sprintf("secret %q", id))
		}
		return resp, nil
	}

	resp, err := client.GetByOrgByScopeBySlug(ctx, &secretv1.SecretScopeSlugRequest{
		Org:   org,
		Scope: scope,
		Slug:  slug,
	})
	if err != nil {
		return nil, domains.RPCError(err, fmt.Sprintf("secret %q (scope=%s) in org %q", slug, scope, org))
	}
	return resp, nil
}

// resolveSecretID resolves identification inputs to a system-assigned
// secret ID string. When an ID is already provided it is returned directly.
// Otherwise the secret is fetched by org+scope+slug and its metadata ID is
// extracted.
func resolveSecretID(ctx context.Context, conn *grpc.ClientConn, id, org string, scope secretv1.SecretSpec_Scope, slug string) (string, error) {
	if id != "" {
		return id, nil
	}

	s, err := resolveSecret(ctx, conn, id, org, scope, slug)
	if err != nil {
		return "", err
	}

	resourceID := s.GetMetadata().GetId()
	if resourceID == "" {
		return "", fmt.Errorf("resolved secret %q (scope=%s) in org %q but it has no ID — this indicates a backend issue", slug, scope, org)
	}
	return resourceID, nil
}

// describeSecret returns a human-readable description for error messages.
func describeSecret(id, org string, scope secretv1.SecretSpec_Scope, slug string) string {
	if id != "" {
		return fmt.Sprintf("secret %q", id)
	}
	return fmt.Sprintf("secret %q (scope=%s) in org %q", slug, scope, org)
}
