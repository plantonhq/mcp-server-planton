package kubernetes

import "github.com/modelcontextprotocol/go-sdk/mcp"

func Register(srv *mcp.Server, serverAddress string) {
	mcp.AddTool(srv, GetKubernetesObjectTool(), GetKubernetesObjectHandler(serverAddress))
	mcp.AddTool(srv, FindKubernetesObjectsByKindTool(), FindKubernetesObjectsByKindHandler(serverAddress))
	mcp.AddTool(srv, FindKubernetesNamespacesTool(), FindKubernetesNamespacesHandler(serverAddress))
	mcp.AddTool(srv, FindKubernetesPodsTool(), FindKubernetesPodsHandler(serverAddress))
	mcp.AddTool(srv, GetKubernetesPodTool(), GetKubernetesPodHandler(serverAddress))
	mcp.AddTool(srv, LookupKubernetesSecretKeyValueTool(), LookupKubernetesSecretKeyValueHandler(serverAddress))
	mcp.AddTool(srv, UpdateKubernetesObjectTool(), UpdateKubernetesObjectHandler(serverAddress))
	mcp.AddTool(srv, DeleteKubernetesObjectTool(), DeleteKubernetesObjectHandler(serverAddress))
}
