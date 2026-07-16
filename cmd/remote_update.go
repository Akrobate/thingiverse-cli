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

		accessToken, err := cmd.Flags().GetString("access_token")
		if err != nil {
			fmt.Println("failed to retrieve access_token flag: %w", err)
		} else {
			fmt.Println("Access token : " + accessToken)
		}

		thing, err := thing.NewThingParams()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		thing.Save()

		return nil
	},
}

func init() {
	remoteUpdateCmd.Flags().String("access_token", "", "Spécifie directement l'access token de Thingiverse")
	remoteCmd.AddCommand(remoteUpdateCmd)
}
