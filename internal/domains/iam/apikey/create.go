package apikey

import (
	"context"
	"fmt"
	"time"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	apikeyv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/apikey/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Create generates a new API key via ApiKeyCommandController.Create.
// The raw key value is returned ONLY in the create response and is never
// persisted or retrievable afterward.
func Create(ctx context.Context, serverAddress, name string, neverExpires bool, expiresAt string) (string, error) {
	spec := &apikeyv1.ApiKeySpec{
		NeverExpires: neverExpires,
	}
	if expiresAt != "" && !neverExpires {
		t, err := time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			return "", fmt.Errorf("'expires_at' must be in RFC 3339 format (e.g. 2026-12-31T23:59:59Z): %w", err)
		}
		spec.ExpiresAt = timestamppb.New(t)
	}

	key := &apikeyv1.ApiKey{
		ApiVersion: "iam.planton.ai/v1",
		Kind:       "ApiKey",
		Metadata:   &apiresource.ApiResourceMetadata{Name: name},
		Spec:       spec,
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := apikeyv1.NewApiKeyCommandControllerClient(conn)
			resp, err := client.Create(ctx, key)
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("API key %q", name))
			}
			return domains.MarshalJSON(resp)
		})
}
