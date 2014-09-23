package hook

type DockerRecord struct {
	Repository struct {
		Name string `json:"repo_name"`
		URL  string `json:"repo_url"`
	} `json:"repository"`
}

func (r DockerRecord) GetURL() string {
	return r.Repository.URL
}

func (r DockerRecord) GetName() string {
	return r.Repository.Name
}
