package hook

import (
	"errors"
	"os"
	"path"
	"strings"
)

// ResolveScript is resolving the target script.
func ResolveScript(dir, name string) (string, error) {
	if path.Ext(name) == "" {
		name += ".sh"
	}
	script := path.Clean(path.Join(dir, name))
	if !strings.HasPrefix(script, dir) {
		return "", errors.New("Invalid script path: " + name)
	}
	if _, err := os.Stat(script); os.IsNotExist(err) {
		return "", errors.New("Script not found: " + script)
	}

	return script, nil
}
