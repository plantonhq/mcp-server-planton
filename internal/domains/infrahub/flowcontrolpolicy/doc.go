// Package flowcontrolpolicy provides the MCP tools for the FlowControlPolicy
// domain, backed by the FlowControlPolicyQueryController and
// FlowControlPolicyCommandController RPCs
// (ai.planton.infrahub.flowcontrolpolicy.v1) on the Planton backend.
//
// Flow control policies govern how stack jobs (Pulumi/Terraform operations)
// execute for a given scope. They control behaviors like manual approval
// gates, lifecycle event disabling, refresh skipping, and preview-before-
// update pauses. Policies are scoped to an organization, environment,
// platform, or individual cloud resource via an ApiResourceSelector.
//
// Three tools are exposed:
//   - apply_flow_control_policy:  create or update a flow control policy
//   - get_flow_control_policy:    retrieve a policy by ID or by selector (kind + id)
//   - delete_flow_control_policy: delete a policy by ID
package flowcontrolpolicy
