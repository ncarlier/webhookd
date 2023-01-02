package test

import (
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/helper"
)

func TestToSnakeCase(t *testing.T) {
	testCases := []struct {
		value    string
		expected string
	}{
		{"hello-world", "hello_world"},
		{"helloWorld", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"Hello/_World", "hello__world"},
		{"Hello/world", "hello_world"},
	}
	for _, tc := range testCases {
		value := helper.ToSnake(tc.value)
		assert.Equal(t, tc.expected, value, "")
	}
}
