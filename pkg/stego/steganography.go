package stego

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"io"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/auyer/steganography"
	"github.com/enrichman/stegosecrets/pkg/file"
)

func EncodeSecretFromFileWithChecksum(secret []byte, inputFile, outputFile string) error {
	err := EncodeSecretFromFile(secret, inputFile, outputFile)
	if err != nil {
		return err
	}

	outputImageFile, err := os.Open(outputFile)
	if err != nil {
		return err
	}
	defer outputImageFile.Close()

	h := sha256.New()
	if _, err := io.Copy(h, outputImageFile); err != nil {
		return err
	}

	checksum := fmt.Sprintf("%x\n", h.Sum(nil))
	file.WriteFile([]byte(checksum), fmt.Sprintf("%s.checksum", outputFile))

	return nil
}

func EncodeSecretFromFile(secret []byte, inputFile, outputFile string) error {
	inputImageFile, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer inputImageFile.Close()

	outputImageFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outputImageFile.Close()

	return EncodeSecret(secret, inputImageFile, outputImageFile)
}

func EncodeSecret(secret []byte, imgIn io.Reader, imgOut io.Writer) error {
	img, format, err := image.Decode(bufio.NewReader(imgIn))
	if err != nil {
		return err
	}
	fmt.Println("format", format)

	sizeOfMessage := steganography.MaxEncodeSize(img)
	fmt.Println("sizeOfMessage", sizeOfMessage)

	w := new(bytes.Buffer)
	err = steganography.Encode(w, img, secret)
	if err != nil {
		return err
	}

	_, err = w.WriteTo(imgOut)
	return err
}

func DecodeSecret(imgIn io.Reader) ([]byte, error) {
	img, format, err := image.Decode(bufio.NewReader(imgIn))
	if err != nil {
		return nil, err
	}
	fmt.Println("format", format)

	sizeOfMessage := steganography.GetMessageSizeFromImage(img)
	fmt.Println("sizeOfMessage", sizeOfMessage)

	msg := steganography.Decode(sizeOfMessage, img)
	return msg, nil
}
