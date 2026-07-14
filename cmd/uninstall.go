package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/configuration"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove configuration file",
	Long: `Remove configuration file

Examples:
  thingiverse-cli uninstall`,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		config, err := configuration.NewConfiguration()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		if askForConfirmation("This will remove your app credentials, and your access token") {
			err = config.RemoveConfigFolder()
			if err != nil {
				fmt.Errorf("failed to remove configuration: %w", err)
			}
			fmt.Println("Configuration removed")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
