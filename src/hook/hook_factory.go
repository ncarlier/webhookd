package hook

import (
	"errors"
)

type Record interface {
	GetURL() string
	GetName() string
}

func RecordFactory(hookname string) (Record, error) {
	switch hookname {
	case "bitbucket":
		return new(BitbucketRecord), nil
	case "github":
		return new(GithubRecord), nil
	case "docker":
		return new(DockerRecord), nil
	default:
		return nil, errors.New("Unknown hookname.")
	}
}
