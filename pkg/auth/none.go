package auth

import (
	"net/http"
)

type noAuth struct {
}

func (c *noAuth) Usage() string {
	return "No Auth. Usage: -auth none"
}

func (c *noAuth) Init(_ bool) {}

func (c *noAuth) ParseParam(_ string) error {
	return nil
}

// NoAuth A Nop Auth middleware
func (c *noAuth) Middleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return h
	}
}
