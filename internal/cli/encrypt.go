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
	keyParts      uint8
	keyThreshold  uint8
	outputDir     string
	imagesDir     string
)

func newEncryptCmd() *cobra.Command {
	encryptCmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt a file/message",
		RunE:  runEncryptCmd,
	}

	encryptCmd.Flags().StringVarP(&cleartextFile, "file", "f", "",
		`The file to encrypt. If not specified a message from STDIN will be read.`)
	encryptCmd.Flags().Uint8VarP(&keyParts, "parts", "p", 0,
		`The number of parts (partial keys) in which the secret will be splitted.
If empty only the master-key will be generated.`)
	encryptCmd.Flags().Uint8VarP(&keyThreshold, "threshold", "t", 0,
		`The minimum number of parts (partial keys) needed to decrypt the secret`)
	encryptCmd.Flags().StringVarP(&outputDir, "output", "o", "out",
		`The output directory where the encoded secret and keys/images will be saved.`)
	encryptCmd.Flags().StringVarP(&imagesDir, "images", "i", "images",
		`The directory where to look for the images where the partial keys will be hidden.
If empty no images will be generated.`)

	return encryptCmd
}

func runEncryptCmd(cmd *cobra.Command, args []string) error {
	if keyThreshold > keyParts {
		return errors.Errorf("threshold %d cannot exceed the parts %d", keyThreshold, keyParts)
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

	logger := log.NewSimpleLogger(cmd.OutOrStdout(), log.NewLevel(silent, verbose))

	encrypter, err := encrypt.NewEncrypter(
		encrypt.WithPartsAndThreshold(keyParts, keyThreshold),
		encrypt.WithOutputDir(outputDir),
		encrypt.WithImagesDir(imagesDir),
		encrypt.WithLogger(logger),
	)
	if err != nil {
		return errors.Wrap(err, "failed creating encrypter")
	}

	err = encrypter.Encrypt(bytes.NewReader(toEncrypt), cleartextFile)
	if err != nil {
		return errors.Wrapf(err, "failed encrypting file '%s'", cleartextFile)
	}

	return nil
}
