package file

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func WriteFileChecksum(filename string) error {
	fileInput, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileInput.Close()

	content, err := io.ReadAll(fileInput)
	if err != nil {
		return err
	}

	return WriteChecksum(content, filename)
}

func WriteChecksum(content []byte, filename string) error {
	h := sha256.New()

	_, err := h.Write(content)
	if err != nil {
		return err
	}

	checksum := fmt.Sprintf("%x\t%s", h.Sum(nil), filepath.Base(filename))

	return WriteFile([]byte(checksum), fmt.Sprintf("%s.checksum", filename))
}

func WriteKey(key []byte, filename string) error {
	base64EncodedKey := base64.StdEncoding.EncodeToString(key)

	return WriteFile([]byte(base64EncodedKey), fmt.Sprintf("%s.key", filename))
}

func WriteFile(content []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)

	return err
}

func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bb, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return bb, nil
}

func ReadKey(filename string) ([]byte, error) {
	encodedKey, err := ReadFile(filename)
	if err != nil {
		return nil, err
	}

	decodedKey, err := base64.StdEncoding.DecodeString(string(encodedKey))
	if err != nil {
		return nil, err
	}

	return decodedKey, nil
}
