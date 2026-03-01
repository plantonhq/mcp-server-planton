# Implementation Plans

Plans created for the MCP Server Gap Completion project. Each plan documents a specific implementation effort.

## Plan Registry

| Plan | Status | Created | Completed | Description |
|------|--------|---------|-----------|-------------|
| `t07-purge-cloud-resource.plan.md` | Completed | 2026-03-01 | 2026-03-01 | Add `purge_cloud_resource` tool to CloudResource domain |
| `t03-organization-full-crud.plan.md` | Completed | 2026-03-01 | 2026-03-01 | Add get, create, update, delete tools to Organization domain |
| `t04-environment-full-crud.plan.md` | Completed | 2026-03-01 | 2026-03-01 | Add get (dual-resolution), create, update, delete tools to Environment domain |
| `t06-stackjob-ai-native-tools.plan.md` | Completed | 2026-03-01 | 2026-03-01 | Add 5 AI-native and diagnostic tools to StackJob domain (IaC resources, stack input, service env status, error recommendation) |
| `t05-connect-domain-credentials.plan.md` | Completed | 2026-03-01 | 2026-03-01 | Implement Connect domain: 22 tools + 2 MCP resources across 5 sub-packages (credential, github, defaultprovider, defaultrunner, runner) for 19 credential types and 3 platform connection types |
| `t08-iam-domain.plan.md` | Completed | 2026-03-01 | 2026-03-01 | Implement IAM bounded context (20 tools across 5 sub-packages) + ProviderConnectionAuthorization (3 tools in connect/providerauth), totaling 23 new tools across 7 phases |
| `t09-t10-t11-remaining-tier2-tools.plan.md` | Completed | 2026-03-01 | 2026-03-01 | Add 8 tools across 3 tasks: delete_infra_pipeline (T09), PromotionPolicy CRUD+query (T10, 4 tools), FlowControlPolicy CRUD (T11, 3 tools) |
| `t12-expand-mcp-resources.plan.md` | Completed | 2026-03-01 | 2026-03-01 | Add api-resource-kinds://catalog MCP resource — platform navigational index (29 kinds, 6 domains). Scope reduced from 5 to 1 resource after discovering 3 already covered by tools, 1 delivered in T05 |
| `t15-mcp-prompts.plan.md` | Completed | 2026-03-01 | 2026-03-01 | Add 5 cross-domain MCP prompts (debug_failed_deployment, assess_change_impact, explore_infrastructure, provision_cloud_resource, manage_access) — first implementation of the third MCP primitive |

---

*Last updated: 2026-03-01*
