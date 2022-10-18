package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/enrichman/stegosecrets/internal/encrypt"
	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/spf13/cobra"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "stego",
		Short: "stego",
		Long:  ``,
	}

	rootCmd.AddCommand(
		newEncryptCmd(),
		newDecryptCmd(),
		newImageCmd(),
	)

	return rootCmd
}

func newEncryptCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "encrypt",
		Short: "encrypt",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fi, err := os.Stdin.Stat()
			if err != nil {
				panic(err)
			}

			if fi.Mode()&os.ModeNamedPipe == 0 {
				fmt.Print("Enter message: ")
			}

			message, _ := bufio.NewReader(os.Stdin).ReadBytes('\n')

			encrypter := encrypt.Encrypter{Parts: 2, Threshold: 2}
			encrypter.Encrypt(bytes.NewReader(message))
		},
	}
}

func bla() {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "20221010_1212_*")
	if err != nil {
		panic(err)
	}
	fmt.Println(tmpDir)
	defer func() {
		os.RemoveAll(tmpDir)
	}()
}

func newImageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "images",
		Short: "images",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			for i := 1; i <= 10; i++ {
				resp, err := http.Get("https://picsum.photos/900/600")
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				bb, err := io.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				file.WriteFile(bb, fmt.Sprintf("images/%d.foobar.jpg", i))
			}
		},
	}
}
