package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ncarlier/webhookd/pkg/strcase"
)

// URLValuesToShellVars convert URL values to shell vars.
func URLValuesToShellVars(q url.Values) []string {
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

func nextRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
