package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
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
		Short: "Download some stock images that can be used to hide the keys",
		RunE:  runImagesCmd,
	}

	imagesCmd.Flags().Uint16Var(&width, "width", 900, "The width of the images")
	imagesCmd.Flags().Uint16Var(&height, "height", 600, "The height of the images")
	imagesCmd.Flags().StringVarP(&output, "output", "o", "images",
		"The output directory where the images will be downloaded")
	imagesCmd.Flags().Uint16VarP(&imagesNum, "num", "n", 10, "The number of images to download")

	return imagesCmd
}

var (
	client                = http.Client{Timeout: 30 * time.Second}
	errInvalidNumOfImages = errors.New("number of images must be at least 1")
)

func runImagesCmd(cmd *cobra.Command, args []string) error {
	if imagesNum == 0 {
		return errInvalidNumOfImages
	}
	// creates the output folder if it doesn't exists
	err := os.MkdirAll(output, 0o755)
	if err != nil {
		return errors.Wrapf(err, "failed creating output images folder '%s'", output)
	}

	bar := progressbar.Default(int64(imagesNum), "Downloading images...")

	for i := 1; i <= int(imagesNum); i++ {
		url := fmt.Sprintf("https://picsum.photos/%d/%d", width, height)

		resp, err := client.Get(url)
		if err != nil {
			return errors.Wrapf(err, "failed http get request to Picsum [%s]", url)
		}
		defer resp.Body.Close()

		bb, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed reading response from Picsum")
		}

		imageFilename := fmt.Sprintf("%s/%03d.jpg", output, i)

		err = file.WriteFile(nil, bb, imageFilename)
		if err != nil {
			return errors.Wrapf(err, "failed writing file '%s'", imageFilename)
		}

		err = bar.Add(1)
		if err != nil {
			fmt.Println("Error adding value to progress bar: ", err)
		}
	}

	err = bar.Finish()
	if err != nil {
		fmt.Println("Error closing progress bar: ", err)
	}

	return nil
}
