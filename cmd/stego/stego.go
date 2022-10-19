package main

import (
	"fmt"
	"os"

	"github.com/enrichman/stegosecrets/internal/cli"
)

func main() {
	fmt.Print()
	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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
