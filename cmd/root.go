package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Akrobate/thingiverse-cli/pkg/configuration"
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

func getAccessToken(cmd *cobra.Command) (string, error) {
	accessToken, err := cmd.Flags().GetString("access_token")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve access_token flag: %w", err)
	}

	if accessToken != "" {
		return accessToken, nil
	}

	config, err := configuration.NewConfiguration()
	if err != nil {
		return "", fmt.Errorf("failed to initialize configuration: %w", err)
	}

	if !config.ConfigurationExists() {
		return "", fmt.Errorf("Unknown client_id, Please run config first")
	}
	config.Load()

	if !config.AccessTokenExists() {
		return "", fmt.Errorf("Access token not in config, run auth")
	}

	return config.AccessToken, nil
}
