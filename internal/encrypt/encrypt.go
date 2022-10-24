package encrypt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/enrichman/stegosecrets/internal/log"
	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/enrichman/stegosecrets/pkg/image"
	sss "github.com/enrichman/stegosecrets/pkg/stego"
)

type Encrypter struct {
	Parts     int
	Threshold int

	Logger log.Logger
}

type OptFunc func(*Encrypter) error

func NewEncrypter(opts ...OptFunc) (*Encrypter, error) {
	enc := &Encrypter{}

	for _, opt := range opts {
		if err := opt(enc); err != nil {
			return nil, err
		}
	}

	if enc.Threshold > enc.Parts {
		return nil, errors.Errorf("threshold %d cannot exceed the parts %d", enc.Threshold, enc.Parts)
	}

	return enc, nil
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

func (e *Encrypter) Encrypt(reader io.Reader, filename string) error {
	e.Logger.Print(fmt.Sprintf("Encrypting '%s'", filename))

	e.Logger.Debug("generateAndSaveMasterKey")
	masterKey, err := e.generateAndSaveMasterKey(filename)
	if err != nil {
		return err
	}

	e.Logger.Debug("encryptAndSaveMessage")
	err = e.encryptAndSaveMessage(masterKey, reader, filename)
	if err != nil {
		return err
	}

	if e.Parts <= 1 {
		e.Logger.Print("No parts provided. Only the master-key will be generated.")
	}

	if e.Parts > 1 {
		err = e.splitAndSaveKey(masterKey)
		if err != nil {
			return err
		}
	}

	return nil
}

const outDirName = "out"

func (e *Encrypter) generateAndSaveMasterKey(filename string) ([]byte, error) {
	masterKey, err := sss.GenerateMasterKey()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outDirName, 0744); err != nil {
		return nil, err
	}

	err = file.WriteKey(masterKey, fmt.Sprintf("%s/%s.enc", outDirName, filename))
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

	err = file.WriteChecksum(message, fmt.Sprintf("%s/%s.enc", outDirName, filename))
	if err != nil {
		return err
	}

	encryptedMessage, err := sss.Encrypt(masterKey, message)
	if err != nil {
		return err
	}

	err = file.WriteFile(encryptedMessage, fmt.Sprintf("%s/%s.enc", outDirName, filename))
	if err != nil {
		return err
	}

	err = file.WriteChecksum(encryptedMessage, fmt.Sprintf("%s/%s.enc", outDirName, filename))
	if err != nil {
		return err
	}
	return nil
}

func (e *Encrypter) splitAndSaveKey(masterKey []byte) error {
	e.Logger.Print(fmt.Sprintf("Splitting key into %d parts (threshold: %d)", e.Parts, e.Threshold))

	parts, err := sss.Split(masterKey, e.Parts, e.Threshold)
	if err != nil {
		return err
	}

	images, err := e.getImages(len(parts))
	if err != nil {
		return err
	}

	if len(images) == 0 {
		e.Logger.Print("No images found.")
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
	if lenImages == 0 {
		return nil, errors.Errorf("no image files in %s dir: run 'stego images' to get some random pics", dir)
	}
	for lenImages < count {
		images = append(images, images...)
		lenImages = len(images)
	}

	return images[:count], nil
}

func (e *Encrypter) saveKeysIntoImages(parts []sss.Part, images []string) error {
	for i, part := range parts {
		partialKeyFilename := fmt.Sprintf("%s/%d", outDirName, i+1)

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
