# GitHub Branch Protection Configuration

This document provides step-by-step instructions for configuring branch protection rules for the `main` branch.

## Prerequisites

- Repository administrator access
- Go to: https://github.com/plantoncloud-inc/mcp-server-planton/settings/branches

## Configuration Steps

### Step 1: Create Branch Protection Rule

1. Navigate to **Settings** → **Branches**
2. Click **Add rule** or **Add branch protection rule**
3. Enter branch name pattern: `main`

### Step 2: Configure Protection Settings

#### Pull Request Requirements

- ☑️ **Require a pull request before merging**
  - ☑️ **Require approvals**: Set to `1`
  - ☑️ **Dismiss stale pull request approvals when new commits are pushed**
  - ☑️ **Require review from Code Owners**

#### Status Checks

- ☑️ **Require status checks to pass before merging**
  - ☑️ **Require branches to be up to date before merging**
  - **Required status checks** (select these after they appear from CI runs):
    - `lint-and-test`
    - `golangci-lint`

#### Conversation Resolution

- ☑️ **Require conversation resolution before merging**

#### Commit Signing (Recommended)

- ☑️ **Require signed commits**

#### Commit History

- ☑️ **Require linear history**

#### Administrator Enforcement

- ☑️ **Do not allow bypassing the above settings**
  - ☑️ **Include administrators**

#### Push Restrictions

- ☑️ **Restrict who can push to matching branches**
  - Add users/teams: `@sureshattaluri`, `@swarupdonepudi`
  - Or use GitHub teams for better management

#### Force Push Settings

- ☑️ **Allow force pushes**: `Specify who can force push`
  - Leave empty to disallow all force pushes (recommended)

#### Deletion Protection

- ☐ **Allow deletions** (leave unchecked)

### Step 3: Save Configuration

Click **Create** or **Save changes** to apply the branch protection rules.

## Verification

After configuration, verify the following:

1. Try to push directly to `main` - should be blocked
2. Create a test PR - should require approval and passing CI
3. Check that status checks are enforced

## Important Notes

- **Status checks** (`lint-and-test`, `golangci-lint`) will only appear in the list after they have run at least once. If you don't see them immediately, merge a PR first, then update the branch protection settings.
- Administrators can optionally bypass these rules, but it's recommended to include them in the enforcement.
- Signed commits provide additional security but require all contributors to set up GPG/SSH signing.

## Troubleshooting

### Status checks not appearing
- Run a PR to trigger the CI workflows
- Wait for the workflows to complete
- Return to branch protection settings and select the checks

### Cannot push to main
- This is expected behavior after protection is enabled
- All changes must go through pull requests

### Need to make emergency changes
- If absolutely necessary, temporarily disable branch protection
- Make the change
- Re-enable protection immediately

## References

- [GitHub Branch Protection Documentation](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)
- [Requiring status checks](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches#require-status-checks-before-merging)
