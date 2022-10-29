package cli

import (
	"bytes"

	"github.com/enrichman/stegosecrets/internal/encrypt"
	"github.com/enrichman/stegosecrets/internal/log"
	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "failed creating encrypter")
	}

	if silent {
		encrypter.Logger = &log.SilentLogger{}
	} else {
		encrypter.Logger = log.NewSimpleLogger(verbose)
	}

	var toEncrypt []byte
	if cleartextFile != "" {
		toEncrypt, err = file.ReadFile(cleartextFile)
	} else {
		cleartextFile = "secret.enc"
		toEncrypt, err = getInputFromStdin()
	}

	if err != nil {
		return errors.Wrapf(err, "failed getting input to encrypt '%s'", cleartextFile)
	}

	err = encrypter.Encrypt(bytes.NewReader(toEncrypt), cleartextFile)
	if err != nil {
		return errors.Wrapf(err, "failed encrypting file '%s'", cleartextFile)
	}

	return nil
}
