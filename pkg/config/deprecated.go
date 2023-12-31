package config

import (
	"flag"
	"log/slog"
)

// OldConfig contain global configuration
type OldConfig struct {
	NbWorkers int    `flag:"nb-workers" desc:"Number of workers to start [DEPRECATED]" default:"2"`
	ScriptDir string `flag:"scripts" desc:"Scripts location [DEPRECATED]" default:"scripts"`
}

// ManageDeprecatedFlags manage configuration legacy
func (c *Config) ManageDeprecatedFlags() {
	// TODO check env variable
	// TODO other legacy parameters?
	// TODO code factorization
	if isFlagPassed("nb-workers") {
		slog.Warn("using deprecated configuration flag", "flag", "nb-workers")
		c.Hook.Workers = c.NbWorkers
	}
	if isFlagPassed("scripts") {
		slog.Warn("using deprecated configuration flag", "flag", "scripts")
		c.Hook.ScriptsDir = c.ScriptDir
	}
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
