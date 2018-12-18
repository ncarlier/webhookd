package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ncarlier/webhookd/pkg/auth"
	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/middleware"
)

// NewRouter creates router with declared routes
func NewRouter(conf *config.Config) *http.ServeMux {
	router := http.NewServeMux()
	authenticator := auth.NewAuthenticator(conf)

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc(conf)
		handler = middleware.Logger(handler)
		handler = middleware.Tracing(nextRequestID)(handler)

		if authenticator != nil {
			handler = middleware.Auth(handler, authenticator)
		}
		router.Handle(route.Path, handler)
	}

	return router
}
