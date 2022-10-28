package cli

import (
	"fmt"
	"github.com/enrichman/stegosecrets/pkg/file"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	width  uint16
	height uint16
	output string
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

	return imagesCmd
}

var client = http.Client{Timeout: 30 * time.Second}

func runImagesCmd(cmd *cobra.Command, args []string) error {
	// creates the output folder if it doesn't exists
	err := os.MkdirAll(output, 0755)
	if err != nil {
		return err
	}

	bar := progressbar.Default(10, "Downloading images...")
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

		err = file.WriteFile(bb, fmt.Sprintf("%s/%d.jpg", output, i))
		if err != nil {
			return err
		}
		err = bar.Add(1)
		if err != nil {
			fmt.Println("Could not add to bar: ", err)
		}
	}
	err = bar.Finish()
	if err != nil {
		fmt.Println("Could not add to bar: ", err)
	}
	return nil
}
