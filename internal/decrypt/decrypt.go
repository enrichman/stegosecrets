package decrypt

import (
	"strings"

	"github.com/enrichman/stegosecrets/pkg/file"
	sss "github.com/enrichman/stegosecrets/pkg/stego"
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
			return nil, err
		}
	}

	return decrypter, nil
}

func WithMasterKeyFile(filename string) OptFunc {
	return func(d *Decrypter) error {
		masterKey, err := file.ReadKey(filename)
		if err != nil {
			return err
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
				return err
			}
		}
		return nil
	}
}

func WithPartialKeyFile(filename string) OptFunc {
	return func(d *Decrypter) error {
		partialKey, err := file.ReadKey(filename)
		if err != nil {
			return err
		}

		d.Parts = append(d.Parts, sss.NewPart(partialKey))
		return nil
	}
}

// TODO fix
func WithPartialKeyImageFile(filename string) OptFunc {
	return func(d *Decrypter) error {
		partialKey, err := file.ReadKey(filename)
		if err != nil {
			return err
		}

		d.Parts = append(d.Parts, sss.NewPart(partialKey))
		return nil
	}
}

func (d *Decrypter) Decrypt(content []byte, filename string) error {

	var key []byte
	var err error

	if len(d.MasterKey) > 0 {
		key = d.MasterKey
	} else {
		key, err = sss.Combine(d.Parts)
		if err != nil {
			return err
		}
	}

	cleartext, err := sss.Decrypt(key, content)
	if err != nil {
		return err
	}

	// TODO check checksum
	return file.WriteFile(cleartext, strings.TrimSuffix(filename, ".enc"))
}
