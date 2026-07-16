package cmd

import (
	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "work with remotes",
}

func init() {
	rootCmd.AddCommand(remoteCmd)
}
