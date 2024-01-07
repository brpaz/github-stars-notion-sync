package syncer

import "time"

// starredRepoCollection is a struct that holds all the starred repos.
// it is useful to have this wrapper instead of using a slice directly, because it allows us to add some helper methods like for example, checking if a repository already exists in the collection
type starredRepoCollection struct {
	TotalCount int
	Repos      []starredRepo
}

// starredRepo holds essential information about a starred repository. This is the information that will be synced to notion and avoid using the raw github.Repository struct, which contains a lot of information that is not needed
type starredRepo struct {
	ID          int64
	Name        string
	Description string
	Language    string
	Topics      []string
	URL         string
	StarredAt   time.Time
}

// newStarredRepoCollection creates a new instance starredRepoCollection
func newStarredRepoCollection() *starredRepoCollection {
	return &starredRepoCollection{
		TotalCount: 0,
		Repos:      make([]starredRepo, 0),
	}
}

// Add adds a new starred repo to the collection
func (c *starredRepoCollection) Add(repo starredRepo) {
	c.Repos = append(c.Repos, repo)
	c.TotalCount++
}

// Contains checks if a repository already exists in the collection (by ID, not by name, because a user can have multiple repositories with the same name (but different IDs)
func (c *starredRepoCollection) Contains(repoID int64) bool {
	for _, repo := range c.Repos {
		if repo.ID == repoID {
			return true
		}
	}

	return false
}
