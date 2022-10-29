package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const AppName = "stego"

var (
	Version = "0.0.0-dev"
	verbose bool
	silent  bool
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   AppName,
		Short: "stego",
		Long:  ``,
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "silent output")

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
		return nil, errors.Wrap(err, "failed reading bytes from stdin")
	}

	return text, nil
}
