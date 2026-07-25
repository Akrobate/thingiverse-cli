package cmd

import (
	"github.com/spf13/cobra"
)

var thingCmd = &cobra.Command{
	Use:   "thing",
	Short: "Work with local declaration of thing",
}

func init() {
	rootCmd.AddCommand(thingCmd)
}
