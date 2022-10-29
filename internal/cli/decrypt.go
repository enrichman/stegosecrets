package cli

import (
	"github.com/enrichman/stegosecrets/internal/decrypt"
	"github.com/enrichman/stegosecrets/pkg/file"
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
		Short: "decrypt",
		Long:  ``,
		RunE:  runDecryptCmd,
	}

	decryptCmd.Flags().StringVarP(&encryptedFile, "file", "f", "", "file")
	decryptCmd.Flags().StringVar(&masterKeyFile, "master-key", "", "masterkey")
	decryptCmd.Flags().StringArrayVar(&keyFiles, "key", []string{}, "keys")
	decryptCmd.Flags().StringArrayVar(&imageFiles, "img", []string{}, "images")

	return decryptCmd
}

func runDecryptCmd(cmd *cobra.Command, args []string) error {
	decrypter, err := buildDecrypter()
	if err != nil {
		return errors.Wrap(err, "failed building decrypter")
	}

	encryptedBytes, err := file.ReadFile(encryptedFile)
	if err != nil {
		return errors.Wrapf(err, "failed reading file '%s'", encryptedFile)
	}

	err = decrypter.Decrypt(encryptedBytes, encryptedFile)
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
