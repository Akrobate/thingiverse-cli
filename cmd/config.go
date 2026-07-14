package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/configuration"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Initialize the configuration file",
	Long: `Initializes the configuration file with clien_id and client_secret.

This command will create the configuration file in your app config directory
You will be prompted for information.

Examples:
  thingiverse-cli config`,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		config, err := configuration.NewConfiguration()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		client_id := askUser("Client id", "")
		client_secret := askUser("Client secret", "")

		config.ClientId = client_id
		config.ClientSecret = client_secret

		err = config.Save()
		if err != nil {
			return fmt.Errorf("Failed to save configuration: %w", err)
		}

		fmt.Println("Credentials saved")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
