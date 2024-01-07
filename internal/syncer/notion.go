package syncer

import (
	"github.com/jomei/notionapi"
)

const (
	databasePropertyTitle       = "Name"
	databasePropertyCreatedTime = "Created time"
	databasePropertyDescription = "Description"
	databasePropertyLanguage    = "Language"
	databasePropertyTopics      = "Topics"
	databasePropertyRepoURL     = "Repository URL"
	databasePropertyRepoID      = "Repository ID"
)

// RequiredProperty represents a required property for the notion database
type RequiredProperty struct {
	PropertyName string
	PropertyType notionapi.PropertyType
}

var requiredProperties = []RequiredProperty{
	{
		PropertyName: databasePropertyCreatedTime,
		PropertyType: notionapi.PropertyTypeCreatedTime,
	},
	{
		PropertyName: databasePropertyDescription,
		PropertyType: notionapi.PropertyTypeRichText,
	},
	{
		PropertyName: databasePropertyLanguage,
		PropertyType: notionapi.PropertyTypeSelect,
	},
	{
		PropertyName: databasePropertyTopics,
		PropertyType: notionapi.PropertyTypeMultiSelect,
	},
	{
		PropertyName: databasePropertyTitle,
		PropertyType: notionapi.PropertyTypeTitle,
	},
	{
		PropertyName: databasePropertyRepoID,
		PropertyType: notionapi.PropertyTypeNumber,
	},
	{
		PropertyName: databasePropertyRepoURL,
		PropertyType: notionapi.PropertyTypeURL,
	},
}

// databasePages is a struct that holds all the existing notion pages.
// It is useful to have this wrapper instead of using a slice directly, because it allows us to add some helper methods
// to the collection itself, like for example, checking if a github repository already exists in the collection
type databasePages struct {
	Pages []notionPage
}

// notionPage is a small representation of a notion page. It holds only the required information for syncing
type notionPage struct {
	ID       string
	Title    string
	GitHubID int64
}

func newDatabasePages() *databasePages {
	return &databasePages{
		Pages: make([]notionPage, 0),
	}
}

// Add adds a new notion page to the collection
func (c *databasePages) Add(page notionPage) {
	c.Pages = append(c.Pages, page)
}

// ContainsRepo checks if a github repository already exists in the collection
func (c *databasePages) ContainsRepo(repoID int64) bool {
	for _, page := range c.Pages {
		if page.GitHubID == repoID {
			return true
		}
	}

	return false
}

// buildCreatePageRequestFromRepo builds a notion page create request from a starred repo object
func buildCreatePageRequestFromRepo(databaseID notionapi.DatabaseID, repo *starredRepo) *notionapi.PageCreateRequest {
	topicsProperty := make([]notionapi.Option, len(repo.Topics))

	for i, topic := range repo.Topics {
		topicsProperty[i] = notionapi.Option{
			Name: topic,
		}
	}

	properties := notionapi.Properties{
		databasePropertyTitle: &notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: &notionapi.Text{
						Content: repo.Name,
					},
				},
			},
		},
		databasePropertyDescription: &notionapi.RichTextProperty{
			RichText: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: &notionapi.Text{
						Content: repo.Description,
					},
				},
			},
		},
		databasePropertyRepoURL: &notionapi.URLProperty{
			URL: repo.URL,
		},
		databasePropertyRepoID: &notionapi.NumberProperty{
			Number: float64(repo.ID),
		},
		databasePropertyTopics: &notionapi.MultiSelectProperty{
			MultiSelect: topicsProperty,
		},
	}

	if repo.Language != "" {
		properties[databasePropertyLanguage] = &notionapi.SelectProperty{
			Select: notionapi.Option{
				Name: repo.Language,
			},
		}
	}

	request := &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: databaseID,
		},
		Properties: properties,
	}

	return request
}
