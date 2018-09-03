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
	debug      bool
	authheader string
}

func (c *basicAuth) Init(debug bool) {
	c.debug = debug
	if debug {
		logger.Warning.Println("\u001B[33mBasic Auth: Debug mode enabled. Might Leak sentitive information in log output.\u001B[0m")
	}
}

func (c *basicAuth) Usage() string {
	return "HTTP Basic Auth. Usage: -auth basic -authparam <username>:<password>[:<realm>]  (example: -auth basic -authparam foo:bar)"
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
				if c.debug {
					logger.Debug.Printf("HTTP Basic Auth: %s:%s PASSED\n", username, password)
				}
				next.ServeHTTP(w, r)
			} else if !ok {
				if c.debug {
					logger.Debug.Println("HTTP Basic Auth: Auth header not present.")
				}
				w.Header().Add("WWW-Authenticate", c.authheader)
				w.WriteHeader(401)
				w.Write([]byte("Authentication required."))
			} else {
				if c.debug {
					logger.Debug.Printf("HTTP Basic Auth: Invalid credentials: %s:%s \n", username, password)
				}
				w.WriteHeader(403)
				w.Write([]byte("Forbidden."))
			}
		})
	}
}
