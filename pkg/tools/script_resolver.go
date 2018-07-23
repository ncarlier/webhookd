package tools

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/ncarlier/webhookd/pkg/logger"
)

// ResolveScript is resolving the target script.
func ResolveScript(dir, name string) (string, error) {
	script := path.Join(dir, fmt.Sprintf("%s.sh", name))
	logger.Debug.Println("Resolving script: ", script, "...")
	if _, err := os.Stat(script); os.IsNotExist(err) {
		return "", errors.New("Script not found: " + script)
	}

	return script, nil
}
