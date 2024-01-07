package sync

import (
	"errors"

	"github.com/spf13/pflag"
)

const (
	FlagGitHubToken      = "github-token"
	FlagNotionToken      = "notion-token"
	FlagNotionDatabaseID = "notion-database-id"
)

var (
	ErrGitHubTokenRequired      = errors.New("github-token is required")
	ErrNotionTokenRequired      = errors.New("notion-token is required")
	ErrNotionDatabaseIDRequired = errors.New("notion-database-id is required")
)

// Flags encapsulates all the options that are required to run the sync command
type Flags struct {
	GitHubToken      string
	NotionToken      string
	NotionDatabaseID string
}

// a map of required flags and their respective error.
var requiredFlags = map[string]error{
	FlagGitHubToken:      ErrGitHubTokenRequired,
	FlagNotionToken:      ErrNotionTokenRequired,
	FlagNotionDatabaseID: ErrNotionDatabaseIDRequired,
}

// validateRequiredFlags validates the flags passed to the sync command
func validateRequiredFlags(flags *pflag.FlagSet) error {
	for flagName, flagErr := range requiredFlags {
		flagValue, err := flags.GetString(flagName)
		if err != nil {
			return err
		}

		if flagValue == "" {
			return flagErr
		}
	}
	return nil
}

// parseFlags parses the flags received in the command and construct an "Options" struct with their values
func parseFlags(flags *pflag.FlagSet) (Flags, error) {
	gitHubToken, err := flags.GetString(FlagGitHubToken)
	if err != nil {
		return Flags{}, err
	}

	notionToken, _ := flags.GetString(FlagNotionToken)

	if err != nil {
		return Flags{}, err
	}

	notionDatabaseID, _ := flags.GetString(FlagNotionDatabaseID)

	if err != nil {
		return Flags{}, err
	}

	return Flags{
		GitHubToken:      gitHubToken,
		NotionToken:      notionToken,
		NotionDatabaseID: notionDatabaseID,
	}, nil
}
