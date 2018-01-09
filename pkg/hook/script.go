package hook

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/ncarlier/webhookd/pkg/logger"
)

var (
	scriptsdir = os.Getenv("APP_SCRIPTS_DIR")
)

// ResolveScript is resolving the target script.
func ResolveScript(p string) (string, error) {
	if scriptsdir == "" {
		scriptsdir = "scripts"
	}

	script := path.Join(scriptsdir, fmt.Sprintf("%s.sh", p))
	logger.Debug.Println("Resolving script: ", script, "...")
	if _, err := os.Stat(script); os.IsNotExist(err) {
		return "", errors.New("Script not found: " + script)
	}

	return script, nil
}
