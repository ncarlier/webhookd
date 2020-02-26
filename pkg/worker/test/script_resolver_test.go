package test

import (
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/worker"
)

func TestResolveScript(t *testing.T) {
	script, err := worker.ResolveScript("../../../scripts", "../scripts/echo")
	assert.Nil(t, err, "")
	assert.Equal(t, "../../../scripts/echo.sh", script, "")
}

func TestNotResolveScript(t *testing.T) {
	_, err := worker.ResolveScript("../../scripts", "foo")
	assert.NotNil(t, err, "")
	assert.Equal(t, "Script not found: ../../scripts/foo.sh", err.Error(), "")
}

func TestResolveBadScript(t *testing.T) {
	_, err := worker.ResolveScript("../../scripts", "../tests/test_simple")
	assert.NotNil(t, err, "")
	assert.Equal(t, "Invalid script path: ../tests/test_simple.sh", err.Error(), "")
}

func TestResolveScriptWithExtension(t *testing.T) {
	_, err := worker.ResolveScript("../../scripts", "node.js")
	assert.NotNil(t, err, "")
	assert.Equal(t, "Script not found: ../../scripts/node.js", err.Error(), "")
}
