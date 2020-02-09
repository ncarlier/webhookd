package config

// Config contain global configuration
type Config struct {
	ListenAddr      string `flag:"listen-addr" desc:"HTTP listen address" default:":8080"`
	TLSListenAddr   string `flag:"tls-listen-addr" desc:"TLS listen address"`
	TLSCertFile     string `flag:"tls-cert-file" desc:"TLS certificate file" default:"server.pem"`
	TLSKeyFile      string `flag:"tls-key-file" desc:"TLS key file" default:"server.key"`
	TLSDomain       string `flag:"tls-domain" desc:"TLS domain name used by ACME"`
	NbWorkers       int    `flag:"nb-workers" desc:"Number of workers to start" default:"2"`
	Debug           bool   `flag:"debug" desc:"Output debug logs" default:"false"`
	Timeout         int    `flag:"timeout" desc:"Maximum hook execution time in second" default:"10"`
	ScriptDir       string `flag:"scripts" desc:"Scripts location" default:"scripts"`
	PasswdFile      string `flag:"passwd-file" desc:"Password file for basic HTTP authentication" default:".htpasswd"`
	LogDir          string `flag:"log-dir" desc:"Hook execution logs location" default:""`
	NotificationURI string `flag:"notification-uri" desc:"Notification URI"`
}
