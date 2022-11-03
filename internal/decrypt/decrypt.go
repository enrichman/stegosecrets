package decrypt

import (
	"os"
	"strings"

	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/enrichman/stegosecrets/pkg/image"
	sss "github.com/enrichman/stegosecrets/pkg/stego"
	"github.com/pkg/errors"
)

type Decrypter struct {
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

		d.Parts = append(d.Parts, sss.NewPart(partialKey))

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

		d.Parts = append(d.Parts, sss.NewPart(partialKey))

		return nil
	}
}

func (d *Decrypter) Decrypt(content []byte, filename string) error {
	var (
		key []byte
		err error
	)

	if len(d.MasterKey) > 0 {
		key = d.MasterKey
	} else {
		key, err = sss.Combine(d.Parts)
		if err != nil {
			return errors.Wrap(err, "failed combining parts")
		}
	}

	cleartext, err := sss.Decrypt(key, content)
	if err != nil {
		return errors.Wrap(err, "failed decrypting content")
	}

	// TODO check checksum
	err = file.WriteFile(cleartext, strings.TrimSuffix(filename, ".enc"))
	if err != nil {
		return errors.Wrap(err, "failed writing decoded file")
	}

	return nil
}
