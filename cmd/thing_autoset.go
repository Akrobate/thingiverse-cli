package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var thinAutosetCmd = &cobra.Command{
	Use:   "autoset",
	Short: "autoset fills de thingyverse file with images or files",
	Long: `autoset fills de thingyverse file with images or files

Examples:
  thingiverse-cli thing autoset images
  thingiverse-cli thing autoset files
  `,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("you must provide exactly one argument: \"files\" or \"images\"")
		}

		if args[0] != "files" && args[0] != "images" {
			return fmt.Errorf("invalid argument %q: must be \"files\" or \"images\"", args[0])
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		outputType := args[0]
		if outputType != "files" && outputType != "images" {
			return fmt.Errorf("invalid argument %q: must be \"files\" or \"images\"", args[0])
		}

		t, err := thing.NewThing()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		if err := t.Load(); err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		return t.AutosetFilesAndImages(".", outputType)
	},
}

func init() {
	thingCmd.AddCommand(thinAutosetCmd)
}
