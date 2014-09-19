package hooks

import (
	"errors"
)

type Record interface {
	GetGitURL() string
	GetName() string
}

func RecordFactory(hookname string) (Record, error) {
	switch hookname {
	case "bitbucket":
		return new(BitbucketRecord), nil
	case "github":
		return new(GithubRecord), nil
	default:
		return nil, errors.New("Unknown hookname.")
	}
}
