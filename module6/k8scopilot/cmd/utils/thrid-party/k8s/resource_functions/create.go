package resource_functions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/k8s"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/llm"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/restmapper"
)

type DeployResource struct {
	ai     *llm.OpenAI
	client *k8s.ClientGo
}

func NewDeployResource(ai *llm.OpenAI, client *k8s.ClientGo) K8sFunc {
	return &DeployResource{ai: ai, client: client}
}

func (d *DeployResource) ResourceFunc(arguments string) (string, error) {
	fmt.Println("generateAndDeployResource")
	params := struct {
		UserInput string `json:"user_input"`
	}{}
	err := json.Unmarshal([]byte(arguments), &params)
	if err != nil {
		return "", fmt.Errorf("failed to parse function call generateAndDeployResource arguments=%s", arguments)
	}
	yaml, err := d.ai.SendMessage("你现在是一个 K8s 资源生成器，根据用户输入生成 K8s YAML，注意除了 YAML 内容以外不要输出任何内容，此外不要把 YAML 放在 ``` 代码快里", params.UserInput)
	if err != nil {
		return "", fmt.Errorf("ChatGPT error: %v", err)
	}
	resources, err := restmapper.GetAPIGroupResources(d.client.DiscoveryClient)
	if err != nil {
		return "", err
	}
	unstructuredObj := &unstructured.Unstructured{}
	_, _, err = scheme.Codecs.UniversalDeserializer().Decode([]byte(yaml), nil, unstructuredObj)
	if err != nil {
		return "", err
	}
	mapper := restmapper.NewDiscoveryRESTMapper(resources)
	groupVersionKind := unstructuredObj.GroupVersionKind()
	mapping, err := mapper.RESTMapping(groupVersionKind.GroupKind(), groupVersionKind.Version)
	if err != nil {
		return "", err
	}
	namespace := unstructuredObj.GetNamespace()
	if namespace == "" {
		namespace = corev1.NamespaceDefault
	}
	_, err = d.client.DynamicClient.Resource(mapping.Resource).Namespace(namespace).Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("YAML:\n%s\n\nDeployment successful", yaml), nil
}

func (d *DeployResource) GetKey() string {
	return "generateAndDeployResource"
}

func (d *DeployResource) GetTool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "generateAndDeployResource",
			Description: "生成K8s YAML并部署资源",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"user_input": {
						Type:        jsonschema.String,
						Description: "用户输入内容要求包含资源类型和镜像",
					},
				},
				Required: []string{"user_input"},
			},
		},
	}
}
