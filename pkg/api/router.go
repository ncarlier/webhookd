package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ncarlier/webhookd/pkg/auth"
	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/middleware"
	"github.com/ncarlier/webhookd/pkg/pubkey"
)

func nextRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// NewRouter creates router with declared routes
func NewRouter(conf *config.Config) *http.ServeMux {
	router := http.NewServeMux()

	// Load authenticator...
	authenticator, err := auth.NewHtpasswdFromFile(conf.PasswdFile)
	if err != nil {
		logger.Debug.Printf("unable to load htpasswd file (\"%s\"): %s\n", conf.PasswdFile, err)
	}

	// Load key store...
	keystore, err := pubkey.NewKeyStore(conf.KeyStoreURI)
	if err != nil {
		logger.Warning.Printf("unable to load key store (\"%s\"): %s\n", conf.KeyStoreURI, err)
	}

	// Register HTTP routes...
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc(conf)
		handler = middleware.Method(handler, route.Methods)
		handler = middleware.Cors(handler)
		if conf.TLSListenAddr != "" {
			handler = middleware.HSTS(handler)
		}
		handler = middleware.Logger(handler)
		handler = middleware.Tracing(nextRequestID)(handler)

		if keystore != nil {
			handler = middleware.HTTPSignature(handler, keystore)
		}
		if authenticator != nil {
			handler = middleware.Auth(handler, authenticator)
		}
		router.Handle(route.Path, handler)
	}

	return router
}
