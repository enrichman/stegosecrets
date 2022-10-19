package encrypt

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/enrichman/stegosecrets/pkg/image"
	sss "github.com/enrichman/stegosecrets/pkg/stego"
)

type Encrypter struct {
	Parts     int
	Threshold int
}

type OptFunc func(*Encrypter) error

func NewEncrypter(opts ...OptFunc) (*Encrypter, error) {
	encrypter := &Encrypter{}

	for _, opt := range opts {
		err := opt(encrypter)
		if err != nil {
			return nil, err
		}
	}

	return encrypter, nil
}

func WithParts(parts int) OptFunc {
	return func(e *Encrypter) error {
		if parts < 0 || parts > 256 {
			return errors.New("invalid parts")
		}
		e.Parts = parts
		return nil
	}
}

func WithThreshold(threshold int) OptFunc {
	return func(e *Encrypter) error {
		if threshold < 0 || threshold > 256 {
			return errors.New("invalid threshold")
		}
		e.Threshold = threshold
		return nil
	}
}

func EncryptFile(filename string) {

}

func (e *Encrypter) Encrypt(reader io.Reader) error {
	masterKey, err := sss.GenerateMasterKey()
	if err != nil {
		return err
	}
	file.WriteKey(masterKey, "out/file.aes")

	message, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	encryptedMessage, err := sss.Encrypt(masterKey, message)
	if err != nil {
		return err
	}

	file.WriteFile(encryptedMessage, "out/file.aes")
	file.WriteChecksum(encryptedMessage, "out/file.aes")

	if e.Parts > 1 {
		parts, err := sss.Split(masterKey, e.Parts, e.Threshold)
		if err != nil {
			return err
		}
		return encodePartsInImages(parts)
	}

	return nil
}

func encodePartsInImages(parts []sss.Part) error {
	dir := "images"
	// get images
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	path, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	images := []fs.DirEntry{}
	for _, f := range files {
		switch filepath.Ext(f.Name()) {
		case ".jpg", ".jpeg", ".png":
			images = append(images, f)
			fmt.Printf(" -  %s\n", f.Name())
		}
	}
	fmt.Printf("found %d images\n", len(images))

	// TODO if parts > len(images) add same image

	fmt.Println("Encrypted parts:")

	for i, part := range parts {
		fmt.Printf(" %d) %s\n", i, part.Base64())

		outName := fmt.Sprintf("out/%d", i)
		if len(images) > 0 {
			imagePath := filepath.Join(path, images[i].Name())
			outName = fmt.Sprintf("%s%s", outName, filepath.Ext(images[i].Name()))
			image.EncodeSecretFromFile(part.Bytes(), imagePath, outName)
			file.WriteFileChecksum(outName)
		}

		file.WriteKey(part.Bytes(), outName)
	}

	return nil
}
