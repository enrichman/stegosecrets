package file

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func WriteFileChecksum(filename string) {
	fileInput, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fileInput.Close()

	content, err := io.ReadAll(fileInput)
	if err != nil {
		panic(err)
	}

	WriteChecksum(content, filename)
}

func WriteChecksum(content []byte, filename string) {
	h := sha256.New()
	h.Write(content)
	checksum := fmt.Sprintf("%x\t%s", h.Sum(nil), filepath.Base(filename))

	WriteFile([]byte(checksum), fmt.Sprintf("%s.checksum", filename))
}

func WriteKey(key []byte, filename string) {
	base64EncodedKey := base64.StdEncoding.EncodeToString(key)
	WriteFile([]byte(base64EncodedKey), fmt.Sprintf("%s.key", filename))
}

func WriteFile(content []byte, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		panic(err)
	}
}

func ReadFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bb, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return bb
}

func ReadKey(filename string) []byte {
	encodedKey := ReadFile(filename)
	decodedKey, err := base64.StdEncoding.DecodeString(string(encodedKey))
	if err != nil {
		panic(err)
	}
	return decodedKey
}
