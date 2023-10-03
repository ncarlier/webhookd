package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/ncarlier/webhookd/pkg/api"
	"github.com/ncarlier/webhookd/pkg/config"

	"golang.org/x/crypto/acme/autocert"
)

func cacheDir() (dir string) {
	if u, _ := user.Current(); u != nil {
		dir = filepath.Join(os.TempDir(), "webhookd-acme-cache-"+u.Username)
		if err := os.MkdirAll(dir, 0o700); err == nil {
			return dir
		}
	}
	return ""
}

// Server is a HTTP server wrapper used to manage TLS
type Server struct {
	self     *http.Server
	tls      bool
	certFile string
	keyFile  string
}

// ListenAndServe start HTTP(s) server
func (s *Server) ListenAndServe() error {
	if s.tls {
		return s.self.ListenAndServeTLS(s.certFile, s.keyFile)
	}
	return s.self.ListenAndServe()
}

// Shutdown stop HTTP(s) server
func (s *Server) Shutdown(ctx context.Context) error {
	s.self.SetKeepAlivesEnabled(false)
	return s.self.Shutdown(ctx)
}

// NewServer create new HTTP(s) server
func NewServer(cfg *config.Config) *Server {
	logger := slog.NewLogLogger(slog.Default().Handler(), slog.LevelError)
	server := &Server{
		tls: cfg.TLS,
		self: &http.Server{
			Addr:     cfg.ListenAddr,
			Handler:  api.NewRouter(cfg),
			ErrorLog: logger,
		},
	}
	if server.tls {
		// HTTPs server
		if cfg.TLSDomain == "" {
			server.certFile = cfg.TLSCertFile
			server.keyFile = cfg.TLSKeyFile
		} else {
			m := &autocert.Manager{
				Cache:      autocert.DirCache(cacheDir()),
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(cfg.TLSDomain),
			}
			server.self.TLSConfig = m.TLSConfig()
			server.certFile = ""
			server.keyFile = ""
		}
	}
	return server
}
