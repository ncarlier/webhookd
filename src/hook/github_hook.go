package hook

import (
	"encoding/json"
	"net/http"
)

type GithubRecord struct {
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"git_url"`
	} `json:"repository"`
}

func (r *GithubRecord) GetURL() string {
	return r.Repository.URL
}

func (r *GithubRecord) GetName() string {
	return r.Repository.Name
}

func (r *GithubRecord) Decode(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&r)
	if err != nil {
		return err
	}
	return nil
}
