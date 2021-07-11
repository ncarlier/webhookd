package middleware

import (
	"net/http"
)

// Methods is a middleware to check that the request use the correct HTTP method
func Methods(methods ...string) Middleware {
	return func(next http.Handler) http.Handler {
		allowedMethods := make(map[string]struct{}, len(methods))
		for _, s := range methods {
			allowedMethods[s] = struct{}{}
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := allowedMethods[r.Method]; ok {
				next.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(405)
			w.Write([]byte("405 Method Not Allowed\n"))
		})
	}
}
