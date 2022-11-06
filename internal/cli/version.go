package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the application version",
		Run:   runVersionCmd,
	}
}

func runVersionCmd(_ *cobra.Command, _ []string) {
	fmt.Printf("%s version %s\n", AppName, Version)
}
