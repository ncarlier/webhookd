package hooks

type GithubRecord struct {
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"git_url"`
	} `json:"repository"`
}

func (r GithubRecord) GetGitURL() string {
	return r.Repository.URL
}

func (r GithubRecord) GetName() string {
	return r.Repository.Name
}
