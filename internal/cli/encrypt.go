package cli

import (
	"bytes"

	"github.com/enrichman/stegosecrets/internal/encrypt"
	"github.com/spf13/cobra"
)

var (
	cleartextFile string
	keyParts      int
	keyThreshold  int
	outputDir     string
)

func newEncryptCmd() *cobra.Command {
	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt",
		Long:  ``,
		RunE:  runEncryptCmd,
	}

	encryptCmd.Flags().StringVarP(&cleartextFile, "file", "f", "", "file")
	encryptCmd.Flags().IntVarP(&keyParts, "parts", "p", 0, "parts")
	encryptCmd.Flags().IntVarP(&keyThreshold, "threshold", "t", 0, "threshold")
	encryptCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output dir")

	return encryptCmd
}

func runEncryptCmd(cmd *cobra.Command, args []string) error {
	encrypter, err := encrypt.NewEncrypter(
		encrypt.WithParts(keyParts),
		encrypt.WithThreshold(keyThreshold),
	)
	if err != nil {
		return err
	}

	toEncrypt, err := getInputFromFileOrStdin(cleartextFile)
	if err != nil {
		return err
	}

	return encrypter.Encrypt(bytes.NewReader(toEncrypt), cleartextFile)
}
