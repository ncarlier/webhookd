package auth

import (
	"net/http"

	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/logger"
)

// Authenticator is a generic interface to validate an HTTP request
type Authenticator interface {
	Validate(r *http.Request) bool
}

// NewAuthenticator creates new authenticator form the configuration
func NewAuthenticator(conf *config.Config) Authenticator {
	authenticator, err := NewHtpasswdFromFile(*conf.PasswdFile)
	if err != nil {
		logger.Debug.Printf("unable to load htpasswd file: \"%s\" (%s)\n", *conf.PasswdFile, err)
		return nil
	}
	return authenticator
}
