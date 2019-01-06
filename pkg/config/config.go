package config

import (
	"flag"
	"os"
	"strconv"
)

// Config contain global configuration
type Config struct {
	ListenAddr      *string
	NbWorkers       *int
	Debug           *bool
	Timeout         *int
	ScriptDir       *string
	PasswdFile      *string
	LogDir          *string
	NotificationURI *string
}

var config = &Config{
	ListenAddr:      flag.String("listen", getEnv("LISTEN_ADDR", ":8080"), "HTTP service address"),
	NbWorkers:       flag.Int("nb-workers", getIntEnv("NB_WORKERS", 2), "The number of workers to start"),
	Debug:           flag.Bool("debug", getBoolEnv("DEBUG", false), "Output debug logs"),
	Timeout:         flag.Int("timeout", getIntEnv("HOOK_TIMEOUT", 10), "Hook maximum delay (in second) before timeout"),
	ScriptDir:       flag.String("scripts", getEnv("SCRIPTS_DIR", "scripts"), "Scripts directory"),
	PasswdFile:      flag.String("passwd", getEnv("PASSWD_FILE", ".htpasswd"), "Password file encoded with htpasswd"),
	LogDir:          flag.String("log-dir", getEnv("LOG_DIR", os.TempDir()), "Webhooks execution log directory"),
	NotificationURI: flag.String("notification-uri", getEnv("NOTIFICATION_URI", ""), "Notification URI"),
}

func init() {
	// set shorthand parameters
	const shorthand = " (shorthand)"
	usage := flag.Lookup("listen").Usage + shorthand
	flag.StringVar(config.ListenAddr, "l", *config.ListenAddr, usage)
	usage = flag.Lookup("debug").Usage + shorthand
	flag.BoolVar(config.Debug, "d", *config.Debug, usage)
	usage = flag.Lookup("passwd").Usage + shorthand
	flag.StringVar(config.PasswdFile, "p", *config.PasswdFile, usage)
}

// Get global configuration
func Get() *Config {
	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv("APP_" + key); ok {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	strValue := getEnv(key, strconv.Itoa(fallback))
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	strValue := getEnv(key, strconv.FormatBool(fallback))
	if value, err := strconv.ParseBool(strValue); err == nil {
		return value
	}
	return fallback
}
