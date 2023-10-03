package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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
		o := &responseObserver{ResponseWriter: w}
		start := time.Now()
		defer func() {
			requestID, ok := r.Context().Value(requestIDKey).(string)
			if !ok {
				requestID = "0"
			}
			logger.LogIf(
				logger.RequestOutputEnabled,
				slog.LevelInfo+1,
				fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
				"ip", getRequestIP(r),
				"time", start.Format("02/Jan/2006:15:04:05 -0700"),
				"duration", time.Since(start).Milliseconds(),
				"status", o.status,
				"bytes", o.written,
				"referer", r.Referer(),
				"ua", r.UserAgent(),
				"reqid", requestID,
			)
		}()
		next.ServeHTTP(o, r)
	})
}

func getRequestIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if comma := strings.Index(ip, ","); comma != -1 {
		ip = ip[0:comma]
	}
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func (o *responseObserver) Flush() {
	flusher, ok := o.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
