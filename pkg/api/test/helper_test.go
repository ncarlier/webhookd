package test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/ncarlier/webhookd/pkg/api"
	"github.com/ncarlier/webhookd/pkg/assert"
)

func TestQueryParamsToShellVars(t *testing.T) {
	tc := url.Values{
		"string": []string{"foo"},
		"list":   []string{"foo", "bar"},
	}
	values := api.HTTPParamsToShellVars(tc)
	assert.ContainsStr(t, "string=foo", values, "")
	assert.ContainsStr(t, "list=foo,bar", values, "")
}

func TestHTTPHeadersToShellVars(t *testing.T) {
	tc := http.Header{
		"Content-Type": []string{"text/plain"},
		"X-Foo-Bar":    []string{"foo", "bar"},
	}
	values := api.HTTPParamsToShellVars(tc)
	assert.ContainsStr(t, "content_type=text/plain", values, "")
	assert.ContainsStr(t, "x_foo_bar=foo,bar", values, "")
}
