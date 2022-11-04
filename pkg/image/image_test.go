package image_test

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	stegoimage "github.com/enrichman/stegosecrets/pkg/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeSecretFromFile(t *testing.T) {
	inJpeg, err := os.CreateTemp(os.TempDir(), "in.jpg")
	assert.Nil(t, err)

	defer os.Remove(filepath.Join(os.TempDir(), "in.jpg"))

	testImage := image.NewRGBA(image.Rect(0, 0, 256, 256))

	var imageBuff bytes.Buffer
	err = jpeg.Encode(&imageBuff, testImage, nil)
	assert.Nil(t, err)

	_, err = inJpeg.Write(imageBuff.Bytes())
	assert.Nil(t, err)

	inJpeg.Close()
	assert.FileExists(t, inJpeg.Name())

	outJpeg := filepath.Join(os.TempDir(), "out.jpg")
	assert.NoFileExists(t, outJpeg)

	secret := []byte("test secret")
	err = stegoimage.EncodeSecretFromFile(secret, inJpeg.Name(), outJpeg)
	assert.Nil(t, err)
	assert.FileExists(t, outJpeg)

	defer os.Remove(filepath.Join(os.TempDir(), "out.jpg"))

	outJpegFile, err := os.Open(outJpeg)
	assert.Nil(t, err)

	decodedSecret, err := stegoimage.DecodeSecret(outJpegFile)
	assert.Nil(t, err)
	assert.Equal(t, secret, decodedSecret)
}

func TestEncodeDecodeSecret(t *testing.T) {
	secret := []byte("test secret")

	testImage := image.NewRGBA(image.Rect(0, 0, 256, 256))

	var imageBuff bytes.Buffer
	err := jpeg.Encode(&imageBuff, testImage, nil)
	require.NoError(t, err)

	var imageOut bytes.Buffer
	err = stegoimage.EncodeSecret(secret, &imageBuff, &imageOut)
	require.NoError(t, err)

	out, err := stegoimage.DecodeSecret(&imageOut)
	require.NoError(t, err)

	require.Equal(t, secret, out)
}
