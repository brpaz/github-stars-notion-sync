package main

import (
	"os"
	"runtime"

	"github.com/brpaz/github-stars-notion-sync/cmd/root"
	"github.com/brpaz/github-stars-notion-sync/cmd/sync"
	versionCmd "github.com/brpaz/github-stars-notion-sync/cmd/version"
	"github.com/brpaz/github-stars-notion-sync/internal/syncer"
	"github.com/google/go-github/v57/github"
	"github.com/jomei/notionapi"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	gitCommit = "none"
	buildDate = "unknown"
)

// Application entrypoint
func main() {
	rootCmd := root.NewCommand(version)
	registerCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initSyncer(flags sync.Flags) (sync.Syncer, error) {
	gitHubClient := github.
		NewClient(nil).
		WithAuthToken(flags.GitHubToken)

	notionClient := notionapi.NewClient(
		notionapi.Token(flags.NotionToken),
	)

	return syncer.New(gitHubClient, notionClient)
}

func registerCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(sync.NewCommand(initSyncer))
	rootCmd.AddCommand(versionCmd.NewCommand(versionCmd.VersionInfo{
		Version:   version,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
	}))
}
