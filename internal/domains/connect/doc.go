// Package connect is the parent package for the Connect bounded context,
// which manages connections to third-party cloud providers and platform
// connection resources (DefaultProviderConnection, DefaultRunnerBinding,
// RunnerRegistration).
//
// Sub-packages:
//   - connection: Generic connection CRUD tools with type dispatch (5 tools + 2 MCP resources)
//   - github: GitHub-specific extras (webhooks, installation info, repository listing)
//   - defaultprovider: Default provider connection management
//   - defaultrunner: Default runner binding management
//   - runner: Runner registration management
//   - providerauth: Provider connection authorization management
package connect
