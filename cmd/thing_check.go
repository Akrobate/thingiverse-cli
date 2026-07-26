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

		t, err := thing.NewThing()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}
		//t.CheckFilesExists()

		t.Load()
		// t.GenerateHashFiles()

		t.CompareAndUpdateFiles(accessToken)
		// for _, item := range *images {
		// 	fmt.Printf("%d \t %s\n", item.Id, item.Name)
		// }

		// t.DeleteAllFilesAndImages(accessToken)
		return nil
	},
}

func init() {
	thingCheckCmd.Flags().String("access_token", "", "Access token for thingiverse")
	thingCmd.AddCommand(thingCheckCmd)
}
