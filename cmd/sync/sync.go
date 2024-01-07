package sync

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// Syncer interface that allows you to sync your github stars with a notion database
type Syncer interface {
	SyncStars(ctx context.Context, databaseID string) error
}

// SyncerInitializer function provides a way to initialize the syncer with the given options
// This abstraction is useful to allow mocking the syncer in tests
type SyncerInitializer func(opts Flags) (Syncer, error)

var ErrSyncerInitializerRequired = errors.New("syncer initializer is required")

// NewCommand returns a new cobra command that allows you to sync your github stars with a notion database
func NewCommand(initializerFn SyncerInitializer) *cobra.Command {
	command := &cobra.Command{
		Use:   "sync",
		Short: "Sync your github stars with a notion database",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateRequiredFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, initializerFn)
		},
	}

	command.Flags().StringP(FlagGitHubToken, "", os.Getenv("GITHUB_TOKEN"), "A github token to authenticate with the github api")
	command.Flags().StringP(FlagNotionToken, "", os.Getenv("NOTION_TOKEN"), "A notion token to authenticate with the notion api")
	command.Flags().StringP(FlagNotionDatabaseID, "", os.Getenv("NOTION_DATABASE_ID"), "The id of the notion database to sync with")

	return command
}

// run executes the sync command
func run(cmd *cobra.Command, syncerInitializer SyncerInitializer) error {
	if syncerInitializer == nil {
		return ErrSyncerInitializerRequired
	}

	ctx := cmd.Context()
	flags, err := parseFlags(cmd.Flags())

	if err != nil {
		return err
	}

	syncerSvc, err := syncerInitializer(flags)

	if err != nil {
		return err
	}

	return syncerSvc.SyncStars(ctx, flags.NotionDatabaseID)
}
