package root_test

import (
	"testing"

	"github.com/brpaz/github-stars-notion-sync/cmd/root"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestNewCommand(t *testing.T) {
	t.Parallel()

	t.Run("instanciates the command", func(t *testing.T) {
		cmd := root.NewCommand("1.0.0")
		require.IsType(t, &cobra.Command{}, cmd)
	})
}
