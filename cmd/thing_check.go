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

	Verify if the fields are filled and if the remote creation would be ok.
	Check if listed files exists

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

		if err := t.CheckThingParams(); err != nil {
			return fmt.Errorf("Bad param in thingiverse configuration file: %w", err)
		}

		t.CheckFilesExists()

		return nil
	},
}

func init() {
	thingCmd.AddCommand(thingCheckCmd)
}
