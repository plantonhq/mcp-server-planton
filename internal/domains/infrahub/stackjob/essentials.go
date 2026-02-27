package stackjob

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	apiresource "github.com/plantonhq/planton/apis/stubs/go/ai/planton/commons/apiresource"
	stackjobv1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/stackjob/v1"
	"google.golang.org/grpc"
)

// CheckEssentials validates that all prerequisites required to run a stack
// job are in place for a given cloud resource kind and owner (org + optional
// env) via the StackJobEssentialsQueryController.Check RPC.
//
// The response contains four preflight checks — iac_module,
// backend_credential, flow_control, and provider_credential — each with
// a passed flag and any errors encountered.
func CheckEssentials(ctx context.Context, serverAddress string, kind, org, env string) (string, error) {
	resolvedKind, err := domains.ResolveKind(kind)
	if err != nil {
		return "", err
	}

	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := stackjobv1.NewStackJobEssentialsQueryControllerClient(conn)
			resp, err := client.Check(ctx, &stackjobv1.CheckStackJobEssentialsInput{
				CloudResourceKind: resolvedKind,
				CloudResourceOwner: &apiresource.CloudResourceOwner{
					Org: org,
					Env: env,
				},
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf(
					"check stack job essentials for kind %q in org %q", kind, org))
			}
			return domains.MarshalJSON(resp)
		})
}
