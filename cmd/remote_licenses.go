package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var remoteLicensesCmd = &cobra.Command{
	Use:   "licenses",
	Short: "list licenses",
	Long: `lists licenses

Examples:
  thingiverse-cli remote licenses
  `,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Printf("%s\t%s\n", "Value", "Description")
		fmt.Println("---")
		for _, item := range thing.SupportedLicenses {
			fmt.Printf("%s\t %s\n", item.Value, item.Description)
		}
		return nil
	},
}

func init() {
	remoteCmd.AddCommand(remoteLicensesCmd)
}
