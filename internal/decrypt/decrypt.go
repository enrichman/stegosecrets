package decrypt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/enrichman/stegosecrets/internal/log"
	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/enrichman/stegosecrets/pkg/image"
	sss "github.com/enrichman/stegosecrets/pkg/stego"
	"github.com/pkg/errors"
)

type Decrypter struct {
	Logger log.Logger

	MasterKey []byte
	Parts     []sss.Part
}

type OptFunc func(*Decrypter) error

func NewDecrypter(opts ...OptFunc) (*Decrypter, error) {
	decrypter := &Decrypter{
		Parts: []sss.Part{},
	}

	for _, opt := range opts {
		err := opt(decrypter)
		if err != nil {
			return nil, errors.Wrap(err, "failed applying options to decrypter")
		}
	}

	return decrypter, nil
}

func WithMasterKeyFile(filename string) OptFunc {
	return func(d *Decrypter) error {
		masterKey, err := file.ReadKey(filename)
		if err != nil {
			return errors.Wrap(err, "failed reading master key file")
		}

		d.MasterKey = masterKey

		return nil
	}
}

func WithPartialKeyFiles(filenames []string) OptFunc {
	return func(d *Decrypter) error {
		for _, filename := range filenames {
			err := WithPartialKeyFile(filename)(d)
			if err != nil {
				return errors.Wrap(err, "failed reading partial key file")
			}
		}

		return nil
	}
}

func WithPartialKeyFile(filename string) OptFunc {
	return func(d *Decrypter) error {
		partialKey, err := file.ReadKey(filename)
		if err != nil {
			return errors.Wrap(err, "failed reading partial key file")
		}

		part, err := sss.NewPartFromContent(partialKey)
		if err != nil {
			return errors.Wrap(err, "failed creating part")
		}

		d.Parts = append(d.Parts, part)

		return nil
	}
}

func WithPartialKeyImageFile(filename string) OptFunc {
	return func(d *Decrypter) error {
		file, err := os.Open(filename)
		if err != nil {
			return errors.Wrapf(err, "failed opening file '%s'", filename)
		}
		defer file.Close()

		partialKey, err := image.DecodeSecret(file)
		if err != nil {
			return errors.Wrap(err, "failed reading partial key image file")
		}

		part, err := sss.NewPartFromContent(partialKey)
		if err != nil {
			return errors.Wrap(err, "failed creating part")
		}

		d.Parts = append(d.Parts, part)

		return nil
	}
}

func (d *Decrypter) Decrypt(filename string) error {
	d.Logger.Print(fmt.Sprintf("Decrypting '%s'", filepath.Base(filename)))

	if len(d.MasterKey) == 0 && len(d.Parts) < 2 {
		return errors.New("at least a master-key or more than one part needs to be specified")
	}

	encryptedFile, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Wrap(err, "file not found")
		}

		return errors.Wrapf(err, "failed opening file '%s'", filename)
	}

	var key []byte

	if len(d.MasterKey) > 0 {
		d.Logger.Print("Decrypting with master-key")
		key = d.MasterKey
	} else {
		d.Logger.Print("Decrypting with partial keys")
		key, err = sss.Combine(d.Parts)
		if err != nil {
			return errors.Wrap(err, "failed combining parts")
		}
	}

	content, err := io.ReadAll(encryptedFile)
	if err != nil {
		return errors.Wrap(err, "failed to read content")
	}

	cleartext, err := sss.Decrypt(key, content)
	if err != nil {
		return errors.Wrap(err, "failed decrypting content")
	}

	// TODO check checksum
	outputFile := strings.TrimSuffix(filename, ".enc")

	err = file.WriteFile(d.Logger, cleartext, outputFile)
	if err != nil {
		return errors.Wrap(err, "failed writing decoded file")
	}

	d.Logger.Print("Decrypted file saved to:", outputFile)

	return nil
}
