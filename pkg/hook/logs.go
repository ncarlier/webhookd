package hook

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ncarlier/webhookd/pkg/helper"
)

// GetLogFile get hook log with its name and id
func GetLogFile(id, name, base string) (*os.File, error) {
	logPattern := filepath.Join(base, fmt.Sprintf("%s_%s_*.txt", helper.ToSnake(name), id))
	files, err := filepath.Glob(logPattern)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		filename := files[len(files)-1]
		return os.Open(filename)
	}
	return nil, nil
}
