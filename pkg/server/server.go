package server

import (
	"context"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/ncarlier/webhookd/pkg/api"
	"github.com/ncarlier/webhookd/pkg/config"
	"github.com/ncarlier/webhookd/pkg/logger"

	"golang.org/x/crypto/acme/autocert"
)

func cacheDir() (dir string) {
	if u, _ := user.Current(); u != nil {
		dir = filepath.Join(os.TempDir(), "webhookd-acme-cache-"+u.Username)
		if err := os.MkdirAll(dir, 0700); err == nil {
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
	} else {
		return s.self.ListenAndServe()
	}
}

// Shutdown stop HTTP(s) server
func (s *Server) Shutdown(ctx context.Context) error {
	s.self.SetKeepAlivesEnabled(false)
	return s.self.Shutdown(ctx)
}

// NewServer create new HTTP(s) server
func NewServer(cfg *config.Config) *Server {
	server := &Server{}
	if cfg.TLSListenAddr == "" {
		// Simple HTTP server
		server.self = &http.Server{
			Addr:     cfg.ListenAddr,
			Handler:  api.NewRouter(cfg),
			ErrorLog: logger.Error,
		}
		server.tls = false
	} else {
		// HTTPs server
		if cfg.TLSDomain == "" {
			server.self = &http.Server{
				Addr:     cfg.TLSListenAddr,
				Handler:  api.NewRouter(cfg),
				ErrorLog: logger.Error,
			}
			server.certFile = cfg.TLSCertFile
			server.keyFile = cfg.TLSKeyFile
		} else {
			m := &autocert.Manager{
				Cache:      autocert.DirCache(cacheDir()),
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(cfg.TLSDomain),
			}
			server.self = &http.Server{
				Addr:      cfg.TLSListenAddr,
				Handler:   api.NewRouter(cfg),
				ErrorLog:  logger.Error,
				TLSConfig: m.TLSConfig(),
			}
			server.certFile = ""
			server.keyFile = ""
		}
		server.tls = true
	}
	return server
}
