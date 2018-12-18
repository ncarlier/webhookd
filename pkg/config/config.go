package config

import (
	"flag"
	"os"
	"strconv"
)

// Config contain global configuration
type Config struct {
	ListenAddr *string
	NbWorkers  *int
	Debug      *bool
	Timeout    *int
	ScriptDir  *string
	PasswdFile *string
}

var config = &Config{
	ListenAddr: flag.String("listen", getEnv("LISTEN_ADDR", ":8080"), "HTTP service address (e.g.address, ':8080')"),
	NbWorkers:  flag.Int("nb-workers", getIntEnv("NB_WORKERS", 2), "The number of workers to start"),
	Debug:      flag.Bool("debug", getBoolEnv("DEBUG", false), "Output debug logs"),
	Timeout:    flag.Int("timeout", getIntEnv("HOOK_TIMEOUT", 10), "Hook maximum delay before timeout (in second)"),
	ScriptDir:  flag.String("scripts", getEnv("SCRIPTS_DIR", "scripts"), "Scripts directory"),
	PasswdFile: flag.String("passwd", getEnv("PASSWD_FILE", ".htpasswd"), "Password file (encoded with htpasswd)"),
}

func init() {
	flag.StringVar(config.ListenAddr, "l", *config.ListenAddr, "HTTP service (e.g address: ':8080')")
	flag.BoolVar(config.Debug, "d", *config.Debug, "Output debug logs")
	flag.StringVar(config.PasswdFile, "p", *config.PasswdFile, "Password file (encoded with htpasswd)")
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
