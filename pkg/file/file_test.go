package file

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WriteFileChecksum(t *testing.T) {
	tmpDir := t.TempDir()
	original := path.Join(tmpDir, "file")

	err := os.WriteFile(original, []byte("content"), 0o600)
	require.NoError(t, err)

	err = WriteFileChecksum(original)
	require.NoError(t, err)

	checksum, err := os.ReadFile(fmt.Sprintf("%s.checksum", original))
	require.NoError(t, err)

	expectedChecksum := []byte("ed7002b439e9ac845f22357d822bac1444730fbdb6016d3ec9432297b9ec9f73\tfile")
	require.Equal(t, expectedChecksum, checksum)
}

func Test_WriteKey(t *testing.T) {
	tmpDir := t.TempDir()
	keyFile := path.Join(tmpDir, "file")

	expectedKey := []byte("test")
	err := WriteKey(expectedKey, keyFile)
	require.NoError(t, err)

	key, err := ReadKey(keyFile + ".key")
	require.NoError(t, err)
	require.Equal(t, expectedKey, key)
}
