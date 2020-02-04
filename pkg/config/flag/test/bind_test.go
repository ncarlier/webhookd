package test

import (
	"flag"
	"testing"
	"time"

	"github.com/ncarlier/webhookd/pkg/assert"
	configflag "github.com/ncarlier/webhookd/pkg/config/flag"
)

type sampleConfig struct {
	Label         string        `flag:"label" desc:"String parameter" default:"foo"`
	Override      string        `flag:"override" desc:"String parameter to override" default:"bar"`
	Count         int           `flag:"count" desc:"Number parameter" default:"2"`
	Debug         bool          `flag:"debug" desc:"Boolean parameter" default:"false"`
	Timer         time.Duration `flag:"timer" desc:"Duration parameter" default:"30s"`
	Array         []string      `flag:"array" desc:"Array parameter" default:"foo,bar"`
	OverrideArray []string      `flag:"override-array" desc:"Array parameter to override" default:"foo"`
}

func TestFlagBinding(t *testing.T) {
	conf := &sampleConfig{}
	err := configflag.Bind(conf, "FOO")
	flag.CommandLine.Parse([]string{"-override", "test", "-override-array", "a", "-override-array", "b"})
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, "foo", conf.Label, "")
	assert.Equal(t, "test", conf.Override, "")
	assert.Equal(t, 2, conf.Count, "")
	assert.Equal(t, false, conf.Debug, "")
	assert.Equal(t, time.Second*30, conf.Timer, "")
	assert.Equal(t, 2, len(conf.Array), "")
	assert.Equal(t, "foo", conf.Array[0], "")
	assert.Equal(t, 2, len(conf.OverrideArray), "")
	assert.Equal(t, "a", conf.OverrideArray[0], "")
}
