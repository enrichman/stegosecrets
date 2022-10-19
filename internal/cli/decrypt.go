package cli

import (
	"github.com/enrichman/stegosecrets/internal/decrypt"
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
		return err
	}

	encryptedBytes, err := getInputFromFileOrStdin(encryptedFile)
	if err != nil {
		return err
	}

	return decrypter.Decrypt(encryptedBytes)
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

	return decrypt.NewDecrypter(decrypterOpts...)
}
