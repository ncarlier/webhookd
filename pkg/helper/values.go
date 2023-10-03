package helper

import (
	"net/url"
	"strings"
)

// GetValueOrAlt get value or alt
func GetValueOrAlt(values url.Values, key, alt string) string {
	if val, ok := values[key]; ok {
		return strings.Join(val, ",")
	}
	return alt
}
