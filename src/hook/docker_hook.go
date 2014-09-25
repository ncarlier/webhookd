package hook

import (
	"encoding/json"
	"net/http"
)

type DockerRecord struct {
	Repository struct {
		Name string `json:"repo_name"`
		URL  string `json:"repo_url"`
	} `json:"repository"`
}

func (r *DockerRecord) GetURL() string {
	return r.Repository.URL
}

func (r *DockerRecord) GetName() string {
	return r.Repository.Name
}

func (r *DockerRecord) Decode(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&r)
	if err != nil {
		return err
	}
	return nil
}
