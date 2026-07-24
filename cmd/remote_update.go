package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var remoteUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "------",
	Long: `-------

Examples:
  thingiverse-cli remote update
  thingiverse-cli remote update --access_token=YOUR_ACCESS_TOKEN
  `,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		accessToken, err := getAccessToken(cmd)
		if err != nil {
			return fmt.Errorf("failed to retrieve access_token: %w", err)
		}

		t, err := thing.NewThing()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		if err := t.Update(accessToken); err != nil {
			return err
		}

		fmt.Println("[OK] Update success")

		return nil
	},
}

func init() {
	remoteUpdateCmd.Flags().String("access_token", "", "Access token for thingiverse")
	remoteCmd.AddCommand(remoteUpdateCmd)
}
