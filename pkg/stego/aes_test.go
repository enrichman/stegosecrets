package stego

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_EncryptDecrypt(t *testing.T) {
	key, err := GenerateMasterKey()
	require.NoError(t, err)

	message := []byte("test message")
	encrypted, err := Encrypt(key, message)
	require.NoError(t, err)

	decrypted, err := Decrypt(key, encrypted)
	require.NoError(t, err)
	require.Equal(t, message, decrypted)
}

func Test_EncryptDecryptEmptyMessage(t *testing.T) {
	key, err := GenerateMasterKey()
	require.NoError(t, err)

	message := []byte{}
	encrypted, err := Encrypt(key, message)
	require.NoError(t, err)

	decrypted, err := Decrypt(key, encrypted)
	require.NoError(t, err)
	require.Equal(t, message, decrypted)
}

func Benchmark_Encrypt(b *testing.B) {
	key, err := GenerateMasterKey()
	if err != nil {
		b.Fatal(err)
	}

	message := []byte("test message")

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := Encrypt(key, message)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_Decrypt(b *testing.B) {
	key, err := GenerateMasterKey()
	if err != nil {
		b.Fatal(err)
	}

	encr, err := Encrypt(key, []byte("test message"))
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := Decrypt(key, encr)
		if err != nil {
			b.Fatal(err)
		}
	}
}
