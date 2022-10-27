package main

import (
	"fmt"
	"os"

	"github.com/enrichman/stegosecrets/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
