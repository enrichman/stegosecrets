package image

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"image"
	"image/jpeg"
	"testing"
)

func Test_EncodeDecodeSecret(t *testing.T) {
	secret := []byte("test secret")

	testImage := image.NewRGBA(image.Rect(0, 0, 256, 256))
	var imageBuff bytes.Buffer
	err := jpeg.Encode(&imageBuff, testImage, nil)
	require.NoError(t, err)

	var imageOut bytes.Buffer
	err = EncodeSecret(secret, &imageBuff, &imageOut)
	require.NoError(t, err)

	out, err := DecodeSecret(&imageOut)
	require.NoError(t, err)

	require.Equal(t, secret, out)
}
