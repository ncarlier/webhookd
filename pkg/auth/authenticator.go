package auth

import (
	"net/http"
)

// Authenticator is a generic interface to validate HTTP request credentials.
// It's returns the authentication result along with the principal (username) if it has one.
type Authenticator interface {
	Validate(r *http.Request) (ok bool, username string)
}
