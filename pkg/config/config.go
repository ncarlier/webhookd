package config

import (
	"fmt"
	"regexp"
)

// Config store root configuration
type Config struct {
	ListenAddr     string             `flag:"listen-addr" desc:"HTTP listen address" default:":8080"`
	PasswdFile     string             `flag:"passwd-file" desc:"Password file for basic HTTP authentication" default:".htpasswd"`
	TruststoreFile string             `flag:"truststore-file" desc:"Truststore used by HTTP signature verifier (.pem or .p12)"`
	Hook           HookConfig         `flag:"hook"`
	Log            LogConfig          `flag:"log"`
	Notification   NotificationConfig `flag:"notification"`
	Static         StaticConfig       `flag:"static"`
	TLS            TLSConfig          `flag:"tls"`
	OldConfig      `flag:""`
}

// HookConfig store Hook execution configuration
type HookConfig struct {
	DefaultExt string `flag:"default-ext" desc:"Default extension for hook scripts" default:"sh"`
	Timeout    int    `flag:"timeout" desc:"Maximum hook execution time in second" default:"10"`
	ScriptsDir string `flag:"scripts" desc:"Scripts location" default:"scripts"`
	LogDir     string `flag:"log-dir" desc:"Hook execution logs location" default:""`
	Workers    int    `flag:"workers" desc:"Number of workers to start" default:"2"`
}

// LogConfig store logger configuration
type LogConfig struct {
	Level   string   `flag:"level" desc:"Log level (debug, info, warn or error)" default:"info"`
	Format  string   `flag:"format" desc:"Log format (json or text)" default:"text"`
	Modules []string `flag:"modules" desc:"Logging modules to activate (http,hook)" default:""`
}

// NotificationConfig store notification configuration
type NotificationConfig struct {
	URI string `flag:"uri" desc:"Notification URI"`
}

// StaticConfig store static assets configuration
type StaticConfig struct {
	Dir  string `flag:"dir" desc:"Static file directory to serve on /static path" default:""`
	Path string `flag:"path" desc:"Path to serve static file directory" default:"/static"`
}

// TLSConfig store TLS configuration
type TLSConfig struct {
	Enabled  bool   `flag:"enabled" desc:"Enable TLS" default:"false"`
	CertFile string `flag:"cert-file" desc:"TLS certificate file (unused if ACME used)" default:"server.pem"`
	KeyFile  string `flag:"key-file" desc:"TLS key file (unused if ACME used)" default:"server.key"`
	Domain   string `flag:"domain" desc:"TLS domain name used by ACME"`
}

// Validate the configuration
func (c *Config) Validate() error {
	if matched, _ := regexp.MatchString(`^/\w+$`, c.Static.Path); !matched {
		return fmt.Errorf("invalid static path: %s", c.Static.Path)
	}
	return nil
}
