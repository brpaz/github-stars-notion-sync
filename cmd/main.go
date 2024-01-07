package main

import (
	"os"

	"github.com/brpaz/github-stars-notion-sync/cmd/root"
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
