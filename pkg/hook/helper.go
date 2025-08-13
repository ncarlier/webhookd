package hook

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ResolveScript is resolving the target script.
func ResolveScript(dir, name, defaultExt string) (string, error) {
	if filepath.Ext(name) == "" {
		name += "." + defaultExt
	}
	script := filepath.Clean(filepath.Join(dir, name))
	if !strings.HasPrefix(script, dir) {
		return "", errors.New("Invalid script path: " + name)
	}
	if _, err := os.Stat(script); os.IsNotExist(err) {
		return "", errors.New("Script not found: " + script)
	}

	return script, nil
}
