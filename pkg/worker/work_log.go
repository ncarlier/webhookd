package worker

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/ncarlier/webhookd/pkg/tools"
)

// getLogDir returns log directory
func getLogDir() string {
	if value, ok := os.LookupEnv("APP_LOG_DIR"); ok {
		return value
	}
	return os.TempDir()
}

func createLogFile(work *WorkRequest) (*os.File, error) {
	logFilename := path.Join(getLogDir(), fmt.Sprintf("%s_%d_%s.txt", tools.ToSnakeCase(work.Name), work.ID, time.Now().Format("20060102_1504")))
	return os.Create(logFilename)
}

// GetLogFile retrieve work log with its name and id
func GetLogFile(id, name string) (*os.File, error) {
	logPattern := path.Join(getLogDir(), fmt.Sprintf("%s_%s_*.txt", tools.ToSnakeCase(name), id))
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
