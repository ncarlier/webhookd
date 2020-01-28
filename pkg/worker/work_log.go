package worker

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ncarlier/webhookd/pkg/tools"
)

// RetrieveLogFile retrieve work log with its name and id
func RetrieveLogFile(id, name, base string) (*os.File, error) {
	logPattern := path.Join(base, fmt.Sprintf("%s_%s_*.txt", tools.ToSnakeCase(name), id))
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
