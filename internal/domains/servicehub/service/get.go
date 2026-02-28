package service

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	servicev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/servicehub/service/v1"
	"google.golang.org/grpc"
)

// Get retrieves a service by ID or by org+slug via the
// ServiceQueryController RPCs.
//
// Two identification paths are supported:
//   - ID path: calls Get(ServiceId) directly.
//   - Slug path: calls GetByOrgBySlug(ApiResourceByOrgBySlugRequest).
func Get(ctx context.Context, serverAddress, id, org, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			svc, err := resolveService(ctx, conn, id, org, slug)
			if err != nil {
				return "", err
			}
			return domains.MarshalJSON(svc)
		})
}

// resolveService fetches the full Service proto by ID or by org+slug.
// Used by operations that need the full resource before acting on it.
func resolveService(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (*servicev1.Service, error) {
	client := servicev1.NewServiceQueryControllerClient(conn)

	if id != "" {
		resp, err := client.Get(ctx, &servicev1.ServiceId{Value: id})
		if err != nil {
			return nil, domains.RPCError(err, fmt.Sprintf("service %q", id))
		}
		return resp, nil
	}

	resp, err := client.GetByOrgBySlug(ctx, &apiresource.ApiResourceByOrgBySlugRequest{
		Org:  org,
		Slug: slug,
	})
	if err != nil {
		return nil, domains.RPCError(err, fmt.Sprintf("service %q in org %q", slug, org))
	}
	return resp, nil
}

// resolveServiceID resolves identification inputs to a system-assigned service
// ID string. When an ID is already provided it is returned directly. Otherwise
// the service is fetched by org+slug and its metadata ID is extracted.
func resolveServiceID(ctx context.Context, conn *grpc.ClientConn, id, org, slug string) (string, error) {
	if id != "" {
		return id, nil
	}

	svc, err := resolveService(ctx, conn, id, org, slug)
	if err != nil {
		return "", err
	}

	resourceID := svc.GetMetadata().GetId()
	if resourceID == "" {
		return "", fmt.Errorf("resolved service %q in org %q but it has no ID — this indicates a backend issue", slug, org)
	}
	return resourceID, nil
}

// describeService returns a human-readable description of the service for
// use in error messages.
func describeService(id, org, slug string) string {
	if id != "" {
		return fmt.Sprintf("service %q", id)
	}
	return fmt.Sprintf("service %q in org %q", slug, org)
}
