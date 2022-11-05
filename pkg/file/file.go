package file

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/enrichman/stegosecrets/internal/log"
	"github.com/pkg/errors"
)

func WriteFileChecksum(logger log.Logger, filename string) error {
	fileInput, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "failed opening file '%s'", filename)
	}
	defer fileInput.Close()

	content, err := io.ReadAll(fileInput)
	if err != nil {
		return errors.Wrapf(err, "failed reading file content '%s'", filename)
	}

	return WriteChecksum(logger, content, filename)
}

func WriteChecksum(logger log.Logger, content []byte, filename string) error {
	h := sha256.New()

	_, err := h.Write(content)
	if err != nil {
		return errors.Wrapf(err, "failed hashing content of '%s' file", filename)
	}

	checksum := fmt.Sprintf("%x\t%s", h.Sum(nil), filepath.Base(filename))

	return WriteFile(logger, []byte(checksum), fmt.Sprintf("%s.checksum", filename))
}

func WriteKey(logger log.Logger, key []byte, filename string) error {
	base64EncodedKey := base64.StdEncoding.EncodeToString(key)

	return WriteFile(logger, []byte(base64EncodedKey), fmt.Sprintf("%s.key", filename))
}

func WriteFile(logger log.Logger, content []byte, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "failed creating file '%s'", filename)
	}

	_, err = f.Write(content)
	if err != nil {
		return errors.Wrapf(err, "failed writing content to file '%s'", filename)
	}

	err = f.Close()
	if err != nil {
		return errors.Wrapf(err, "failed writing content to file '%s'", filename)
	}

	if logger != nil {
		logger.Debug("Created file:", filename)
	}

	return nil
}

func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed opening file '%s'", filename)
	}
	defer file.Close()

	bb, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading file '%s'", filename)
	}

	return bb, nil
}

func ReadKey(filename string) ([]byte, error) {
	encodedKey, err := ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed reading key file '%s'", filename)
	}

	decodedKey, err := base64.StdEncoding.DecodeString(string(encodedKey))
	if err != nil {
		return nil, errors.Wrapf(err, "failed decoding file '%s' from base64", filename)
	}

	return decodedKey, nil
}

func Check(filename, checksumFilename string) error {
	content, err := ReadFile(filename)
	if err != nil {
		return err
	}

	h := sha256.New()

	_, err = h.Write(content)
	if err != nil {
		return errors.Wrapf(err, "failed hashing content of '%s' file", filename)
	}

	checksumFileContent, err := ReadFile(checksumFilename)
	if err != nil {
		return err
	}

	checksumToVerify := strings.Split(string(checksumFileContent), "\t")[0]
	checksumContent := fmt.Sprintf("%x", h.Sum(nil))

	if checksumToVerify != checksumContent {
		return errors.New("failed checksum verification")
	}

	return nil
}
