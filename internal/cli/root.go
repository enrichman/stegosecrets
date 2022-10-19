package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/enrichman/stegosecrets/pkg/file"
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

func getInputFromFileOrStdin(filename string) ([]byte, error) {
	// if filename is specified try to get input from file
	if filename != "" {
		return file.ReadFile(filename)
	}

	// else read from stdin
	return getInputFromStdin()
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
