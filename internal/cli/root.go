package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "stego",
		Short: "stego",
		Long:  ``,
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
