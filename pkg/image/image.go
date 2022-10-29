package image

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/jpeg" // enable decoding for jpeg images.
	_ "image/png"  // enable decoding for png images.
	"io"
	"os"

	"github.com/auyer/steganography"
	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/pkg/errors"
)

func EncodeSecretFromFileWithChecksum(secret []byte, inputFile, outputFile string) error {
	err := EncodeSecretFromFile(secret, inputFile, outputFile)
	if err != nil {
		return errors.Wrap(err, "failed encoding secret from file")
	}

	outputImageFile, err := os.Open(outputFile)
	if err != nil {
		return errors.Wrapf(err, "failed opening output file '%s'", outputFile)
	}
	defer outputImageFile.Close()

	h := sha256.New()
	if _, err := io.Copy(h, outputImageFile); err != nil {
		return errors.Wrapf(err, "failed copying hash of output file '%s'", outputFile)
	}

	checksum := fmt.Sprintf("%x\n", h.Sum(nil))
	checksumFilename := fmt.Sprintf("%s.checksum", outputFile)

	err = file.WriteFile([]byte(checksum), checksumFilename)
	if err != nil {
		return errors.Wrapf(err, "failed writing checksum file '%s'", checksumFilename)
	}

	return nil
}

func EncodeSecretFromFile(secret []byte, inputFile, outputFile string) error {
	inputImageFile, err := os.Open(inputFile)
	if err != nil {
		return errors.Wrapf(err, "failed opening input file '%s'", inputFile)
	}
	defer inputImageFile.Close()

	outputImageFile, err := os.Create(outputFile)
	if err != nil {
		return errors.Wrapf(err, "failed creating output file '%s'", outputFile)
	}
	defer outputImageFile.Close()

	return EncodeSecret(secret, inputImageFile, outputImageFile)
}

func EncodeSecret(secret []byte, imgIn io.Reader, imgOut io.Writer) error {
	img, format, err := image.Decode(bufio.NewReader(imgIn))
	if err != nil {
		return errors.Wrapf(err, "failed decoding '%s' image", format)
	}

	// TODO this should be checked with the secret to see if the image is big enough
	// sizeOfMessage := steganography.MaxEncodeSize(img)

	w := new(bytes.Buffer)

	err = steganography.Encode(w, img, secret)
	if err != nil {
		return errors.Wrap(err, "failed encoding secret into image")
	}

	_, err = w.WriteTo(imgOut)

	return errors.Wrap(err, "failed writing out image")
}

func DecodeSecret(imgIn io.Reader) ([]byte, error) {
	img, format, err := image.Decode(bufio.NewReader(imgIn))
	if err != nil {
		return nil, errors.Wrapf(err, "failed decoding '%s' image", format)
	}

	sizeOfMessage := steganography.GetMessageSizeFromImage(img)
	msg := steganography.Decode(sizeOfMessage, img)

	return msg, nil
}
