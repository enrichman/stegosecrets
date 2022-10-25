package cli

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/spf13/cobra"
)

var (
	width  uint16
	height uint16
)

func newImagesCmd() *cobra.Command {
	imagesCmd := &cobra.Command{
		Use:   "images",
		Short: "images",
		Long:  ``,
		RunE:  runImagesCmd,
	}

	imagesCmd.Flags().Uint16Var(&width, "width", 900, "width")
	imagesCmd.Flags().Uint16Var(&height, "height", 600, "height")

	return imagesCmd
}

var client = http.Client{Timeout: 30 * time.Second}

func runImagesCmd(cmd *cobra.Command, args []string) error {
	for i := 1; i <= 10; i++ {
		resp, err := client.Get(fmt.Sprintf("https://picsum.photos/%d/%d", width, height))
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
