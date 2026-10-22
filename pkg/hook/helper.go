package hook

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ResolveScript is resolving the target script.
func ResolveScript(dir, name, defaultExt string) (string, error) {
	// Reject decoded "." and ".." path segments
	for _, part := range strings.Split(name, "/") {
		if part == "." || part == ".." {
			return "", errors.New("invalid script path: " + name)
		}
	}

	scriptPath := path.Join(dir, name)
	// Validate script path
	scriptPath = path.Clean(scriptPath)

	rel, err := filepath.Rel(dir, scriptPath)
	if err != nil || rel == ".." || strings.HasPrefix(rel, "../") || strings.HasPrefix(rel, "..\\") {
		return "", errors.New("invalid script path: " + name)
	}
	// Check if the script exists
	if _, err := os.Stat(scriptPath); errors.Is(err, os.ErrNotExist) {
		// Try to add the default extension if not provided
		if path.Ext(name) == "" {
			scriptPathWithExt := scriptPath + "." + defaultExt
			if _, err := os.Stat(scriptPathWithExt); errors.Is(err, os.ErrNotExist) {
				return "", errors.New("script not found: " + scriptPath + "[." + defaultExt + "]")
			} else {
				return scriptPathWithExt, nil
			}
		}
		return "", errors.New("script not found: " + scriptPath)
	} else {
		return scriptPath, nil
	}
}
