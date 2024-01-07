package syncer

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/jomei/notionapi"

	"github.com/brpaz/github-stars-notion-sync/internal/log"
)

var (
	ErrNilGithubClient = errors.New("github client cannot be nil")
	ErrNilNotionClient = errors.New("notion client cannot be nil")
)

const (
	githubReposPerPage = 100
	notionPagesPerPage = 50
)

type Syncer struct {
	github *github.Client
	notion *notionapi.Client
}

// New creates a new Syncer instance with the given github and notion clients
func New(githubClient *github.Client, notionClient *notionapi.Client) (*Syncer, error) {
	if githubClient == nil {
		return nil, ErrNilGithubClient
	}

	if notionClient == nil {
		return nil, ErrNilNotionClient
	}

	return &Syncer{
		github: githubClient,
		notion: notionClient,
	}, nil
}

// SyncStars syncs the github stars with the notion database
// This method will:
// 1. Get all the starred repos from github
// 2. Get all the pages from the notion database
// 3. Compare the two lists and create/delete the notion pages accordingly
func (s *Syncer) SyncStars(ctx context.Context, notionDatabaseID string) error {
	log.Info(ctx, "starting syncer")

	databaseID := notionapi.DatabaseID(notionDatabaseID)
	notionDatabase, err := s.notion.Database.Get(ctx, databaseID)
	if err != nil {
		return fmt.Errorf("error getting notion database: %w", err)
	}

	// ensure that the notion database has the required fields.
	// this is critical to ensure that the syncer works as expected.
	if err := s.validateDatabaseFields(notionDatabase); err != nil {
		return fmt.Errorf("error validating notion database: %w", err)
	}

	log.Info(ctx, "fetching pages from notion database. Depending on the size of the database, this might take a while.")
	notionPages, err := s.getPagesFromNotionDatabase(ctx, databaseID)
	if err != nil {
		return fmt.Errorf("error getting notion pages: %w", err)
	}

	log.Info(ctx, fmt.Sprintf("found %d pages in notion", len(notionPages.Pages)))

	log.Info(ctx, "fetching starred repos from github. Depending on the number of starred repos, this might take a while.")
	starredRepos, err := s.fetchGitHubStarredRepos(ctx)
	if err != nil {
		return fmt.Errorf("error getting starred repos: %w", err)
	}

	log.Info(ctx, fmt.Sprintf("found %d starred repos in github", len(starredRepos.Repos)))

	if err := s.doSync(ctx, databaseID, notionPages, starredRepos); err != nil {
		return fmt.Errorf("error syncing notion database: %w", err)
	}

	return nil
}

func (s *Syncer) validateDatabaseFields(database *notionapi.Database) error {
	for _, requiredProperty := range requiredProperties {
		if _, ok := database.Properties[requiredProperty.PropertyName]; !ok {
			return fmt.Errorf("notion database is missing required property %s", requiredProperty.PropertyName)
		}

		if notionapi.PropertyType(database.Properties[requiredProperty.PropertyName].GetType()) != requiredProperty.PropertyType {
			return fmt.Errorf("notion database property %s is of type %s, but should be %s", requiredProperty.PropertyName, database.Properties[requiredProperty.PropertyName].GetType(), requiredProperty.PropertyType)
		}
	}

	return nil
}

// fetchGitHubStarredRepos returns a collection of starred repos from github
func (s *Syncer) fetchGitHubStarredRepos(ctx context.Context) (*starredRepoCollection, error) {
	starredRepos := newStarredRepoCollection()

	opt := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: githubReposPerPage},
	}

	for {
		repos, resp, err := s.github.Activity.ListStarred(ctx, "", opt)
		if err != nil {
			return starredRepos, err
		}

		for _, repo := range repos {
			starredRepos.Add(starredRepo{
				ID:          repo.Repository.GetID(),
				Name:        repo.Repository.GetName(),
				Description: repo.Repository.GetDescription(),
				URL:         repo.Repository.GetHTMLURL(),
				Topics:      repo.Repository.Topics,
				Language:    repo.Repository.GetLanguage(),
				StarredAt:   repo.StarredAt.Time,
			})
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return starredRepos, nil
}

// getPagesFromNotionDatabase returns a collection of pages from the specified notion database
func (s *Syncer) getPagesFromNotionDatabase(ctx context.Context, databaseID notionapi.DatabaseID) (*databasePages, error) {
	pages := newDatabasePages()
	cursor := notionapi.Cursor("")

	for {
		resp, err := s.notion.Database.Query(ctx, databaseID, &notionapi.DatabaseQueryRequest{
			PageSize:    notionPagesPerPage,
			StartCursor: cursor,
		})
		if err != nil {
			return pages, err
		}

		for _, result := range resp.Results {
			titleProperty := result.Properties[databasePropertyTitle].(*notionapi.TitleProperty)
			repoIDProperty := result.Properties[databasePropertyRepoID].(*notionapi.NumberProperty)

			pages.Add(notionPage{
				ID:       result.ID.String(),
				Title:    titleProperty.Title[0].PlainText,
				GitHubID: int64(repoIDProperty.Number),
			})
		}

		if !resp.HasMore {
			break
		}

		cursor = resp.NextCursor
	}

	return pages, nil
}

// doSync compares the notion pages with the starred repos and creates/deletes the notion pages accordingly
func (s *Syncer) doSync(ctx context.Context, databaseID notionapi.DatabaseID, notionPages *databasePages, starredRepos *starredRepoCollection) error {
	pagesToCreate := make([]starredRepo, 0)
	pagesToDelete := make([]notionPage, 0)

	// find the pages that need to be created (i.e. starred repos that are not in the notion database)
	for _, repo := range starredRepos.Repos {
		if !notionPages.ContainsRepo(repo.ID) {
			pagesToCreate = append(pagesToCreate, repo)
		}
	}

	// find the pages that need to be deleted (i.e. notion pages whose repo is not starred anymore)
	for _, page := range notionPages.Pages {
		if !starredRepos.Contains(page.GitHubID) {
			pagesToDelete = append(pagesToDelete, page)
		}
	}

	log.Info(ctx, fmt.Sprintf("found %d pages to create", len(pagesToCreate)))

	for _, repo := range pagesToCreate {
		if err := s.createNotionPage(ctx, databaseID, &repo); err != nil {
			log.Error(ctx, "error creating notion page", log.String("repo", repo.Name), log.String("error", err.Error()))
			continue
		}

		log.Info(ctx, "notion page created", log.String("repo", repo.Name))
	}

	log.Info(ctx, fmt.Sprintf("found %d pages to delete", len(pagesToDelete)))
	for _, page := range pagesToDelete {
		if err := s.deleteNotionPage(ctx, notionapi.PageID(page.ID)); err != nil {
			log.Error(ctx, "error deleting notion page", log.String("page", page.Title), log.String("error", err.Error()))
			continue
		}

		log.Info(ctx, "notion page deleted", log.String("page", page.Title))
	}

	return nil
}

func (s *Syncer) createNotionPage(ctx context.Context, databaseID notionapi.DatabaseID, repo *starredRepo) error {
	request := buildCreatePageRequestFromRepo(databaseID, repo)
	_, err := s.notion.Page.Create(ctx, request)

	return err
}

func (s *Syncer) deleteNotionPage(ctx context.Context, pageID notionapi.PageID) error {
	_, err := s.notion.Page.Update(ctx, pageID, &notionapi.PageUpdateRequest{
		Archived: true,
	})

	return err
}
