package middleware

import (
	"fmt"
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
			addr := r.RemoteAddr
			if i := strings.LastIndex(addr, ":"); i != -1 {
				addr = addr[:i]
			}
			logger.Info.Printf(
				"%s - - [%s] %q %d %d %q %q %q",
				addr,
				start.Format("02/Jan/2006:15:04:05 -0700"),
				fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
				o.status,
				o.written,
				r.Referer(),
				r.UserAgent(),
				fmt.Sprintf("REQID=%s", requestID),
			)
		}()
		next.ServeHTTP(o, r)
	})
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
