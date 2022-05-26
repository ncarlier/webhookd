package middleware

import (
	"net/http"
)

// Cors is a middleware to enabling CORS on HTTP requests
func Cors(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if r.Method != "OPTIONS" {
			inner.ServeHTTP(w, r)
		}
	})
}
