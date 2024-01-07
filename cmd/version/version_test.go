package version_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/brpaz/github-stars-notion-sync/cmd/version"
)

func TestNewVersionCmd(t *testing.T) {

	t.Parallel()

	t.Run("Prints version information", func(t *testing.T) {
		mockVersionInfo := version.VersionInfo{
			Version:   "1.0.0",
			GitCommit: "123456",
			BuildDate: "2023-12-10T15:09:00Z",
			GoVersion: "1.21.1",
		}

		versionCmd := version.NewCommand(mockVersionInfo)
		out := bytes.NewBuffer([]byte{})
		versionCmd.SetOut(out)

		err := versionCmd.RunE(versionCmd, []string{})
		require.NoError(t, err)

		expected := "Version: 1.0.0\nGit commit: 123456\nBuild date: 2023-12-10T15:09:00Z\nGo version: 1.21.1\n"

		assert.Equal(t, expected, out.String())
	})
}
