package cli

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/spf13/cobra"
)

var (
	width     uint16
	height    uint16
	output    string
	imagesNum uint16
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
	imagesCmd.Flags().StringVarP(&output, "output", "o", "images", "output directory")
	imagesCmd.Flags().Uint16VarP(&imagesNum, "num", "n", 10, "number of images")

	return imagesCmd
}

var client = http.Client{Timeout: 30 * time.Second}

func runImagesCmd(cmd *cobra.Command, args []string) error {
	if imagesNum == 0 {
		return errors.New("number of images must be at least 1")
	}
	// creates the output folder if it doesn't exists
	err := os.MkdirAll(output, 0755)
	if err != nil {
		return err
	}

	for imagesNum != 0 {
		resp, err := client.Get(fmt.Sprintf("https://picsum.photos/%d/%d", width, height))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bb, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = file.WriteFile(bb, fmt.Sprintf("%s/%d.jpg", output, imagesNum))
		if err != nil {
			return err
		}
		imagesNum--
	}

	return nil
}
