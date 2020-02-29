package api

import (
	"net/http"

	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/middleware"
)

// HandlerFunc custom function handler
type HandlerFunc func(conf *config.Config) http.Handler

// Route is the structure of an HTTP route definition
type Route struct {
	Path        string
	HandlerFunc HandlerFunc
	Middlewares []middleware.Middleware
}

func route(path string, handler HandlerFunc, middlewares ...middleware.Middleware) Route {
	return Route{
		Path:        path,
		HandlerFunc: handler,
		Middlewares: middlewares,
	}
}

// Routes is a list of Route
type Routes []Route

var routes = Routes{
	route("/", index, middleware.Methods("GET", "POST")),
	route("/healthz", healthz, middleware.Methods("GET")),
	route("/varz", varz, middleware.Methods("GET")),
}
