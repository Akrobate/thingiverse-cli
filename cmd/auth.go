package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/configuration"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate generating the token",
	Long: `Launch the authentication process.

This command will show the url for processing the OAuth2 procedure
You will be prompter to pastle you app token

Examples:
  thingiverse-cli auth`,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		config, err := configuration.NewConfiguration()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		if !config.ConfigurationExists() {
			return fmt.Errorf("Unknown client_id, Please run config first")
		}
		config.Load()

		fmt.Println("Please enter this URL in your browser, Authenticate, pastle here your access token")
		fmt.Println("")
		fmt.Println(config.GenerateConnectionUrl())
		fmt.Println("")
		access_token := askUser("Access token ?", "")

		if access_token == "" {
			return fmt.Errorf("Error, bad token")
		}

		config.AccessToken = access_token
		config.Save()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
