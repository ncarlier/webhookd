package config_test

import (
	"flag"
	"testing"

	"github.com/ncarlier/webhookd/pkg/assert"
	"github.com/ncarlier/webhookd/pkg/config"
)

func TestFlagBuilder(t *testing.T) {
	flag.Parse()
	conf := &config.Config{}
	err := config.HydrateFromFlags(conf)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, ":8080", conf.ListenAddr, "")
	assert.Equal(t, 2, conf.NbWorkers, "")
	assert.Equal(t, 10, conf.Timeout, "")
	assert.Equal(t, "scripts", conf.ScriptDir, "")
	assert.Equal(t, false, conf.Debug, "")
}
