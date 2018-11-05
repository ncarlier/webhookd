package tools_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/tools"
)

func TestToSnakeCase(t *testing.T) {
	testCases := []struct {
		value    string
		expected string
	}{
		{"hello-world", "helloworld"},
		{"helloWorld", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"Hello/World", "hello__world"},
		{"Hello/world", "hello_world"},
	}
	for _, tc := range testCases {
		value := tools.ToSnakeCase(tc.value)
		assert.Equal(t, tc.expected, value, "")
	}
}

func TestQueryParamsToShellVars(t *testing.T) {
	tc := url.Values{
		"string": []string{"foo"},
		"list":   []string{"foo", "bar"},
	}
	values := tools.QueryParamsToShellVars(tc)
	assert.ContainsStr(t, "string=foo", values, "")
	assert.ContainsStr(t, "list=foo,bar", values, "")
}

func TestHTTPHeadersToShellVars(t *testing.T) {
	tc := http.Header{
		"Content-Type": []string{"text/plain"},
		"X-Foo-Bar":    []string{"foo", "bar"},
	}
	values := tools.HTTPHeadersToShellVars(tc)
	assert.ContainsStr(t, "content_type=text/plain", values, "")
	assert.ContainsStr(t, "x_foo_bar=foo,bar", values, "")
}
