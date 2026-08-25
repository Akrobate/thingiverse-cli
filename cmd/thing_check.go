package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var thingCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the local file",
	Long: `Check the local file

Examples:
  thingiverse-cli thing check
  `,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		t, err := thing.NewThing()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		if err := t.Load(); err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if err := t.CheckThingParamsoad(); err != nil {
			return fmt.Errorf("Bad param in thingiverse configuration file: %w", err)
		}

		return nil
	},
}

func init() {
	thingCheckCmd.Flags().String("access_token", "", "Access token for thingiverse")
	thingCheckCmd.Flags().Bool("dry_run", false, "Preview operation")
	thingCheckCmd.Flags().Bool("debug", false, "Show debug")
	thingCmd.AddCommand(thingCheckCmd)
}
