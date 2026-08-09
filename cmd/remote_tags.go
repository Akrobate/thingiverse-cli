package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var remoteTagsCmd = &cobra.Command{
	Use:   "tags [searchString]",
	Short: "search tags",
	Long: `lists tags with count

Examples:
  thingiverse-cli remote tags search_this
  thingiverse-cli remote tags search_this --access_token=YOUR_ACCESS_TOKEN
  `,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		accessToken, err := getAccessToken(cmd)
		if err != nil {
			return fmt.Errorf("failed to retrieve access_token: %w", err)
		}

		searchString := ""
		if len(args) > 0 {
			searchString = args[0]
		}

		resp, err := thing.TagSearch(searchString, accessToken)
		if err != nil {
			return fmt.Errorf("Failed to Tags: %w", err)
		}

		fmt.Printf("%s\t%s\n", "Tag count", "Tag")
		fmt.Println("---")
		for _, item := range *resp {
			fmt.Printf("%d\t%s\n", item.TagCount, item.Tag)
		}

		return nil
	},
}

func init() {
	remoteTagsCmd.Flags().String("access_token", "", "Access token for thingiverse")
	remoteCmd.AddCommand(remoteTagsCmd)
}
