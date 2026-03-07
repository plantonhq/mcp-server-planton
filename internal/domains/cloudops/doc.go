// Package cloudops provides shared utilities for the CloudOps MCP tools.
//
// CloudOps tools expose live cloud infrastructure operations (list pods,
// describe EC2 instances, etc.) through the Planton control plane. All
// operations route through the CloudOpsRequestContext, which the control
// plane uses to resolve credentials and authorize access without requiring
// raw cloud credentials from the MCP client.
//
// Sub-packages:
//   - kubernetes: Kubernetes object, pod, namespace, and secret operations
//   - aws: EC2, VPC, subnet, security group, availability zone, and S3 operations
//   - gcp: Compute instance and storage bucket operations
//   - azure: Virtual machine and blob container operations
package cloudops
