/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/k8s"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/k8s/resource_functions"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/llm"
	"github.com/sashabaranov/go-openai"
	"os"

	"github.com/spf13/cobra"
)

// chatgptCmd represents the chatgpt command
var chatgptCmd = &cobra.Command{
	Use:   "chatgpt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		startChat()
	},
}

func startChat() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("我是K8s Copilot, 请问有什么我可以帮你的")
	for {
		fmt.Print("> ")
		if scanner.Scan() {
			input := scanner.Text()
			if input == "exit" {
				fmt.Println("再见")
				break
			}

			if input == "" {
				continue
			}
			resp := processInput(input)
			fmt.Println(resp)
		}
	}
}

func processInput(input string) string {
	client, err := llm.GetOpenAI()
	if err != nil {
		return err.Error()
	}
	return functionCalling(input, client)
}

func functionCalling(input string, client *llm.OpenAI) string {
	dialogue := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		},
	}
	if functionCall == nil {
		return "functionCall is nil"
	}
	response, err := client.Client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT4o,
		Messages: dialogue,
		Tools:    functionCall.Tools,
	})
	if err != nil || len(response.Choices) != 1 {
		return fmt.Sprintf("Completion error1: err:%v len(choice):%v\n", err, len(response.Choices))
	}
	message := response.Choices[0].Message
	if len(message.ToolCalls) != 1 {
		return fmt.Sprintf("Completion error2: len(toolcalls): %v\n", len(message.ToolCalls))
	}
	dialogue = append(dialogue, message)
	function := message.ToolCalls[0].Function
	result, err := functionCall.CallFunc(function.Name, function.Arguments)
	if err != nil {
		return err.Error()
	}
	return result
}

func getFunctionCall() (*resource_functions.FunctionCall, error) {
	ai, err := llm.GetOpenAI()
	if err != nil {
		return nil, err
	}

	clientGo, err := k8s.GetClientGo()
	if err != nil {
		return nil, err
	}

	funcCall := resource_functions.NewFunctionCall(ai, clientGo)
	err = funcCall.Init()
	if err != nil {
		return nil, err
	}
	return funcCall, nil
}

func init() {
	askCmd.AddCommand(chatgptCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatgptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatgptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
