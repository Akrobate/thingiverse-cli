package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var remoteGetCmd = &cobra.Command{
	Use:   "get",
	Short: "thingiverse-cli remote get thingiverse_id",
	Long: `thingiverse-cli remote get thingiverse_id

Examples:
  thingiverse-cli remote get 123456
  thingiverse-cli remote get 123456 --access_token=YOUR_ACCESS_TOKEN
  `,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		accessToken, err := getAccessToken(cmd)
		if err != nil {
			return fmt.Errorf("failed to retrieve access_token: %w", err)
		}

		resp, err := thing.GetApi(args[0], accessToken)
		data, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			fmt.Printf("Erreur d'affichage : %v\n", err)
		} else {
			fmt.Println(string(data))
		}

		return nil
	},
}

func init() {
	remoteGetCmd.Flags().String("access_token", "", "Access token for thingiverse")
	remoteCmd.AddCommand(remoteGetCmd)
}
