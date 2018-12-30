package worker

import (
	"strconv"
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/logger"
	"github.com/ncarlier/webhookd/pkg/model"
)

func printWorkMessages(work *model.WorkRequest) {
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
	work := model.NewWorkRequest("test", script, payload, args, 5)
	assert.NotNil(t, work, "")
	printWorkMessages(work)
	err := run(work)
	assert.Nil(t, err, "")
	assert.Equal(t, work.Status, model.Success, "")
	assert.Equal(t, work.GetLogContent("notify:"), "OK\n", "")

	// Test that we can retrieve log file afterward
	id := strconv.FormatUint(work.ID, 10)
	logFile, err := RetrieveLogFile(id, "test")
	defer logFile.Close()
	assert.Nil(t, err, "Log file should exists")
	assert.NotNil(t, logFile, "Log file should be retrieve")
}

func TestWorkRunnerWithError(t *testing.T) {
	logger.Init("debug")
	script := "../../tests/test_error.sh"
	work := model.NewWorkRequest("test", script, "", []string{}, 5)
	assert.NotNil(t, work, "")
	printWorkMessages(work)
	err := run(work)
	assert.NotNil(t, err, "")
	assert.Equal(t, work.Status, model.Error, "")
	assert.Equal(t, "exit status 1", err.Error(), "")
}

func TestWorkRunnerWithTimeout(t *testing.T) {
	logger.Init("debug")
	script := "../../tests/test_timeout.sh"
	work := model.NewWorkRequest("test", script, "", []string{}, 1)
	assert.NotNil(t, work, "")
	printWorkMessages(work)
	err := run(work)
	assert.NotNil(t, err, "")
	assert.Equal(t, work.Status, model.Error, "")
	assert.Equal(t, "signal: killed", err.Error(), "")
}
