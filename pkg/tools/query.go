package tools

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase convert string to snakecase.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// QueryParamsToShellVars convert URL query parameters to shell vars.
func QueryParamsToShellVars(q url.Values) []string {
	var params []string
	for k, v := range q {
		var buf bytes.Buffer
		buf.WriteString(ToSnakeCase(k))
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(strings.Join(v[:], ",")))
		params = append(params, buf.String())
	}
	return params
}
