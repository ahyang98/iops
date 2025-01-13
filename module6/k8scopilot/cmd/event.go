/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/k8s"
	"github.com/ahyang98/k8scopilot/cmd/utils/thrid-party/llm"
	"github.com/sashabaranov/go-openai"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
)

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("event called")
		eventsAndLogs, err := getPodEventsAndLogs()
		if err != nil {
			return
		}
		// fmt.Printf("event and logs %v", eventsAndLogs)
		result, err := send2ChatGPT(eventsAndLogs)
		if err != nil {
			return
		}
		fmt.Println(result)
	},
}

func send2ChatGPT(podInfo map[string][]string) (string, error) {
	openAI, err := llm.GetOpenAI()
	if err != nil {
		return "", err
	}

	combinedInfo := "找到以下Pod Warning事件及其日志: \n\n"
	for name, info := range podInfo {
		combinedInfo += fmt.Sprintf("Pod 名字: %s\n", name)
		for _, line := range info {
			combinedInfo += line + "\n"
		}
		combinedInfo += "\n"
	}
	fmt.Printf("send2ChatGPT combinedInfo %v", combinedInfo)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "您是一位Kubernetes专家，你要帮助用户诊断多个Pod的问题",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("以下是多个Pod Event事件及其日志: \n%s\n请主要针对Pod Log给出实质性、可操作的建议", combinedInfo),
		},
	}

	resp, err := openAI.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4oMini,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func getPodEventsAndLogs() (map[string][]string, error) {
	client, err := k8s.GetClientGo()
	if err != nil {
		return nil, err
	}
	result := make(map[string][]string)
	events, err := client.ClientSet.CoreV1().Events(corev1.NamespaceDefault).List(context.Background(), metav1.ListOptions{FieldSelector: "type=Warning"})
	if err != nil {
		return nil, err
	}
	for _, event := range events.Items {
		fmt.Printf("getPodEventsAndLogs message %v \n", event.Message)
		if event.InvolvedObject.Kind == "Pod" {
			podEventLog := handlePodEvent(client, event)
			if podEventLog != nil {
				fmt.Printf("getPodEventsAndLogs podEventLog %v \n", podEventLog)
				result[event.InvolvedObject.Name] = podEventLog
			}
		}
	}
	fmt.Printf("getPodEventsAndLogs result %v \n", result)
	return result, nil
}

func handlePodEvent(client *k8s.ClientGo, event corev1.Event) []string {
	result := make([]string, 0)
	podName := event.InvolvedObject.Name
	namespace := event.InvolvedObject.Namespace
	message := event.Message
	logOptions := &corev1.PodLogOptions{}
	request := client.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, logOptions)
	podLogs, err := request.Stream(context.Background())
	if err != nil {
		return nil
	}
	defer podLogs.Close()
	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(podLogs)
	if err != nil {
		return nil
	}
	result = append(result, fmt.Sprintf("Event Message: %s", message))
	result = append(result, fmt.Sprintf("Namespce:%s", namespace))
	result = append(result, fmt.Sprintf("Logs:%s", buf.String()))
	return result
}

func init() {
	analyzeCmd.AddCommand(eventCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// eventCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// eventCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
