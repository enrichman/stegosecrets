package decrypt

import (
	"github.com/enrichman/stegosecrets/pkg/file"
	sss "github.com/enrichman/stegosecrets/pkg/stego"
)

type Decrypter struct {
	MasterKey []byte
}

func (d *Decrypter) Decrypt(content []byte) {
	cleartext, err := sss.Decrypt(d.MasterKey, content)
	if err != nil {
		panic(err)
	}

	file.WriteFile(cleartext, "out/clear.txt")

	// check checksum
}
