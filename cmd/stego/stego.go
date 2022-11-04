package main

import (
	"fmt"
	"os"

	"github.com/enrichman/stegosecrets/internal/cli"
)

func main() {
	rootCmd := cli.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(rootCmd.ErrOrStderr(), "❌", err.Error())
		os.Exit(1)
	}
}
