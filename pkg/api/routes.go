package api

import (
	"github.com/ncarlier/webhookd/pkg/auth"
	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/middleware"
	"github.com/ncarlier/webhookd/pkg/truststore"
)

var commonMiddlewares = middleware.Middlewares{
	middleware.Cors,
	middleware.Logger,
	middleware.Tracing(nextRequestID),
}

func buildMiddlewares(conf *config.Config) middleware.Middlewares {
	var middlewares = commonMiddlewares
	if conf.TLS {
		middlewares = middlewares.UseAfter(middleware.HSTS)
	}

	// Load trust store...
	ts, err := truststore.New(conf.TrustStoreFile)
	if err != nil {
		logger.Warning.Printf("unable to load trust store (\"%s\"): %s\n", conf.TrustStoreFile, err)
	}
	if ts != nil {
		middlewares = middlewares.UseAfter(middleware.Signature(ts))
	}

	// Load authenticator...
	authenticator, err := auth.NewHtpasswdFromFile(conf.PasswdFile)
	if err != nil {
		logger.Debug.Printf("unable to load htpasswd file (\"%s\"): %s\n", conf.PasswdFile, err)
	}
	if authenticator != nil {
		middlewares = middlewares.UseAfter(middleware.AuthN(authenticator))
	}
	return middlewares
}

func routes(conf *config.Config) Routes {
	middlewares := buildMiddlewares(conf)
	staticPath := conf.StaticPath + "/"
	return Routes{
		route(
			"/",
			index,
			middlewares...,
		),
		route(
			staticPath,
			static(staticPath),
			middlewares.UseBefore(middleware.Methods("GET"))...,
		),
		route(
			"/healthz",
			healthz,
			commonMiddlewares.UseBefore(middleware.Methods("GET"))...,
		),
		route(
			"/varz",
			varz,
			middlewares.UseBefore(middleware.Methods("GET"))...,
		),
	}
}
