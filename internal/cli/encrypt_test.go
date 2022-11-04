package cli_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/enrichman/stegosecrets/internal/cli"
	"github.com/stretchr/testify/assert"
)

func setupTest(tb testing.TB) func(tb testing.TB) {
	tb.Helper()
	tb.Log("setup test")

	return func(tb testing.TB) {
		tb.Helper()
		tb.Log("teardown test")

		assert.DirExists(tb, "out")
		assert.Nil(tb, os.RemoveAll("out"))
	}
}

func TestEncryptCmd_NoInput(t *testing.T) {
	rootCmd := cli.NewRootCmd()
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})

	rootCmd.SetArgs([]string{"encrypt"})

	assert.NoDirExists(t, "out")

	err := rootCmd.Execute()
	assert.NotNil(t, err)
	assert.NoDirExists(t, "out")
}

func TestEncryptCmd_Stdin(t *testing.T) {
	teardown := setupTest(t)
	defer teardown(t)

	rootCmd := cli.NewRootCmd()

	outAndErr := &bytes.Buffer{}
	rootCmd.SetOut(outAndErr)
	rootCmd.SetErr(outAndErr)

	rootCmd.SetIn(strings.NewReader("hello\n"))

	rootCmd.SetArgs([]string{"encrypt"})

	assert.NoDirExists(t, "out")

	err := rootCmd.Execute()
	assert.Nil(t, err, outAndErr)

	assert.DirExists(t, "out")
	assert.FileExists(t, "out/secret.enc")
	assert.FileExists(t, "out/secret.enc.key")
	assert.FileExists(t, "out/secret.checksum")
	assert.FileExists(t, "out/secret.enc.checksum")
}
