// Package kubernetes provides MCP tools for Kubernetes CloudOps operations.
// Tools query and mutate Kubernetes resources (objects, pods, namespaces, secrets)
// through the Planton control plane using either cloud resource or provider connection access.
package kubernetes

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	k8sns "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/kubernetes/v1/namespace"
	k8sobject "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/kubernetes/v1/object"
	k8spod "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/kubernetes/v1/pod"
	k8ssecret "github.com/plantonhq/mcp-server-planton/gen/go/ai/planton/cloudops/provider/kubernetes/v1/secret"
	"github.com/plantonhq/mcp-server-planton/internal/domains"
	cloudopsctx "github.com/plantonhq/mcp-server-planton/internal/domains/cloudops"
	"google.golang.org/grpc"
)

type kubernetesContextInput struct {
	Org               string `json:"org"                          jsonschema:"required,Organization slug."`
	Env               string `json:"env,omitempty"                jsonschema:"Environment slug. Required for cloud_resource access mode."`
	CloudResourceKind string `json:"cloud_resource_kind,omitempty" jsonschema:"Cloud resource kind (PascalCase, e.g. 'KubernetesDeployment'). Use with cloud_resource_slug for cloud resource access mode."`
	CloudResourceSlug string `json:"cloud_resource_slug,omitempty" jsonschema:"Cloud resource slug. Use with cloud_resource_kind for cloud resource access mode."`
	Connection        string `json:"connection,omitempty"        jsonschema:"Provider connection slug for direct access. Mutually exclusive with cloud resource fields."`
}

type GetKubernetesObjectInput struct {
	kubernetesContextInput
	Namespace  string `json:"namespace" jsonschema:"required,Namespace of the object. Empty for cluster-scoped resources."`
	ApiVersion string `json:"api_version" jsonschema:"required,API version of the resource (e.g. apps/v1, v1)."`
	Kind       string `json:"kind"       jsonschema:"required,Kind of the resource (e.g. Deployment, Service)."`
	Name       string `json:"name"      jsonschema:"required,Name of the resource."`
}

func GetKubernetesObjectTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_kubernetes_object",
		Description: "Get a specific Kubernetes object by namespace, apiVersion, kind, and name. " +
			"Returns the full object manifest. Use when you need to inspect a Deployment, ConfigMap, Service, or any other Kubernetes resource.",
	}
}

func GetKubernetesObjectHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetKubernetesObjectInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetKubernetesObjectInput) (*mcp.CallToolResult, any, error) {
		if input.ApiVersion == "" {
			return nil, nil, fmt.Errorf("'api_version' is required")
		}
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}
		if input.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				req := &k8sobject.GetKubernetesObjectRequest{
					Context:    opsCtx,
					Namespace:  input.Namespace,
					ApiVersion: input.ApiVersion,
					Kind:       input.Kind,
					Name:       input.Name,
				}
				resp, err := k8sobject.NewKubernetesObjectQueryControllerClient(conn).Get(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "kubernetes object")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type FindKubernetesObjectsByKindInput struct {
	kubernetesContextInput
	Namespace string `json:"namespace" jsonschema:"required,Namespace to search in."`
	Kind      string `json:"kind"       jsonschema:"required,Kind to search for (e.g. Deployment, ConfigMap)."`
}

func FindKubernetesObjectsByKindTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "find_kubernetes_objects_by_kind",
		Description: "Find all Kubernetes objects of a specific kind in a namespace. " +
			"Use to list Deployments, Pods, ConfigMaps, or any other resource kind in a given namespace.",
	}
}

func FindKubernetesObjectsByKindHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *FindKubernetesObjectsByKindInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *FindKubernetesObjectsByKindInput) (*mcp.CallToolResult, any, error) {
		if input.Namespace == "" {
			return nil, nil, fmt.Errorf("'namespace' is required")
		}
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				req := &k8sobject.FindKubernetesObjectsByKindRequest{
					Context:   opsCtx,
					Namespace: input.Namespace,
					Kind:      input.Kind,
				}
				resp, err := k8sobject.NewKubernetesObjectQueryControllerClient(conn).FindByKind(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "kubernetes objects by kind")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type FindKubernetesNamespacesInput struct {
	kubernetesContextInput
	LabelSelector string `json:"label_selector,omitempty" jsonschema:"Optional Kubernetes label selector to filter namespaces."`
}

func FindKubernetesNamespacesTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "find_kubernetes_namespaces",
		Description: "List namespaces accessible through the context. " +
			"In cloud resource mode returns the single namespace for that resource; in provider connection mode returns all accessible namespaces, optionally filtered by label selector.",
	}
}

func FindKubernetesNamespacesHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *FindKubernetesNamespacesInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *FindKubernetesNamespacesInput) (*mcp.CallToolResult, any, error) {
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				req := &k8sns.FindKubernetesNamespacesRequest{
					Context:       opsCtx,
					LabelSelector: input.LabelSelector,
				}
				resp, err := k8sns.NewKubernetesNamespaceQueryControllerClient(conn).Find(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "kubernetes namespaces")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type FindKubernetesPodsInput struct {
	kubernetesContextInput
	Namespace      string `json:"namespace"                jsonschema:"required,Namespace to search for pods in."`
	PodManager     string `json:"pod_manager,omitempty"    jsonschema:"Optional name of the pod manager (Deployment, StatefulSet) to filter by."`
	PodManagerKind string `json:"pod_manager_kind,omitempty" jsonschema:"Optional kind of the pod manager (e.g. Deployment, StatefulSet)."`
}

func FindKubernetesPodsTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "find_kubernetes_pods",
		Description: "Find pods in a namespace, optionally filtered by manager (Deployment, StatefulSet, etc.). " +
			"Use to list pods for a specific workload or all pods in a namespace.",
	}
}

func FindKubernetesPodsHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *FindKubernetesPodsInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *FindKubernetesPodsInput) (*mcp.CallToolResult, any, error) {
		if input.Namespace == "" {
			return nil, nil, fmt.Errorf("'namespace' is required")
		}
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				req := &k8spod.FindKubernetesPodsRequest{
					Context:        opsCtx,
					Namespace:      input.Namespace,
					PodManager:     input.PodManager,
					PodManagerKind: input.PodManagerKind,
				}
				resp, err := k8spod.NewKubernetesPodQueryControllerClient(conn).Find(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "kubernetes pods")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type GetKubernetesPodInput struct {
	kubernetesContextInput
	Namespace string `json:"namespace" jsonschema:"required,Namespace of the pod."`
	Name      string `json:"name"      jsonschema:"required,Name of the pod."`
}

func GetKubernetesPodTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "get_kubernetes_pod",
		Description: "Get details of a specific pod by namespace and name. " +
			"Returns pod status, conditions, container states, and events. Use when debugging pod issues.",
	}
}

func GetKubernetesPodHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *GetKubernetesPodInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *GetKubernetesPodInput) (*mcp.CallToolResult, any, error) {
		if input.Namespace == "" {
			return nil, nil, fmt.Errorf("'namespace' is required")
		}
		if input.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				req := &k8spod.GetKubernetesPodRequest{
					Context:   opsCtx,
					Namespace: input.Namespace,
					Name:      input.Name,
				}
				resp, err := k8spod.NewKubernetesPodQueryControllerClient(conn).Get(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "kubernetes pod")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type LookupKubernetesSecretKeyValueInput struct {
	kubernetesContextInput
	Namespace  string `json:"namespace"    jsonschema:"required,Namespace of the secret."`
	SecretName string `json:"secret_name" jsonschema:"required,Name of the Kubernetes Secret resource."`
	Key        string `json:"key"          jsonschema:"required,Key within the secret's data map to look up."`
}

func LookupKubernetesSecretKeyValueTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "lookup_kubernetes_secret_key_value",
		Description: "Look up a specific key's value in a Kubernetes Secret. " +
			"Returns the decoded value. Use when you need to read a password, token, or other secret data from a Secret resource.",
	}
}

func LookupKubernetesSecretKeyValueHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *LookupKubernetesSecretKeyValueInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *LookupKubernetesSecretKeyValueInput) (*mcp.CallToolResult, any, error) {
		if input.Namespace == "" {
			return nil, nil, fmt.Errorf("'namespace' is required")
		}
		if input.SecretName == "" {
			return nil, nil, fmt.Errorf("'secret_name' is required")
		}
		if input.Key == "" {
			return nil, nil, fmt.Errorf("'key' is required")
		}
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				req := &k8ssecret.LookupKubernetesSecretKeyValueRequest{
					Context:    opsCtx,
					Namespace:  input.Namespace,
					SecretName: input.SecretName,
					Key:        input.Key,
				}
				resp, err := k8ssecret.NewKubernetesSecretQueryControllerClient(conn).LookupKeyValue(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "kubernetes secret key value")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type UpdateKubernetesObjectInput struct {
	kubernetesContextInput
	Namespace  string `json:"namespace"   jsonschema:"required,Namespace of the object. Empty for cluster-scoped resources."`
	YamlBase64 string `json:"yaml_base64" jsonschema:"required,Base64-encoded YAML manifest of the resource to update."`
}

func UpdateKubernetesObjectTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "update_kubernetes_object",
		Description: "Update a Kubernetes object from a base64-encoded YAML manifest. " +
			"The YAML must contain apiVersion, kind, and metadata.name. Use to apply changes to Deployments, ConfigMaps, or any other resource.",
	}
}

func UpdateKubernetesObjectHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *UpdateKubernetesObjectInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *UpdateKubernetesObjectInput) (*mcp.CallToolResult, any, error) {
		if input.YamlBase64 == "" {
			return nil, nil, fmt.Errorf("'yaml_base64' is required")
		}
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				req := &k8sobject.UpdateKubernetesObjectRequest{
					Context:    opsCtx,
					Namespace:  input.Namespace,
					YamlBase64: input.YamlBase64,
				}
				resp, err := k8sobject.NewKubernetesObjectCommandControllerClient(conn).Update(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "kubernetes object update")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}

type DeleteKubernetesObjectInput struct {
	kubernetesContextInput
	Namespace  string `json:"namespace"   jsonschema:"required,Namespace of the object. Empty for cluster-scoped resources."`
	ApiVersion string `json:"api_version" jsonschema:"required,API version of the resource to delete (e.g. apps/v1, v1)."`
	Kind       string `json:"kind"        jsonschema:"required,Kind of the resource to delete (e.g. Deployment, ConfigMap)."`
	Name       string `json:"name"        jsonschema:"required,Name of the resource to delete."`
}

func DeleteKubernetesObjectTool() *mcp.Tool {
	return &mcp.Tool{
		Name: "delete_kubernetes_object",
		Description: "Delete a specific Kubernetes object by namespace, apiVersion, kind, and name. " +
			"Use with caution; the object and its dependent resources may be removed.",
	}
}

func DeleteKubernetesObjectHandler(serverAddress string) func(context.Context, *mcp.CallToolRequest, *DeleteKubernetesObjectInput) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, _ *mcp.CallToolRequest, input *DeleteKubernetesObjectInput) (*mcp.CallToolResult, any, error) {
		if input.ApiVersion == "" {
			return nil, nil, fmt.Errorf("'api_version' is required")
		}
		if input.Kind == "" {
			return nil, nil, fmt.Errorf("'kind' is required")
		}
		if input.Name == "" {
			return nil, nil, fmt.Errorf("'name' is required")
		}
		opsCtx, err := cloudopsctx.BuildContext(input.Org, input.Env, input.CloudResourceKind, input.CloudResourceSlug, input.Connection)
		if err != nil {
			return nil, nil, err
		}
		text, err := domains.WithConnection(ctx, serverAddress,
			func(ctx context.Context, conn *grpc.ClientConn) (string, error) {
				req := &k8sobject.DeleteKubernetesObjectRequest{
					Context:    opsCtx,
					Namespace:  input.Namespace,
					ApiVersion: input.ApiVersion,
					Kind:       input.Kind,
					Name:       input.Name,
				}
				resp, err := k8sobject.NewKubernetesObjectCommandControllerClient(conn).Delete(ctx, req)
				if err != nil {
					return "", domains.RPCError(err, "kubernetes object delete")
				}
				return domains.MarshalJSON(resp)
			})
		if err != nil {
			return nil, nil, err
		}
		return domains.TextResult(text)
	}
}
