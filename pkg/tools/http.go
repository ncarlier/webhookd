package tools

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

// ToSnakeCase convert string to snakecase.
func ToSnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	result := strings.Replace(string(out), "/", "_", -1)
	return strings.Replace(result, "-", "", -1)
}

// QueryParamsToShellVars convert URL query parameters to shell vars.
func QueryParamsToShellVars(q url.Values) []string {
	var params []string
	for k, v := range q {
		var buf bytes.Buffer
		value, err := url.QueryUnescape(strings.Join(v[:], ","))
		if err != nil {
			continue
		}
		buf.WriteString(ToSnakeCase(k))
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
		buf.WriteString(ToSnakeCase(k))
		buf.WriteString("=")
		buf.WriteString(value)
		params = append(params, buf.String())
	}
	return params
}
