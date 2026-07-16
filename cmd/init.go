package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the configuration file",
	Long: `Initializes the configuration file with clien_id and client_secret.

This command will create the configuration file in your app config directory
You will be prompted for information.

Examples:
  thingiverse-cli config`,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		thing, err := thing.NewThingParams()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		thing.Save()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
