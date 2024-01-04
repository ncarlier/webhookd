package assert

import (
	"testing"
)

// Nil assert that an object is nil
func Nil(t *testing.T, actual interface{}, message string) {
	if message == "" {
		message = "Nil assertion failed"
	}
	if actual != nil {
		t.Fatalf("%s - actual: %s", message, actual)
	}
}

// NotNil assert that an object is not nil
func NotNil(t *testing.T, actual interface{}, message string) {
	if message == "" {
		message = "Not nil assertion failed"
	}
	if actual == nil {
		t.Fatalf("%s - actual: nil", message)
	}
}

// Equal assert that an object is equal to an expected value
func Equal[K comparable](t *testing.T, expected, actual K, message string) {
	if message == "" {
		message = "Equal assertion failed"
	}
	if actual != expected {
		t.Fatalf("%s - expected: %v, actual: %v", message, expected, actual)
	}
}

// NotEqual assert that an object is not equal to an expected value
func NotEqual[K comparable](t *testing.T, expected, actual K, message string) {
	if message == "" {
		message = "Not equal assertion failed"
	}
	if actual == expected {
		t.Fatalf("%s - unexpected: %v, actual: %v", message, expected, actual)
	}
}

// ContainsStr assert that an array contains an expected value
func Contains[K comparable](t *testing.T, expected K, array []K, message string) {
	if message == "" {
		message = "Array don't contains expected value"
	}
	for _, str := range array {
		if str == expected {
			return
		}
	}
	t.Fatalf("%s - array: %v, expected value: %v", message, array, expected)
}

// True assert that an expression is true
func True(t *testing.T, expression bool, message string) {
	if message == "" {
		message = "Expression is not true"
	}
	if !expression {
		t.Fatalf("%s : %v", message, expression)
	}
}
