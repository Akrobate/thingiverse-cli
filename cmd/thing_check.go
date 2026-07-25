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
		//return nil

		t.Load()
		fmt.Println("IMAGES")
		images, err := thing.GetImagesAPI(t.Id, accessToken)
		if err != nil {
			return err
		}
		for _, item := range *images {
			fmt.Printf("%d\t %s\n", item.Id, item.Name)
		}

		fmt.Println("FILES")
		files, err := thing.GetFilesAPI(t.Id, accessToken)
		if err != nil {
			return err
		}
		for _, item := range *files {
			fmt.Printf("%d\t %s\t %s\t %d\n", item.Id, item.Name, item.Hash, item.DefaultImage.Id)
		}
		fmt.Println("------------------------------------------")

		t.DeleteAllFilesAndImages(accessToken)
		return nil

		// return t.UploadFiles(accessToken)
	},
}

func init() {
	thingCheckCmd.Flags().String("access_token", "", "Access token for thingiverse")
	thingCmd.AddCommand(thingCheckCmd)
}
