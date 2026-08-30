package cmd

import (
	"fmt"

	"github.com/Akrobate/thingiverse-cli/pkg/thing"
	"github.com/spf13/cobra"
)

var remoteUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update thing on thingiverse",
	Long: `Update thing on thingiverse, images and files, meta data, or all

Examples:
  thingiverse-cli remote update
  thingiverse-cli remote update --access_token=YOUR_ACCESS_TOKEN
  `,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("you must provide exactly one argument: \"info\", \"files\" or \"all\"")
		}

		if args[0] != "files" && args[0] != "info" && args[0] != "all" {
			return fmt.Errorf("invalid argument %q: must be \"info\",  \"files\" or \"all\"", args[0])
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		updateType := args[0]
		accessToken, err := getAccessToken(cmd)

		if err != nil {
			return fmt.Errorf("failed to retrieve access_token: %w", err)
		}
		dryRun, err := cmd.Flags().GetBool("dry_run")
		if err != nil {
			return fmt.Errorf("failed to retrieve dry_run flag: %w", err)
		}

		if dryRun {
			fmt.Println("[MODE] DRY RUN")
		}

		debug, err := cmd.Flags().GetBool("debug")

		if err != nil {
			return fmt.Errorf("failed to retrieve debug flag: %w", err)
		}
		t, err := thing.NewThing()
		if err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		if updateType == "info" || updateType == "all" {
			if err := t.Update(accessToken); err != nil {
				return fmt.Errorf("[ERROR] Update info: %w", err)
			}
			fmt.Println("[OK] Update info success")
		}

		if updateType == "files" || updateType == "all" {
			if err := t.CompareAndUpdateFiles(accessToken, dryRun, debug); err != nil {
				return err
			}

			if err := t.DeleteAndUpdateAllImages(accessToken, dryRun, debug); err != nil {
				return err
			}

			fmt.Println("Applying gallery order...")
			if err := t.UpdateOrderFilesAndImage(accessToken); err != nil {
				return err
			}

			fmt.Println("[OK] Update files, images, order success")
		}

		return nil
	},
}

func init() {
	remoteUpdateCmd.Flags().String("access_token", "", "Access token for thingiverse")
	remoteUpdateCmd.Flags().Bool("dry_run", false, "Preview operation")
	remoteUpdateCmd.Flags().Bool("debug", false, "Show debug")
	remoteCmd.AddCommand(remoteUpdateCmd)
}
