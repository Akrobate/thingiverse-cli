package main

import (
	"os"

	"github.com/Akrobate/thingiverse-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
