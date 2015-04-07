package api_test

import (
	"github.com/ncarlier/webhookd/api"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server   *httptest.Server
	reader   io.Reader
)

func init() {
	server = httptest.NewServer(api.Handlers())
}

func assertHook(t *testing.T, url string, json string, expectedStatus int) {
	reader = strings.NewReader(json)
	request, err := http.NewRequest("POST", url, reader)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != expectedStatus {
		t.Errorf("Status expected: %d, Actual status: %d", expectedStatus, res.StatusCode)
	}
}

func TestBadHook(t *testing.T) {
	url := fmt.Sprintf("%s/bad/echo", server.URL)
	json := `{"foo": "bar"}`
	assertHook(t, url, json, 404)
}


func TestGitlabHook(t *testing.T) {
	url := fmt.Sprintf("%s/gitlab/echo", server.URL)

	json := `{
		"object_kind": "push",
		"before": "95790bf891e76fee5e1747ab589903a6a1f80f22",
		"after": "da1560886d4f094c3e6c9ef40349f7d38b5d27d7",
		"ref": "refs/heads/master",
		"user_email": "john@example.com",
		"project_id": 15,
		"repository": {
			"name": "Diaspora",
			"url": "git@example.com:mike/diasporadiaspora.git",
			"description": "",
			"git_http_url":"http://example.com/mike/diaspora.git",
			"git_ssh_url":"git@example.com:mike/diaspora.git"
		}
	}`

	assertHook(t, url, json, 200)
}

func TestGithubHook(t *testing.T) {
	url := fmt.Sprintf("%s/github/echo", server.URL)

	json := `{
		"repository": {
			"id": 20000106,
			"name": "public-repo",
			"full_name": "baxterthehacker/public-repo",
			"html_url": "https://github.com/baxterthehacker/public-repo",
			"description": "",
			"url": "https://github.com/baxterthehacker/public-repo",
			"git_url": "git://github.com/baxterthehacker/public-repo.git",
			"ssh_url": "git@github.com:baxterthehacker/public-repo.git",
			"homepage": null
		}
	}`

	assertHook(t, url, json, 200)
}

func TestDockerHook(t *testing.T) {
	url := fmt.Sprintf("%s/docker/echo", server.URL)

	json := `{
		"repository":{
			"status":"Active",
			"description":"my docker repo that does cool things",
			"full_description":"This is my full description",
			"repo_url":"https://registry.hub.docker.com/u/username/reponame/",
			"owner":"username",
			"name":"reponame",
			"namespace":"username",
			"repo_name":"username/reponame"
		}
	}`

	assertHook(t, url, json, 200)
}

