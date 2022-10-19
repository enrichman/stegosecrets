package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/enrichman/stegosecrets/internal/decrypt"
	"github.com/enrichman/stegosecrets/pkg/file"
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
		return err
	}

	// _, err := sss.Combine(parts)
	// if err != nil {
	// 	panic(err)
	// }

	var encryptedBytes []byte
	if encryptedFile != "" {
		encryptedBytes, err = getEncryptedInputFromFile(encryptedFile)
		if err != nil {
			return err
		}
	} else {
		encryptedBytes, err = getEncryptedInputFromStdin()
		if err != nil {
			return err
		}
	}

	return decrypter.Decrypt(encryptedBytes)
}

func getEncryptedInputFromFile(filename string) ([]byte, error) {
	return file.ReadFile(filename)
}

func getEncryptedInputFromStdin() ([]byte, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		// TODO ?
		os.Exit(1)
	}

	input := []byte{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("scanning")
		scanner.Scan()
		// Holds the string that was scanned
		text := scanner.Bytes()
		if len(text) == 0 {
			break
		}
		input = append(input, text...)
	}

	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}

	return input, nil
}
