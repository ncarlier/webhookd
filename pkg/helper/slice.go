package helper

import "strings"

// ContainsFold checks if array contains string (case-insensitive)
func ContainsFold(arr []string, str string) bool {
	if len(arr) == 0 {
		return false
	}

	for _, s := range arr {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}
