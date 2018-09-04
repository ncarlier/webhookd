package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ncarlier/webhookd/pkg/logger"
)

type basicAuth struct {
	username   string
	password   string
	authheader string
}

func (c *basicAuth) Init(_ bool) {}

func (c *basicAuth) Usage() string {
	return "HTTP Basic Auth. Usage: -auth basic -authparam <username>:<password>[:<realm>]  (example: -auth basic -auth-param foo:bar)"
}

func (c *basicAuth) ParseParam(param string) error {
	res := strings.Split(param, ":")
	realm := "Authentication required."
	switch len(res) {
	case 3:
		realm = res[2]
		fallthrough
	case 2:
		c.username, c.password = res[0], res[1]
		c.authheader = fmt.Sprintf("Basic realm=\"%s\"", realm)
		return nil
	}
	return errors.New("Invalid Auth param")

}

// BasicAuth HTTP Basic Auth implementation
func (c *basicAuth) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if username, password, ok := r.BasicAuth(); ok && username == c.username && password == c.password {
				logger.Info.Printf("HTTP Basic Auth: %s PASSED\n", username)
				next.ServeHTTP(w, r)
			} else if !ok {
				logger.Debug.Println("HTTP Basic Auth: Auth header not present.")
				w.Header().Add("WWW-Authenticate", c.authheader)
				w.WriteHeader(401)
				w.Write([]byte("Authentication required."))
			} else {
				logger.Warning.Printf("HTTP Basic Auth: Invalid credentials for username %s\n", username)
				w.WriteHeader(403)
				w.Write([]byte("Forbidden."))
			}
		})
	}
}
