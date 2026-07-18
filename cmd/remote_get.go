package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var remoteGetCmd = &cobra.Command{
	Use:   "get",
	Short: "get remote thing",
	Long: `get remote thing

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

		fmt.Println(args)

		resp, err := thing.Get(args[0], accessToken)

		fmt.Println("------ ERROR -------")
		fmt.Println(err)
		fmt.Println("------ /ERROR -------")
		fmt.Println(resp)
		fmt.Printf("%#v\n", resp)
		return nil
	},
}

func init() {
	remoteGetCmd.Flags().String("access_token", "", "Access token for thingiverse")
	remoteCmd.AddCommand(remoteGetCmd)
}
