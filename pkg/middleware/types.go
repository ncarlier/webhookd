package middleware

import "net/http"

// Middleware function definition
type Middleware func(inner http.Handler) http.Handler

// Middlewares list
type Middlewares []Middleware

// UseBefore insert a middleware at the begining of the middleware chain
func (ms Middlewares) UseBefore(m Middleware) Middlewares {
	return append([]Middleware{m}, ms...)
}

// UseAfter add a middleware at the end of the middleware chain
func (ms Middlewares) UseAfter(m Middleware) Middlewares {
	return append(ms, m)
}
