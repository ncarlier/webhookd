package middleware

import (
	"net/http"
	"time"

	"github.com/ncarlier/webhookd/pkg/logger"
)

type key int

const (
	requestIDKey key = 0
)

// Logger is a middleware to log HTTP request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "unknown"
			}
			logger.Info.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}
