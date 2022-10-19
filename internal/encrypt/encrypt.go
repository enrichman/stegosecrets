package encrypt

import (
	"errors"
	"fmt"
	"io"
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

func (e *Encrypter) Encrypt(reader io.Reader, filename string) error {
	masterKey, err := e.generateAndSaveMasterKey(filename)
	if err != nil {
		return err
	}

	err = e.encryptAndSaveMessage(masterKey, reader, filename)
	if err != nil {
		return err
	}

	if e.Parts > 1 {
		err = e.splitAndSaveKey(masterKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Encrypter) generateAndSaveMasterKey(filename string) ([]byte, error) {
	masterKey, err := sss.GenerateMasterKey()
	if err != nil {
		return nil, err
	}

	err = file.WriteKey(masterKey, "out/"+filename)
	if err != nil {
		return nil, err
	}

	return masterKey, nil
}

func (e *Encrypter) encryptAndSaveMessage(masterKey []byte, reader io.Reader, filename string) error {
	message, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = file.WriteChecksum(message, "out/"+filename)
	if err != nil {
		return err
	}

	encryptedMessage, err := sss.Encrypt(masterKey, message)
	if err != nil {
		return err
	}

	err = file.WriteFile(encryptedMessage, fmt.Sprintf("out/%s.enc", filename))
	if err != nil {
		return err
	}

	err = file.WriteChecksum(encryptedMessage, fmt.Sprintf("out/%s.enc", filename))
	if err != nil {
		return err
	}
	return nil
}

func (e *Encrypter) splitAndSaveKey(masterKey []byte) error {
	parts, err := sss.Split(masterKey, e.Parts, e.Threshold)
	if err != nil {
		return err
	}

	images, err := e.getImages(len(parts))
	if err != nil {
		return err
	}

	err = e.saveKeysIntoImages(parts, images)
	if err != nil {
		return err
	}

	return nil
}

func (e *Encrypter) getImages(count int) ([]string, error) {
	dir := "images"

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	images := make([]string, 0, count)

	for _, file := range files {
		switch filepath.Ext(file.Name()) {
		case ".jpg", ".jpeg", ".png":
			images = append(images, filepath.Join(dir, file.Name()))
		}
		if len(images) >= count {
			break
		}
	}

	// TODO we can improve this
	lenImages := len(images)
	for lenImages < count {
		images = append(images, images...)
		lenImages = len(images)
	}

	return images[:count], nil
}

func (e *Encrypter) saveKeysIntoImages(parts []sss.Part, images []string) error {
	for i, part := range parts {
		partialKeyFilename := fmt.Sprintf("out/%d", i+1)

		// write .key file
		err := file.WriteKey(part.Bytes(), partialKeyFilename)
		if err != nil {
			return err
		}

		// if the images are available hide the key inside them
		if len(images) > 0 {
			imageOutName := fmt.Sprintf("%s%s", partialKeyFilename, filepath.Ext(images[i]))

			err := image.EncodeSecretFromFile(part.Bytes(), images[i], imageOutName)
			if err != nil {
				return err
			}

			err = file.WriteFileChecksum(imageOutName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
