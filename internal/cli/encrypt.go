package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/enrichman/stegosecrets/internal/encrypt"
	"github.com/spf13/cobra"
)

func newEncryptCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt",
		Long:  ``,
		RunE:  runEncryptCmd,
	}
}

func runEncryptCmd(cmd *cobra.Command, args []string) error {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	// no pipe
	if fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Print("Enter message: ")
	}

	message, _ := bufio.NewReader(os.Stdin).ReadBytes('\n')

	encrypter := encrypt.Encrypter{Parts: 2, Threshold: 2}
	encrypter.Encrypt(bytes.NewReader(message))

	return nil
}
