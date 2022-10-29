package encrypt_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/enrichman/stegosecrets/internal/encrypt"
	"github.com/stretchr/testify/assert"
)

func TestNewEncrypter_WithPartsThreshold(t *testing.T) {
	type args struct {
		parts     int
		threshold int
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
			name: "invalid threshold",
			args: args{
				threshold: -1,
			},
			wantErr: true,
		},
		{
			name: "invalid parts",
			args: args{
				parts: 300,
			},
			wantErr: true,
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
				encrypt.WithParts(tc.args.parts),
				encrypt.WithThreshold(tc.args.threshold),
			)

			if tc.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, encrypter)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, encrypter)
			}
		})
	}
}

func TestNewEncrypter_WithOutputDir(t *testing.T) {
	// it should create the output dir
	t.Run("non existing dir", func(t *testing.T) {
		dir := filepath.Join(os.TempDir(), "non-existing-dir")
		assert.Nil(t, os.RemoveAll(dir))

		_, err := os.Stat(dir)
		assert.ErrorIs(t, err, fs.ErrNotExist)

		encrypter, err := encrypt.NewEncrypter(
			encrypt.WithOutputDir(dir),
		)

		assert.Nil(t, err)
		assert.NotNil(t, encrypter)

		_, err = os.Stat(dir)
		assert.Nil(t, err)
	})
}

func TestNewEncrypter_WithImagesDir(t *testing.T) {
	t.Run("set images dir", func(t *testing.T) {
		imagesDir := "random"

		encrypter, err := encrypt.NewEncrypter(
			encrypt.WithImagesDir(imagesDir),
		)

		assert.Nil(t, err)
		assert.NotNil(t, encrypter)

		absoluteImagesDir, err := filepath.Abs(imagesDir)
		assert.Nil(t, err)
		assert.Equal(t, absoluteImagesDir, encrypter.ImagesDir)
	})
}
