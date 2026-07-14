package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "thingiverse-cli",
	Short: "Cli for management of things on thingiverse",
	Long: `Cli for management of things on thingiverse


`,
	Version: "0.0.1",
}

func Execute() error {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd.Execute()
}

func askUser(promt string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(promt + ": ")
	if defaultValue != "" {
		fmt.Printf("(" + defaultValue + ") ")
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		input = defaultValue
	}
	return input
}

func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [Y/n]: ", prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		input = strings.TrimSpace(strings.ToLower(input))

		if input == "" || input == "y" || input == "yes" {
			return true
		}

		if input == "n" || input == "no" {
			return false
		}

		fmt.Println("Invalid input, please answer with Yes or No")
	}
}
