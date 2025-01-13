package resource_functions

import (
	"fmt"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/k8s"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/llm"
	"github.com/sashabaranov/go-openai"
)

type FunctionCall struct {
	functions map[string]func(string) (string, error)
	Tools     []openai.Tool
	ai        *llm.OpenAI
	client    *k8s.ClientGo
}

func NewFunctionCall(ai *llm.OpenAI, client *k8s.ClientGo) *FunctionCall {
	instance := &FunctionCall{
		functions: make(map[string]func(string) (string, error)),
		Tools:     make([]openai.Tool, 0),
	}
	return instance
}

func (c *FunctionCall) Init() error {
	funcs := make([]K8sFunc, 0)
	funcs = append(funcs, NewDeployResource(c.ai, c.client))
	funcs = append(funcs, NewQueryResource(c.client))
	funcs = append(funcs, NewDeleteResource(c.client))

	for _, sFunc := range funcs {
		c.functions[sFunc.GetKey()] = sFunc.ResourceFunc
		c.Tools = append(c.Tools, sFunc.GetTool())
	}
	return nil
}

func (c *FunctionCall) CallFunc(name, arguments string) (string, error) {
	f, ok := c.functions[name]
	if !ok {
		return "", fmt.Errorf("function %s is not supported", name)
	}
	return f(arguments)
}
