package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/enrichman/stegosecrets/internal/tui"
	"github.com/spf13/cobra"
)

const AppName = "stego"

var Version = "0.0.0-dev"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   AppName,
		Short: "stego",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			tui.Run()
		},
	}

	rootCmd.AddCommand(
		newEncryptCmd(),
		newDecryptCmd(),
		newImagesCmd(),
		newVersionCmd(),
	)

	return rootCmd
}

func getInputFromStdin() ([]byte, error) {
	fmt.Print("Enter text: ")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	return text, nil
}
