package api

import (
	"net/http"

	"github.com/ncarlier/webhookd/pkg/config"
)

// HandlerFunc custom function handler
type HandlerFunc func(conf *config.Config) http.Handler

// Route is the structure of an HTTP route definition
type Route struct {
	Method      string
	Path        string
	HandlerFunc HandlerFunc
}

// Routes is a list of Route
type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/",
		index,
	},
	Route{
		"GET",
		"/healtz",
		healthz,
	},
	Route{
		"GET",
		"/varz",
		varz,
	},
}
