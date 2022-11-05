package cli

import (
	"github.com/enrichman/stegosecrets/internal/decrypt"
	"github.com/enrichman/stegosecrets/internal/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	encryptedFile string
	masterKeyFile string
	keyFiles      []string
	imageFiles    []string
)

func newDecryptCmd() *cobra.Command {
	decryptCmd := &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt a file with the provided master-key, keys or images",
		RunE:  runDecryptCmd,
	}

	decryptCmd.Flags().StringVarP(&encryptedFile, "file", "f", "", "The file to decrypt")
	decryptCmd.Flags().StringVar(&masterKeyFile, "master-key", "", `The master-key used to decrypt the file.
If provided keys or images will be ignored`)
	decryptCmd.Flags().StringArrayVar(&keyFiles, "key", []string{}, "The files containing the partial keys")
	decryptCmd.Flags().StringArrayVar(&imageFiles, "img", []string{}, "The image files containing the hidden partial keys")

	return decryptCmd
}

func runDecryptCmd(cmd *cobra.Command, args []string) error {
	if encryptedFile == "" {
		return errors.New("missing file to decrypt. Use -f/--file flag")
	}

	decrypter, err := buildDecrypter()
	if err != nil {
		return errors.Wrap(err, "failed building decrypter")
	}

	loggerLevel := log.NewLevel(silent, verbose)
	decrypter.Logger = log.NewSimpleLogger(cmd.OutOrStdout(), loggerLevel)

	err = decrypter.Decrypt(encryptedFile)
	if err != nil {
		return errors.Wrapf(err, "failed decrypting file '%s'", encryptedFile)
	}

	return nil
}

func buildDecrypter() (*decrypt.Decrypter, error) {
	decrypterOpts := []decrypt.OptFunc{}

	if masterKeyFile != "" {
		decrypterOpts = append(decrypterOpts, decrypt.WithMasterKeyFile(masterKeyFile))
	}

	for _, filename := range keyFiles {
		decrypterOpts = append(decrypterOpts, decrypt.WithPartialKeyFile(filename))
	}

	for _, filename := range imageFiles {
		decrypterOpts = append(decrypterOpts, decrypt.WithPartialKeyImageFile(filename))
	}

	decrypter, err := decrypt.NewDecrypter(decrypterOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating decrypter")
	}

	return decrypter, nil
}
