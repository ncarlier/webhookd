package api

import (
	"net/http"

	"github.com/ncarlier/webhookd/pkg/config"
)

// NewRouter creates router with declared routes
func NewRouter(conf *config.Config) *http.ServeMux {
	router := http.NewServeMux()

	// Register HTTP routes...
	for _, route := range routes(conf) {
		handler := route.HandlerFunc(conf)
		for _, mw := range route.Middlewares {
			handler = mw(handler)
		}
		router.Handle(route.Path, handler)
	}

	return router
}
