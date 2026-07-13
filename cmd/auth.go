package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate generating the token",
	Long: `Launch the authentication process.

This command will show the url for processing the OAuth2 procedure
You will be prompter to pastle you app token

Examples:
  thingiverse-cli auth`,
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
	rootCmd.AddCommand(authCmd)
}
