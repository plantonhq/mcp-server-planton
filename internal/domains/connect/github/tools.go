package github

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/commons/apiresource"
	githubconnectionv1 "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/connect/githubconnection/v1"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	"google.golang.org/grpc"
)

// ---------------------------------------------------------------------------
// configure_github_webhook
// ---------------------------------------------------------------------------

type ConfigureWebhookInput struct {
	GithubConnectionID string `json:"github_connection_id" jsonschema:"required,ID of the GithubConnection to use for webhook configuration."`
	RepoName           string `json:"repo_name" jsonschema:"required,Full repository name (e.g. 'owner/repo')."`
}

func ConfigureWebhookTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "configure_github_webhook",
		Description: "Configure a webhook on a GitHub repository using a GitHub connection. " +
			"Returns the webhook ID on success.",
	}
}

func ConfigureWebhookHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ConfigureWebhookInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ConfigureWebhookInput) (*mcp.CallToolResult, any, error) {
		if input.GithubConnectionID == "" {
			return nil, nil, fmt.Errorf("'github_connection_id' is required")
		}
		if input.RepoName == "" {
			return nil, nil, fmt.Errorf("'repo_name' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := githubconnectionv1.NewGithubCommandControllerClient(conn)
				resp, err := client.ConfigureWebhook(ctx, &githubconnectionv1.GithubRepoWebhookRequest{
					GithubConnectionId: input.GithubConnectionID,
					RepoName:           input.RepoName,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("webhook for repo %q", input.RepoName))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// remove_github_webhook
// ---------------------------------------------------------------------------

type RemoveWebhookInput struct {
	GithubConnectionID string `json:"github_connection_id" jsonschema:"required,ID of the GithubConnection used for the webhook."`
	RepoName           string `json:"repo_name" jsonschema:"required,Full repository name (e.g. 'owner/repo')."`
}

func RemoveWebhookTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "remove_github_webhook",
		Description: "Remove a webhook from a GitHub repository. " +
			"Requires the GitHub connection ID and repository name.",
	}
}

func RemoveWebhookHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *RemoveWebhookInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *RemoveWebhookInput) (*mcp.CallToolResult, any, error) {
		if input.GithubConnectionID == "" {
			return nil, nil, fmt.Errorf("'github_connection_id' is required")
		}
		if input.RepoName == "" {
			return nil, nil, fmt.Errorf("'repo_name' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := githubconnectionv1.NewGithubCommandControllerClient(conn)
				_, err := client.RemoveWebhook(ctx, &githubconnectionv1.GithubRepoWebhookRequest{
					GithubConnectionId: input.GithubConnectionID,
					RepoName:           input.RepoName,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("webhook for repo %q", input.RepoName))
				}
				return fmt.Sprintf("Webhook removed from repository %q", input.RepoName), nil
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_github_installation_info
// ---------------------------------------------------------------------------

type GetInstallationInfoInput struct {
	InstallationID int64 `json:"installation_id" jsonschema:"required,GitHub App installation ID."`
}

func GetInstallationInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_github_installation_info",
		Description: "Get information about a GitHub App installation by installation ID. " +
			"Returns account type, account ID, and avatar URL.",
	}
}

func GetInstallationInfoHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInstallationInfoInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInstallationInfoInput) (*mcp.CallToolResult, any, error) {
		if input.InstallationID == 0 {
			return nil, nil, fmt.Errorf("'installation_id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := githubconnectionv1.NewGithubQueryControllerClient(conn)
				resp, err := client.GetInstallationInfo(ctx, &githubconnectionv1.GithubAppInstallationId{
					Value: input.InstallationID,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("GitHub installation %d", input.InstallationID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// list_github_repositories
// ---------------------------------------------------------------------------

type ListRepositoriesInput struct {
	GithubConnectionID string `json:"github_connection_id" jsonschema:"required,ID of the GithubConnection to list repositories for."`
	SearchText         string `json:"search_text,omitempty" jsonschema:"Optional text to filter repositories by name."`
}

func ListRepositoriesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "list_github_repositories",
		Description: "List GitHub repositories accessible via a GitHub connection. " +
			"Optionally filter by search text. Returns repository names and metadata.",
	}
}

func ListRepositoriesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *ListRepositoriesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *ListRepositoriesInput) (*mcp.CallToolResult, any, error) {
		if input.GithubConnectionID == "" {
			return nil, nil, fmt.Errorf("'github_connection_id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := githubconnectionv1.NewGithubQueryControllerClient(conn)
				resp, err := client.FindGithubRepositories(ctx, &githubconnectionv1.FindGithubRepositoriesInput{
					GithubConnectionId: input.GithubConnectionID,
					SearchText:         input.SearchText,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("repositories for connection %q", input.GithubConnectionID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

// ---------------------------------------------------------------------------
// get_github_installation_token
// ---------------------------------------------------------------------------

type GetInstallationTokenInput struct {
	GithubConnectionID string `json:"github_connection_id" jsonschema:"required,ID of the GithubConnection to generate an installation token for."`
}

func GetInstallationTokenTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_github_installation_token",
		Description: "Generate a short-lived GitHub App installation token for a GitHub connection. " +
			"The token can be used for Git operations and API calls. Tokens expire after 1 hour.",
	}
}

func GetInstallationTokenHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetInstallationTokenInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetInstallationTokenInput) (*mcp.CallToolResult, any, error) {
		if input.GithubConnectionID == "" {
			return nil, nil, fmt.Errorf("'github_connection_id' is required")
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				client := githubconnectionv1.NewGithubQueryControllerClient(conn)
				resp, err := client.GetInstallationToken(ctx, &apiresource.ApiResourceId{
					Value: input.GithubConnectionID,
				})
				if err != nil {
					return "", domains.RPCError(err, fmt.Sprintf("installation token for connection %q", input.GithubConnectionID))
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
