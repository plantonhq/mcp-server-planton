package variable

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	variablev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/configmanager/variable/v1"
	"google.golang.org/grpc"
)

// Resolve retrieves a variable's plain string value via the
// VariableQueryController.Resolve RPC.
//
// Unlike Get which returns the full Variable resource with metadata and spec,
// Resolve returns only the resolved value as a string. This is convenient for
// operational lookups where the agent just needs the value.
func Resolve(ctx context.Context, serverAddress, org string, scope variablev1.VariableSpec_Scope, slug string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := variablev1.NewVariableQueryControllerClient(conn)
			resp, err := client.Resolve(ctx, &variablev1.ResolveVariableRequest{
				Org:   org,
				Scope: scope,
				Slug:  slug,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("variable %q (scope=%s) in org %q", slug, scope, org))
			}
			return resp.GetValue(), nil
		})
}
