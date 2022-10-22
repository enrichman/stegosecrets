package encrypt_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/enrichman/stegosecrets/internal/encrypt"
)

func TestNewEncrypter(t *testing.T) {
	t.Run("new encrypter with valid options", func(t *testing.T) {
		threshold := 3
		parts := 5

		e, err := encrypt.NewEncrypter(encrypt.WithThreshold(threshold), encrypt.WithParts(parts))

		assert.NotNil(t, e)
		assert.Nil(t, err)
	})
	t.Run("new encrypter with invalid threshold", func(t *testing.T) {
		threshold := -1

		e, err := encrypt.NewEncrypter(encrypt.WithThreshold(threshold))

		assert.Nil(t, e)
		assert.EqualError(t, err, "invalid threshold")
	})
	t.Run("new encrypter with invalid parts", func(t *testing.T) {
		parts := 300

		e, err := encrypt.NewEncrypter(encrypt.WithParts(parts))

		assert.Nil(t, e)
		assert.EqualError(t, err, "invalid parts")
	})
	t.Run("new encrypter with invalid threshold and parts pair", func(t *testing.T) {
		threshold := 5
		parts := 3

		e, err := encrypt.NewEncrypter(encrypt.WithParts(parts), encrypt.WithThreshold(threshold))

		assert.Nil(t, e)
		assert.EqualError(t, err, fmt.Sprintf("threshold %d cannot exceed the parts %d", threshold, parts))
	})
}
