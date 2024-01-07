package syncer_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/brpaz/github-stars-notion-sync/internal/syncer"
	"github.com/google/go-github/v57/github"
	"github.com/h2non/gock"
	"github.com/jomei/notionapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	notionAPIURL = "https://api.notion.com"
	githubAPIURL = "https://api.github.com"
)

func loadFixture(t *testing.T, path string) []byte {
	t.Helper()

	// Get the absolute path to the project root
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get the caller information")
	}

	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

	jsonFilePath := filepath.Join(projectRoot, "internal", "testdata", path)
	jsonData, err := os.ReadFile(jsonFilePath)

	require.NoError(t, err)
	return jsonData
}

func TestSyncer_New(t *testing.T) {
	t.Parallel()

	t.Run("should return error if github client is nil", func(t *testing.T) {
		t.Parallel()
		syncerSvc, err := syncer.New(nil, notionapi.NewClient(""))

		assert.Error(t, err)
		assert.Equal(t, err, syncer.ErrNilGithubClient)
		assert.Nil(t, syncerSvc)
	})

	t.Run("should return error if notion client is nil", func(t *testing.T) {
		t.Parallel()
		syncerSvc, err := syncer.New(github.NewClient(nil), nil)
		assert.Error(t, err)
		assert.Equal(t, err, syncer.ErrNilNotionClient)
		assert.Nil(t, syncerSvc)
	})

	t.Run("should return syncer instance if all arguments are valid", func(t *testing.T) {
		t.Parallel()
		syncerSvc, err := syncer.New(github.NewClient(nil), &notionapi.Client{})
		assert.NoError(t, err)
		assert.NotNil(t, syncerSvc)
	})
}

func TestSyncer_SyncStars(t *testing.T) {
	syncerSvc, err := syncer.New(github.NewClient(nil), notionapi.NewClient(""))
	require.NoError(t, err)

	t.Run("should return error if notion database does not exist", func(t *testing.T) {
		mockDatabaseID := "7099f06c-95ef-46a2-a47d-672cf8e4a8a4"
		notionAPIResponse := loadFixture(t, path.Join("notionapi", "database_not_found_response.json"))
		defer gock.Off()

		gock.New(notionAPIURL).
			Get(fmt.Sprintf("/v1/databases/%s", mockDatabaseID)).
			Reply(404).
			JSON(notionAPIResponse)

		err := syncerSvc.SyncStars(context.Background(), mockDatabaseID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Could not find database with ID")
		assert.True(t, gock.IsDone())
	})

	t.Run("should return error with invalid notion credentials", func(t *testing.T) {
		notionAPIResponse := loadFixture(t, path.Join("notionapi", "invalid_credentials_response.json"))
		mockDatabaseID := "7099f06c-95ef-46a2-a47d-672cf8e4a8a4"

		defer gock.Off()

		gock.New(notionAPIURL).
			Get(fmt.Sprintf("/v1/databases/%s", mockDatabaseID)).
			Reply(401).
			JSON(notionAPIResponse)

		err := syncerSvc.SyncStars(context.Background(), mockDatabaseID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API token is invalid")
		assert.True(t, gock.IsDone())
	})

	/**
	* This test should sync 3 pages:
	* - 1 page that exists in the database but not in github (will be deleted)
	* - 3 page that exists in github but not in the database (will be created)
	 */
	t.Run("successfully syncs stars", func(t *testing.T) {
		notionGetDatabaseResponse := loadFixture(t, path.Join("notionapi", "get_database_response.json"))
		notionDatabasePagesResponse := loadFixture(t, path.Join("notionapi", "get_database_pages_response.json"))
		githubStarredReposResponse := loadFixture(t, path.Join("githubapi", "get_starred_repos_response.json"))

		createPage1Request := loadFixture(t, path.Join("notionapi", "create_page_1_request.json"))
		createPage2Request := loadFixture(t, path.Join("notionapi", "create_page_2_request.json"))
		createPage3Request := loadFixture(t, path.Join("notionapi", "create_page_3_request.json"))

		createPageResponse := loadFixture(t, path.Join("notionapi", "create_page_response.json"))
		mockDatabaseID := "705baa92-0ea9-4a4f-bb97-4916d1cb45bc"

		defer gock.Off()

		gock.New(notionAPIURL).
			Get(fmt.Sprintf("/v1/databases/%s", mockDatabaseID)).
			Reply(200).
			JSON(notionGetDatabaseResponse)

		gock.New(notionAPIURL).
			Post(fmt.Sprintf("/v1/databases/%s/query", mockDatabaseID)).
			Reply(200).
			JSON(notionDatabasePagesResponse)

		gock.New(githubAPIURL).
			Get("/user/starred").
			Reply(200).
			JSON(githubStarredReposResponse)

		gock.New(notionAPIURL).
			Post("/v1/pages").
			BodyString(string(createPage1Request)).
			Reply(200).
			JSON(createPageResponse)

		gock.New(notionAPIURL).
			Post("/v1/pages").
			BodyString(string(createPage2Request)).
			Reply(200).
			JSON(createPageResponse)

		gock.New(notionAPIURL).
			Post("/v1/pages").
			BodyString(string(createPage3Request)).
			Reply(200).
			JSON(createPageResponse)

		gock.New(notionAPIURL).
			Patch("/v1/pages/9ef240ab-18de-4808-92ee-22f6dce028e9").
			BodyString(`{"properties":null,"archived":true}`).
			Reply(200).
			JSON(createPageResponse)

		err := syncerSvc.SyncStars(context.Background(), mockDatabaseID)

		assert.NoError(t, err)
		assert.True(t, gock.IsDone())
	})
}
