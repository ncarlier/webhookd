package middleware

import (
	"net/http"

	"github.com/ncarlier/webhookd/pkg/auth"
)

const xWebAuthUser = "X-WebAuth-User"

// AuthN is a middleware to checks HTTP request credentials
func AuthN(authenticator auth.Authenticator) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Del(xWebAuthUser)
			if ok, username := authenticator.Validate(r); ok {
				w.Header().Set(xWebAuthUser, username)
				next.ServeHTTP(w, r)
				return
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Ah ah ah, you didn't say the magic word"`)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
		})
	}
}
