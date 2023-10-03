package test

import (
	"log/slog"
	"os"
	"strconv"
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/hook"
)

func printJobMessages(job *hook.Job) {
	go func() {
		for {
			msg, open := <-job.MessageChan
			if !open {
				break
			}
			slog.Info(string(msg))
		}
	}()
}

func TestHookJob(t *testing.T) {
	req := &hook.Request{
		Name:    "test_simple",
		Method:  "GET",
		Payload: "{\"foo\": \"bar\"}",
		Args: []string{
			"name=foo",
			"user_agent=test",
		},
		Timeout:   5,
		BaseDir:   "../test",
		OutputDir: os.TempDir(),
	}
	job, err := hook.NewHookJob(req)
	assert.Nil(t, err, "")
	assert.NotNil(t, job, "")
	printJobMessages(job)
	err = job.Run()
	assert.Nil(t, err, "")
	assert.Equal(t, job.Status(), hook.Success, "")
	assert.Equal(t, job.Logs("notify:"), "OK\n", "")

	// Test that we can retrieve log file afterward
	id := strconv.FormatUint(job.ID(), 10)
	logFile, err := hook.Logs(id, "test", os.TempDir())
	assert.Nil(t, err, "Log file should exists")
	defer logFile.Close()
	assert.NotNil(t, logFile, "Log file should be retrieve")
}

func TestWorkRunnerWithError(t *testing.T) {
	req := &hook.Request{
		Name:      "test_error",
		Method:    "POST",
		Payload:   "",
		Args:      []string{},
		Timeout:   5,
		BaseDir:   "../test",
		OutputDir: os.TempDir(),
	}
	job, err := hook.NewHookJob(req)
	assert.Nil(t, err, "")
	assert.NotNil(t, job, "")
	printJobMessages(job)
	err = job.Run()
	assert.NotNil(t, err, "")
	assert.Equal(t, job.Status(), hook.Error, "")
	assert.Equal(t, "exit status 1", err.Error(), "")
}

func TestWorkRunnerWithTimeout(t *testing.T) {
	req := &hook.Request{
		Name:      "test_timeout",
		Method:    "POST",
		Payload:   "",
		Args:      []string{},
		Timeout:   1,
		BaseDir:   "../test",
		OutputDir: os.TempDir(),
	}
	job, err := hook.NewHookJob(req)
	assert.Nil(t, err, "")
	assert.NotNil(t, job, "")
	printJobMessages(job)
	err = job.Run()
	assert.NotNil(t, err, "")
	assert.Equal(t, job.Status(), hook.Error, "")
	assert.Equal(t, "signal: killed", err.Error(), "")
}
