package middleware

import (
	"net/http"
)

// Method is a middleware to check that the request use the correct HTTP method
func Method(inner http.Handler, methods []string) http.Handler {
	allowedMethods := make(map[string]struct{}, len(methods))
	for _, s := range methods {
		allowedMethods[s] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := allowedMethods[r.Method]; ok {
			inner.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(405)
		w.Write([]byte("405 Method Not Allowed\n"))
		return
	})
}
