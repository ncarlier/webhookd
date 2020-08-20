package api

import (
	"net/http"

	"github.com/ncarlier/webhookd/pkg/auth"
	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/middleware"
	"github.com/ncarlier/webhookd/pkg/pubkey"
)

var commonMiddlewares = []middleware.Middleware{
	middleware.Cors,
	middleware.Logger,
	middleware.Tracing(nextRequestID),
}

// NewRouter creates router with declared routes
func NewRouter(conf *config.Config) *http.ServeMux {
	router := http.NewServeMux()

	var middlewares = commonMiddlewares
	if conf.TLS {
		middlewares = append(middlewares, middleware.HSTS)
	}

	// Load trust store...
	trustStore, err := pubkey.NewTrustStore(conf.TrustStoreFile)
	if err != nil {
		logger.Warning.Printf("unable to load trust store (\"%s\"): %s\n", conf.TrustStoreFile, err)
	}
	if trustStore != nil {
		middlewares = append(middlewares, middleware.HTTPSignature(trustStore))
	}

	// Load authenticator...
	authenticator, err := auth.NewHtpasswdFromFile(conf.PasswdFile)
	if err != nil {
		logger.Debug.Printf("unable to load htpasswd file (\"%s\"): %s\n", conf.PasswdFile, err)
	}
	if authenticator != nil {
		middlewares = append(middlewares, middleware.AuthN(authenticator))
	}

	// Register HTTP routes...
	for _, route := range routes {
		handler := route.HandlerFunc(conf)
		for _, mw := range route.Middlewares {
			handler = mw(handler)
		}
		for _, mw := range middlewares {
			if route.Path == "/healthz" {
				continue
			}
			handler = mw(handler)
		}
		router.Handle(route.Path, handler)
	}

	return router
}
