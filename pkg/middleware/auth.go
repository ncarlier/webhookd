package middleware

import (
	"net/http"

	"github.com/ncarlier/webhookd/pkg/auth"
)

// Auth is a middleware to checks HTTP request credentials
func Auth(inner http.Handler, authn auth.Authenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authn.Validate(r) {
			inner.ServeHTTP(w, r)
			return
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="Ah ah ah, you didn't say the magic word"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
	})
}
