package hook

import (
	"encoding/json"
	"net/http"
)

type GitlabRecord struct {
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"git_ssh_url"`
	} `json:"repository"`
}

func (r *GitlabRecord) GetURL() string {
	return r.Repository.URL
}

func (r *GitlabRecord) GetName() string {
	return r.Repository.Name
}

func (r *GitlabRecord) Decode(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&r)
	if err != nil {
		return err
	}
	return nil
}
