package auth

import "net/http"

// Method an interface describing an authentication method
type Method interface {
	// Called after ParseParam method.
	// auth.Method should initialize itself here and get ready to receive requests.
	// Logger has been initialized so it is safe to call logger methods here.
	Init(debug bool)
	// Return Method Usage Info
	Usage() string
	// Parse the parameter passed through the -authparam flag
	// Logger is not initialized at this state so do NOT call logger methods
	// If the parameter is unacceptable, return an error and main should exit
	ParseParam(string) error
	// Return a middleware to handle connections.
	Middleware() func(http.Handler) http.Handler
}

var (
	// AvailableMethods Returns a map of available auth methods
	AvailableMethods = map[string]Method{
		"none":  new(noAuth),
		"basic": new(basicAuth),
	}
)
