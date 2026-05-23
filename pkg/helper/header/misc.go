package header

import (
	"net/http"
	"strings"
)

// FilterHeaders filters the given HTTP headers based on a list of allowed header names.
// If the allowed list contains "*", it returns all headers unmodified.
func FilterHeaders(headers http.Header, allowedHeaders []string) http.Header {
	for _, h := range allowedHeaders {
		if h == "*" {
			return headers
		}
	}

	filteredHeaders := make(http.Header)
	for k, v := range headers {
		for _, allowed := range allowedHeaders {
			if strings.EqualFold(k, allowed) {
				filteredHeaders[k] = v
				break
			}
		}
	}
	return filteredHeaders
}
