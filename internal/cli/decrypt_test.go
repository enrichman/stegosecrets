package cli_test

import (
	"bytes"
	"testing"

	"github.com/enrichman/stegosecrets/internal/cli"
	"github.com/enrichman/stegosecrets/pkg/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAssetsDir = "../../test/assets/p5t3/"

func TestDecryptCmd_NoInput(t *testing.T) {
	rootCmd := cli.NewRootCmd()
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetErr(&bytes.Buffer{})

	rootCmd.SetArgs([]string{"decrypt"})

	err := rootCmd.Execute()
	assert.Error(t, err)
}

func TestDecryptCmd(t *testing.T) {
	tt := []struct {
		name           string
		args           []string
		wantExecuteErr bool
		wantCheckErr   bool
	}{
		{
			name: "happy path",
			args: []string{
				"--key", testAssetsDir + "001.key",
				"--key", testAssetsDir + "002.key",
				"--key", testAssetsDir + "003.key",
			},
		},
		{
			name: "not enough parts",
			args: []string{
				"--key", testAssetsDir + "001.key",
				"--key", testAssetsDir + "002.key",
			},
			wantExecuteErr: true,
		},
		{
			name: "decode with images",
			args: []string{
				"--img", testAssetsDir + "001.jpg",
				"--img", testAssetsDir + "002.jpg",
				"--img", testAssetsDir + "003.jpg",
			},
		},
		{
			name: "decode with master-key",
			args: []string{
				"--master-key", testAssetsDir + "secret.enc.key",
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rootCmd := cli.NewRootCmd()

			outAndErr := &bytes.Buffer{}
			rootCmd.SetOut(outAndErr)
			rootCmd.SetErr(outAndErr)

			cmdArgs := append([]string{
				"decrypt",
				"-f", testAssetsDir + "secret.enc",
			}, tc.args...)

			rootCmd.SetArgs(cmdArgs)

			err := rootCmd.Execute()
			if tc.wantExecuteErr {
				require.Error(t, err, outAndErr)
			} else {
				require.NoError(t, err, outAndErr)
			}

			err = file.Check(testAssetsDir+"secret", testAssetsDir+"secret.checksum")
			if tc.wantCheckErr {
				require.Error(t, err, outAndErr)
			} else {
				require.NoError(t, err, outAndErr)
			}
		})
	}
}
