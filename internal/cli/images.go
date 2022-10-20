package cli

import (
	"fmt"
	"io"
	"net/http"

	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/spf13/cobra"
)

func newImagesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "images",
		Short: "images",
		Long:  ``,
		RunE:  runImagesCmd,
	}
}

func runImagesCmd(cmd *cobra.Command, args []string) error {
	// TODO add flags for w/h, number of images, folder..
	for i := 1; i <= 10; i++ {
		resp, err := http.Get("https://picsum.photos/900/600")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bb, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = file.WriteFile(bb, fmt.Sprintf("images/%d.jpg", i))
		if err != nil {
			return err
		}
	}

	return nil
}
