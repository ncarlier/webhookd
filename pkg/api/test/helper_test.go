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
	assert.Contains(t, "string=foo", values, "")
	assert.Contains(t, "list=foo,bar", values, "")
}

func TestHTTPHeadersToShellVars(t *testing.T) {
	tc := http.Header{
		"Content-Type": []string{"text/plain"},
		"X-Foo-Bar":    []string{"foo", "bar"},
	}
	values := api.HTTPParamsToShellVars(tc)
	assert.Contains(t, "content_type=text/plain", values, "")
	assert.Contains(t, "x_foo_bar=foo,bar", values, "")
}
