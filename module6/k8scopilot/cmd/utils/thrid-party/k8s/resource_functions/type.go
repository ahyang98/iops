package resource_functions

import "github.com/sashabaranov/go-openai"

type K8sFunc interface {
	ResourceFunc(arguments string) (string, error)
	GetKey() string
	GetTool() openai.Tool
}
