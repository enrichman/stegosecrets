package encrypt

import (
	"encoding/base64"
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
	Parts     uint8
	Threshold uint8
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

	if enc.OutputDir == "" {
		if err := WithOutputDir("out")(enc); err != nil {
			return nil, err
		}
	}

	if enc.ImagesDir == "" {
		if err := WithImagesDir("images")(enc); err != nil {
			return nil, err
		}
	}

	return enc, nil
}

func WithPartsAndThreshold(parts, threshold uint8) OptFunc {
	return func(e *Encrypter) error {
		if threshold > parts {
			return errors.Errorf("threshold %d cannot exceed parts %d", threshold, parts)
		}

		e.Parts = parts
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

func WithLogger(logger log.Logger) OptFunc {
	return func(e *Encrypter) error {
		e.Logger = logger

		return nil
	}
}

func (e *Encrypter) Encrypt(reader io.Reader, filename string) error {
	e.Logger.Print(fmt.Sprintf("🔒 Encrypting '%s'", filename))

	masterKey, err := e.generateAndSaveMasterKey(filename)
	if err != nil {
		return errors.Wrapf(err, "failed generating and saving master key '%s'", filename)
	}

	if e.Parts <= 1 {
		e.Logger.Print("No parts provided. Only the master-key will be generated.")
	}

	e.Logger.Debug("Generated master-key:", base64.StdEncoding.EncodeToString(masterKey))

	err = e.encryptAndSaveMessage(masterKey, reader, filename)
	if err != nil {
		return errors.Wrapf(err, "failed encrypting and saving message '%s'", filename)
	}

	if e.Parts > 1 {
		err = e.splitAndSaveKey(masterKey)
		if err != nil {
			return errors.Wrap(err, "failed splitting and saving master key")
		}
	}

	e.Logger.Print("Encrypted files and keys saved to:", e.OutputDir)

	return nil
}

func (e *Encrypter) generateAndSaveMasterKey(filename string) ([]byte, error) {
	masterKey, err := sss.GenerateMasterKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed generating master key")
	}

	encFilename := filepath.Join(e.OutputDir, fmt.Sprintf("%s.enc", filename))

	err = file.WriteKey(e.Logger, masterKey, encFilename)
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

	err = file.WriteChecksum(e.Logger, message, filepath.Join(e.OutputDir, filename))
	if err != nil {
		return errors.Wrap(err, "failed writing checksum file of original message")
	}

	encryptedMessage, err := sss.Encrypt(masterKey, message)
	if err != nil {
		return errors.Wrap(err, "failed encrypting message")
	}

	encryptedFilename := filepath.Join(e.OutputDir, fmt.Sprintf("%s.enc", filename))

	err = file.WriteFile(e.Logger, encryptedMessage, encryptedFilename)
	if err != nil {
		return errors.Wrap(err, "failed writing encoded file")
	}

	err = file.WriteChecksum(e.Logger, encryptedMessage, encryptedFilename)
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

	e.Logger.Debug("Partial keys:")

	for i, p := range parts {
		e.Logger.Debug(fmt.Sprintf("%d) %s", i+1, p.Base64()))
	}

	images, err := e.getImages(len(parts))
	if err != nil {
		e.Logger.Print("failed getting images")
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

		// if we have sufficient amount of images, we're done, early return
		if len(images) == count {
			return images, nil
		}
	}

	if len(images) == 0 {
		return nil, errors.Errorf("no image files in %s dir: run 'stego images' to get some random pics", e.ImagesDir)
	}

	// if we don't have sufficient amount of images, fill up with images we have
	i := 0
	for len(images) < count {
		images = append(images, images[i])
		i++
	}

	return images, nil
}

func (e *Encrypter) saveKeysIntoImages(parts []sss.Part, images []string) error {
	if len(images) == 0 {
		e.Logger.Print("No images found.")
	}

	for i, part := range parts {
		partialKeyFilename := filepath.Join(e.OutputDir, fmt.Sprintf("%03d", i+1))

		e.Logger.Print(fmt.Sprintf("🔑 Writing partial key %03d", i+1))

		// write .key file
		err := file.WriteKey(e.Logger, part.Bytes(), partialKeyFilename)
		if err != nil {
			return errors.Wrapf(err, "failed writing key file '%s'", partialKeyFilename)
		}

		// if the images are available hide the key inside them
		if len(images) > 0 {
			imageOutName := partialKeyFilename + filepath.Ext(images[i])

			e.Logger.Debug(fmt.Sprintf("Writing partial key %03d into image", i+1))

			err := image.EncodeSecretFromFile(part.Bytes(), images[i], imageOutName)
			if err != nil {
				return errors.Wrapf(err, "failed encoding secret into image file '%s'", imageOutName)
			}

			e.Logger.Debug(fmt.Sprintf("Writing partial key %03d checksum", i+1))

			err = file.WriteFileChecksum(e.Logger, imageOutName)
			if err != nil {
				return errors.Wrapf(err, "failed writing checksum file '%s'", imageOutName)
			}
		}
	}

	return nil
}
