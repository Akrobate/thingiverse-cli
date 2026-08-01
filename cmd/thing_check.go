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

		accessToken, err := getAccessToken(cmd)
		if err != nil {
			return fmt.Errorf("failed to retrieve access_token: %w", err)
		}

		dryRun, err := cmd.Flags().GetBool("dry_run")
		if err != nil {
			return fmt.Errorf("failed to retrieve dry_run flag: %w", err)
		}

		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return fmt.Errorf("failed to retrieve debug flag: %w", err)
		}

		t, err := thing.NewThing()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		t.Load()
		return t.CompareAndUpdateFiles(accessToken, dryRun, debug)
	},
}

func init() {
	thingCheckCmd.Flags().String("access_token", "", "Access token for thingiverse")
	thingCheckCmd.Flags().Bool("dry_run", false, "Preview operation")
	thingCheckCmd.Flags().Bool("debug", false, "Show debug")
	thingCmd.AddCommand(thingCheckCmd)
}
