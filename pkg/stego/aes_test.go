package stego_test

import (
	"testing"

	"github.com/enrichman/stegosecrets/pkg/stego"
	"github.com/stretchr/testify/require"
)

func Test_EncryptDecrypt(t *testing.T) {
	key, err := stego.GenerateMasterKey()
	require.NoError(t, err)

	message := []byte("test message")
	encrypted, err := stego.Encrypt(key, message)
	require.NoError(t, err)

	decrypted, err := stego.Decrypt(key, encrypted)
	require.NoError(t, err)
	require.Equal(t, message, decrypted)
}

func GenerateMasterKey() {
	panic("unimplemented")
}

func Test_EncryptDecryptEmptyMessage(t *testing.T) {
	key, err := stego.GenerateMasterKey()
	require.NoError(t, err)

	message := []byte{}
	encrypted, err := stego.Encrypt(key, message)
	require.NoError(t, err)

	decrypted, err := stego.Decrypt(key, encrypted)
	require.NoError(t, err)
	require.Equal(t, message, decrypted)
}

func Benchmark_Encrypt(b *testing.B) {
	key, err := stego.GenerateMasterKey()
	if err != nil {
		b.Fatal(err)
	}

	message := []byte("test message")

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := stego.Encrypt(key, message)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Decrypt(b *testing.B) {
	key, err := stego.GenerateMasterKey()
	if err != nil {
		b.Fatal(err)
	}

	encr, err := stego.Encrypt(key, []byte("test message"))
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := stego.Decrypt(key, encr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func FuzzEncryptDecrypt(f *testing.F) {
	kk, err := stego.GenerateMasterKey()
	require.NoError(f, err)

	f.Add(kk, []byte(`message`))
	f.Fuzz(func(t *testing.T, key []byte, message []byte) {
		encrypted, err := stego.Encrypt(key, message)
		if err != nil {
			require.Nil(t, encrypted)

			return
		}

		require.NoError(t, err)

		decrypted, err := stego.Decrypt(key, encrypted)
		if err != nil {
			require.Nil(t, decrypted)

			return
		}

		require.NoError(t, err)
		require.Equal(t, message, decrypted)
	})
}
