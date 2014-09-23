package hook

import (
	"fmt"
)

type BitbucketRecord struct {
	Repository struct {
		Slug  string `json:"slug"`
		Owner string `json:"owner"`
	} `json:"repository"`
}

func (r BitbucketRecord) GetURL() string {
	return fmt.Sprintf("git@bitbucket.org:%s/%s.git", r.Repository.Owner, r.Repository.Owner)
}

func (r BitbucketRecord) GetName() string {
	return r.Repository.Slug
}
