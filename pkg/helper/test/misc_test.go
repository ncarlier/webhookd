package test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/ncarlier/webhookd/pkg/helper/header"
)

func TestFilterHeaders(t *testing.T) {
	testCases := []struct {
		name           string
		headers        http.Header
		allowedHeaders []string
		expected       http.Header
	}{
		{
			name: "allow all headers",
			headers: http.Header{
				"X-Foo": []string{"bar"},
				"Y-Bar": []string{"baz"},
			},
			allowedHeaders: []string{"*"},
			expected: http.Header{
				"X-Foo": []string{"bar"},
				"Y-Bar": []string{"baz"},
			},
		},
		{
			name: "filter specific header",
			headers: http.Header{
				"X-Foo": []string{"bar"},
				"Y-Bar": []string{"baz"},
			},
			allowedHeaders: []string{"X-Foo"},
			expected: http.Header{
				"X-Foo": []string{"bar"},
			},
		},
		{
			name: "case insensitive filter",
			headers: http.Header{
				"X-Foo-Bar": []string{"baz"},
				"Y-Bar":     []string{"foo"},
			},
			allowedHeaders: []string{"x-foo-bar"},
			expected: http.Header{
				"X-Foo-Bar": []string{"baz"},
			},
		},
		{
			name: "no allowed headers",
			headers: http.Header{
				"X-Foo": []string{"bar"},
			},
			allowedHeaders: []string{},
			expected:       http.Header{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := header.FilterHeaders(tc.headers, tc.allowedHeaders)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
