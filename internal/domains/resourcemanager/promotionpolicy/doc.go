// Package promotionpolicy provides the MCP tools for the PromotionPolicy
// domain, backed by the PromotionPolicyQueryController and
// PromotionPolicyCommandController RPCs
// (ai.planton.resourcemanager.promotionpolicy.v1) on the Planton backend.
//
// Promotion policies define the DAG of environments through which deployments
// are promoted (e.g. dev -> staging -> production). Policies are scoped to an
// organization or to the platform as a default. The "which" query resolves
// the effective policy for any given scope with inheritance (org-specific
// takes precedence over the platform default).
//
// Four tools are exposed:
//   - apply_promotion_policy:  create or update a promotion policy
//   - get_promotion_policy:    retrieve a policy by ID or by selector (kind + id)
//   - which_promotion_policy:  resolve the effective policy for a given scope
//   - delete_promotion_policy: delete a policy by ID
package promotionpolicy
