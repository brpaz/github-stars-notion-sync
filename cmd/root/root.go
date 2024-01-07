package root

import "github.com/spf13/cobra"

// NewCommand creates a new instance of the root command
func NewCommand(version string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "github-stars-notion-sync",
		Short:   "Sync your github stars with a notion database",
		Version: version,
	}

	return rootCmd
}
