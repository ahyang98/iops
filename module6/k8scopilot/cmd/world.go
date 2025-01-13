/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/ahyang98/k8scopilot/cmd/utils"

	"github.com/spf13/cobra"
)

// worldCmd represents the world command
var worldCmd = &cobra.Command{
	Use:   "world",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig, _ := utils.GetConfig().Get("kubeconfig")
		fmt.Println("kubeconfig", kubeconfig)
		namespace, _ := utils.GetConfig().Get("namespace")
		fmt.Println("namespace", namespace)
		fmt.Println("source", source)
	},
	Deprecated: "this cmd is deprecated.",
}

var source string

func init() {
	helloCmd.AddCommand(worldCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// worldCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// worldCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	worldCmd.Flags().StringVarP(&source, "source", "s", "world", "The source of the message")
}
