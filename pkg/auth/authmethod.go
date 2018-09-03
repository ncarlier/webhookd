package auth

import "net/http"

// Method an interface describing an authentication method
type Method interface {
	Init(debug bool)
	Usage() string
	ParseParam(string) error
	Middleware() func(http.Handler) http.Handler
}

var (
	// AvailableMethods Returns a map of available auth methods
	AvailableMethods = map[string]Method{
		"none":  new(noAuth),
		"basic": new(basicAuth),
	}
)
