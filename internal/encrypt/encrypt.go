package encrypt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/enrichman/stegosecrets/internal/log"
	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/enrichman/stegosecrets/pkg/image"
	sss "github.com/enrichman/stegosecrets/pkg/stego"
	"github.com/pkg/errors"
)

type Encrypter struct {
	Parts     int
	Threshold int
	OutputDir string
	ImagesDir string

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

func WithOutputDir(outputDir string) OptFunc {
	return func(e *Encrypter) error {
		absDir, err := filepath.Abs(outputDir)
		if err != nil {
			return errors.Wrap(err, "error getting absolute path for output directory")
		}

		if err := os.MkdirAll(absDir, 0o744); err != nil {
			return errors.Wrap(err, "failed creating output filede")
		}

		e.OutputDir = absDir

		return nil
	}
}

func WithImagesDir(imagesDir string) OptFunc {
	return func(e *Encrypter) error {
		absDir, err := filepath.Abs(imagesDir)
		if err != nil {
			return errors.Wrap(err, "error getting absolute path for images directory")
		}

		e.ImagesDir = absDir

		return nil
	}
}

func (e *Encrypter) Encrypt(reader io.Reader, filename string) error {
	e.Logger.Print(fmt.Sprintf("Encrypting '%s'", filename))

	e.Logger.Debug("generateAndSaveMasterKey")

	masterKey, err := e.generateAndSaveMasterKey(filename)
	if err != nil {
		return errors.Wrapf(err, "failed generating and saving master key '%s'", filename)
	}

	e.Logger.Debug("encryptAndSaveMessage")

	err = e.encryptAndSaveMessage(masterKey, reader, filename)
	if err != nil {
		return errors.Wrapf(err, "failed encrypting and saving message '%s'", filename)
	}

	if e.Parts <= 1 {
		e.Logger.Print("No parts provided. Only the master-key will be generated.")
	}

	if e.Parts > 1 {
		err = e.splitAndSaveKey(masterKey)
		if err != nil {
			return errors.Wrap(err, "failed splitting and saving master key")
		}
	}

	return nil
}

func (e *Encrypter) generateAndSaveMasterKey(filename string) ([]byte, error) {
	masterKey, err := sss.GenerateMasterKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed generating master key")
	}

	err = file.WriteKey(masterKey, filepath.Join(e.OutputDir, fmt.Sprintf("%s.enc", filename)))
	if err != nil {
		return nil, errors.Wrap(err, "failed writing key file")
	}

	return masterKey, nil
}

func (e *Encrypter) encryptAndSaveMessage(masterKey []byte, reader io.Reader, filename string) error {
	message, err := io.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "failed reading message")
	}

	err = file.WriteChecksum(message, filepath.Join(e.OutputDir, filename))
	if err != nil {
		return errors.Wrap(err, "failed writing checksum file of original message")
	}

	encryptedMessage, err := sss.Encrypt(masterKey, message)
	if err != nil {
		return errors.Wrap(err, "failed encrypting message")
	}

	encryptedFilename := filepath.Join(e.OutputDir, fmt.Sprintf("%s.enc", filename))

	err = file.WriteFile(encryptedMessage, encryptedFilename)
	if err != nil {
		return errors.Wrap(err, "failed writing encoded file")
	}

	err = file.WriteChecksum(encryptedMessage, encryptedFilename)
	if err != nil {
		return errors.Wrap(err, "failed writing checksum file")
	}

	return nil
}

func (e *Encrypter) splitAndSaveKey(masterKey []byte) error {
	e.Logger.Print(fmt.Sprintf("Splitting key into %d parts (threshold: %d)", e.Parts, e.Threshold))

	parts, err := sss.Split(masterKey, e.Parts, e.Threshold)
	if err != nil {
		return errors.Wrap(err, "failed splitting masterkey")
	}

	images, err := e.getImages(len(parts))
	if err != nil {
		e.Logger.Print("failed getting images")
	}

	if len(images) == 0 {
		e.Logger.Print("No images found.")
	}

	err = e.saveKeysIntoImages(parts, images)
	if err != nil {
		return errors.Wrap(err, "failed saving keys into images")
	}

	return nil
}

func (e *Encrypter) getImages(count int) ([]string, error) {
	files, err := os.ReadDir(e.ImagesDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading images folder '%s'", e.ImagesDir)
	}

	images := make([]string, 0, count)

	for _, file := range files {
		switch filepath.Ext(file.Name()) {
		case ".jpg", ".jpeg", ".png":
			images = append(images, filepath.Join(e.ImagesDir, file.Name()))
		}

		if len(images) >= count {
			break
		}
	}

	// TODO we can improve this
	lenImages := len(images)
	if lenImages == 0 {
		return nil, errors.Errorf("no image files in %s dir: run 'stego images' to get some random pics", e.ImagesDir)
	}

	for lenImages < count {
		images = append(images, images...)
		lenImages = len(images)
	}

	return images[:count], nil
}

func (e *Encrypter) saveKeysIntoImages(parts []sss.Part, images []string) error {
	for i, part := range parts {
		partialKeyFilename := filepath.Join(e.OutputDir, fmt.Sprintf("%03d", i+1))

		// write .key file
		err := file.WriteKey(part.Bytes(), partialKeyFilename)
		if err != nil {
			return errors.Wrapf(err, "failed writing key file '%s'", partialKeyFilename)
		}

		// if the images are available hide the key inside them
		if len(images) > 0 {
			imageOutName := partialKeyFilename + filepath.Ext(images[i])

			err := image.EncodeSecretFromFile(part.Bytes(), images[i], imageOutName)
			if err != nil {
				return errors.Wrapf(err, "failed encoding secret into image file '%s'", imageOutName)
			}

			err = file.WriteFileChecksum(imageOutName)
			if err != nil {
				return errors.Wrapf(err, "failed writing checksum file '%s'", imageOutName)
			}
		}
	}

	return nil
}
