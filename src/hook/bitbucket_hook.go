package hook

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (r BitbucketRecord) Decode(req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	if payload, ok := req.PostForm["payload"]; ok {
		err := json.Unmarshal([]byte(payload[0]), &r)
		if err != nil {
			return err
		}
	} else {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&r)
		if err != nil {
			return err
		}
	}
	return nil
}
