package llm

import (
	"context"
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	Client  *openai.Client
	context context.Context
}

var ai *OpenAI

func GetOpenAI() (*OpenAI, error) {
	var err error = nil
	if ai == nil {
		ai, err = newOpenAIClient()
	}
	return ai, err
}

func newOpenAIClient() (*OpenAI, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("open ai key is not set")
	}
	config := openai.DefaultConfig(apiKey)
	client := openai.NewClientWithConfig(config)
	return &OpenAI{
		Client:  client,
		context: context.Background(),
	}, nil
}

func (o *OpenAI) SendMessage(prompt, content string) (string, error) {
	request := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: prompt,
			},
			{
				Role:    "user",
				Content: content,
			},
		},
	}
	response, err := o.Client.CreateChatCompletion(o.context, request)
	if err != nil {
		return "", err
	}
	if len(response.Choices) == 0 {
		return "", errors.New("no response from openai")
	}
	return response.Choices[0].Message.Content, nil
}
