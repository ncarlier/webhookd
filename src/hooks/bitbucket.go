package hooks

import (
	"fmt"
)

type BitbucketRecord struct {
	Repository struct {
		Slug string `json:"slug"`
		Name string `json:"name"`
		URL  string `json:"absolute_url"`
	} `json:"repository"`
	BaseURL string `json:"canon_url"`
	User    string `json:"user"`
}

func (r BitbucketRecord) GetGitURL() string {
	return fmt.Sprintf("%s%s", r.BaseURL, r.Repository.URL)
}

func (r BitbucketRecord) GetName() string {
	return r.Repository.Name
}
