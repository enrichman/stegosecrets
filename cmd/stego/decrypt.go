package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/enrichman/stegosecrets/internal/decrypt"
	"github.com/enrichman/stegosecrets/pkg/file"
	sss "github.com/enrichman/stegosecrets/pkg/stego"
	"github.com/spf13/cobra"
)

var (
	encryptedFile string
	keyFiles      []string
	imageFiles    []string
)

func newDecryptCmd() *cobra.Command {
	decryptCmd := &cobra.Command{
		Use:   "decrypt",
		Short: "decrypt",
		Long:  ``,
		Run:   runDecryptCmd,
	}

	decryptCmd.PersistentFlags().StringVarP(&encryptedFile, "file", "f", "", "file")
	decryptCmd.PersistentFlags().StringArrayVar(&keyFiles, "key", []string{}, "keys")
	decryptCmd.PersistentFlags().StringArrayVar(&imageFiles, "img", []string{}, "images")

	return decryptCmd
}

func runDecryptCmd(cmd *cobra.Command, args []string) {

	input := []byte{}
	if encryptedFile != "" {
		getEncryptedInputFromFile(input)
	} else {
		getEncryptedInputFromStdin(input)
	}

	masterKey := file.ReadKey("out/file.aes.key")

	key1 := file.ReadKey("out/0.input.jpg.key")
	part1 := sss.NewPart(key1)
	key2 := file.ReadKey("out/1.input.png.key")
	part2 := sss.NewPart(key2)

	_, err := sss.Combine([]sss.Part{part1, part2})
	if err != nil {
		panic(err)
	}

	decrypter := decrypt.Decrypter{MasterKey: masterKey}

	encFile, err := os.Open("out/file.aes")
	if err != nil {
		panic(err)
	}
	defer encFile.Close()

	encryptedBytes, err := io.ReadAll(encFile)
	if err != nil {
		panic(err)
	}

	decrypter.Decrypt(encryptedBytes)
}

func getEncryptedInputFromFile(input []byte) {

}

func getEncryptedInputFromStdin(input []byte) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		os.Exit(1)
	}

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

	fmt.Println(input)

	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
}
