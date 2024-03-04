package config

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ncarlier/webhookd/pkg/helper"
)

// OldConfig contain global configuration
type OldConfig struct {
	NbWorkers int    `flag:"nb-workers" desc:"Number of workers to start [DEPRECATED]" default:"2"`
	Scripts   string `flag:"scripts" desc:"Scripts location [DEPRECATED]" default:"scripts"`
}

// ManageDeprecatedFlags manage legacy configuration
func (c *Config) ManageDeprecatedFlags(prefix string) {
	if isUsingDeprecatedConfigParam(prefix, "nb-workers") {
		c.Hook.Workers = c.NbWorkers
	}
	if isUsingDeprecatedConfigParam(prefix, "scripts") {
		c.Hook.ScriptsDir = c.Scripts
	}
}

func isUsingDeprecatedConfigParam(prefix, flagName string) bool {
	envVar := helper.ToScreamingSnake(prefix + "_" + flagName)
	switch {
	case isFlagPassed(flagName):
		slog.Warn("using deprecated configuration flag", "flag", flagName)
		return true
	case isEnvExists(envVar):
		slog.Warn("using deprecated configuration environment variable", "variable", envVar)
		return true
	default:
		return false
	}
}

func isEnvExists(name string) bool {
	_, exists := os.LookupEnv(name)
	return exists
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
