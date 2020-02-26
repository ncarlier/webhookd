package api

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"

	"github.com/ncarlier/webhookd/pkg/strcase"
)

// QueryParamsToShellVars convert URL query parameters to shell vars.
func QueryParamsToShellVars(q url.Values) []string {
	var params []string
	for k, v := range q {
		var buf bytes.Buffer
		value, err := url.QueryUnescape(strings.Join(v[:], ","))
		if err != nil {
			continue
		}
		buf.WriteString(strcase.ToSnake(k))
		buf.WriteString("=")
		buf.WriteString(value)
		params = append(params, buf.String())
	}
	return params
}

// HTTPHeadersToShellVars convert HTTP headers to shell vars.
func HTTPHeadersToShellVars(h http.Header) []string {
	var params []string
	for k, v := range h {
		var buf bytes.Buffer
		value, err := url.QueryUnescape(strings.Join(v[:], ","))
		if err != nil {
			continue
		}
		buf.WriteString(strcase.ToSnake(k))
		buf.WriteString("=")
		buf.WriteString(value)
		params = append(params, buf.String())
	}
	return params
}
