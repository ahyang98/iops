package resource_functions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/k8s"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type QueryResource struct {
	client *k8s.ClientGo
}

func NewQueryResource(client *k8s.ClientGo) K8sFunc {
	return &QueryResource{client: client}
}

func (q *QueryResource) GetKey() string {
	return "queryResource"
}

func (q *QueryResource) GetTool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "queryResource",
			Description: "查询K8s资源",
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
				},
			},
		},
	}
}

func (q *QueryResource) ResourceFunc(arguments string) (string, error) {
	fmt.Println("queryResource")
	params := struct {
		Namespace    string `json:"namespace"`
		ResourceType string `json:"resource_type"`
	}{}
	err := json.Unmarshal([]byte(arguments), &params)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	gvr, err := k8s.GetGVR(params.ResourceType)
	if err != nil {
		return "", err
	}
	list, err := q.client.DynamicClient.Resource(gvr).Namespace(params.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	result := ""
	for _, item := range list.Items {
		result += fmt.Sprintf("found %s:%s\n", params.ResourceType, item.GetName())
	}
	return result, nil
}

func (q *QueryResource) Register() {
	//TODO implement me
	panic("implement me")
}
