package cli

import (
	"bytes"
	"path/filepath"

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
	imagesDir     string
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
	encryptCmd.Flags().StringVarP(&outputDir, "output", "o", "out", "output dir")
	encryptCmd.Flags().StringVarP(&imagesDir, "images", "i", "images", "images dir")

	return encryptCmd
}

func runEncryptCmd(cmd *cobra.Command, args []string) error {
	var logger log.Logger
	if silent {
		logger = &log.SilentLogger{}
	} else {
		logger = log.NewSimpleLogger(cmd.OutOrStdout(), verbose)
	}

	var (
		toEncrypt []byte
		err       error
	)

	if cleartextFile != "" {
		toEncrypt, err = file.ReadFile(cleartextFile)
		cleartextFile = filepath.Base(cleartextFile)
	} else {
		toEncrypt, err = getInputFromStdin(cmd)
		cleartextFile = "secret"
	}

	if err != nil {
		return errors.Wrapf(err, "failed getting input to encrypt '%s'", cleartextFile)
	}

	encrypter, err := encrypt.NewEncrypter(
		encrypt.WithParts(keyParts),
		encrypt.WithThreshold(keyThreshold),
		encrypt.WithOutputDir(outputDir),
		encrypt.WithImagesDir(imagesDir),
	)
	if err != nil {
		return errors.Wrap(err, "failed creating encrypter")
	}

	encrypter.Logger = logger

	err = encrypter.Encrypt(bytes.NewReader(toEncrypt), cleartextFile)
	if err != nil {
		return errors.Wrapf(err, "failed encrypting file '%s'", cleartextFile)
	}

	return nil
}
