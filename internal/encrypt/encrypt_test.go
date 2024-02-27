package encrypt_test

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/enrichman/stegosecrets/internal/encrypt"
	"github.com/enrichman/stegosecrets/internal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEncrypter_WithPartsThreshold(t *testing.T) {
	type args struct {
		parts     uint8
		threshold uint8
	}

	tt := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid options",
			args: args{
				parts:     5,
				threshold: 3,
			},
			wantErr: false,
		},
		{
			name: "invalid threshold and parts",
			args: args{
				parts:     3,
				threshold: 5,
			},
			wantErr: true,
		},
		{
			name: "invalid output",
			args: args{
				parts:     3,
				threshold: 5,
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			encrypter, err := encrypt.NewEncrypter(
				encrypt.WithPartsAndThreshold(tc.args.parts, tc.args.threshold),
			)

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, encrypter)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, encrypter)
			}
		})
	}
}

func TestNewEncrypter_WithOutputDir(t *testing.T) {
	// it should create the output dir
	t.Run("non existing dir", func(t *testing.T) {
		dir := filepath.Join(os.TempDir(), "non-existing-dir")
		require.NoError(t, os.RemoveAll(dir))

		_, err := os.Stat(dir)
		require.ErrorIs(t, err, fs.ErrNotExist)

		encrypter, err := encrypt.NewEncrypter(
			encrypt.WithOutputDir(dir),
		)

		require.NoError(t, err)
		assert.NotNil(t, encrypter)

		_, err = os.Stat(dir)
		assert.NoError(t, err)
	})
}

func TestNewEncrypter_WithImagesDir(t *testing.T) {
	t.Run("set images dir", func(t *testing.T) {
		imagesDir := "random"

		encrypter, err := encrypt.NewEncrypter(
			encrypt.WithImagesDir(imagesDir),
		)

		require.NoError(t, err)
		assert.NotNil(t, encrypter)

		absoluteImagesDir, err := filepath.Abs(imagesDir)
		require.NoError(t, err)
		assert.Equal(t, absoluteImagesDir, encrypter.ImagesDir)
	})
}

func TestEncrypt(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "out-*")
	require.NoError(t, err)
	assert.DirExists(t, tmpDir)

	encrypter, err := encrypt.NewEncrypter(
		encrypt.WithPartsAndThreshold(5, 2),
		encrypt.WithOutputDir(tmpDir),
		encrypt.WithImagesDir("../../test/assets/p5t3"),
		encrypt.WithLogger(log.NewSimpleLogger(io.Discard, log.None)),
	)
	require.NoError(t, err)

	err = encrypter.Encrypt(strings.NewReader("hello world!"), "secret")
	require.NoError(t, err)

	assert.DirExists(t, tmpDir)
	assert.FileExists(t, tmpDir+"/secret.enc")
	assert.FileExists(t, tmpDir+"/secret.enc.key")
	assert.FileExists(t, tmpDir+"/secret.checksum")
	assert.FileExists(t, tmpDir+"/secret.enc.checksum")

	for i := 1; i <= 5; i++ {
		assert.FileExists(t, fmt.Sprintf("%s/%03d.png", tmpDir, i))
		assert.FileExists(t, fmt.Sprintf("%s/%03d.png.checksum", tmpDir, i))
		assert.FileExists(t, fmt.Sprintf("%s/%03d.key", tmpDir, i))
	}

	err = os.RemoveAll(tmpDir)
	require.NoError(t, err)
}
