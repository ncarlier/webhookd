package metric

import (
	"expvar"
	"runtime"
	"time"
)

var startTime = time.Now().UTC()

func goroutines() interface{} {
	return runtime.NumGoroutine()
}

// uptime is an expvar.Func compliant wrapper for uptime info.
func uptime() interface{} {
	uptime := time.Since(startTime)
	return int64(uptime)
}

var stats = expvar.NewMap("hookstats")

var (
	// Requests count the number of request
	Requests expvar.Int
	// RequestsFailed count the number of failed request
	RequestsFailed expvar.Int
)

func init() {
	stats.Set("requests", &Requests)
	stats.Set("requests_failed", &RequestsFailed)
	expvar.Publish("goroutines", expvar.Func(goroutines))
	expvar.Publish("uptime", expvar.Func(uptime))
}
