package sync_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brpaz/github-stars-notion-sync/cmd/sync"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockSyncer struct {
	mock.Mock
}

func (m *MockSyncer) SyncStars(ctx context.Context, databaseID string) error {
	args := m.Called(ctx, databaseID)
	return args.Error(0)
}

// reset environment variables to make sure the tests run in a clean environment
func resetEnv(t *testing.T) {
	t.Helper()
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("NOTION_TOKEN", "")
	t.Setenv("NOTION_DATABASE_ID", "")
}

func TestNewCommand(t *testing.T) {
	t.Parallel()

	t.Run("instanciates the command", func(t *testing.T) {
		cmd := sync.NewCommand(nil)
		require.IsType(t, &cobra.Command{}, cmd)
	})
}

func TestRun_WithMissingArgs_ReturnsError(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "should return error if github token is not provided",
			args:     []string{"--notion-token", "123", "--notion-database-id", "123"},
			expected: "github-token is required",
		},
		{
			name:     "should return error if notion token is not provided",
			args:     []string{"--github-token", "123", "--notion-database-id", "123"},
			expected: "notion-token is required",
		},
		{
			name:     "should return error if notion database id is not provided",
			args:     []string{"--github-token", "123", "--notion-token", "123"},
			expected: "notion-database-id is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resetEnv(t)
			cmd := sync.NewCommand(nil)
			cmd.SetArgs(tc.args)

			err := cmd.Execute()

			require.Error(t, err)
			require.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestRun_WithValidArgs(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		mockSyncer := &MockSyncer{}
		cmd := sync.NewCommand(func(opts sync.Flags) (sync.Syncer, error) {
			return mockSyncer, nil
		})
		cmd.SetArgs([]string{"--github-token", "123", "--notion-token", "123", "--notion-database-id", "123"})

		mockSyncer.On("SyncStars", context.Background(), "123").Return(nil)
		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		mockSyncer := &MockSyncer{}
		cmd := sync.NewCommand(func(opts sync.Flags) (sync.Syncer, error) {
			return mockSyncer, nil
		})
		cmd.SetArgs([]string{"--github-token", "123", "--notion-token", "123", "--notion-database-id", "123"})

		mockSyncer.On("SyncStars", context.Background(), "123").Return(errors.New("some-error"))
		err := cmd.Execute()

		require.Error(t, err)
	})
}
