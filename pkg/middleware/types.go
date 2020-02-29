package middleware

import "net/http"

// Middleware function definition
type Middleware func(inner http.Handler) http.Handler
