package auth

import (
	"net/http"
)

// Authenticator is a generic interface to validate an HTTP request
type Authenticator interface {
	Validate(r *http.Request) bool
}
