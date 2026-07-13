package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "thingiverse-cli",
	Short: "Cli for management of things on thingiverse",
	Long: `Cli for management of things on thingiverse


`,
	Version: "0.0.1",
}

func Execute() error {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd.Execute()
}
