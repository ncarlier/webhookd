package tools

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// ResolveScript is resolving the target script.
func ResolveScript(dir, name string) (string, error) {
	script := path.Clean(path.Join(dir, fmt.Sprintf("%s.sh", name)))
	if !strings.HasPrefix(script, dir) {
		return "", errors.New("Invalid script path: " + script)
	}
	if _, err := os.Stat(script); os.IsNotExist(err) {
		return "", errors.New("Script not found: " + script)
	}

	return script, nil
}
