package worker

import (
	"strconv"
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/logger"
)

func printWorkMessages(work *WorkRequest) {
	go func() {
		for {
			msg, open := <-work.MessageChan
			if !open {
				break
			}
			logger.Info.Println(string(msg))
		}
	}()
}

func TestWorkRunner(t *testing.T) {
	logger.Init("debug")
	script := "../../tests/test_simple.sh"
	args := []string{
		"name=foo",
		"user_agent=test",
	}
	payload := "{\"foo\": \"bar\"}"
	work := NewWorkRequest("test", script, payload, args, 5)
	assert.NotNil(t, work, "")
	printWorkMessages(work)
	err := run(work)
	assert.Nil(t, err, "")
	assert.Equal(t, work.Status, Success, "")

	// Test that log file is ok
	id := strconv.FormatUint(work.ID, 10)
	logFile, err := GetLogFile(id, "test")
	defer logFile.Close()
	assert.Nil(t, err, "Log file should exists")
	assert.NotNil(t, logFile, "Log file should be retrieve")
}

func TestWorkRunnerWithError(t *testing.T) {
	logger.Init("debug")
	script := "../../tests/test_error.sh"
	work := NewWorkRequest("test", script, "", []string{}, 5)
	assert.NotNil(t, work, "")
	printWorkMessages(work)
	err := run(work)
	assert.NotNil(t, err, "")
	assert.Equal(t, work.Status, Error, "")
	assert.Equal(t, "exit status 1", err.Error(), "")
}

func TestWorkRunnerWithTimeout(t *testing.T) {
	logger.Init("debug")
	script := "../../tests/test_timeout.sh"
	work := NewWorkRequest("test", script, "", []string{}, 1)
	assert.NotNil(t, work, "")
	printWorkMessages(work)
	err := run(work)
	assert.NotNil(t, err, "")
	assert.Equal(t, work.Status, Error, "")
	assert.Equal(t, "signal: killed", err.Error(), "")
}
