package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
