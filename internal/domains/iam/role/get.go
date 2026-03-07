package role

import (
	"context"
	"fmt"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	iamrolev1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/iam/iamrole/v1"
	"google.golang.org/grpc"
)

// Get retrieves an IAM role by ID via IamRoleQueryController.Get.
func Get(ctx context.Context, serverAddress, roleID string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := iamrolev1.NewIamRoleQueryControllerClient(conn)
			resp, err := client.Get(ctx, &iamrolev1.IamRoleId{Value: roleID})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("IAM role %q", roleID))
			}
			return domains.MarshalJSON(resp)
		})
}

// ListForResourceKind retrieves all IAM roles available for a given
// resource kind and principal type, via
// IamRoleQueryController.FindByApiResourceKindAndPrincipalType.
func ListForResourceKind(ctx context.Context, serverAddress, resourceKind, principalType string) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			client := iamrolev1.NewIamRoleQueryControllerClient(conn)
			resp, err := client.FindByApiResourceKindAndPrincipalType(ctx, &iamrolev1.ApiResourceKindAndPrincipalTypeInput{
				ResourceKind:  resourceKind,
				PrincipalType: principalType,
			})
			if err != nil {
				return "", domains.RPCError(err, fmt.Sprintf("IAM roles for resource kind %q and principal type %q", resourceKind, principalType))
			}
			return domains.MarshalJSON(resp)
		})
}
