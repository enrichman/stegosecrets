package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/enrichman/stegosecrets/pkg/file"
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

func runImagesCmd(cmd *cobra.Command, args []string) error {
	checkIfDirectoryExistsAndCreate(output)

	for i := 1; i <= 10; i++ {
		resp, err := http.Get(fmt.Sprintf("https://picsum.photos/%d/%d", width, height))
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
	}

	return nil
}

func checkIfDirectoryExistsAndCreate(folder string) bool {
	if err := os.Mkdir(folder, 0755); os.IsExist(err) {
		return true
	}
	return false
}
