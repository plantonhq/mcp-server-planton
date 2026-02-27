package cloudresource

import (
	"context"

	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudresourcev1 "github.com/plantonhq/planton/apis/stubs/go/ai/planton/infrahub/cloudresource/v1"
	"google.golang.org/grpc"
)

// ListLocks retrieves lock information for a cloud resource via the
// CloudResourceLockController.ListLocks RPC.
//
// The response includes whether the resource is currently locked, details
// about the lock holder (workflow ID, acquired timestamp, TTL remaining),
// and any workflows queued for the lock.
//
// The caller must validate the ResourceIdentifier before calling this function.
func ListLocks(ctx context.Context, serverAddress string, id ResourceIdentifier) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveResourceID(ctx, conn, id)
			if err != nil {
				return "", err
			}

			desc := describeIdentifier(id)
			client := cloudresourcev1.NewCloudResourceLockControllerClient(conn)
			info, err := client.ListLocks(ctx, &cloudresourcev1.CloudResourceId{Value: resourceID})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(info)
		})
}

// RemoveLocks removes all locks (active lock and wait queue) for a cloud
// resource via the CloudResourceLockController.RemoveLocks RPC.
//
// WARNING: Removing locks on a resource with an active stack job may cause
// IaC state corruption. Callers should verify that no stack jobs are in
// progress before invoking this operation.
//
// The caller must validate the ResourceIdentifier before calling this function.
func RemoveLocks(ctx context.Context, serverAddress string, id ResourceIdentifier) (string, error) {
	return domains.WithConnection(ctx, serverAddress,
		func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
			resourceID, err := resolveResourceID(ctx, conn, id)
			if err != nil {
				return "", err
			}

			desc := describeIdentifier(id)
			client := cloudresourcev1.NewCloudResourceLockControllerClient(conn)
			resp, err := client.RemoveLocks(ctx, &cloudresourcev1.CloudResourceId{Value: resourceID})
			if err != nil {
				return "", domains.RPCError(err, desc)
			}
			return domains.MarshalJSON(resp)
		})
}
