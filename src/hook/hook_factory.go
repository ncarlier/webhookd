package hook

import (
	"errors"
	"net/http"
)

type Record interface {
	GetURL() string
	GetName() string
	Decode(r *http.Request) error
}

func RecordFactory(hookname string) (Record, error) {
	switch hookname {
	case "bitbucket":
		return new(BitbucketRecord), nil
	case "github":
		return new(GithubRecord), nil
	case "gitlab":
		return new(GitlabRecord), nil
	case "docker":
		return new(DockerRecord), nil
	default:
		return nil, errors.New("Unknown hookname: " + hookname)
	}
}
