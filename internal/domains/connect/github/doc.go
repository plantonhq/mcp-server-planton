// Package github provides MCP tools for GitHub-specific operations that go
// beyond standard credential CRUD. These tools use the GithubCommandController
// and GithubQueryController RPCs from the githubcredential/v1 proto package.
//
// Five tools are exposed:
//   - configure_github_webhook:      set up a webhook on a GitHub repository
//   - remove_github_webhook:         remove a webhook from a GitHub repository
//   - get_github_installation_info:  retrieve GitHub App installation details
//   - list_github_repositories:      search repositories accessible via a GitHub credential
//   - get_github_installation_token: obtain a short-lived installation token
package github
