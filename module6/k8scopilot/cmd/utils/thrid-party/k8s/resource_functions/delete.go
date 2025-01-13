package resource_functions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ahyang98/k8scopilot/cmd/utils"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/k8s"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteResource struct {
	client *k8s.ClientGo
}

func NewDeleteResource(client *k8s.ClientGo) K8sFunc {
	return &DeleteResource{client: client}
}

func (d *DeleteResource) ResourceFunc(arguments string) (string, error) {
	fmt.Println("deleteResource")
	params := struct {
		Namespace    string `json:"namespace"`
		ResourceType string `json:"resource_type"`
		ResourceName string `json:"resource_name"`
	}{}
	err := json.Unmarshal([]byte(arguments), &params)
	if err != nil {
		return "", err
	}

	namespace := params.Namespace
	if namespace == "" {
		namespace = corev1.NamespaceDefault
	}

	if !utils.AskForConfirmation(fmt.Sprintf("Are you sure you want to delete %s/%s in namespace %s?", params.ResourceType, params.ResourceName, namespace)) {
		fmt.Println("Deletion cancelled.")
		return "Deletion cancelled", nil
	}

	client, err := k8s.GetClientGo()
	if err != nil {
		return "", err
	}
	gvr, err := k8s.GetGVR(params.ResourceType)
	if err != nil {
		return "", err
	}
	err = client.DynamicClient.Resource(gvr).Namespace(namespace).Delete(context.Background(), params.ResourceName, metav1.DeleteOptions{})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("deleted Resource %s/%s in namespace %s.", params.ResourceType, params.ResourceName, namespace), nil
}

func (d *DeleteResource) GetKey() string {
	return "deleteResource"
}

func (d *DeleteResource) GetTool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "deleteResource",
			Description: "删除K8s资源",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"namespace": {
						Type:        jsonschema.String,
						Description: "资源所在命名空间",
					},
					"resource_type": {
						Type:        jsonschema.String,
						Description: "K8s标准资源类型，例如pod、service、deployment",
					},
					"resource_name": {
						Type:        jsonschema.String,
						Description: "删除的资源名称",
					},
				},
			},
		},
	}
}
