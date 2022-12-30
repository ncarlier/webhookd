package middleware

import (
	"net"
	"net/http"
)

const xForwardedFor = "X-Forwarded-For"

func getIP(req *http.Request) string {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}
	return ip
}

// XFF is a middleware to identifying the originating IP address using X-Forwarded-For header
func XFF(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xff := r.Header.Get(xForwardedFor)
		if xff == "" {
			r.Header.Set(xForwardedFor, getIP(r))
		}
		inner.ServeHTTP(w, r)
	})
}
