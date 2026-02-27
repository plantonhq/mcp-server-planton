package cloudresource

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	cloudResourceAPIVersion = "infra-hub.planton.ai/v1"
	cloudResourceKindConst  = "CloudResource"
)

// Apply creates or updates a cloud resource via the CloudResourceCommandController.Apply RPC.
// The CloudResource proto must be fully constructed before calling this function.
func Apply(ctx context.Context, serverAddress string, cr *cloudresourcev1.CloudResource) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := cloudresourcev1.NewCloudResourceCommandControllerClient(conn)
			result, err := client.Apply(ctx, cr)
			if err != nil {
				desc := describeResource(cr)
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(result)
		})
}

// buildCloudResource assembles a CloudResource proto from the parsed components.
//
//   - cloudObject: the original cloud_object map (used for metadata extraction)
//   - kindStr: the PascalCase kind string (e.g. "AwsEksCluster")
//   - normalizedObject: the spec-validated structpb.Struct from the generated parser
func buildCloudResource(cloudObject map[string]any, kindStr string, normalizedObject *structpb.Struct) (*cloudresourcev1.CloudResource, error) {
	kind, err := resolveKind(kindStr)
	if err != nil {
		return nil, err
	}

	metadata, err := extractMetadata(cloudObject)
	if err != nil {
		return nil, err
	}

	return &cloudresourcev1.CloudResource{
		ApiVersion: cloudResourceAPIVersion,
		Kind:       cloudResourceKindConst,
		Metadata:   metadata,
		Spec: &cloudresourcev1.CloudResourceSpec{
			Kind:        kind,
			CloudObject: normalizedObject,
		},
	}, nil
}

// describeResource returns a human-readable description for error messages.
func describeResource(cr *cloudresourcev1.CloudResource) string {
	md := cr.GetMetadata()
	if md == nil {
		return "cloud resource"
	}
	name := md.GetName()
	if name == "" {
		name = md.GetSlug()
	}
	return fmt.Sprintf("cloud resource %q in org %q env %q", name, md.GetOrg(), md.GetEnv())
}
