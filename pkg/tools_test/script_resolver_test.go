package tools_test

import (
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/tools"
)

func TestResolveScript(t *testing.T) {
	script, err := tools.ResolveScript("../../scripts", "../scripts/echo")
	assert.Nil(t, err, "")
	assert.Equal(t, "../../scripts/echo.sh", script, "")
}

func TestNotResolveScript(t *testing.T) {
	_, err := tools.ResolveScript("../../scripts", "foo")
	assert.NotNil(t, err, "")
	assert.Equal(t, "Script not found: ../../scripts/foo.sh", err.Error(), "")
}

func TestResolveBadScript(t *testing.T) {
	_, err := tools.ResolveScript("../../scripts", "../tests/test_simple")
	assert.NotNil(t, err, "")
	assert.Equal(t, "Invalid script path: ../tests/test_simple", err.Error(), "")
}
