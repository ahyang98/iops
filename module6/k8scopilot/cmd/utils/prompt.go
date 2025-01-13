package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AskForConfirmation prompts the user for a yes/no confirmation
func AskForConfirmation(message string) bool {
	fmt.Printf("%s (yes/no): ", message)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "yes" || response == "y"
}
