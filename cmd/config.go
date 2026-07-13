package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Initialize the configuration file",
	Long: `Initializes the configuration file with clien_id and client_secret.

This command will create the configuration file in your app config directory
You will be prompted for information.

Examples:
  thingiverse-cli config`,
	Args: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		// mgr, err := manager.NewManager()
		// if err != nil {
		// 	return fmt.Errorf("failed to initialize manager: %w", err)
		// }

		// dir, err := os.Getwd()
		// base := filepath.Base(dir)

		// var pkg manager.Package
		// pkg.Name = askUser("package name", base)
		// pkg.Version = askUser("version", "1.0.0")
		// pkg.Description = askUser("description", "")
		// pkg.Repository = askUser("repository", "")
		// pkg.Author = askUser("author", "")

		// if err := mgr.Init(&pkg); err != nil {
		// 	return fmt.Errorf("failed to install package: %w", err)
		// }

		client_id := askUser("Client id", "")
		client_secret := askUser("Client secret", "")

		fmt.Println("Client id : " + client_id)
		fmt.Println("Client secret : " + client_secret)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
