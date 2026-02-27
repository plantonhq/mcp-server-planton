package cloudresource

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"google.golang.org/grpc"
)

// ResolveValueReferences resolves all valueFrom references in a cloud
// resource's spec via the CloudResourceQueryController.ResolveValueFromReferences
// RPC.
//
// The server loads the resource from the database, walks its specification to
// find all valueFrom references, resolves them to concrete values, and returns
// the fully transformed cloud resource as YAML along with resolution status,
// errors, and diagnostics.
//
// The caller must validate inputs before calling this function: kind must be a
// valid CloudResourceKind enum value, and the ResourceIdentifier must identify
// exactly one resource.
func ResolveValueReferences(ctx context.Context, serverAddress string, kind cloudresourcekind.CloudResourceKind, id ResourceIdentifier) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveResourceID(ctx, conn, id)
			if err != nil {
				return "", err
			}

			desc := describeIdentifier(id)
			client := cloudresourcev1.NewCloudResourceQueryControllerClient(conn)
			resp, err := client.ResolveValueFromReferences(ctx, &cloudresourcev1.ResolveValueFromReferencesRequest{
				CloudResourceKind: kind,
				CloudResourceId:   resourceID,
			})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(resp)
		})
}
