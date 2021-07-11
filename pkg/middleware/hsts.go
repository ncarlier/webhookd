package middleware

import (
	"net/http"
)

// HSTS is a middleware to enabling HSTS on HTTP requests
func HSTS(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=15768000 ; includeSubDomains")
		inner.ServeHTTP(w, r)
	})
}
