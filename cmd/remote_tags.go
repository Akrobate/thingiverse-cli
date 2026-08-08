package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var remoteTagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "list tags for find the category id",
	Long: `lists tags with child categories, each category is provided with count items

Examples:
  thingiverse-cli remote tags
  thingiverse-cli remote tags --access_token=YOUR_ACCESS_TOKEN
  `,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		accessToken, err := getAccessToken(cmd)
		if err != nil {
			return fmt.Errorf("failed to retrieve access_token: %w", err)
		}

		resp, err := thing.CategorySearch(accessToken)
		if err != nil {
			return fmt.Errorf("Failed to Tags: %w", err)
		}

		fmt.Printf("%s\t%s\n", "Id", "Name (count)")
		fmt.Println("---")
		for _, item := range *resp {
			fmt.Printf("%d\t%s (%d)\n", item.Id, item.Name, item.Count)

			subcategory_resp, err := thing.SubCategorySearch(item.Slug, accessToken)
			if err != nil {
				return fmt.Errorf("Failed to CategorySearch: %w", err)
			}

			for _, item := range subcategory_resp.Children {
				fmt.Printf("%d\t - %s (%d)\n", item.Id, item.Name, item.Count)
			}
		}

		return nil
	},
}

func init() {
	remoteTagsCmd.Flags().String("access_token", "", "Access token for thingiverse")
	remoteCmd.AddCommand(remoteTagsCmd)
}
