package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	GitCommit = "none"
	BuildDate = "unknown"
)

// VersionInfo stores information about the application version
type VersionInfo struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
}

// NewCommand creates a new comamnd that prints information about the application version, including the git commit, build date and go runtime version.
func NewCommand(versionInfo VersionInfo) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the application version",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Version: %s\n", versionInfo.Version)
			fmt.Fprintf(out, "Git commit: %s\n", versionInfo.GitCommit)
			fmt.Fprintf(out, "Build date: %s\n", versionInfo.BuildDate)
			fmt.Fprintf(out, "Go version: %s\n", versionInfo.GoVersion)

			return nil
		},
	}
}
